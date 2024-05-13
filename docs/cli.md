# CLI Summary

## Commands and Flags

### knowledge

- **askdir**: Retrieve sources for a query from a dataset generated from a directory
  - Flags:
    - `--auto-migrate`: Auto migrate database (Env: `KNOW_DB_AUTO_MIGRATE`, default: `true`)
    - `-c, --concurrency`: Number of concurrent ingestion processes (Env: `KNOW_INGEST_CONCURRENCY`, default: 10)
    - `--dsn`: Server database connection string (Env: `KNOW_DB_DSN`, default: `sqlite://$XDG_DATA_HOME/gptscript/knowledge/knowledge.db`)
    - `-h, --help`: Help for askdir
    - `--ignore-extensions`: Comma-separated list of file extensions to ignore (Env: `KNOW_INGEST_IGNORE_EXTENSIONS`)
    - `--openai-api-base`: OpenAI API base (Env: `OPENAI_BASE_URL`, default: `https://api.openai.com/v1`)
    - `--openai-api-key`: OpenAI API key (Env: `OPENAI_API_KEY`, default: `sk-foo`)
    - `--openai-api-type`: OpenAI API type (Env: `OPENAI_API_TYPE`, default: `OPEN_AI`)
    - `--openai-api-version`: OpenAI API version (for Azure) (Env: `OPENAI_API_VERSION`, default: `2024-02-01`)
    - `--openai-azure-deployment`: Azure OpenAI deployment name (Env: `OPENAI_AZURE_DEPLOYMENT`)
    - `--openai-embedding-model`: OpenAI Embedding model (Env: `OPENAI_EMBEDDING_MODEL`, default: `text-embedding-ada-002`)
    - `-p, --path`: Path to the directory to query (Env: `KNOWLEDGE_CLIENT_ASK_DIR_PATH`, default: `./knowledge`)
    - `-r, --recursive`: Recursively ingest directories (Env: `KNOW_INGEST_RECURSIVE`)
    - `--server`: URL of the Knowledge API Server (Env: `KNOW_SERVER_URL`)
    - `-k, --top-k`: Number of sources to retrieve (Env: `KNOWLEDGE_CLIENT_ASK_DIR_TOP_K`, default: 5)
    - `--vector-dbpath`: VectorDBPath to the vector database (Env: `KNOW_VECTOR_DB_PATH`, default: `$XDG_DATA_HOME/gptscript/knowledge/vector.db`)

- **completion**: Generate the autocompletion script for the specified shell
  - Flags:
    - `-h, --help`: Help for completion

- **create-dataset**: Create a new dataset
  - Flags:
    - `--auto-migrate`: Auto migrate database (Env: `KNOW_DB_AUTO_MIGRATE`, default: `true`)
    - `--dsn`: Server database connection string (Env: `KNOW_DB_DSN`, default: `sqlite://$XDG_DATA_HOME/gptscript/knowledge/knowledge.db`)
    - `--embed-dim`: Embedding dimension (Env: `KNOWLEDGE_CLIENT_CREATE_DATASET_EMBED_DIM`, default: 1536)
    - `-h, --help`: Help for create-dataset
    - `--openai-api-base`: OpenAI API base (Env: `OPENAI_BASE_URL`, default: `https://api.openai.com/v1`)
    - `--openai-api-key`: OpenAI API key (Env: `OPENAI_API_KEY`, default: `sk-foo`)
    - `--openai-api-type`: OpenAI API type (Env: `OPENAI_API_TYPE`, default: `OPEN_AI`)
    - `--openai-api-version`: OpenAI API version (for Azure) (Env: `OPENAI_API_VERSION`, default: `2024-02-01`)
    - `--openai-azure-deployment`: Azure OpenAI deployment name (Env: `OPENAI_AZURE_DEPLOYMENT`)
    - `--openai-embedding-model`: OpenAI Embedding model (Env: `OPENAI_EMBEDDING_MODEL`, default: `text-embedding-ada-002`)
    - `--server`: URL of the Knowledge API Server (Env: `KNOW_SERVER_URL`)
    - `--vector-dbpath`: VectorDBPath to the vector database (Env: `KNOW_VECTOR_DB_PATH`, default: `$XDG_DATA_HOME/gptscript/knowledge/vector.db`)

- **delete-dataset**: Delete a dataset
  - Flags:
    - `--auto-migrate`: Auto migrate database (Env: `KNOW_DB_AUTO_MIGRATE`, default: `true`)
    - `--dsn`: Server database connection string (Env: `KNOW_DB_DSN`, default: `sqlite://$XDG_DATA_HOME/gptscript/knowledge/knowledge.db`)
    - `-h, --help`: Help for delete-dataset
    - `--openai-api-base`: OpenAI API base (Env: `OPENAI_BASE_URL`, default: `https://api.openai.com/v1`)
    - `--openai-api-key`: OpenAI API key (Env: `OPENAI_API_KEY`, default: `sk-foo`)
    - `--openai-api-type`: OpenAI API type (Env: `OPENAI_API_TYPE`, default: `OPEN_AI`)
    - `--openai-api-version`: OpenAI API version (for Azure) (Env: `OPENAI_API_VERSION`, default: `2024-02-01`)
    - `--openai-azure-deployment`: Azure OpenAI deployment name (Env: `OPENAI_AZURE_DEPLOYMENT`)
    - `--openai-embedding-model`: OpenAI Embedding model (Env: `OPENAI_EMBEDDING_MODEL`, default: `text-embedding-ada-002`)
    - `--server`: URL of the Knowledge API Server (Env: `KNOW_SERVER_URL`)
    - `--vector-dbpath`: VectorDBPath to the vector database (Env: `KNOW_VECTOR_DB_PATH`, default: `$XDG_DATA_HOME/gptscript/knowledge/vector.db`)

- **get-dataset**: Get a dataset
  - Flags:
    - `--auto-migrate`: Auto migrate database (Env: `KNOW_DB_AUTO_MIGRATE`, default: `true`)
    - `--dsn`: Server database connection string (Env: `KNOW_DB_DSN`, default: `sqlite://$XDG_DATA_HOME/gptscript/knowledge/knowledge.db`)
    - `-h, --help`: Help for get-dataset
    - `--openai-api-base`: OpenAI API base (Env: `OPENAI_BASE_URL`, default: `https://api.openai.com/v1`)
    - `--openai-api-key`: OpenAI API key (Env: `OPENAI_API_KEY`, default: `sk-foo`)
    - `--openai-api-type`: OpenAI API type (Env: `OPENAI_API_TYPE`, default: `OPEN_AI`)
    - `--openai-api-version`: OpenAI API version (for Azure) (Env: `OPENAI_API_VERSION`, default: `2024-02-01`)
    - `--openai-azure-deployment`: Azure OpenAI deployment name (Env: `OPENAI_AZURE_DEPLOYMENT`)
    - `--openai-embedding-model`: OpenAI Embedding model (Env: `OPENAI_EMBEDDING_MODEL`, default: `text-embedding-ada-002`)
    - `--server`: URL of the Knowledge API Server (Env: `KNOW_SERVER_URL`)
    - `--vector-dbpath`: VectorDBPath to the vector database (Env: `KNOW_VECTOR_DB_PATH`, default: `$XDG_DATA_HOME/gptscript/knowledge/vector.db`)

- **ingest**: Ingest a file/directory into a dataset (non-recursive)
  - Flags:
    - `--auto-migrate`: Auto migrate database (Env: `KNOW_DB_AUTO_MIGRATE`, default: `true`)
    - `-c, --concurrency`: Number of concurrent ingestion processes (Env: `KNOW_INGEST_CONCURRENCY`, default: 10)
    - `-d, --dataset`: Target Dataset ID (Env: `KNOW_TARGET_DATASET`, default: `default`)
    - `--dsn`: Server database connection string (Env: `KNOW_DB_DSN`, default: `sqlite://$XDG_DATA_HOME/gptscript/knowledge/knowledge.db`)
    - `-h, --help`: Help for ingest
    - `--ignore-extensions`: Comma-separated list of file extensions to ignore (Env: `KNOW_INGEST_IGNORE_EXTENSIONS`)
    - `--openai-api-base`: OpenAI API base (Env: `OPENAI_BASE_URL`, default: `https://api.openai.com/v1`)
    - `--openai-api-key`: OpenAI API key (Env: `OPENAI_API_KEY`, default: `sk-foo`)
    - `--openai-api-type`: OpenAI API type (Env: `OPENAI_API_TYPE`, default: `OPEN_AI`)
    - `--openai-api-version`: OpenAI API version (for Azure) (Env: `OPENAI_API_VERSION`, default: `2024-02-01`)
    - `--openai-azure-deployment`: Azure OpenAI deployment name (Env: `OPENAI_AZURE_DEPLOYMENT`)
    - `--openai-embedding-model`: OpenAI Embedding model (Env: `OPENAI_EMBEDDING_MODEL`, default: `text-embedding-ada-002`)
    - `-r, --recursive`: Recursively ingest directories (Env: `KNOW_INGEST_RECURSIVE`)
    - `--server`: URL of the Knowledge API Server (Env: `KNOW_SERVER_URL`)
    - `--textsplitter-chunk-overlap`: Textsplitter Chunk Overlap (Env: `KNOW_TEXTSPLITTER_CHUNK_OVERLAP`, default: 256)
    - `--textsplitter-chunk-size`: Textsplitter Chunk Size (Env: `KNOW_TEXTSPLITTER_CHUNK_SIZE`, default: 1024)
    - `--textsplitter-encoding-name`: Textsplitter Encoding Name (Env: `KNOW_TEXTSPLITTER_ENCODING_NAME`, default: `cl100k_base`)
    - `--textsplitter-model-name`: Textsplitter Model Name (Env: `KNOW_TEXTSPLITTER_MODEL_NAME`, default: `gpt-4`)
    - `--vector-dbpath`: VectorDBPath to the vector database (Env: `KNOW_VECTOR_DB_PATH`, default: `$XDG_DATA_HOME/gptscript/knowledge/vector.db`)

- **list-datasets**: List existing datasets
  - Flags:
    - `--auto-migrate`: Auto migrate database (Env: `KNOW_DB_AUTO_MIGRATE`, default: `true`)
    - `--dsn`: Server database connection string (Env: `KNOW_DB_DSN`, default: `sqlite://$XDG_DATA_HOME/gptscript/knowledge/knowledge.db`)
    - `-h, --help`: Help for list-datasets
    - `--openai-api-base`: OpenAI API base (Env: `OPENAI_BASE_URL`, default: `https://api.openai.com/v1`)
    - `--openai-api-key`: OpenAI API key (Env: `OPENAI_API_KEY`, default: `sk-foo`)
    - `--openai-api-type`: OpenAI API type (Env: `OPENAI_API_TYPE`, default: `OPEN_AI`)
    - `--openai-api-version`: OpenAI API version (for Azure) (Env: `OPENAI_API_VERSION`, default: `2024-02-01`)
    - `--openai-azure-deployment`: Azure OpenAI deployment name (Env: `OPENAI_AZURE_DEPLOYMENT`)
    - `--openai-embedding-model`: OpenAI Embedding model (Env: `OPENAI_EMBEDDING_MODEL`, default: `text-embedding-ada-002`)
    - `--server`: URL of the Knowledge API Server (Env: `KNOW_SERVER_URL`)
    - `--vector-dbpath`: VectorDBPath to the vector database (Env: `KNOW_VECTOR_DB_PATH`, default: `$XDG_DATA_HOME/gptscript/knowledge/vector.db`)

- **retrieve**: Retrieve sources for a query from a dataset
  - Flags:
    - `--auto-migrate`: Auto migrate database (Env: `KNOW_DB_AUTO_MIGRATE`, default: `true`)
    - `-d, --dataset`: Target Dataset ID (Env: `KNOW_TARGET_DATASET`, default: `default`)
    - `--dsn`: Server database connection string (Env: `KNOW_DB_DSN`, default: `sqlite://$XDG_DATA_HOME/gptscript/knowledge/knowledge.db`)
    - `-h, --help`: Help for retrieve
    - `--openai-api-base`: OpenAI API base (Env: `OPENAI_BASE_URL`, default: `https://api.openai.com/v1`)
    - `--openai-api-key`: OpenAI API key (Env: `OPENAI_API_KEY`, default: `sk-foo`)
    - `--openai-api-type`: OpenAI API type (Env: `OPENAI_API_TYPE`, default: `OPEN_AI`)
    - `--openai-api-version`: OpenAI API version (for Azure) (Env: `OPENAI_API_VERSION`, default: `2024-02-01`)
    - `--openai-azure-deployment`: Azure OpenAI deployment name (Env: `OPENAI_AZURE_DEPLOYMENT`)
    - `--openai-embedding-model`: OpenAI Embedding model (Env: `OPENAI_EMBEDDING_MODEL`, default: `text-embedding-ada-002`)
    - `--server`: URL of the Knowledge API Server (Env: `KNOW_SERVER_URL`)
    - `-k, --top-k`: Number of sources to retrieve (Env: `KNOWLEDGE_CLIENT_RETRIEVE_TOP_K`, default: 5)
    - `--vector-dbpath`: VectorDBPath to the vector database (Env: `KNOW_VECTOR_DB_PATH`, default: `$XDG_DATA_HOME/gptscript/knowledge/vector.db`)

- **server**: Start the Knowledge API Server
  - Flags:
    - `--auto-migrate`: Auto migrate database (Env: `KNOW_DB_AUTO_MIGRATE`, default: `true`)
    - `--dsn`: Server database connection string (Env: `KNOW_DB_DSN`, default: `sqlite://$XDG_DATA_HOME/gptscript/knowledge/knowledge.db`)
    - `-h, --help`: Help for server
    - `--openai-api-base`: OpenAI API base (Env: `OPENAI_BASE_URL`, default: `https://api.openai.com/v1`)
    - `--openai-api-key`: OpenAI API key (Env: `OPENAI_API_KEY`, default: `sk-foo`)
    - `--openai-api-type`: OpenAI API type (Env: `OPENAI_API_TYPE`, default: `OPEN_AI`)
    - `--openai-api-version`: OpenAI API version (for Azure) (Env: `OPENAI_API_VERSION`, default: `2024-02-01`)
    - `--openai-azure-deployment`: Azure OpenAI deployment name (Env: `OPENAI_AZURE_DEPLOYMENT`)
    - `--openai-embedding-model`: OpenAI Embedding model (Env: `OPENAI_EMBEDDING_MODEL`, default: `text-embedding-ada-002`)
    - `--server-apibase`: Server API base (Env: `KNOW_SERVER_API_BASE`, default: `/v1`)
    - `--server-port`: Server port (Env: `KNOW_SERVER_PORT`, default: `8000`)
    - `--server-url`: Server URL (Env: `KNOW_SERVER_URL`, default: `http://localhost`)
    - `--vector-dbpath`: VectorDBPath to the vector database (Env: `KNOW_VECTOR_DB_PATH`, default: `$XDG_DATA_HOME/gptscript/knowledge/vector.db`)

- **version**: Display the version of the Knowledge tool
  - Flags:
    - `-h, --help`: Help for version
