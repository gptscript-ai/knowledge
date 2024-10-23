---
title: "knowledge askdir"
---
## knowledge askdir

Retrieve sources for a query from a dataset generated from a directory

```
knowledge askdir [--path <path>] <query> [flags]
```

### Options

```
      --auto-migrate string               Auto migrate database ($KNOW_DB_AUTO_MIGRATE) (default "true")
      --concurrency int                   Number of concurrent ingestion processes ($KNOW_INGEST_CONCURRENCY) (default 10)
  -c, --config-file string                Path to the configuration file ($KNOW_CONFIG_FILE)
      --dedupe-func string                Name of the deduplication function to use ($KNOW_INGEST_DEDUPE_FUNC)
      --embedding-model-provider string   Embedding model provider ($KNOW_EMBEDDING_MODEL_PROVIDER) (default "openai")
      --err-on-unsupported-file           Error on unsupported file types ($KNOW_INGEST_ERR_ON_UNSUPPORTED_FILE)
      --flow string                       Flow name ($KNOW_FLOW)
      --flows-file string                 Path to a YAML/JSON file containing ingestion/retrieval flows ($KNOW_FLOWS_FILE) (default "blueprint:default")
  -h, --help                              help for askdir
      --ignore-extensions string          Comma-separated list of file extensions to ignore ($KNOW_INGEST_IGNORE_EXTENSIONS)
      --ignore-file string                Path to a .gitignore style file ($KNOW_INGEST_IGNORE_FILE)
      --include-hidden                    Include hidden files and directories ($KNOW_INGEST_INCLUDE_HIDDEN)
      --index-dsn string                  Index Database Connection string (relational DB) (default "sqlite://$XDG_DATA_HOME/gptscript/knowledge/knowledge.db") ($KNOW_INDEX_DSN)
  -w, --keyword strings                   Keywords that retrieved documents must contain ($KNOW_RETRIEVE_KEYWORDS)
      --no-create-dataset                 Do NOT create the dataset if it doesn't exist ($KNOW_INGEST_NO_CREATE_DATASET)
      --no-prune                          Do not prune deleted files ($KNOW_ASKDIR_NO_PRUNE)
      --no-recursive                      Don't recursively ingest directories ($KNOW_NO_INGEST_RECURSIVE)
  -p, --path string                       Path to the directory to query ($KNOWLEDGE_CLIENT_ASK_DIR_PATH) (default ".")
      --server string                     URL of the Knowledge API Server ($KNOW_SERVER_URL)
  -k, --top-k int                         Number of sources to retrieve ($KNOWLEDGE_CLIENT_ASK_DIR_TOP_K) (default 10)
      --vector-dsn string                 DSN to the vector database (default "chromem:$XDG_DATA_HOME/gptscript/knowledge/vector.db") ($KNOW_VECTOR_DSN)
```

### SEE ALSO

* [knowledge](knowledge.md)	 - 

