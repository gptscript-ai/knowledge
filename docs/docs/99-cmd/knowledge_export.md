---
title: "knowledge export"
---
## knowledge export

Export one or more datasets as an archive (zip)

```
knowledge export <dataset-id> [<dataset-id>...] [flags]
```

### Options

```
  -a, --all                               Export all datasets ($KNOWLEDGE_CLIENT_EXPORT_DATASETS_ALL)
      --auto-migrate string               Auto migrate database ($KNOW_DB_AUTO_MIGRATE) (default "true")
  -c, --config-file string                Path to the configuration file ($KNOW_CONFIG_FILE)
      --dsn string                        Server database connection string (default "sqlite://$XDG_DATA_HOME/gptscript/knowledge/knowledge.db") ($KNOW_DB_DSN)
      --embedding-model-provider string   Embedding model provider ($KNOW_EMBEDDING_MODEL_PROVIDER) (default "openai")
  -h, --help                              help for export
      --output string                     Output path ($KNOWLEDGE_CLIENT_EXPORT_DATASETS_OUTPUT) (default ".")
      --server string                     URL of the Knowledge API Server ($KNOW_SERVER_URL)
      --vector-dbpath string              VectorDBPath to the vector database (default "$XDG_DATA_HOME/gptscript/knowledge/vector.db") ($KNOW_VECTOR_DB_PATH)
```

### SEE ALSO

* [knowledge](knowledge.md)	 - 

