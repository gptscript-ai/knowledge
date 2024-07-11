# Transformers

## Available Transformers


### extra_metadata

Add extra metadata to every document.

**Options** 

- `metadata`

### filter_markdown_docs_no_content

Filter out markdown-formatted documents that do not have any content (e.g. only whitespace/newline/headings).

### keywords

Let the LLM extract some keywords from the document content and add them to the metadata (field `keywords`).

**Options**

- `NumKeywords`
- `LLM`