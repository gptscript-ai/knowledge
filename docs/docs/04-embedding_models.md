---
title: Embedding Models
---

# Embedding Models

## Generate Embeddings

Embeddings are automatically generated when ingesting a document.
Currently, this is part of the job of the vector store implementation (chromem-go).


## Choosing an Embedding Model

At the moment, knowledge only support OpenAI API compatible model provider endpoints, which you can configure using the following environment variables:

| Variable                 | Default                     | Notes                                 |
|--------------------------|-----------------------------|---------------------------------------|
| `OPENAI_BASE_URL`        | `https://api.openai.com/v1` | ---                                   |
| `OPENAI_API_KEY`         | `sk-foo`                    | -!-                                   |
| `OPENAI_EMBEDDING_MODEL` | `text-embedding-ada-002`    | ---                                   |
| `OPENAI_API_VERSION`     | `2024-02-01`                | for Azure                             |
| `OPENAI_API_TYPE`        | `OPEN_AI`                   | one of `OPEN_AI`, `AZURE`, `AZURE_AD` |

As you can see above, knowledge defaults to the `text-embedding-ada-002` model from OpenAI.
[Below](#example-configurations) you can see how to configure some other providers/models.

### Example Configurations

Here are some example configurations that we tested with the knowledge tool to confirm that they're working as expected.
This is no judgement on the quality of the models or how well they work with the knowledge tool in terms of retrieval accuracy.

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

::: note

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