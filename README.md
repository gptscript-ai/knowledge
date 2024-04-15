<p align="center">
  <img src="src/static/img/icon.png" />
</p>

# knowledge-retrieval-api

Standalone Knowledge Retrieval API Server to be used with Rubra

## Development

- Run in development mode (hot-reloading): `make run-dev` (Requires `docker` and `compose`)
- Dependency Management: `uv`
- Linting & Formatting: `ruff`

## File Types

Currently, the following file types are supported for ingestion via llama-index' [`SimpleDirectoryReader`](https://docs.llamaindex.ai/en/stable/module_guides/loading/simpledirectoryreader.html#supported-file-types) interface:

- `.csv` - comma-separated values
- `.docx` - Microsoft Word
- `.epub` - EPUB ebook format
- `.hwp` - Hangul Word Processor
- `.ipynb` - Jupyter Notebook
- `.jpeg`, .jpg - JPEG image
- `.mbox` - MBOX email archive
- `.md` - Markdown
- `.mp3, .mp4` - audio and video
- `.pdf` - Portable Document Format
- `.png` - Portable Network Graphics
- `.ppt, .pptm, .pptx` - Microsoft PowerPoint

## Examples

You can use the GPTScript example in the examples/ directory to test the ingestion and querying parts of the API.
The GPTScript will do the following:

1. Ingest the llama2 Paper located as `examples/data/llama2.pdf` (only if it hasn't been ingested before)
2. Query the Dataset to tell us something about the topics "Truthfulness, Toxicity, and Bias"

The returned response should contain a reference to the source page.

Just run this from the repository root:

```bash
make run-dev # if you haven't already

# Create the dataset
curl -X 'POST' \
  'http://localhost:8000/datasets/create' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "name": "llama2",
  "embed_dim": 0
}'

# Run the GPTScript example
gptscript examples/example.gpt
```
