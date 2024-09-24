from answer_factuality import answer_factuality
from answer_faithfulness import answer_faithfulness
from answer_relevance import answer_relevancy

def main(input_dict: dict):
    try:
        questions = input_dict["questions"]
        answers = input_dict["answers"]
        contexts = input_dict["contexts"]
        ground_truths = input_dict["ground_truths"]
    except Exception as e:
        print(e)
    factuality_score = answer_factuality(answers=answers, ground_truths=ground_truths)
    faithfulness_score = answer_faithfulness(answers=answers, contexts=contexts)
    relevance_score = answer_relevancy(answers=answers, questions=questions)
    print(factuality_score)
    print(faithfulness_score)
    print(relevance_score)


if __name__ == "__main__":
    pass
    # main()