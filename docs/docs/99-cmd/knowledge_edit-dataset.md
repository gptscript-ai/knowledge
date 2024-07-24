---
title: "knowledge edit-dataset"
---
## knowledge edit-dataset

Edit an existing dataset

```
knowledge edit-dataset <dataset-id> [flags]
```

### Options

```
      --auto-migrate string               Auto migrate database ($KNOW_DB_AUTO_MIGRATE) (default "true")
  -c, --config-file string                Path to the configuration file ($KNOW_CONFIG_FILE)
      --dsn string                        Server database connection string (default "sqlite://$XDG_DATA_HOME/gptscript/knowledge/knowledge.db") ($KNOW_DB_DSN)
      --embedding-model-provider string   Embedding model provider ($KNOW_EMBEDDING_MODEL_PROVIDER) (default "openai")
  -h, --help                              help for edit-dataset
      --replace-metadata strings          replace metadata with key-value pairs (existing metadata will be removed) ($KNOWLEDGE_CLIENT_EDIT_DATASET_REPLACE_METADATA)
      --reset-metadata                    reset metadata to default (empty) ($KNOWLEDGE_CLIENT_EDIT_DATASET_RESET_METADATA)
      --server string                     URL of the Knowledge API Server ($KNOW_SERVER_URL)
      --update-metadata strings           update metadata key-value pairs (existing metadata will be updated/preserved) ($KNOWLEDGE_CLIENT_EDIT_DATASET_UPDATE_METADATA)
      --vector-dbpath string              VectorDBPath to the vector database (default "$XDG_DATA_HOME/gptscript/knowledge/vector.db") ($KNOW_VECTOR_DB_PATH)
```

### SEE ALSO

* [knowledge](knowledge.md)	 - 

