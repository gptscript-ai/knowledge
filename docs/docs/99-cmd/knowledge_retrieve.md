---
title: "knowledge retrieve"
---
## knowledge retrieve

Retrieve sources for a query from a dataset

```
knowledge retrieve [--dataset <dataset-id>] <query> [flags]
```

### Options

```
      --archive string                    Path to the archive file ($KNOWLEDGE_CLIENT_RETRIEVE_ARCHIVE)
      --auto-migrate string               Auto migrate database ($KNOW_DB_AUTO_MIGRATE) (default "true")
  -c, --config-file string                Path to the configuration file ($KNOW_CONFIG_FILE)
  -d, --dataset string                    Target Dataset ID ($KNOW_TARGET_DATASET) (default "default")
      --dsn string                        Server database connection string (default "sqlite://$XDG_DATA_HOME/gptscript/knowledge/knowledge.db") ($KNOW_DB_DSN)
      --embedding-model-provider string   Embedding model provider ($KNOW_EMBEDDING_MODEL_PROVIDER) (default "openai")
      --flow string                       Flow name ($KNOW_FLOW)
      --flows-file string                 Path to a YAML/JSON file containing ingestion/retrieval flows ($KNOW_FLOWS_FILE)
  -h, --help                              help for retrieve
      --server string                     URL of the Knowledge API Server ($KNOW_SERVER_URL)
  -k, --top-k int                         Number of sources to retrieve ($KNOWLEDGE_CLIENT_RETRIEVE_TOP_K) (default 10)
      --vector-dbpath string              VectorDBPath to the vector database (default "$XDG_DATA_HOME/gptscript/knowledge/vector.db") ($KNOW_VECTOR_DB_PATH)
```

### SEE ALSO

* [knowledge](knowledge.md)	 - 

