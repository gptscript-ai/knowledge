---
title: Config File
---

# Config File

:::warning

    The config file format is subject to change as it's still in development.

:::

:::note

    This global configuration file is independent from the [flow configuration files](11-flows/01-overview.md#flow-config-file---flows-file).

:::

## Usage

Using the config file is as simple as passing `-c <path>` or `--config-file <path>` to the knowledge CLI on [supported commands](99-cmd/knowledge.md).
You may as well use the `KNOW_CONFIG_FILE` environment variable to set the path to the config file.

## Configuration Overview

Here we try to capture all supported configuration items in one example.

:::note

    You can write the config in YAML or JSON format. 
    You can find some example config files in the [GitHub repository](https://github.com/gptscript-ai/knowledge/blob/main/examples/configfiles).    

:::

```yaml
embeddings:
  providers:
  - name: my-cohere
    type: cohere
    config:
      apiKey: "${COHERE_API_KEY}" # environment variables are expanded when reading the config file
      model: "embed-english-v2.0"
  - name: myopenai
    type: openai
    config:
      apiKey: "${OPENAI_API_KEY}"
      embeddingEndpoint: "/some-custom-endpoint" # anything that's not the default /embeddings
  - name: foobar
    type: vertex
    config:
      apiKey: "${GOOGLE_API_KEY}"
      project: "acorn-io"
      # apiEndpoint: https://us-central1-aiplatform.googleapis.com
      model: "text-embedding-004"
```

### Sections

- `embeddings`: See [Embedding Models](05-embedding_models.md) for more details.
  - Select a provider using the command line flag `--embedding-model-provider` or the environment variable `KNOW_EMBEDDING_MODEL_PROVIDER` (default: `openai`).
  - **Note**: If a provider is selected but not specified in the config file, we'll assume that it's a standard provider configured via standard environment variables.
    - E.g. you select `vertex`, but that name is not configured, so we default to `type=vertex` and use the `VERTEX_*` environment variables to configure a standard Google Vertex AI provider.