---
title: "knowledge import"
---
## knowledge import

Import one or more datasets from an archive (zip) (default: all datasets)

### Synopsis

Import one or more datasets from an archive (zip) (default: all datasets).
## IMPORTANT: Embedding functions
   Embedding functions are not part of exported knowledge base archives, so you'll have to know the embedding function used to import the archive.
   This primarily concerns the choice of the embeddings provider (model) and the embedding dimension.
   Note: This is only relevant if you plan to add more documents to the dataset after importing it.


```
knowledge import <path> [<dataset-id>...] [flags]
```

### Options

```
      --auto-migrate string              Auto migrate database ($KNOW_DB_AUTO_MIGRATE) (default "true")
      --dsn string                       Server database connection string (default "sqlite://$XDG_DATA_HOME/gptscript/knowledge/knowledge.db") ($KNOW_DB_DSN)
  -h, --help                             help for import
      --openai-api-base string           OpenAI API base ($OPENAI_BASE_URL) (default "https://api.openai.com/v1")
      --openai-api-key string            OpenAI API key (not required if used with clicky-chats) ($OPENAI_API_KEY) (default "sk-foo")
      --openai-api-type string           OpenAI API type (OPEN_AI, AZURE, AZURE_AD) ($OPENAI_API_TYPE) (default "OPEN_AI")
      --openai-api-version string        OpenAI API version (for Azure) ($OPENAI_API_VERSION) (default "2024-02-01")
      --openai-azure-deployment string   Azure OpenAI deployment name (overrides openai-embedding-model, if set) ($OPENAI_AZURE_DEPLOYMENT)
      --openai-embedding-model string    OpenAI Embedding model ($OPENAI_EMBEDDING_MODEL) (default "text-embedding-ada-002")
      --openai-model string              OpenAI model ($OPENAI_MODEL) (default "gpt-4")
      --server string                    URL of the Knowledge API Server ($KNOW_SERVER_URL)
      --vector-dbpath string             VectorDBPath to the vector database (default "$XDG_DATA_HOME/gptscript/knowledge/vector.db") ($KNOW_VECTOR_DB_PATH)
```

### SEE ALSO

* [knowledge](knowledge.md)	 - 

