package main

import (
	"context"
	"database/sql"
	"log"
	"os"
)

func main() {
	file, err := os.Open("people.csv")
	if err != nil {
		log.Fatal("unable to open file", err)
	}

	db, err := sql.Open("pgx", "postgres://user:super-secret@localhost:5432/people?sslmode=disable")
	if err != nil {
		log.Fatal("unable to open DB", err)
	}

	if err := importData(context.Background(), file, db); err != nil {
		log.Fatal("failed to import data", err)
	}
}
