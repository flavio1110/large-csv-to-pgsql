## Importing data to Postgres using pgx without loading in memory

Example of importing data into the database without loading the entire file into memory.
Implementing the `CopyFromSource` of `pgx`.

### Pre-requisites to Run:
1. [Docker installed](https://www.docker.com/products/docker-desktop/)
2. Run postgres
```go
docker compose up -d
```

### Run App to import the file "people.csv" into the database
```go
go run .
```

### Check the imported data
```sql
select * from people p 
```