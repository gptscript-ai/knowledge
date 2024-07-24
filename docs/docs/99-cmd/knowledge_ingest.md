---
title: "knowledge ingest"
---
## knowledge ingest

Ingest a file/directory into a dataset

```
knowledge ingest [--dataset <dataset-id>] <path> [flags]
```

### Options

```
      --auto-migrate string                 Auto migrate database ($KNOW_DB_AUTO_MIGRATE) (default "true")
      --concurrency int                     Number of concurrent ingestion processes ($KNOW_INGEST_CONCURRENCY) (default 10)
  -c, --config-file string                  Path to the configuration file ($KNOW_CONFIG_FILE)
  -d, --dataset string                      Target Dataset ID ($KNOW_TARGET_DATASET) (default "default")
      --dsn string                          Server database connection string (default "sqlite://$XDG_DATA_HOME/gptscript/knowledge/knowledge.db") ($KNOW_DB_DSN)
      --embedding-model-provider string     Embedding model provider ($KNOW_EMBEDDING_MODEL_PROVIDER) (default "openai")
      --flow string                         Flow name ($KNOW_FLOW)
      --flows-file string                   Path to a YAML/JSON file containing ingestion/retrieval flows ($KNOW_FLOWS_FILE)
  -h, --help                                help for ingest
      --ignore-extensions string            Comma-separated list of file extensions to ignore ($KNOW_INGEST_IGNORE_EXTENSIONS)
      --no-recursive                        Don't recursively ingest directories ($KNOW_NO_INGEST_RECURSIVE)
      --server string                       URL of the Knowledge API Server ($KNOW_SERVER_URL)
      --textsplitter-chunk-overlap int      Textsplitter Chunk Overlap ($KNOW_TEXTSPLITTER_CHUNK_OVERLAP) (default 256)
      --textsplitter-chunk-size int         Textsplitter Chunk Size ($KNOW_TEXTSPLITTER_CHUNK_SIZE) (default 1024)
      --textsplitter-encoding-name string   Textsplitter Encoding Name ($KNOW_TEXTSPLITTER_ENCODING_NAME) (default "cl100k_base")
      --textsplitter-model-name string      Textsplitter Model Name ($KNOW_TEXTSPLITTER_MODEL_NAME) (default "gpt-4")
      --vector-dbpath string                VectorDBPath to the vector database (default "$XDG_DATA_HOME/gptscript/knowledge/vector.db") ($KNOW_VECTOR_DB_PATH)
```

### SEE ALSO

* [knowledge](knowledge.md)	 - 

