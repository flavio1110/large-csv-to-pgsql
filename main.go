package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
)

var fileName = "people.csv"

// Connection string based on the DB created using the compose file
var connString = "postgres://user:super-secret@localhost:5432/people?sslmode=disable"

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		log.Fatal("Provide the command: gen, import-stream, or import-read-all")
	}

	switch args[0] {
	case "gen":
		generateLargeFile()
	case "import-stream":
		importFile(importWithStream, "stream")
	case "import-read-all":
		importFile(importReadAll, "read-all")
	default:
		log.Fatal("Provide the command: gen, import-stream, or import-read-all")
	}
}

func generateLargeFile() {
	f, err := os.Create(fileName)
	if err != nil {
		log.Fatal("Failed to create file:", err.Error())
	}
	defer f.Close()

	// header
	f.WriteString("first name,last name,city\n")

	// 1M rows ~16MB
	for i := 0; i < 1_000_000; i++ {
		f.WriteString("John,Doe,New York\n")
	}
}

type importer func(ctx context.Context, csvStream io.Reader, db *sql.DB) error

func importFile(imp importer, id string) {
	defer profileMemory(id)
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal("unable to open file", err)
	}
	defer file.Close()

	db, err := sql.Open("pgx", connString)
	if err != nil {
		log.Fatal("unable to open DB", err)
	}
	defer db.Close()
	defer cleanupTable(db)

	if err := imp(context.Background(), file, db); err != nil {
		log.Fatal("failed to import data ", err)
	}

}

func cleanupTable(db *sql.DB) {
	_, err := db.Exec("delete from people")
	if err != nil {
		log.Fatal("Failed to clean up table", err)
	}
}

func profileMemory(id string) {
	fmt.Printf("===> %s", id)
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("\n\tAlloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	//fmt.Printf("\tFrees = %v MiB", bToMb(m.Frees))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v", m.NumGC)
	fmt.Printf("\tObjs = %v\n", m.Mallocs)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
