# Knowledge API

Standalone Knowledge API Server to be used with GPTScript and GPTStudio

## Run

Just run this from the repository root:

```bash
make run # it will be available at http://localhost:8000
```

## Supported File Types

- `.md`
- `.txt`
- `.html`
- `.pdf`
- `.ipynb`
- `.csv`
- `.docx`
- `.rtf`
- `.odt`

## OpenAPI / Swagger

The API is documented using OpenAPI 2.0 (Swagger), automatically generated using [`swaggo/swag`](https://github.com/swaggo/swag) (`make openapi`).
