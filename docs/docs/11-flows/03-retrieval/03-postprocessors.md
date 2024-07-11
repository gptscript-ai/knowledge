# Postprocessors

## Available Postprocessors

### extra_metadata

Add some extra metadata to every document.

**Options**

- `metadata`

### filter_markdown_docs_no_content

Shouldn't be required anymore at this stage in your pipeline, but still: drop any document that doesn't have any content (considering markdown syntax).

### keywords

Use the LLM to extract keywords from the content of the documents and add them to the metadata (field: `keywords`).

**Options:**

- `NumKeywords`
- `LLM`

### similarity

Drop any document that doesn't have a similarity score above a certain threshold.

**Options**

- `Threshold`

### content_substring_filter

Drop any document where the content contains or doesn't contain some substrings.

**Options**

- `Contains`
- `NotContains`

### content_filter

Drop any document where the content doesn't match a certain criteria.
The LLM judges whether the content matches the criteria or not.

**Options**

- `Question`
- `Include`
- `LLM`

### cohere_rerank

Use Cohere's reranking API to rerank the documents based on their relevance regarding the input query.

**Options**

- `ApiKey`
- `Model`
- `TopN`