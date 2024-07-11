# Retrievers

## Available Retrievers

### basic

The default retriever. No fuzz, just similarity search.

**Options**

- `TopK`

Example: [examples/advanced_config.yaml](https://github.com/gptscript-ai/knowledge/blob/main/examples/advanced_config.yaml)

### subquery

A relict from the past. It's a subquery retriever that uses the LLM to generate subqueries.
As of now, you can achieve the same with a query modifier combined with a normal retriever.
There are probably pros and cons for either way.

**Options**

- `Model`
- `Limit`
- `TopK`

Example: [examples/subquery_retriever.yaml](https://github.com/gptscript-ai/knowledge/blob/main/examples/subquery_retriever.yaml)

### routing

The routing retriever lets the LLM choose which datasets it wants to query.
It presents the LLM with a list of available datasets along with their metadata (description, etc. if present).

**Options**

- `Model`
- `AvailableDatasets`
- `TopK`

Example: [examples/subquery_retriever.yaml](https://github.com/gptscript-ai/knowledge/blob/main/examples/subquery_retriever.yaml)