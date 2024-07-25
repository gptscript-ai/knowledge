---
title: "knowledge import"
---
## knowledge import

Import one or more datasets from an archive (zip) (default: all datasets)

### Synopsis

Import one or more datasets from an archive (zip) (default: all datasets).
## IMPORTANT: Embedding functions
   Embedding functions are not part of exported knowledge base archives, so you'll have to know the embedding function used to import the archive.
   This primarily concerns the choice of the embeddings provider (model).
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

