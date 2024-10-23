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
      --embedding-model-provider string   Embedding model provider ($KNOW_EMBEDDING_MODEL_PROVIDER) (default "openai")
  -h, --help                              help for export
      --index-dsn string                  Index Database Connection string (relational DB) (default "sqlite://$XDG_DATA_HOME/gptscript/knowledge/knowledge.db") ($KNOW_INDEX_DSN)
      --output string                     Output path ($KNOWLEDGE_CLIENT_EXPORT_DATASETS_OUTPUT) (default ".")
      --server string                     URL of the Knowledge API Server ($KNOW_SERVER_URL)
      --vector-dsn string                 DSN to the vector database (default "chromem:$XDG_DATA_HOME/gptscript/knowledge/vector.db") ($KNOW_VECTOR_DSN)
```

### SEE ALSO

* [knowledge](knowledge.md)	 - 

