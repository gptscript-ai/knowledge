---
title: "Sharing Datasets"
---

# Sharing Datasets

Share your Knowledge!
Datasets are not limited to live on your machine. You can share them with others.

You can even export a dataset and distribute it with another GPTScript tool to add tool-specific Knowledge to the mix.
Let's say you're creating a GPTScript tool for a specific software. Then you can create a dataset and feed it with all the information you have about that software.
Then export the dataset and bundle it with your tool, so it always has access to this knowledge base without having to look it up every time.

## Exporting a Dataset

To export a dataset, you can use the `export` command. This will create a `.zip` file that you can share with others.

<details>
<summary>Create a Dataset if you don't have one already</summary>

```bash
knowledge create-dataset my-dataset
knowledge ingest --dataset my-dataset docs/docs/10-datasets/01-managing.md
```

</details>

```bash
knowledge export my-dataset --output my-dataset.zip
```

:::note

1. You can export multiple datasets by specifiying multiple Dataset IDs or by using the `--all` flag.
2. If you leave out the `--output my-dataset.zip` flag, Knowledge will create an export file with a name like `knowledge-export-2024-07-11-10-13-08.zip` in the current directory.

:::

## Importing a Dataset

:::warning

Importing a Dataset works just fine, but there's a culprit when you want to **ingest additional content into an imported dataset**: You'll have to use the exact same embedding function as the original dataset.
The Embedding function is part of the Vector Database implementation and defines how the content is transformed into a vector representation.
Currently, this is defined solely based on the model provider configuration, so it's fairly simple to replicate - you just have to use the same model (`$OPENAI_EMBEDDING_MODEL`) for it to work.

:::

To import a dataset, you can use the `import` command. This will import the dataset from the `.zip` file.

<details>
<summary>Delete the test dataset to prove that it's working</summary>

```bash
knowledge delete-dataset my-dataset
```

</details>

```bash
knowledge import my-dataset.zip
```

:::note

Import only selected datasets from the `.zip` file by specifying the Dataset IDs (no flags required).

:::

## Using an exported Dataset without importing it

You want to leverage the knowledge from an exported dataset, but have no need of extending it with additional data?
In that case you don't even have to import it into your Knowledge datastore. Just use the `--archive` flag.

```bash
knowledge retrieve --archive my-dataset.zip "How do I add metadata to a dataset?"
```

The `--archive` flag is available for the following commands:

- `retrieve`
- `get-dataset`
- `list-datasets`

