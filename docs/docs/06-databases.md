---
title: Index & Vector Databases
---

# Index & Vector Databases

## Index Database

The index database is an additional (relational) metadata database which keeps track of all datasets and ingested files and their relationships.
It enables some extra convenience features but does not store the actual data (content & embeddings).
The current implementation uses **SQLite** by default, which is fully embedded and does not require any additional setup.

You can configure it by setting a database connection string via the `KNOW_INDEX_DSN` environment variable. 
The following options are available:

- [SQLite](https://www.sqlite.org/) (default): `KNOW_INDEX_DSN="sqlite:///home/me/mysqlite.db"`
- [Postgres](https://www.postgresql.org/): `KNOW_INDEX_DSN="postgres://knowledge:knowledge@localhost:5432/knowledge?sslmode=disable"`


## Vector Database

The vector database is the main storage for the content and embeddings of the ingested documents along with some metadata (e.g. source file information).
The current implementation uses [**chromem-go**](https://github.com/philippgille/chromem-go) by default, which is fully embedded and does not require any additional setup.

You can configure it by setting a database connection string via the `KNOW_VECTOR_DSN` environment variable. 
The following options are available:

- [Chromem-Go](https://github.com/philippgille/chromem-go) (default): `KNOW_VECTOR_DSN="chromem:///path/to/directory"` (Note: we're using a customized fork of chromem-go, so some details may differ from the original project)
- [PGVector](https://github.com/pgvector/pgvector): `KNOW_VECTOR_DSN="pgvector://knowledge:knowledge@localhost:5432/knowledge?sslmode=disable"`
- [SQLite-Vec](https://github.com/asg017/sqlite-vec): `KNOW_VECTOR_DSN="sqlite-vec:///home/me/mysqlite.db"`