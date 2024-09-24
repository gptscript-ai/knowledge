import numpy as np
import json
from typing import List, Union
from rag_eval_utils import llm_call


FAITHFULNESS_PROMPT = """
Given an answer and a context, your task is to:
1. break the answer down to a series of factual statements
2. judge the faithfulness of these statements based on the given context. For each statement you must return verdict as 1 if the statement can be directly inferred based on the context or 0 if the statement can not be directly inferred based on the context.
Return your output as a list of json dict. in each dict contains 2 field: `statement` field with the factual statement, and `verdict` field with value either 0 or 1, integer. Do not include '```json'.

ANSWER:
{answer}

CONTEXT:
{context}
"""


def answer_faithfulness( answers: List[str], contexts: List[List[str]]):
    """
    This metric assesses the factual consistency of the generated answer with respect to the provided context.
    It is computed by comparing the answer to the retrieved context and is scaled between 0 and 1, with higher values indicating better consistency.
    This function breaks each answer down to a list of factual statements, each statement is considered verdict if it can be directly inferred based on the context, with value either 0 or 1.
    Finally the faithfulness K-precision of the answer is calculated by: sum(verdict of each statement) / (num_of_statements_in_the_answer)

    Args:
        answers (List[str]): _description_
        contexts (List[List[str]]): _description_
    """
    assert len(answers) == len(contexts)
    res = []
    for a, ctxs in zip( answers, contexts):
        c = "\n".join(ctxs)
        c_faithfulness = llm_call(FAITHFULNESS_PROMPT.format( answer=a, context=c))
        try:
            list_of_faithfulness = json.loads(c_faithfulness)
            c_verdicts = []
            for l in list_of_faithfulness:
                c_verdicts.append(l["verdict"])
            score = sum(c_verdicts) / len(c_verdicts) if len(c_verdicts) > 0 else 0
            res.append(score)
        except Exception as e:
            print(f"fail to parse faithfulness, Error: {e}")
        
    return res



if __name__ == "__main__":
    
    questions= ['When was the first super bowl?', 'Who won the most super bowls?']
    answers = ['The first superbowl was held on Jan 15, 1967', 'The most super bowls have been won by The New England Patriots']
    contexts= [["green bay packers has 1 superbowl", "new england has won 6 superbowls", 'The First AFL–NFL World Championship Game was an American football game played on January 15, 1967, at the Los Angeles Memorial Coliseum in Los Angeles,'], 
    ['The Green Bay Packers...Green Bay, Wisconsin.','The Packers compete...Football Conference']]
    res = answer_faithfulness(answers, contexts)
    
    print(res)