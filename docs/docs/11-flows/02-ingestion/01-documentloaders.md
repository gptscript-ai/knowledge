# Documentloaders

The available document loaders define which file types can be ingested into the system.

## Available Document Loaders

### plaintext

### markdown

### html

### pdf

**Options**

- `Password`
- `StartPage`
- `MaxPages`
- `Source`
- `NumThread`

### csv

**Options**

- `Separator`
- `LazyQuotes`
- `Columns`

### notebook

Support for Jupyter Notebooks (.ipynb)

**Options**

- `IncludeOutputs`
- `Traceback`
- `MaxOutputLength`

### document

Support for docx, rtf, odt