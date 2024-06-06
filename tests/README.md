# E2E benchmark testing 

This folder contains all the tests data and test cases for RAG benchmark testing. 

The testing data is from https://github.com/h2oai/enterprise-h2ogpte/blob/main/rag_benchmark/e2e_df.csv.

To run test:

```commandline
python -m venv venv
pip install -r requiremens.txt
pytest test.py
```