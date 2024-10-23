---
title: "knowledge delete-file"
---
## knowledge delete-file

Delete a file from a dataset

```
knowledge delete-file <file-id|file-abs-path> [flags]
```

### Options

```
      --auto-migrate string               Auto migrate database ($KNOW_DB_AUTO_MIGRATE) (default "true")
  -c, --config-file string                Path to the configuration file ($KNOW_CONFIG_FILE)
  -d, --dataset string                    Target Dataset ID ($KNOWLEDGE_CLIENT_DELETE_FILE_DATASET) (default "default")
      --embedding-model-provider string   Embedding model provider ($KNOW_EMBEDDING_MODEL_PROVIDER) (default "openai")
  -h, --help                              help for delete-file
      --index-dsn string                  Index Database Connection string (relational DB) (default "sqlite://$XDG_DATA_HOME/gptscript/knowledge/knowledge.db") ($KNOW_INDEX_DSN)
      --server string                     URL of the Knowledge API Server ($KNOW_SERVER_URL)
      --vector-dsn string                 DSN to the vector database (default "chromem:$XDG_DATA_HOME/gptscript/knowledge/vector.db") ($KNOW_VECTOR_DSN)
```

### SEE ALSO

* [knowledge](knowledge.md)	 - 

