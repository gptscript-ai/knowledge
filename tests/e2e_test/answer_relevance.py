import numpy as np
import json
from typing import List, Union
from rag_eval_utils import get_embedding, llm_call


ANSWER_RELEVANCY_PROMPT = """Generate a potential question for the given answer, and identify if answer is noncommittal. Give noncommittal as 1 if the answer is noncommittal and 0 if the answer is committal. A noncommittal answer is one that is evasive, vague, or ambiguous. For example, "I don't know" or "I'm not sure" are noncommittal answers.
Generate your output in json format strictly, includes `question` and `noncommittal` field, don't include ```json.

THE GIVEN ANSWER:
{given_answer}
"""


def calculate_similarity(question: str, generated_questions: List[str]):
    """This is a 1VN similarity calculate, question is a string while generated questions is a list of N items.
    """
    question_vec = np.asarray(get_embedding(question))
    gen_question_vec = np.asarray(
        get_embedding(generated_questions)
    )
    norm = np.linalg.norm(gen_question_vec, axis=1) * np.linalg.norm(
        question_vec, axis=1
    )
    return (
        np.dot(gen_question_vec, question_vec.T).reshape(
            -1,
        )
        / norm
    )

def answer_relevancy(questions:List[str], answers: List[str]):
    """
    The Answer Relevancy metric evaluates how closely the generated answer aligns with the given question.
    Lower scores are given to answers that are incomplete or contain unnecessary information, while higher scores reflect greater relevancy.
    This metric is calculated based on both the `question` and the `answer`.
    
    Equation:
    similarity_score.mean() * (1-noncommittal)
    
    The part before * is the mean cosine similarity of the original question to an (number of) artifical questions, which where generated (reverse engineered) based on the answer
    the part after is a binary weight, either 0 or 1, For example `I don't know` is a noncommittal answer, so `(1-noncommittal)` will be 0
    
    Required Args:
    question: the origin question
    answer: the answer generated by llm
    """
    assert len(questions) == len(answers)
    res = []
    for q,a in zip(questions, answers):
        _generated_question = llm_call(ANSWER_RELEVANCY_PROMPT.format(given_answer=a))
        try:
            print(_generated_question)
            _generated_question = json.loads(_generated_question)
            generated_question = _generated_question["question"]
            noncommittal = int(_generated_question["noncommittal"])
            
        # TODO: it might be a good idea to generate more than 1 potential question
            
            similarity_score = calculate_similarity(q, generated_question)
            weighted_score = similarity_score.mean() * (1-noncommittal)
        except Exception as e:
            print(f"Error parsing generated questions : {e}")
            weighted_score = 0
        
        res.append(weighted_score)
    return res    
        
    
if __name__ == "__main__":
    
    questions = ["who win the first superbowl?", 
           "Who are the main characters in Reunion under the Stars story", 
           "What is the main story line for Reunion under the stars story",
           ]
    answers = ["The first superbowl was held on Jan 15, 1967", 
           """The main characters in "Reunion Under the Stars" are Aarav and Priya. Aarav is a scientist who moved to the United States, and Priya is a famous artist. They both hail from India and share a deep, enduring friendship that remains strong despite the distance and time apart.""",
    """The main storyline of "Reunion Under the Stars" revolves around the heartwarming reunion of two childhood friends, Aarav and Priya, who both hail from India. Aarav has moved to the United States to become a scientist, while Priya has become a famous artist. Despite the years and distance that have separated them, their friendship remains strong.The story begins with Aarav and Priya meeting again in a bustling city park in the U.S. They share a joyful and emotional reunion, reminiscing about their childhood in India and catching up on each other's lives. They explore the city together, enjoying Indian food and sharing stories about their respective journeys and achievements.As the day turns into evening, they sit under the stars and make a promise to maintain their friendship despite the distance. They vow to meet more often and keep the essence of their homeland alive in their hearts. The story highlights the enduring bond of friendship that transcends time and distance, celebrating the magic of their enduring connection.""",
           ]
    res = answer_relevancy(questions, answers)
    print(res)