---
title: "knowledge import"
---
## knowledge import

Import one or more datasets from an archive (zip) (default: all datasets)

### Synopsis

Import one or more datasets from an archive (zip) (default: all datasets).
## IMPORTANT: Embedding functions
   When someone first ingests some data into a dataset, the embedding provider configured at that time will be attached to the dataset.
   Upon subsequent ingestion actions, the same embedding provider must be used to ensure that the embeddings are consistent.
   Most of the times, the only field that has to be the same is the model, as that defines the dimensionality usually.
   Note: This is only relevant if you plan to add more documents to the dataset after importing it.


```
knowledge import <path> [<dataset-id>...] [flags]
```

### Options

```
      --auto-migrate string               Auto migrate database ($KNOW_DB_AUTO_MIGRATE) (default "true")
  -c, --config-file string                Path to the configuration file ($KNOW_CONFIG_FILE)
      --dsn string                        Server database connection string (default "sqlite://$XDG_DATA_HOME/gptscript/knowledge/knowledge.db") ($KNOW_DB_DSN)
      --embedding-model-provider string   Embedding model provider ($KNOW_EMBEDDING_MODEL_PROVIDER) (default "openai")
  -h, --help                              help for import
      --server string                     URL of the Knowledge API Server ($KNOW_SERVER_URL)
      --vector-dbpath string              VectorDBPath to the vector database (default "$XDG_DATA_HOME/gptscript/knowledge/vector.db") ($KNOW_VECTOR_DB_PATH)
```

### SEE ALSO

* [knowledge](knowledge.md)	 - 

