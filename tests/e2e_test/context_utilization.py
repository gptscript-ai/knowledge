import numpy as np
import json
from typing import List, Union
from rag_eval_utils import llm_call



CONTEXT_PROMPT = """
Given question, answer and context, Verify if the context was useful in arriving at the given answer, if useful, return 1, if not, return 0. You should only return either 0 or 1.

QUESTION:
{question}

ANSWER:
{answer}

CONTEXT:
{context}
"""


def context_utilization(questions:List[str], answers: List[str], contexts: List[List[str]]):
    """
    Context utilization is a metric that assesses whether context chunks are utilized in answers. It is calculated using the question, answer, and contexts, with values ranging from 0 to 1, where higher scores represent greater precision.
    Note: replace answer with ground_truth, this metrics then becomes `context precision`.

    Args:
        questions (List[str]): _description_
        answers (List[str]): _description_
        contexts (List[List[str]]): _description_
    """
    assert len(questions) == len(answers)
    assert len(questions) == len(contexts)
    res = []
    for q, a, ctxs in zip(questions, answers, contexts):
        c_util_list = []
        for c in ctxs:
            c_useful = int(llm_call(CONTEXT_PROMPT.format(question=q, answer=a, context=c))) # this should return either 0 or 1?
            c_util_list.append(c_useful)
        
        score = sum(c_util_list) / len(c_util_list)
        res.append(score)
    return res



if __name__ == "__main__":
    
    questions= ['When was the first super bowl?', 'Who won the most super bowls?']
    answers = ['The first superbowl was held on Jan 15, 1967', 'The most super bowls have been won by The New England Patriots']
    contexts= [["green bay packers has 1 superbowl", "new england has won 6 superbowls", 'The First AFL–NFL World Championship Game was an American football game played on January 15, 1967, at the Los Angeles Memorial Coliseum in Los Angeles,'], 
    ['The Green Bay Packers...Green Bay, Wisconsin.','The Packers compete...Football Conference']]
    res = context_utilization(questions, answers, contexts)
    print(res)