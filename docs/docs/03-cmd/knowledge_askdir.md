---
title: "knowledge askdir"
---
## knowledge askdir

Retrieve sources for a query from a dataset generated from a directory

```
knowledge askdir [--path <path>] <query> [flags]
```

### Options

```
      --auto-migrate string              Auto migrate database ($KNOW_DB_AUTO_MIGRATE) (default "true")
  -c, --concurrency int                  Number of concurrent ingestion processes ($KNOW_INGEST_CONCURRENCY) (default 10)
      --dsn string                       Server database connection string (default "sqlite://$XDG_DATA_HOME/gptscript/knowledge/knowledge.db") ($KNOW_DB_DSN)
      --flow string                      Flow name ($KNOW_FLOW)
      --flows-file string                Path to a YAML/JSON file containing ingestion/retrieval flows ($KNOW_FLOWS_FILE)
  -h, --help                             help for askdir
      --ignore-extensions string         Comma-separated list of file extensions to ignore ($KNOW_INGEST_IGNORE_EXTENSIONS)
      --no-recursive                     Don't recursively ingest directories ($KNOW_NO_INGEST_RECURSIVE)
      --openai-api-base string           OpenAI API base ($OPENAI_BASE_URL) (default "https://api.openai.com/v1")
      --openai-api-key string            OpenAI API key (not required if used with clicky-chats) ($OPENAI_API_KEY) (default "sk-foo")
      --openai-api-type string           OpenAI API type (OPEN_AI, AZURE, AZURE_AD) ($OPENAI_API_TYPE) (default "OPEN_AI")
      --openai-api-version string        OpenAI API version (for Azure) ($OPENAI_API_VERSION) (default "2024-02-01")
      --openai-azure-deployment string   Azure OpenAI deployment name (overrides openai-embedding-model, if set) ($OPENAI_AZURE_DEPLOYMENT)
      --openai-embedding-model string    OpenAI Embedding model ($OPENAI_EMBEDDING_MODEL) (default "text-embedding-ada-002")
      --openai-model string              OpenAI model ($OPENAI_MODEL) (default "gpt-4")
  -p, --path string                      Path to the directory to query ($KNOWLEDGE_CLIENT_ASK_DIR_PATH) (default ".")
      --server string                    URL of the Knowledge API Server ($KNOW_SERVER_URL)
  -k, --top-k int                        Number of sources to retrieve ($KNOWLEDGE_CLIENT_ASK_DIR_TOP_K) (default 10)
      --vector-dbpath string             VectorDBPath to the vector database (default "$XDG_DATA_HOME/gptscript/knowledge/vector.db") ($KNOW_VECTOR_DB_PATH)
```

### SEE ALSO

* [knowledge](knowledge.md)	 - 

