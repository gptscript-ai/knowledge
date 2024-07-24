---
title: "knowledge get-dataset"
---
## knowledge get-dataset

Get a dataset

```
knowledge get-dataset <dataset-id> [flags]
```

### Options

```
      --archive string                    Path to the archive file ($KNOWLEDGE_CLIENT_GET_DATASET_ARCHIVE)
      --auto-migrate string               Auto migrate database ($KNOW_DB_AUTO_MIGRATE) (default "true")
  -c, --config-file string                Path to the configuration file ($KNOW_CONFIG_FILE)
      --dsn string                        Server database connection string (default "sqlite://$XDG_DATA_HOME/gptscript/knowledge/knowledge.db") ($KNOW_DB_DSN)
      --embedding-model-provider string   Embedding model provider ($KNOW_EMBEDDING_MODEL_PROVIDER) (default "openai")
  -h, --help                              help for get-dataset
      --no-docs                           Do not include documents in output (way less verbose) ($KNOWLEDGE_CLIENT_GET_DATASET_NO_DOCS)
      --server string                     URL of the Knowledge API Server ($KNOW_SERVER_URL)
      --vector-dbpath string              VectorDBPath to the vector database (default "$XDG_DATA_HOME/gptscript/knowledge/vector.db") ($KNOW_VECTOR_DB_PATH)
```

### SEE ALSO

* [knowledge](knowledge.md)	 - 

