[![Lint Build Test](https://github.com/flavio1110/large-csv-to-pgsql/actions/workflows/ci.yaml/badge.svg?branch=main)](https://github.com/flavio1110/large-csv-to-pgsql/actions/workflows/ci.yaml)

## Importing data to Postgres using pgx without loading in memory

Example of importing data into the database without loading the entire file into memory.
Implementing the `CopyFromSource` of `pgx`.

Comparing the memory consumption for each approach importing a file with ~16MB (1M rows).
> You can read more about the meaning of each metric on <https://golang.org/pkg/runtime/#MemStats>.

| Approach/metric  |     |  TotalAlloc |     |          Sys |
|------------------|-----|------------:|-----|-------------:|
| Stream file      |     |      61 MiB |     |       12 Mib |
| Read entire file |     |      84 MiB |     |       58 Mib |
|                  |     | **+37.70%** |     | **+346.15%** |

Run the tests yourself:

### Pre-requisites to Run:
1. [Docker installed](https://www.docker.com/products/docker-desktop/)
2. Run postgres

```shell
make create-db
```
3. Generate test file to run the tests
It will generate a CSV file with ~16MB.

```shell
make gen-file
```
### Run both versions
It will run the application with each implementation variation.
It will populate the DB with the generated file and clean the table afterwards.
It will also output the overall memory consumption.
```shell
make import-both 
```

### Run benchmark both versions
It will benchmark each implementation variation including memory comsumption.
It will populate the DB with the generated file and clean the table afterwards.
It will also output the overall memory consumption.
```shell
make benchmark 
```