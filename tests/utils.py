from gptscript.command import exec
from gptscript.tool import Tool
import json
from openai import OpenAI, Stream
from openai.types.chat import ChatCompletion, ChatCompletionChunk
import time
import os


def run_test(question, answer, dataset, judge_client):
    tools = [
        Tool(
            tools=["retrieve"],
            instructions="""

                    {}
                    """.format(question),
        ),
        Tool(
            name="retrieve",
            description="Retrieve information from a Knowledge Base Dataset",
            args={"query": "Query to be executed against the Knowledge Base Dataset"},
            instructions="""

                #!knowledge retrieve -d {} -k 20 ${{query}}
                """.format(dataset)),
    ]

    os.environ['GPTSCRIPT_INTERNAL_SYSTEM_PROMPT'] = ('You are an expert in understanding context and extracting information')
    resp = exec(tools)

    if isinstance(resp, (list, tuple)):
        returned_answer = ''.join(resp)
    else:
        returned_answer = resp
    returned_answer = returned_answer.replace('\n', '')

    match, reason = judge_final_answer(judge_client, False, returned_answer, answer)
    assert match is True, """
        Question: {},
        Returned answer: {},
        Expected answer: {},
        Reason: {}
        """.format(question, returned_answer, answer, reason)


def judge_final_answer(judge_client: OpenAI, stream: bool, final_answer: str, final_answer_should: str) -> (bool, str):
    if not final_answer_should:
        return None, True

    judge_response = judge_client.chat.completions.create(
        model='gpt-4-turbo-preview',
        response_format={
            "type": "json_object",
        },
        messages=[{
            "role": "system",
            "content": """When given JSON objects that conform to the following JSONSchema:
{
    "name": "judge",
    "type": "object",
    "properties": {
        "final_answer": {
            "type": "string",
            "description": "An answer to judge for correctness."
        },
        "final_answer_should": {
            "type": "string",
            "description": "The constraints that final_answer must completely satisfy to be considered correct."
        }
    },
    "required": [
        "final_answer",
        "final_answer_should"
    ]
}

Determine if `final_answer` satisfies the constraints described by `final_answer_should`.
`final_answer` is considered correct if and only if it satisfies the constraints described by `final_answer_should`.
If `final_answer_should` mentioned `or` condition, then if `final_answer` meet one of the condition it should be considered correct.

After making a determination, respond with a JSON object that conforms to the following JSONSchema:

{
    "name": "ruling",
    "type": "object",
    "properties": {
        "correct": {
            "type": "boolean",
            "description": "Set to true if and only if the answer is considered correct."
        },
        "reasoning": {
            "type": "string",
            "description": "A brief summary of the reasoning used to come to the determination."
        }
    },
    "required": [
        "correct",
        "reasoning"
    ]
}

Your responses are concise and include only the json object described above.
"""
        }, {
            "role": "user",
            "content": json.dumps({
                "final_answer": final_answer,
                "final_answer_should": final_answer_should,
            })
        }],
        stream=stream

    )

    judge_completion = to_chat_completion(judge_response)
    judge_message = judge_completion.choices[0].message.content
    judge_ruling = json.loads(judge_message)

    try:
        return judge_ruling['correct'], judge_ruling['reasoning']
    except KeyError as e:
        raise ValueError(f"Failed to judge final answer. Judge response missing key: {e}")


def to_chat_completion(response: ChatCompletion | Stream[ChatCompletionChunk]) -> ChatCompletion:
    if isinstance(response, ChatCompletion):
        return response

    id = ""
    finish_reason = ""
    model = ""
    system_fingerprint = ""
    choices = {}
    tool_calls = {}

    for chunk in response:
        if chunk is None:
            continue

        id = id or chunk.id
        model = model or chunk.model
        system_fingerprint = system_fingerprint or chunk.system_fingerprint

        for choice in chunk.choices or []:
            if choice is None or choice.delta is None:
                continue

            if choice.finish_reason is not None and finish_reason != "":
                finish_reason = choice.finish_reason

            if choice.index not in choices:
                choices[choice.index] = {
                    "index": choice.index,
                    "message": {
                        "role": choice.delta.role or "",
                        "content": choice.delta.content or "",
                    },
                    "finish_reason": choice.finish_reason
                }
            else:
                choices[choice.index]["message"]["content"] += choice.delta.content or ""
                choices[choice.index]["finish_reason"] = choice.finish_reason

            if choice.delta.tool_calls is None:
                continue

            for response_call in choice.delta.tool_calls:
                if choice.index not in tool_calls:
                    tool_calls[choice.index] = {}

                arguments = ""
                name = ""
                if response_call.function is not None:
                    arguments = response_call.function.arguments or ""
                    name = response_call.function.name or ""

                if response_call.index not in tool_calls[choice.index]:
                    tool_calls[choice.index][response_call.index] = {
                        "id": response_call.id,
                        "type": "function",
                        "function": {
                            "arguments": arguments,
                            "name": name,
                        }
                    }
                else:
                    tool_calls[choice.index][response_call.index]["function"]["arguments"] += arguments

    for index, choice in choices.items():
        if index not in tool_calls:
            continue

        if "message" not in choices[index]:
            continue

        choices[index]["message"]["tool_calls"] = [tool_calls[index][call_index] for call_index in
                                                   sorted(tool_calls[index].keys())]

    return ChatCompletion(
        id=id,
        created=int(time.time()),
        object="chat.completion",
        model=model,
        choices=[choices[index] for index in sorted(choices.keys())] if choices is not None else [],
        system_fingerprint=system_fingerprint
    )