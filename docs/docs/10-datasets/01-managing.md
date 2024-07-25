---
title: "Managing Datasets"
---

# Managing Datasets

Knowledge lets you create multiple datasets.
Each dataset is a separate knowledge base that you can feed files into.
In the Vector Database, each dataset is stored as a separate collection of vectors (vector space) that Knowledge will search against.

## Dataset Lifecycle

The Knowledge tool has a `default` dataset that will be used whenever no dataset is specified via flag or environment variable.
When using `askdir`, a dataset will be created automatically with an ID generated from the absolute path to the target directory.
You can also create datasets manually using the `create-dataset` command - see below.

### Create a Dataset

```bash
# knowledge create-dataset <dataset-id> [flags]
knowledge create-dataset my-dataset
```

<details>
<summary>Expected Output</summary>

```bash
$ knowledge create-dataset my-dataset
2024/07/11 09:45:27 INFO Created dataset id=my-dataset
Created dataset "my-dataset"
```

</details>


### Ingest a file into a specific Dataset

```bash
knowledge ingest --dataset my-dataset docs/docs/10-datasets/01-managing.md
```

<details>
<summary>Expected Output</summary>

```bash
$ knowledge ingest --dataset my-dataset docs/docs/10-datasets/01-managing.md
2024/07/11 09:50:53 INFO Ingested document filename=01-managing.md count=6 absolute_path=<path>/docs/docs/10-datasets/01-managing.md
Ingested 1 files from "docs/docs/10-datasets/01-managing.md" into dataset "my-dataset"
```

</details>

### Get existing Datasets

If you just want an overview of what datasets exist in your Knowledge datastore, list them out like this:

```bash
knowledge list-datasets
```

<details>
<summary>Expected Output</summary>

```bash
$ knowledge list-datasets                                          
[{"id":"default"},{"id":"my-dataset"}]
```

</details>

If you need more details and even want to see, what files/documents are inside a dataset, get it specifically:

```bash
knowledge get-dataset my-dataset
```

<details>
<summary>Expected Output</summary>

```bash
$ knowledge get-dataset my-dataset                                          
{"id":"my-dataset","Files":[{"id":"4c107370-3f5a-11ef-8438-8cf8c5751845","dataset":"my-dataset","Documents":[{"id":"0ee6aa14-3b60-4f76-9a1d-55b5f6d105a5","dataset":"my-dataset","file_id":"4c107370-3f5a-11ef-8438-8cf8c5751845"},{"id":"f993e6c6-c122-426c-900a-589cf3037993","dataset":"my-dataset","file_id":"4c107370-3f5a-11ef-8438-8cf8c5751845"},{"id":"3142d3ec-4548-42b4-950c-a272002fbe60","dataset":"my-dataset","file_id":"4c107370-3f5a-11ef-8438-8cf8c5751845"},{"id":"e3ea75df-f3a2-4daf-9ce8-2f7ba16ad130","dataset":"my-dataset","file_id":"4c107370-3f5a-11ef-8438-8cf8c5751845"},{"id":"37d71728-a78c-4c82-8e51-654c5e224d9a","dataset":"my-dataset","file_id":"4c107370-3f5a-11ef-8438-8cf8c5751845"},{"id":"159bdf1f-697c-4e0d-bffe-f7472b03c5b5","dataset":"my-dataset","file_id":"4c107370-3f5a-11ef-8438-8cf8c5751845"}],"name":"01-managing.md","absolute_path":"<path>/docs/docs/10-datasets/01-managing.md","size":2006,"modified_at":"2024-07-11T09:50:31.981617707+02:00"}]}
```

</details>

:::note

Too many details in the `get-dataset` output? Use the `--no-docs` flag to hide the information about included documents.

:::

### Modify a Dataset

You can modify existing Datasets, e.g. to provide more information about them in the form of [metadata](#dataset-metadata).
(In fact, at the moment, editing Dataset metadata is the only thing you can do here)

```bash
knowledge edit-dataset my-dataset --update-metadata description="Technical Documentation about the GPTScript Knowledge Tool"
```

<details>
<summary>Expected Output</summary>

```bash
$ knowledge edit-dataset my-dataset --update-metadata description="Technical Documentation about the GPTScript Knowledge Tool"
Updated dataset:
 {"id":"my-dataset","Files":null,"metadata":{"description":"Technical Documentation about the GPTScript Knowledge Tool"}}
```

</details>

### EOL: Delete a Dataset

The Dataset got dirty, or you just don't need it anymore? Get rid of it!

```bash
knowledge delete-dataset my-dataset
```

<details>
<summary>Expected Output</summary>

```bash
$ knowledge delete-dataset my-dataset
2024/07/11 10:01:11 INFO Deleting dataset id=my-dataset
Deleted dataset "my-dataset"
```

</details>

## Dataset Metadata

Dataset can have arbitrary metadata in the form of key-value pairs (map, dict, whatever you want to call it).
This metadata serves two purposes (at the moment):

1. Help humans to understand what's inside a dataset
2. Help an LLM to understand what's inside a dataset
   - This is quite interesting, e.g. when used with the **routing retriever**, which presents the list of available datasets to the LLM and the LLM decides, which dataset(s) to query to get a meaningful answer.