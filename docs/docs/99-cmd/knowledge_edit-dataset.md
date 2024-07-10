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
      --auto-migrate string              Auto migrate database ($KNOW_DB_AUTO_MIGRATE) (default "true")
      --dsn string                       Server database connection string (default "sqlite://$XDG_DATA_HOME/gptscript/knowledge/knowledge.db") ($KNOW_DB_DSN)
  -h, --help                             help for edit-dataset
      --openai-api-base string           OpenAI API base ($OPENAI_BASE_URL) (default "https://api.openai.com/v1")
      --openai-api-key string            OpenAI API key (not required if used with clicky-chats) ($OPENAI_API_KEY) (default "sk-foo")
      --openai-api-type string           OpenAI API type (OPEN_AI, AZURE, AZURE_AD) ($OPENAI_API_TYPE) (default "OPEN_AI")
      --openai-api-version string        OpenAI API version (for Azure) ($OPENAI_API_VERSION) (default "2024-02-01")
      --openai-azure-deployment string   Azure OpenAI deployment name (overrides openai-embedding-model, if set) ($OPENAI_AZURE_DEPLOYMENT)
      --openai-embedding-model string    OpenAI Embedding model ($OPENAI_EMBEDDING_MODEL) (default "text-embedding-ada-002")
      --openai-model string              OpenAI model ($OPENAI_MODEL) (default "gpt-4")
      --replace-metadata strings         replace metadata with key-value pairs (existing metadata will be removed) ($KNOWLEDGE_CLIENT_EDIT_DATASET_REPLACE_METADATA)
      --reset-metadata                   reset metadata to default (empty) ($KNOWLEDGE_CLIENT_EDIT_DATASET_RESET_METADATA)
      --server string                    URL of the Knowledge API Server ($KNOW_SERVER_URL)
      --update-metadata strings          update metadata key-value pairs (existing metadata will be updated/preserved) ($KNOWLEDGE_CLIENT_EDIT_DATASET_UPDATE_METADATA)
      --vector-dbpath string             VectorDBPath to the vector database (default "$XDG_DATA_HOME/gptscript/knowledge/vector.db") ($KNOW_VECTOR_DB_PATH)
```

### SEE ALSO

* [knowledge](knowledge.md)	 - 

