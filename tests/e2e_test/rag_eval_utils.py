import openai
from openai import OpenAI
from typing import Union, List
client = OpenAI()

def llm_call(sys_prompt, response_format=None):
    completion = client.chat.completions.create(
    model="gpt-4o",
    messages=[
        {"role": "system", "content": sys_prompt},
    ],
    temperature=0.1,
    )

    return completion.choices[0].message.content

def get_embedding(text_inputs : Union[List[str], str]):
    if isinstance(text_inputs, str):
        text_inputs = [text_inputs]
    response = client.embeddings.create(
    input=text_inputs,
    model="text-embedding-3-small"
    )

    return [d.embedding for d in response.data]