---
title: Overview
---

# Ingestion and Retrieval Flows

Knowledge lets you configure how data is ingested and how it is retrieved back at querytime using flows.
Flows are a series of steps that can be configured via simple YAML files - so-called Flow Files or Flow Configs.

## Ingestion Flows

An Ingestion Flow consists of 3 main parts (all of them are optional and have basic defaults):

1. **Documentloader**: The loader defines how input data/files are parsed and transformed to "LLM-readable" text.
2. **Textsplitter**: The textsplitter defines how the text coming from the loader is split into smaller parts (documents).
3. **Transformers**: Transformers can be used to modify the documents coming out of the textsplitter. They can e.g. add metadata, remove irrelevant documents or generate summaries for every document.

All the documents yielded by that flow will be sent to the Embeddings Model (globally defined, not yet configurable via flow file) and then stored in the vector database.

## Retrieval Flows

A Retrieval Flow consists of 3 main parts (all of them are optional and have basic defaults):

1. **QueryModifiers**: QueryModifiers can be used to modify the input query before it is used for retrieval. They can even generate additional subqueries.
2. **Retriever**: The retriever defines how documents are retrieved from the vector database. E.g. one can use a routing retriever which automatically selects the most suitable dataset for the search.
3. **Postprocessors**: Postprocessors can be used to modify the retrieved documents. They can e.g. filter out irrelevant documents or sort the documents by relevance.

### Retrieval Flow Architecture

![Retrieval Flow Architecture](/img/retrieval_flows.png)

## Example Flow Configs

You can find some example flow files in our [GitHub repository](https://github.com/gptscript-ai/knowledge/tree/main/examples).