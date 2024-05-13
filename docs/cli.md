# CLI Summary

Note: Generated with gptscript: `gptscript docs/clidocs.gpt`

## knowledge

Usage:

    knowledge [flags]
    knowledge [command]

Available Commands:

    askdir          Retrieve sources for a query from a dataset generated from a directory
    completion      Generate the autocompletion script for the specified shell
    create-dataset  Create a new dataset
    delete-dataset  Delete a dataset
    get-dataset     Get a dataset
    help            Help about any command
    ingest          Ingest a file/directory into a dataset (non-recursive)
    list-datasets   List existing datasets
    retrieve        Retrieve sources for a query from a dataset
    server          
    version         

Flags:

    -h, --help   help for knowledge

Use "knowledge [command] --help" for more information about a command.

## knowledge askdir

Retrieve sources for a query from a dataset generated from a directory

Usage:

    knowledge askdir [--path <path>] <query> [flags]

Flags:

        --auto-migrate string             Auto migrate database ($KNOW_DB_AUTO_MIGRATE) (default "true")
    -c, --concurrency int                 Number of concurrent ingestion processes ($KNOW_INGEST_CONCURRENCY) (default 10)
        --dsn string                      Server database connection string (default "sqlite://$XDG_DATA_HOME/gptscript/knowledge/knowledge.db") ($KNOW_DB_DSN)
    -h, --help                            help for askdir
        --ignore-extensions string        Comma-separated list of file extensions to ignore ($KNOW_INGEST_IGNORE_EXTENSIONS)
        --openai-api-base string          OpenAI API base ($OPENAI_BASE_URL) (default "https://api.openai.com/v1")
        --openai-api-key string           OpenAI API key (not required if used with clicky-chats) ($OPENAI_API_KEY) (default "sk-foo")
        --openai-api-type string          OpenAI API type (OPEN_AI, AZURE, AZURE_AD) ($OPENAI_API_TYPE) (default "OPEN_AI")
        --openai-api-version string       OpenAI API version (for Azure) ($OPENAI_API_VERSION) (default "2024-02-01")
        --openai-embedding-model string   OpenAI Embedding model ($OPENAI_EMBEDDING_MODEL) (default "text-embedding-ada-002")
    -p, --path string                     Path to the directory to query ($KNOWLEDGE_CLIENT_ASK_DIR_PATH) (default "./knowledge")
    -r, --recursive                       Recursively ingest directories ($KNOW_INGEST_RECURSIVE)
        --server string                   URL of the Knowledge API Server ($KNOW_SERVER_URL)
    -k, --top-k int                       Number of sources to retrieve ($KNOWLEDGE_CLIENT_ASK_DIR_TOP_K) (default 5)
        --vector-dbpath string            VectorDBPath to the vector database (default "$XDG_DATA_HOME/gptscript/knowledge/vector.db") ($KNOW_VECTOR_DB_PATH)

## knowledge completion

Generate the autocompletion script for knowledge for the specified shell.
See each sub-command's help for details on how to use the generated script.

Usage:

    knowledge completion [command]

Available Commands:

    bash        Generate the autocompletion script for bash
    fish        Generate the autocompletion script for fish
    powershell  Generate the autocompletion script for powershell
    zsh         Generate the autocompletion script for zsh

Flags:

    -h, --help   help for completion

Use "knowledge completion [command] --help" for more information about a command.

## knowledge create-dataset

Create a new dataset

Usage:

    knowledge create-dataset <dataset-id> [flags]

Flags:
  
        --auto-migrate string             Auto migrate database ($KNOW_DB_AUTO_MIGRATE) (default "true")
        --dsn string                      Server database connection string (default "sqlite://$XDG_DATA_HOME/gptscript/knowledge/knowledge.db") ($KNOW_DB_DSN)
        --embed-dim int                   Embedding dimension ($KNOWLEDGE_CLIENT_CREATE_DATASET_EMBED_DIM) (default 1536)
    -h, --help                            help for create-dataset
        --openai-api-base string          OpenAI API base ($OPENAI_BASE_URL) (default "https://api.openai.com/v1")
        --openai-api-key string           OpenAI API key (not required if used with clicky-chats) ($OPENAI_API_KEY) (default "sk-foo")
        --openai-api-type string          OpenAI API type (OPEN_AI, AZURE, AZURE_AD) ($OPENAI_API_TYPE) (default "OPEN_AI")
        --openai-api-version string       OpenAI API version (for Azure) ($OPENAI_API_VERSION) (default "2024-02-01")
        --openai-embedding-model string   OpenAI Embedding model ($OPENAI_EMBEDDING_MODEL) (default "text-embedding-ada-002")
        --server string                   URL of the Knowledge API Server ($KNOW_SERVER_URL)
        --vector-dbpath string            VectorDBPath to the vector database (default "$XDG_DATA_HOME/gptscript/knowledge/vector.db") ($KNOW_VECTOR_DB_PATH)

## knowledge delete-dataset

Delete a dataset

Usage:

    knowledge delete-dataset <dataset-id> [flags]

Flags:

        --auto-migrate string             Auto migrate database ($KNOW_DB_AUTO_MIGRATE) (default "true")
        --dsn string                      Server database connection string (default "sqlite://$XDG_DATA_HOME/gptscript/knowledge/knowledge.db") ($KNOW_DB_DSN)
    -h, --help                            help for delete-dataset
        --openai-api-base string          OpenAI API base ($OPENAI_BASE_URL) (default "https://api.openai.com/v1")
        --openai-api-key string           OpenAI API key (not required if used with clicky-chats) ($OPENAI_API_KEY) (default "sk-foo")
        --openai-api-type string          OpenAI API type (OPEN_AI, AZURE, AZURE_AD) ($OPENAI_API_TYPE) (default "OPEN_AI")
        --openai-api-version string       OpenAI API version (for Azure) ($OPENAI_API_VERSION) (default "2024-02-01")
        --openai-embedding-model string   OpenAI Embedding model ($OPENAI_EMBEDDING_MODEL) (default "text-embedding-ada-002")
        --server string                   URL of the Knowledge API Server ($KNOW_SERVER_URL)
        --vector-dbpath string            VectorDBPath to the vector database (default "$XDG_DATA_HOME/gptscript/knowledge/vector.db") ($KNOW_VECTOR_DB_PATH)

## knowledge get-dataset

Get a dataset

Usage:

    knowledge get-dataset <dataset-id> [flags]

Flags:

        --auto-migrate string             Auto migrate database ($KNOW_DB_AUTO_MIGRATE) (default "true")
        --dsn string                      Server database connection string (default "sqlite://$XDG_DATA_HOME/gptscript/knowledge/knowledge.db") ($KNOW_DB_DSN)
    -h, --help                            help for get-dataset
        --openai-api-base string          OpenAI API base ($OPENAI_BASE_URL) (default "https://api.openai.com/v1")
        --openai-api-key string           OpenAI API key (not required if used with clicky-chats) ($OPENAI_API_KEY) (default "sk-foo")
        --openai-api-type string          OpenAI API type (OPEN_AI, AZURE, AZURE_AD) ($OPENAI_API_TYPE) (default "OPEN_AI")
        --openai-api-version string       OpenAI API version (for Azure) ($OPENAI_API_VERSION) (default "2024-02-01")
        --openai-embedding-model string   OpenAI Embedding model ($OPENAI_EMBEDDING_MODEL) (default "text-embedding-ada-002")
        --server string                   URL of the Knowledge API Server ($KNOW_SERVER_URL)
        --vector-dbpath string            VectorDBPath to the vector database (default "$XDG_DATA_HOME/gptscript/knowledge/vector.db") ($KNOW_VECTOR_DB_PATH)

## knowledge ingest

Ingest a file/directory into a dataset (non-recursive)

Usage:

    knowledge ingest [--dataset <dataset-id>] <path> [flags]

Flags:

        --auto-migrate string                 Auto migrate database ($KNOW_DB_AUTO_MIGRATE) (default "true")
    -c, --concurrency int                     Number of concurrent ingestion processes ($KNOW_INGEST_CONCURRENCY) (default 10)
    -d, --dataset string                      Target Dataset ID ($KNOW_TARGET_DATASET) (default "default")
        --dsn string                          Server database connection string (default "sqlite://$XDG_DATA_HOME/gptscript/knowledge/knowledge.db") ($KNOW_DB_DSN)
    -h, --help                                help for ingest
        --ignore-extensions string            Comma-separated list of file extensions to ignore ($KNOW_INGEST_IGNORE_EXTENSIONS)
        --openai-api-base string              OpenAI API base ($OPENAI_BASE_URL) (default "https://api.openai.com/v1")
        --openai-api-key string               OpenAI API key (not required if used with clicky-chats) ($OPENAI_API_KEY) (default "sk-foo")
        --openai-api-type string              OpenAI API type (OPEN_AI, AZURE, AZURE_AD) ($OPENAI_API_TYPE) (default "OPEN_AI")
        --openai-api-version string           OpenAI API version (for Azure) ($OPENAI_API_VERSION) (default "2024-02-01")
        --openai-embedding-model string       OpenAI Embedding model ($OPENAI_EMBEDDING_MODEL) (default "text-embedding-ada-002")
    -r, --recursive                           Recursively ingest directories ($KNOW_INGEST_RECURSIVE)
        --server string                       URL of the Knowledge API Server ($KNOW_SERVER_URL)
        --textsplitter-chunk-overlap int      Textsplitter Chunk Overlap ($KNOW_TEXTSPLITTER_CHUNK_OVERLAP) (default 256)
        --textsplitter-chunk-size int         Textsplitter Chunk Size ($KNOW_TEXTSPLITTER_CHUNK_SIZE) (default 1024)
        --textsplitter-encoding-name string   Textsplitter Encoding Name ($KNOW_TEXTSPLITTER_ENCODING_NAME) (default "cl100k_base")
        --textsplitter-model-name string      Textsplitter Model Name ($KNOW_TEXTSPLITTER_MODEL_NAME) (default "gpt-4")
        --vector-dbpath string                VectorDBPath to the vector database (default "$XDG_DATA_HOME/gptscript/knowledge/vector.db") ($KNOW_VECTOR_DB_PATH)

## knowledge list-datasets

List existing datasets

Usage:

    knowledge list-datasets [flags]

Flags:

        --auto-migrate string             Auto migrate database ($KNOW_DB_AUTO_MIGRATE) (default "true")
        --dsn string                      Server database connection string (default "sqlite://$XDG_DATA_HOME/gptscript/knowledge/knowledge.db") ($KNOW_DB_DSN)
    -h, --help                            help for list-datasets
        --openai-api-base string          OpenAI API base ($OPENAI_BASE_URL) (default "https://api.openai.com/v1")
        --openai-api-key string           OpenAI API key (not required if used with clicky-chats) ($OPENAI_API_KEY) (default "sk-foo")
        --openai-api-type string          OpenAI API type (OPEN_AI, AZURE, AZURE_AD) ($OPENAI_API_TYPE) (default "OPEN_AI")
        --openai-api-version string       OpenAI API version (for Azure) ($OPENAI_API_VERSION) (default "2024-02-01")
        --openai-embedding-model string   OpenAI Embedding model ($OPENAI_EMBEDDING_MODEL) (default "text-embedding-ada-002")
        --server string                   URL of the Knowledge API Server ($KNOW_SERVER_URL)
        --vector-dbpath string            VectorDBPath to the vector database (default "$XDG_DATA_HOME/gptscript/knowledge/vector.db") ($KNOW_VECTOR_DB_PATH)

## knowledge retrieve

Retrieve sources for a query from a dataset

Usage:

    knowledge retrieve [--dataset <dataset-id>] <query> [flags]

Flags:

        --auto-migrate string             Auto migrate database ($KNOW_DB_AUTO_MIGRATE) (default "true")
    -d, --dataset string                  Target Dataset ID ($KNOW_TARGET_DATASET) (default "default")
        --dsn string                      Server database connection string (default "sqlite://$XDG_DATA_HOME/gptscript/knowledge/knowledge.db") ($KNOW_DB_DSN)
    -h, --help                            help for retrieve
        --openai-api-base string          OpenAI API base ($OPENAI_BASE_URL) (default "https://api.openai.com/v1")
        --openai-api-key string           OpenAI API key (not required if used with clicky-chats) ($OPENAI_API_KEY) (default "sk-foo")
        --openai-api-type string          OpenAI API type (OPEN_AI, AZURE, AZURE_AD) ($OPENAI_API_TYPE) (default "OPEN_AI")
        --openai-api-version string       OpenAI API version (for Azure) ($OPENAI_API_VERSION) (default "2024-02-01")
        --openai-embedding-model string   OpenAI Embedding model ($OPENAI_EMBEDDING_MODEL) (default "text-embedding-ada-002")
        --server string                   URL of the Knowledge API Server ($KNOW_SERVER_URL)
    -k, --top-k int                       Number of sources to retrieve ($KNOWLEDGE_CLIENT_RETRIEVE_TOP_K) (default 5)
        --vector-dbpath string            VectorDBPath to the vector database (default "$XDG_DATA_HOME/gptscript/knowledge/vector.db") ($KNOW_VECTOR_DB_PATH)

## knowledge server

Usage:

    knowledge server [flags]

Flags:

        --auto-migrate string             Auto migrate database ($KNOW_DB_AUTO_MIGRATE) (default "true")
        --dsn string                      Server database connection string (default "sqlite://$XDG_DATA_HOME/gptscript/knowledge/knowledge.db") ($KNOW_DB_DSN)
    -h, --help                            help for server
        --openai-api-base string          OpenAI API base ($OPENAI_BASE_URL) (default "https://api.openai.com/v1")
        --openai-api-key string           OpenAI API key (not required if used with clicky-chats) ($OPENAI_API_KEY) (default "sk-foo")
        --openai-api-type string          OpenAI API type (OPEN_AI, AZURE, AZURE_AD) ($OPENAI_API_TYPE) (default "OPEN_AI")
        --openai-api-version string       OpenAI API version (for Azure) ($OPENAI_API_VERSION) (default "2024-02-01")
        --openai-embedding-model string   OpenAI Embedding model ($OPENAI_EMBEDDING_MODEL) (default "text-embedding-ada-002")
        --server-apibase string           Server API base ($KNOW_SERVER_API_BASE) (default "/v1")
        --server-port string              Server port ($KNOW_SERVER_PORT) (default "8000")
        --server-url string               Server URL ($KNOW_SERVER_URL) (default "http://localhost")
        --vector-dbpath string            VectorDBPath to the vector database (default "$XDG_DATA_HOME/gptscript/knowledge/vector.db") ($KNOW_VECTOR_DB_PATH)

## knowledge version

Usage:

    knowledge version [flags]

Flags:

    -h, --help   help for version
