---
title: "knowledge server"
---
## knowledge server



```
knowledge server [flags]
```

### Options

```
      --auto-migrate string              Auto migrate database ($KNOW_DB_AUTO_MIGRATE) (default "true")
      --dsn string                       Server database connection string (default "sqlite://$XDG_DATA_HOME/gptscript/knowledge/knowledge.db") ($KNOW_DB_DSN)
  -h, --help                             help for server
      --openai-api-base string           OpenAI API base ($OPENAI_BASE_URL) (default "https://api.openai.com/v1")
      --openai-api-key string            OpenAI API key (not required if used with clicky-chats) ($OPENAI_API_KEY) (default "sk-foo")
      --openai-api-type string           OpenAI API type (OPEN_AI, AZURE, AZURE_AD) ($OPENAI_API_TYPE) (default "OPEN_AI")
      --openai-api-version string        OpenAI API version (for Azure) ($OPENAI_API_VERSION) (default "2024-02-01")
      --openai-azure-deployment string   Azure OpenAI deployment name (overrides openai-embedding-model, if set) ($OPENAI_AZURE_DEPLOYMENT)
      --openai-embedding-model string    OpenAI Embedding model ($OPENAI_EMBEDDING_MODEL) (default "text-embedding-ada-002")
      --openai-model string              OpenAI model ($OPENAI_MODEL) (default "gpt-4")
      --server-apibase string            Server API base ($KNOW_SERVER_API_BASE) (default "/v1")
      --server-port string               Server port ($KNOW_SERVER_PORT) (default "8000")
      --server-url string                Server URL ($KNOW_SERVER_URL) (default "http://localhost")
      --vector-dbpath string             VectorDBPath to the vector database (default "$XDG_DATA_HOME/gptscript/knowledge/vector.db") ($KNOW_VECTOR_DB_PATH)
```

### SEE ALSO

* [knowledge](knowledge.md)	 - 

