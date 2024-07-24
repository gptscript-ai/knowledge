---
title: Embedding Models
---

# Embedding Models

## Generate Embeddings

Embeddings are automatically generated when ingesting a document.
Currently, this is part of the job of the vector store implementation ([chromem-go](https://github.com/philippgille/chromem-go)).


## Choosing an Embedding Model Provider

The knowledge tool supports multiple embedding model providers, which you can configure via the [global config file](04-configfile.md#configuration-overview) or via environment variables.
You can choose which of your configured providers to use by setting the `KNOW_EMBEDDING_MODEL_PROVIDER` environment variable or using the `--embedding-model-provider` flag.

:::note

    The default selected provider is **OpenAI**

:::

### [OpenAI](https://openai.com/) + [Azure](https://ai.azure.com/)

OpenAI and Azure are configured via a single provider configuration to make the configuration similar to the one used by GPTScript.

| Environment Variable      | Config Key       | Default                     | Notes                                 |
|---------------------------|------------------|-----------------------------|---------------------------------------|
| `OPENAI_BASE_URL`         | `baseURL`        | `https://api.openai.com/v1` | ---                                   |
| `OPENAI_API_KEY`          | `apiKey`         | `sk-foo`                    | **required**                          |
| `OPENAI_EMBEDDING_MODEL`  | `embeddingModel` | `text-embedding-ada-002`    | ---                                   |
| `OPENAI_API_TYPE`         | `apiType`        | `OPEN_AI`                   | one of `OPEN_AI`, `AZURE`, `AZURE_AD` |
| `OPENAI_API_VERSION`      | `apiVersion`     | `2024-02-01`                | for **Azure**                         |
| `OPENAI_AZURE_DEPLOYMENT` |                  | ``                          | **required** for **Azure**            |


#### OpenAI Compatible Providers

We have some first-class supported providers that are compatible with OpenAI API, that you can find on this page.
If yours is not in the list, you can still try to configure it using the OpenAI provider configuration as shown below for LM-Studio and Ollama.

<details>
<summary id="example-configurations-lm-studio"><strong>LM-Studio</strong></summary>

LM-Studio failed to return any embeddings if requested concurrently, so we set the parallel threads to 1.
This may change in the future. Tested with LM-Studio v0.2.27.


```dotenv
export OPENAI_BASE_URL=http://localhost:1234/v1
export OPENAI_API_KEY=lm-studio
export OPENAI_EMBEDDING_MODEL="CompendiumLabs/bge-large-en-v1.5-gguf"
export VS_CHROMEM_EMBEDDING_PARALLEL_THREAD="1"
```

:::note

    Running with VS_CHROMEM_EMBEDDING_PARALLEL_THREAD="1" may be really really slow for a large amount of files (or just really large files).

:::

</details>

<details>
<summary id="example-configurations-ollama"><strong>Ollama</strong></summary>

Tested with Ollama v0.2.6 (pre-release that introduced OpenAI API compatibility).


```dotenv
export OPENAI_BASE_URL=http://localhost:11434/v1
export OPENAI_EMBEDDING_MODEL="mxbai-embed-large"
```

</details>

### [Cohere](https://cohere.com/)

| Environment Variable | Config Key | Default              | Notes        |
|----------------------|------------|----------------------|--------------|
| `COHERE_API_KEY`     | `apiKey`   | ---                  | **required** |
| `COHERE_MODEL`       | `model`    | `embed-english-v3.0` | ---          |

### [Jina](https://jina.ai/)

| Environment Variable | Config Key | Default                      | Notes        |
|----------------------|------------|------------------------------|--------------|
| `JINA_API_KEY`       | `apiKey`   | ---                          | **required** |
| `JINA_MODEL`         | `model`    | `jina-embeddings-v2-base-en` | ---          |

### [LocalAI](https://localai.io/)

| Environment Variable | Config Key | Default              | Notes |
|----------------------|------------|----------------------|-------|
| `LOCALAI_MODEL`      | `model`    | `bert-cpp-minilm-v6` | ---   |

### [Mistral](https://mistral.ai/)

| Environment Variable | Config Key | Default | Notes        |
|----------------------|------------|---------|--------------|
| `MISTRAL_API_KEY`    | `apiKey`   | ---     | **required** |

### [Mixedbread](https://www.mixedbread.ai/)

| Environment Variable | Config Key | Default            | Notes        |
|----------------------|------------|--------------------|--------------|
| `MIXEDBREAD_API_KEY` | `apiKey`   | ---                | **required** |
| `MIXEDBREAD_MODEL`   | `model`    | `all-MiniLM-L6-v2` | ---          |

### [Ollama](https://ollama.com/)

Requires Ollama v0.2.6 or later.

| Environment Variable | Config Key | Default                     | Notes |
|----------------------|------------|-----------------------------|-------|
| `OLLAMA_BASE_URL`    | `baseURL`  | `http://localhost:11434/v1` | ---   |
| `MIXEDBREAD_MODEL`   | `model`    | `mxbai-embed-large`         | ---   |