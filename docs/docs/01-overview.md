---
title: Overview
slug: /
---

[![Discord](https://img.shields.io/discord/1204558420984864829?label=Discord)](https://discord.gg/9sSf4UyAMC)

<img alt="Knowledge Logo" src="img/logo.png" style={{width: 400}}/>

GPTScript **Knowledge** is a tool to enhance [GPTScript](https://github.com/gptscript-ai/gptscript) with Retrieval Augmented Generation (RAG).
It covers the whole pipeline to ingest data into a local or remote server and retrieve it back when queried.
It enables you to organize (especially non-public) data in knowledge bases (called datasets) that can even be shared with others.

Here's an overview of what the tool does:

1. **Ingest** data into a dataset
2. **Enhance** ingested data with metadata to boost retrieval accuracy
3. **Retrieve** data from one or more datasets to enrich LLM responses with factual information from your data sources
    - You may as well directly ask questions about a directory of files on your laptop
4. **Manage** multiple datasets to build distinct knowledge bases for different use cases
5. **Share** datasets with others to collaborate on knowledge bases and to distribute GPTScript tools with "built-in" knowledge

## Installation

Knowledge is distributed as a single binary for Linux, MacOS, and Windows.
You can download the latest release from the [GitHub releases page](https://github.com/gptscript-ai/knowledge/releases).

## Quickstart

We're going to use all the default settings and run Knowledge in standalone mode. For more information, see the [Usage](/usage) page.

The commands below are using files from this repository.

1. Export your OpenAI API key:

   ```bash
   export OPENAI_API_KEY=sk-foo # replace with your OpenAI API key
   ```

2. Ingest all files (recursively) in the `docs/docs/` directory:

   ```bash
   knowledge ingest docs/docs                       
   ```
   
   <details>
   <summary>Expected Output</summary>
   
   ```bash
   2024/07/10 17:46:34 INFO Created dataset id=default
   2024/07/10 17:46:35 INFO Ingested document filename=01-overview.md count=5 absolute_path=<path>/docs/docs/01-overview.md
   2024/07/10 17:46:35 INFO Ingested document filename=knowledge_import.md count=7 absolute_path=<path>/docs/docs/99-cmd/knowledge_import.md
   2024/07/10 17:46:35 INFO Ingested document filename=knowledge_ingest.md count=6 absolute_path=<path>/docs/docs/99-cmd/knowledge_ingest.md
   2024/07/10 17:46:35 INFO Ingested document filename=knowledge_create-dataset.md count=5 absolute_path=<path>/docs/docs/99-cmd/knowledge_create-dataset.md
   2024/07/10 17:46:35 INFO Ingested document filename=03-transformers.md count=1 absolute_path=<path>/docs/docs/11-flows/02-ingestion/03-transformers.md
   2024/07/10 17:46:35 INFO Ingested document filename=01-overview.md count=6 absolute_path=<path>/docs/docs/11-flows/01-overview.md
   2024/07/10 17:46:35 INFO Ingested document filename=knowledge_version.md count=4 absolute_path=<path>/docs/docs/99-cmd/knowledge_version.md
   2024/07/10 17:46:35 INFO Ingested document filename=knowledge_server.md count=5 absolute_path=<path>/docs/docs/99-cmd/knowledge_server.md
   2024/07/10 17:46:35 INFO Ingested document filename=knowledge_retrieve.md count=6 absolute_path=<path>/docs/docs/99-cmd/knowledge_retrieve.md
   2024/07/10 17:46:35 INFO Ingested document filename=02-usage.md count=7 absolute_path=<path>/docs/docs/02-usage.md
   2024/07/10 17:46:35 INFO Ingested document filename=01-querymodifiers.md count=1 absolute_path=<path>/docs/docs/11-flows/03-retrieval/01-querymodifiers.md
   2024/07/10 17:46:35 INFO Ingested document filename=knowledge_get-dataset.md count=5 absolute_path=<path>/docs/docs/99-cmd/knowledge_get-dataset.md
   2024/07/10 17:46:35 INFO Ingested document filename=02-textsplitters.md count=1 absolute_path=<path>/docs/docs/11-flows/02-ingestion/02-textsplitters.md
   2024/07/10 17:46:35 INFO Ingested document filename=01-documentloaders.md count=2 absolute_path=<path>/docs/docs/11-flows/02-ingestion/01-documentloaders.md
   2024/07/10 17:46:35 INFO Ingested document filename=knowledge_delete-dataset.md count=5 absolute_path=<path>/docs/docs/99-cmd/knowledge_delete-dataset.md
   2024/07/10 17:46:35 INFO Ingested document filename=knowledge_export.md count=5 absolute_path=<path>/docs/docs/99-cmd/knowledge_export.md
   2024/07/10 17:46:35 INFO Ingested document filename=03-architecture.md count=6 absolute_path=<path>/docs/docs/03-architecture.md
   2024/07/10 17:46:35 INFO Ingested document filename=02-retrievers.md count=1 absolute_path=<path>/docs/docs/11-flows/03-retrieval/02-retrievers.md
   2024/07/10 17:46:35 INFO Ingested document filename=knowledge_list-datasets.md count=5 absolute_path=<path>/docs/docs/99-cmd/knowledge_list-datasets.md
   2024/07/10 17:46:35 INFO Ingested document filename=knowledge.md count=4 absolute_path=<path>/docs/docs/99-cmd/knowledge.md
   2024/07/10 17:46:35 INFO Ingested document filename=knowledge_edit-dataset.md count=6 absolute_path=<path>/docs/docs/99-cmd/knowledge_edit-dataset.md
   2024/07/10 17:46:35 INFO Ingested document filename=knowledge_askdir.md count=6 absolute_path=<path>/docs/docs/99-cmd/knowledge_askdir.md
   2024/07/10 17:46:36 INFO Ingested document filename=03-postprocessors.md count=1 absolute_path=<path>/docs/docs/11-flows/03-retrieval/03-postprocessors.md
   2024/07/10 17:46:38 INFO Ingested document filename=_category_.json count=1 absolute_path=<path>/docs/docs/11-flows/03-retrieval/_category_.json
   2024/07/10 17:46:38 INFO Ingested document filename=_category_.json count=1 absolute_path=<path>/docs/docs/11-flows/02-ingestion/_category_.json
   2024/07/10 17:46:38 INFO Ingested document filename=_category_.json count=1 absolute_path=<path>/docs/docs/11-flows/_category_.json
   2024/07/10 17:46:38 INFO Ingested document filename=_category_.json count=1 absolute_path=<path>/docs/docs/10-datasets/_category_.json
   2024/07/10 17:46:38 INFO Ingested document filename=_category_.json count=1 absolute_path=<path>/docs/docs/99-cmd/_category_.json
   Ingested 28 files from "docs/docs" into dataset "default"
   ```
   </details>

3. Ask questions about the ingested data:

   ```bash
   knowledge retrieve "Which vector database is used by knowledge?"
   ```
   
    <details>
    <summary>Expected Output</summary>
    
    ```bash
   Retrieved the following 1 source collections for the original query "Which vector database is used by knowledge?": 
   {"Which vector database is used by knowledge?":[
    {"content":"# Using Knowledge\n## 2. Server Mode\nIn server mode, Knowledge uses a separate server for the Vector Database and the Document Database.\nThis mode is useful when you want to share the data with multiple clients or when you want to use a more powerful server for the Vector Database.","metadata":{"absPath":"<path>/docs/docs/02-usage.md","filename":"02-usage.md"},"similarity_score":0.8550703},
    {"content":"# Knowledge Architecture\n## 4. Vector Database\nThe vector database is the main storage for the embeddings of the ingested documents along with some metadata (e.g. source file information).\nThe current implementation uses [**chromem-go**](https://github.com/philippgille/chromem-go).\nIt's fully embedded and does not require any additional setup.","metadata":{"absPath":"<path>/docs/docs/03-architecture.md","filename":"03-architecture.md"},"similarity_score":0.85070354},
    {"content":"# Using Knowledge\n## 1. Standalone Mode (Default)\nIn standalone mode, Knowledge makes use of an embedded database and embedded Vector Database which the client connects to directly.\nThis is the default and most simple mode of operation and is useful for local usage and offers a great integration with GPTScript.","metadata":{"absPath":"<path>/docs/docs/02-usage.md","filename":"02-usage.md"},"similarity_score":0.8459537},
    {"content":"# Knowledge Architecture\n## 3. Index Database\nThe index database is an additional (relational) metadata database which keeps track of all datasets and ingested files and their relationships.\nIt enables some extra convenience features but does not store the actual data (embeddings).\nThe current implementation uses **SQLite**.\nIt's fully embedded and does not require any additional setup.","metadata":{"absPath":"<path>/docs/docs/03-architecture.md","filename":"03-architecture.md"},"similarity_score":0.8069764},
    {"content":"## knowledge edit-dataset\n### Options\n--update-metadata strings          update metadata key-value pairs (existing metadata will be updated/preserved) ($KNOWLEDGE_CLIENT_EDIT_DATASET_UPDATE_METADATA)\n      --vector-dbpath string             VectorDBPath to the vector database (default \"$XDG_DATA_HOME/gptscript/knowledge/vector.db\") ($KNOW_VECTOR_DB_PATH)\n```","metadata":{"absPath":"<path>/docs/docs/99-cmd/knowledge_edit-dataset.md","filename":"knowledge_edit-dataset.md"},"similarity_score":0.79661065},
    {"content":"# Knowledge Architecture\n![Knowledge Architecture](/img/knowledge_architecture.png)\nKnowledge consists of the following components:","metadata":{"absPath":"<path>/docs/docs/03-architecture.md","filename":"03-architecture.md"},"similarity_score":0.7884838},
    {"content":"## knowledge retrieve\n### Options\n--openai-model string              OpenAI model ($OPENAI_MODEL) (default \"gpt-4\")\n      --server string                    URL of the Knowledge API Server ($KNOW_SERVER_URL)\n  -k, --top-k int                        Number of sources to retrieve ($KNOWLEDGE_CLIENT_RETRIEVE_TOP_K) (default 10)\n      --vector-dbpath string             VectorDBPath to the vector database (default \"$XDG_DATA_HOME/gptscript/knowledge/vector.db\") ($KNOW_VECTOR_DB_PATH)\n```","metadata":{"absPath":"<path>/docs/docs/99-cmd/knowledge_retrieve.md","filename":"knowledge_retrieve.md"},"similarity_score":0.7874399},
    {"content":"# Knowledge Architecture\n## 1. Knowledge Client\nThe Knowledge Client is the main interface to interact with your knowledge bases.\nIn standalone mode, it makes direct use of embedded databases. It's running fully locally.\nIt's also the default entrypoint for the CLI.","metadata":{"absPath":"<path>/docs/docs/03-architecture.md","filename":"03-architecture.md"},"similarity_score":0.7814405},
    {"content":"## knowledge version\n### SEE ALSO\n- [knowledge](knowledge.md)\t -","metadata":{"absPath":"<path>/docs/docs/99-cmd/knowledge_version.md","filename":"knowledge_version.md"},"similarity_score":0.7807098},
    {"content":"# Ingestion and Retrieval Flows\nKnowledge lets you configure how data is ingested and how it is retrieved back at querytime using flows.\nFlows are a series of steps that can be configured via simple YAML files - so-called Flow Files or Flow Configs.","metadata":{"absPath":"<path>/docs/docs/11-flows/01-overview.md","filename":"01-overview.md"},"similarity_score":0.77821434}]}
    ```
   :::note

   The output is tailored for LLM readability, not necessarily meant to be read by humans.
   
   :::
    </details>

## Quickstart with [GPTScript](https://github.com/gptscript-ai/gptscript)

The default exported [GPTScript](https://github.com/gptscript-ai/gptscript) tool for knowledge is based on `knowledge askdir` (i.e. ask questions about contents of a directory).
The target directory defaults to your GPTScript workspace and can be configured using the `GPTSCRIPT_WORKSPACE` environment variable or the `--workspace` flag.

To mirror exactly what we did above, you can run the following command (assuming you have GPTScript installed):

```bash
gptscript --disable-cache --workspace docs/docs github.com/gptscript-ai/knowledge '{"query": "Which vector database is used by knowledge?"}'
```

<details>
<summary>Expected Output</summary>

```bash
18:17:41 started  [main] [input={"query": "Which vector database is used by knowledge?"}]
18:17:41 started  [context: github.com/gptscript-ai/context/workspace]
18:17:41 sent     [context: github.com/gptscript-ai/context/workspace]
18:17:41 ended    [context: github.com/gptscript-ai/context/workspace] [output=The current working directory is: \"<path>\"\nThe workspac...]
18:17:41 sent     [main]
2024/07/10 18:17:41 INFO Created dataset id=default
2024/07/10 18:17:41 INFO Created dataset id=fb0a2f2abeed1a53fb948b221bb968d901df310a
2024/07/10 18:17:42 INFO Ingested document filename=01-documentloaders.md count=2 absolute_path=<path>/docs/docs/11-flows/02-ingestion/01-documentloaders.md
2024/07/10 18:17:42 INFO Ingested document filename=_category_.json count=1 absolute_path=<path>/docs/docs/11-flows/03-retrieval/_category_.json
2024/07/10 18:17:42 INFO Ingested document filename=_category_.json count=1 absolute_path=<path>/docs/docs/10-datasets/_category_.json
2024/07/10 18:17:42 INFO Ingested document filename=knowledge_create-dataset.md count=5 absolute_path=<path>/docs/docs/99-cmd/knowledge_create-dataset.md
2024/07/10 18:17:42 INFO Ingested document filename=03-transformers.md count=1 absolute_path=<path>/docs/docs/11-flows/02-ingestion/03-transformers.md
2024/07/10 18:17:42 INFO Ingested document filename=01-overview.md count=6 absolute_path=<path>/docs/docs/01-overview.md
2024/07/10 18:17:42 INFO Ingested document filename=knowledge_import.md count=7 absolute_path=<path>/docs/docs/99-cmd/knowledge_import.md
2024/07/10 18:17:42 INFO Ingested document filename=02-textsplitters.md count=1 absolute_path=<path>/docs/docs/11-flows/02-ingestion/02-textsplitters.md
2024/07/10 18:17:42 INFO Ingested document filename=03-architecture.md count=6 absolute_path=<path>/docs/docs/03-architecture.md
2024/07/10 18:17:42 INFO Ingested document filename=knowledge_retrieve.md count=6 absolute_path=<path>/docs/docs/99-cmd/knowledge_retrieve.md
2024/07/10 18:17:42 INFO Ingested document filename=02-usage.md count=7 absolute_path=<path>/docs/docs/02-usage.md
2024/07/10 18:17:42 INFO Ingested document filename=knowledge_edit-dataset.md count=6 absolute_path=<path>/docs/docs/99-cmd/knowledge_edit-dataset.md
2024/07/10 18:17:42 INFO Ingested document filename=_category_.json count=1 absolute_path=<path>/docs/docs/11-flows/_category_.json
2024/07/10 18:17:42 INFO Ingested document filename=01-overview.md count=6 absolute_path=<path>/docs/docs/11-flows/01-overview.md
2024/07/10 18:17:42 INFO Ingested document filename=03-postprocessors.md count=1 absolute_path=<path>/docs/docs/11-flows/03-retrieval/03-postprocessors.md
2024/07/10 18:17:42 INFO Ingested document filename=knowledge_server.md count=5 absolute_path=<path>/docs/docs/99-cmd/knowledge_server.md
2024/07/10 18:17:43 INFO Ingested document filename=knowledge_delete-dataset.md count=5 absolute_path=<path>/docs/docs/99-cmd/knowledge_delete-dataset.md
2024/07/10 18:17:43 INFO Ingested document filename=knowledge_ingest.md count=6 absolute_path=<path>/docs/docs/99-cmd/knowledge_ingest.md
2024/07/10 18:17:43 INFO Ingested document filename=knowledge.md count=4 absolute_path=<path>/docs/docs/99-cmd/knowledge.md
2024/07/10 18:17:43 INFO Ingested document filename=knowledge_export.md count=5 absolute_path=<path>/docs/docs/99-cmd/knowledge_export.md
2024/07/10 18:17:43 INFO Ingested document filename=02-retrievers.md count=1 absolute_path=<path>/docs/docs/11-flows/03-retrieval/02-retrievers.md
2024/07/10 18:17:43 INFO Ingested document filename=knowledge_list-datasets.md count=5 absolute_path=<path>/docs/docs/99-cmd/knowledge_list-datasets.md
2024/07/10 18:17:43 INFO Ingested document filename=_category_.json count=1 absolute_path=<path>/docs/docs/99-cmd/_category_.json
2024/07/10 18:17:43 INFO Ingested document filename=01-querymodifiers.md count=1 absolute_path=<path>/docs/docs/11-flows/03-retrieval/01-querymodifiers.md
2024/07/10 18:17:43 INFO Ingested document filename=_category_.json count=1 absolute_path=<path>/docs/docs/11-flows/02-ingestion/_category_.json
2024/07/10 18:17:43 INFO Ingested document filename=knowledge_askdir.md count=6 absolute_path=<path>/docs/docs/99-cmd/knowledge_askdir.md
2024/07/10 18:17:43 INFO Ingested document filename=knowledge_version.md count=4 absolute_path=<path>/docs/docs/99-cmd/knowledge_version.md
2024/07/10 18:17:43 INFO Ingested document filename=knowledge_get-dataset.md count=5 absolute_path=<path>/docs/docs/99-cmd/knowledge_get-dataset.md
         content  [1] content | Retrieved the following 1 source collections for the query "Which vector database is used by knowledge?" from path "<path>/docs/docs": {"Which vector database is used by knowledge?":[{"content":"# Using Knowledge\n## 2. Server Mode\nIn server mode, Knowledge uses a separate server for the Vector Database and the Document Database.\nThis mode is useful when you want to share the data with multiple clients or when you want to use a more powerful server for the Vector Database.","metadata":{"absPath":"<path>/docs/docs/02-usage.md","filename":"02-usage.md"},"similarity_score":0.8550703},{"content":"# Knowledge Architecture\n## 4. Vector Database\nThe vector database is the main storage for the embeddings of the ingested documents along with some metadata (e.g. source file information).\nThe current implementation uses [**chromem-go**](https://github.com/philippgille/chromem-go).\nIt's fully embedded and does not require any additional setup.","metadata":{"absPath":"<path>/docs/docs/03-architecture.md","filename":"03-architecture.md"},"similarity_score":0.85083103},{"content":"# Using Knowledge\n## 1. Standalone Mode (Default)\nIn standalone mode, Knowledge makes use of an embedded database and embedded Vector Database which the client connects to directly.\nThis is the default and most simple mode of operation and is useful for local usage and offers a great integration with GPTScript.","metadata":{"absPath":"<path>/docs/docs/02-usage.md","filename":"02-usage.md"},"similarity_score":0.84598154},{"content":"# Knowledge Architecture\n## 3. Index Database\nThe index database is an additional (relational) metadata database which keeps track of all datasets and ingested files and their relationships.\nIt enables some extra convenience features but does not store the actual data (embeddings).\nThe current implementation uses **SQLite**.\nIt's fully embedded and does not require any additional setup.","metadata":{"absPath":"<path>/docs/docs/03-architecture.md","filename":"03-architecture.md"},"similarity_score":0.8069764},{"content":"## knowledge edit-dataset\n### Options\n--update-metadata strings          update metadata key-value pairs (existing metadata will be updated/preserved) ($KNOWLEDGE_CLIENT_EDIT_DATASET_UPDATE_METADATA)\n      --vector-dbpath string             VectorDBPath to the vector database (default \"$XDG_DATA_HOME/gptscript/knowledge/vector.db\") ($KNOW_VECTOR_DB_PATH)\n```","metadata":{"absPath":"<path>/docs/docs/99-cmd/knowledge_edit-dataset.md","filename":"knowledge_edit-dataset.md"},"similarity_score":0.79661065},{"content":"# Knowledge Architecture\n![Knowledge Architecture](/img/knowledge_architecture.png)\nKnowledge consists of the following components:","metadata":{"absPath":"<path>/docs/docs/03-architecture.md","filename":"03-architecture.md"},"similarity_score":0.78852946},{"content":"## knowledge retrieve\n### Options\n--openai-model string              OpenAI model ($OPENAI_MODEL) (default \"gpt-4\")\n      --server string                    URL of the Knowledge API Server ($KNOW_SERVER_URL)\n  -k, --top-k int                        Number of sources to retrieve ($KNOWLEDGE_CLIENT_RETRIEVE_TOP_K) (default 10)\n      --vector-dbpath string             VectorDBPath to the vector database (default \"$XDG_DATA_HOME/gptscript/knowledge/vector.db\") ($KNOW_VECTOR_DB_PATH)\n```","metadata":{"absPath":"<path>/docs/docs/99-cmd/knowledge_retrieve.md","filename":"knowledge_retrieve.md"},"similarity_score":0.7874399},{"content":"# Knowledge Architecture\n## 1. Knowledge Client\nThe Knowledge Client is the main interface to interact with your knowledge bases.\nIn standalone mode, it makes direct use of embedded databases. It's running fully locally.\nIt's also the default entrypoint for the CLI.","metadata":{"absPath":"<path>/docs/docs/03-architecture.md","filename":"03-architecture.md"},"similarity_score":0.78156495},{"content":"## knowledge version\n### SEE ALSO\n- [knowledge](knowledge.md)\t -","metadata":{"absPath":"<path>/docs/docs/99-cmd/knowledge_version.md","filename":"knowledge_version.md"},"similarity_score":0.7807098},{"content":"# Ingestion and Retrieval Flows\nKnowledge lets you configure how data is ingested and how it is retrieved back at querytime using flows.\nFlows are a series of steps that can be configured via simple YAML files - so-called Flow Files or Flow Configs.","metadata":{"absPath":"<path>/docs/docs/11-flows/01-overview.md","filename":"01-overview.md"},"similarity_score":0.77821434}]}
         content  [1] content | 
         content  [2] content | The current working directory is: "<path>"
         content  [2] content | The workspace directory is "<path>/docs/docs". Use "/ho ...
         content  [2] content | The workspace contains the following files and directories:
         content  [2] content |   01-overview.md
         content  [2] content |   02-usage.md
         content  [2] content |   03-architecture.md
         content  [2] content |   10-datasets/
         content  [2] content |   11-flows/
         content  [2] content |   99-cmd/
         content  [2] content |   docs/
         content  [2] content | Always use absolute paths to interact with files in the workspace
         content  [2] content | 
         content  [2] content | 
18:17:43 ended    [main] [output=Retrieved the following 1 source collections for the query \"Which vector database is used by knowled...]

INPUT:

{"query": "Which vector database is used by knowledge?"}

OUTPUT:

Retrieved the following 1 source collections for the query "Which vector database is used by knowledge?" from path "<path>/docs/docs": {"Which vector database is used by knowledge?":[{"content":"# Using Knowledge\n## 2. Server Mode\nIn server mode, Knowledge uses a separate server for the Vector Database and the Document Database.\nThis mode is useful when you want to share the data with multiple clients or when you want to use a more powerful server for the Vector Database.","metadata":{"absPath":"<path>/docs/docs/02-usage.md","filename":"02-usage.md"},"similarity_score":0.8550703},{"content":"# Knowledge Architecture\n## 4. Vector Database\nThe vector database is the main storage for the embeddings of the ingested documents along with some metadata (e.g. source file information).\nThe current implementation uses [**chromem-go**](https://github.com/philippgille/chromem-go).\nIt's fully embedded and does not require any additional setup.","metadata":{"absPath":"<path>/docs/docs/03-architecture.md","filename":"03-architecture.md"},"similarity_score":0.85083103},{"content":"# Using Knowledge\n## 1. Standalone Mode (Default)\nIn standalone mode, Knowledge makes use of an embedded database and embedded Vector Database which the client connects to directly.\nThis is the default and most simple mode of operation and is useful for local usage and offers a great integration with GPTScript.","metadata":{"absPath":"<path>/docs/docs/02-usage.md","filename":"02-usage.md"},"similarity_score":0.84598154},{"content":"# Knowledge Architecture\n## 3. Index Database\nThe index database is an additional (relational) metadata database which keeps track of all datasets and ingested files and their relationships.\nIt enables some extra convenience features but does not store the actual data (embeddings).\nThe current implementation uses **SQLite**.\nIt's fully embedded and does not require any additional setup.","metadata":{"absPath":"<path>/docs/docs/03-architecture.md","filename":"03-architecture.md"},"similarity_score":0.8069764},{"content":"## knowledge edit-dataset\n### Options\n--update-metadata strings          update metadata key-value pairs (existing metadata will be updated/preserved) ($KNOWLEDGE_CLIENT_EDIT_DATASET_UPDATE_METADATA)\n      --vector-dbpath string             VectorDBPath to the vector database (default \"$XDG_DATA_HOME/gptscript/knowledge/vector.db\") ($KNOW_VECTOR_DB_PATH)\n```","metadata":{"absPath":"<path>/docs/docs/99-cmd/knowledge_edit-dataset.md","filename":"knowledge_edit-dataset.md"},"similarity_score":0.79661065},{"content":"# Knowledge Architecture\n![Knowledge Architecture](/img/knowledge_architecture.png)\nKnowledge consists of the following components:","metadata":{"absPath":"<path>/docs/docs/03-architecture.md","filename":"03-architecture.md"},"similarity_score":0.78852946},{"content":"## knowledge retrieve\n### Options\n--openai-model string              OpenAI model ($OPENAI_MODEL) (default \"gpt-4\")\n      --server string                    URL of the Knowledge API Server ($KNOW_SERVER_URL)\n  -k, --top-k int                        Number of sources to retrieve ($KNOWLEDGE_CLIENT_RETRIEVE_TOP_K) (default 10)\n      --vector-dbpath string             VectorDBPath to the vector database (default \"$XDG_DATA_HOME/gptscript/knowledge/vector.db\") ($KNOW_VECTOR_DB_PATH)\n```","metadata":{"absPath":"<path>/docs/docs/99-cmd/knowledge_retrieve.md","filename":"knowledge_retrieve.md"},"similarity_score":0.7874399},{"content":"# Knowledge Architecture\n## 1. Knowledge Client\nThe Knowledge Client is the main interface to interact with your knowledge bases.\nIn standalone mode, it makes direct use of embedded databases. It's running fully locally.\nIt's also the default entrypoint for the CLI.","metadata":{"absPath":"<path>/docs/docs/03-architecture.md","filename":"03-architecture.md"},"similarity_score":0.78156495},{"content":"## knowledge version\n### SEE ALSO\n- [knowledge](knowledge.md)\t -","metadata":{"absPath":"<path>/docs/docs/99-cmd/knowledge_version.md","filename":"knowledge_version.md"},"similarity_score":0.7807098},{"content":"# Ingestion and Retrieval Flows\nKnowledge lets you configure how data is ingested and how it is retrieved back at querytime using flows.\nFlows are a series of steps that can be configured via simple YAML files - so-called Flow Files or Flow Configs.","metadata":{"absPath":"<path>/docs/docs/11-flows/01-overview.md","filename":"01-overview.md"},"similarity_score":0.77821434}]}
```
</details>

### Human-readable with GPTScript

The above didn't really look nice - It just returned what the Knowledge tool yielded.
Now let's actually use that as part of a GPTScript such that the Knowledge tool output will be used by the LLM to generate a pretty answer for us.

1. Create a GPTScript script that leverages the knowledge tool (see [examples/quickstart.gpt](examples/quickstart.gpt)):
    
   ```yaml
   # examples/quickstart.gpt
   
   tools: github.com/gptscript-ai/knowledge
   
   Answer any question using context retrieved from knowledge.
   ```
2. Run the GPTScript:

   ```bash
   gptscript --workspace docs/docs examples/quickstart.gpt "Which vector database is used by knowledge?"
   ```
   
   :::note

   The final output of this should be: `The vector database used by Knowledge is [**chromem-go**](https://github.com/philippgille/chromem-go).`

   :::
   
   <details>
   <summary>Expected Output</summary>
   
   ```bash
   18:23:42 started  [main] [input=Which vector database is used by knowledge?]
   18:23:42 sent     [main]
            content  [1] content | Waiting for model response...
            content  [1] content | <tool call> knowledge -> {"query":"vector database used by knowledge"}
   18:23:43 started  [https://raw.githubusercontent.com/gptscript-ai/knowledge/215c1002be8d793944648b02ad450f91049e83d5/tool.gpt(2)] [input={"query":"vector database used by knowledge"}]
   18:23:43 started  [context: github.com/gptscript-ai/context/workspace]
   18:23:43 sent     [context: github.com/gptscript-ai/context/workspace]
   18:23:43 ended    [context: github.com/gptscript-ai/context/workspace] [output=The current working directory is: \"<path>\"\nThe workspac...]
   18:23:43 sent     [https://raw.githubusercontent.com/gptscript-ai/knowledge/215c1002be8d793944648b02ad450f91049e83d5/tool.gpt(2)]
   2024/07/10 18:23:43 INFO Ignoring duplicate document filename=03-architecture.md absolute_path=<path>/docs/docs/03-architecture.md
   2024/07/10 18:23:43 INFO Ignoring duplicate document filename=_category_.json absolute_path=<path>/docs/docs/10-datasets/_category_.json
   2024/07/10 18:23:43 INFO Ignoring duplicate document filename=02-usage.md absolute_path=<path>/docs/docs/02-usage.md
   2024/07/10 18:23:43 INFO Ignoring duplicate document filename=01-overview.md absolute_path=<path>/docs/docs/11-flows/01-overview.md
   2024/07/10 18:23:43 INFO Ignoring duplicate document filename=01-querymodifiers.md absolute_path=<path>/docs/docs/11-flows/03-retrieval/01-querymodifiers.md
   2024/07/10 18:23:43 INFO Ignoring duplicate document filename=knowledge.md absolute_path=<path>/docs/docs/99-cmd/knowledge.md
   2024/07/10 18:23:43 INFO Ignoring duplicate document filename=_category_.json absolute_path=<path>/docs/docs/11-flows/03-retrieval/_category_.json
   2024/07/10 18:23:43 INFO Ignoring duplicate document filename=01-documentloaders.md absolute_path=<path>/docs/docs/11-flows/02-ingestion/01-documentloaders.md
   2024/07/10 18:23:43 INFO Ignoring duplicate document filename=knowledge_delete-dataset.md absolute_path=<path>/docs/docs/99-cmd/knowledge_delete-dataset.md
   2024/07/10 18:23:43 INFO Ignoring duplicate document filename=knowledge_server.md absolute_path=<path>/docs/docs/99-cmd/knowledge_server.md
   2024/07/10 18:23:43 INFO Ignoring duplicate document filename=knowledge_retrieve.md absolute_path=<path>/docs/docs/99-cmd/knowledge_retrieve.md
   2024/07/10 18:23:43 INFO Ignoring duplicate document filename=knowledge_list-datasets.md absolute_path=<path>/docs/docs/99-cmd/knowledge_list-datasets.md
   2024/07/10 18:23:43 INFO Ignoring duplicate document filename=02-textsplitters.md absolute_path=<path>/docs/docs/11-flows/02-ingestion/02-textsplitters.md
   2024/07/10 18:23:43 INFO Ignoring duplicate document filename=knowledge_version.md absolute_path=<path>/docs/docs/99-cmd/knowledge_version.md
   2024/07/10 18:23:43 INFO Ignoring duplicate document filename=_category_.json absolute_path=<path>/docs/docs/11-flows/02-ingestion/_category_.json
   2024/07/10 18:23:43 INFO Ignoring duplicate document filename=03-transformers.md absolute_path=<path>/docs/docs/11-flows/02-ingestion/03-transformers.md
   2024/07/10 18:23:43 INFO Ignoring duplicate document filename=knowledge_create-dataset.md absolute_path=<path>/docs/docs/99-cmd/knowledge_create-dataset.md
   2024/07/10 18:23:43 INFO Ignoring duplicate document filename=knowledge_edit-dataset.md absolute_path=<path>/docs/docs/99-cmd/knowledge_edit-dataset.md
   2024/07/10 18:23:43 INFO Ignoring duplicate document filename=03-postprocessors.md absolute_path=<path>/docs/docs/11-flows/03-retrieval/03-postprocessors.md
   2024/07/10 18:23:43 INFO Ignoring duplicate document filename=knowledge_import.md absolute_path=<path>/docs/docs/99-cmd/knowledge_import.md
   2024/07/10 18:23:43 INFO Ignoring duplicate document filename=_category_.json absolute_path=<path>/docs/docs/11-flows/_category_.json
   2024/07/10 18:23:43 INFO Ignoring duplicate document filename=knowledge_ingest.md absolute_path=<path>/docs/docs/99-cmd/knowledge_ingest.md
   2024/07/10 18:23:43 INFO Ignoring duplicate document filename=knowledge_export.md absolute_path=<path>/docs/docs/99-cmd/knowledge_export.md
   2024/07/10 18:23:43 INFO Ignoring duplicate document filename=knowledge_get-dataset.md absolute_path=<path>/docs/docs/99-cmd/knowledge_get-dataset.md
   2024/07/10 18:23:43 INFO Ignoring duplicate document filename=02-retrievers.md absolute_path=<path>/docs/docs/11-flows/03-retrieval/02-retrievers.md
   2024/07/10 18:23:43 INFO Ignoring duplicate document filename=knowledge_askdir.md absolute_path=<path>/docs/docs/99-cmd/knowledge_askdir.md
   2024/07/10 18:23:43 INFO Ignoring duplicate document filename=_category_.json absolute_path=<path>/docs/docs/99-cmd/_category_.json
   2024/07/10 18:23:44 INFO Ingested document filename=01-overview.md count=29 absolute_path=<path>/docs/docs/01-overview.md
            content  [2] content | Retrieved the following 1 source collections for the query "vector database used by knowledge" from  ...
            content  [2] content | 
            content  [3] content | The current working directory is: "<path>"
            content  [3] content | The workspace directory is "<path>/docs/docs". Use "/ho ...
            content  [3] content | The workspace contains the following files and directories:
            content  [3] content |   01-overview.md
            content  [3] content |   02-usage.md
            content  [3] content |   03-architecture.md
            content  [3] content |   10-datasets/
            content  [3] content |   11-flows/
            content  [3] content |   99-cmd/
            content  [3] content |   docs/
            content  [3] content | Always use absolute paths to interact with files in the workspace
            content  [3] content | 
            content  [3] content | 
   18:23:44 ended    [https://raw.githubusercontent.com/gptscript-ai/knowledge/215c1002be8d793944648b02ad450f91049e83d5/tool.gpt(2)] [output=Retrieved the following 1 source collections for the query \"vector database used by knowledge\" from...]
   18:23:44 continue [main]
   18:23:44 sent     [main]
            content  [1] content | Waiting for model response...
            content  [1] content | The vector database used by Knowledge is [**chromem-go**](https://github.com/philippgille/chromem-go).
   18:23:45 ended    [main] [output=The vector database used by Knowledge is [**chromem-go**](https://github.com/philippgille/chromem-go...]
   18:23:45 usage    [total=2194] [prompt=2147] [completion=47]
   
   INPUT:
   
   Which vector database is used by knowledge?
   
   OUTPUT:
   
   The vector database used by Knowledge is [**chromem-go**](https://github.com/philippgille/chromem-go).
   ```
   
   </details>