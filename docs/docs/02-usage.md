---
title: Usage
---

# Using Knowledge

The knowledge tool itself has two modes of operation: Standalone and Server Mode - Check the sections below to learn more about them.

Both modes are configured the same way, via environment variables or command line flags:

## Configuration

### Embedding Model Provider (must have)

The model provider is the provider of the embeddings model that is used to encode ingested documents.
Currently, we only support **OpenAI** and **Azure OpenAI** via the following flags / environment variables:

```bash
--openai-api-base string           OpenAI API base ($OPENAI_BASE_URL) (default "https://api.openai.com/v1")
--openai-api-key string            OpenAI API key ($OPENAI_API_KEY) (default "sk-foo")
--openai-api-type string           OpenAI API type (OPEN_AI, AZURE, AZURE_AD) ($OPENAI_API_TYPE) (default "OPEN_AI")
--openai-api-version string        OpenAI API version (for Azure) ($OPENAI_API_VERSION) (default "2024-02-01")
--openai-azure-deployment string   Azure OpenAI deployment name (overrides openai-embedding-model, if set) ($OPENAI_AZURE_DEPLOYMENT)
--openai-embedding-model string    OpenAI Embedding model ($OPENAI_EMBEDDING_MODEL) (default "text-embedding-ada-002")
```

Those are persistent flags, so they can be set on any knowledge subcommand.


## 1. Standalone Mode (Default)

In standalone mode, Knowledge makes use of an embedded database and embedded Vector Database which the client connects to directly.
This is the default and most simple mode of operation and is useful for local usage and offers a great integration with GPTScript.

### Run the Client

Any `knowledge` command (except for `knowledge server`) will use the standalone client mode, if no `KNOW_SERVER_URL` environment variable is set.

## 2. Server Mode

In server mode, Knowledge uses a separate server for the Vector Database and the Document Database.
This mode is useful when you want to share the data with multiple clients or when you want to use a more powerful server for the Vector Database.

### Run the Server

```bash
knowledge server
```

## Ingestion

To ingest a document, you can use the `knowledge ingest` command:

```bash
knowledge ingest --dataset my-dataset ./path/to/my-document.txt
```

:::note

    By default, the dataset will be created if it doesn't exist.
    If you don't want that, you can use the `--no-create-dataset` flag.

:::

### Ignoring Files

You can ignore files by providing an ignore file, similar to `.gitignore`:

```bash
knowledge ingest --dataset my-dataset --ignore-file .knowledgeignore ./path/to/my-documents
```

Here's an example ignore file which basically tells knowledge to only consider Markdown files and nothing else:

```gitignore
# Ignore everything
*

# Except Markdown files in any directory
!**/*.md   
```

:::note

    Alternatively, you can use the `--ignore-extensions` flag to ignore files with specific extensions.

    ```bash
    knowledge ingest --dataset my-dataset --ignore-extensions=.txt ./path/to/my-documents
    ```

:::


### Remote Files

You can ingest remote files by providing a URL - Currently only Git Repositories are supported:

```bash
knowledge ingest --dataset my-dataset https://github.com/gptscript-ai/knowledge
```

You may also specify a Commit Hash, Tag or Branch that you want to check out:

```bash
knowledge ingest --dataset my-dataset https://github.com/gptscript-ai/knowledge@v0.3.0
```

:::note

    Here, it's advisable to use a [ignore file](#ignoring-files) to avoid ingesting all the git metadata and potentially present vendor files.

:::