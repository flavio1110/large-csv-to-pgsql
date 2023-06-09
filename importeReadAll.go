package main

import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
)

func importReadAll(ctx context.Context, csvStream io.Reader, db *sql.DB) error {
	conn, err := db.Conn(ctx)
	if err != nil {
		return fmt.Errorf("failed to get Db connection fro sql.DB instance: %w", err)
	}

	rows, err := getRowsToBeInserted(csvStream)
	if err != nil {
		return fmt.Errorf("failed to get rows to insert: %w", err)
	}

	if err = conn.Raw(func(driverConn any) error {
		pgxConn := driverConn.(*stdlib.Conn).Conn()
		_, err := pgxConn.CopyFrom(ctx, pgx.Identifier{"people"}, peopleColumns, pgx.CopyFromRows(rows))
		if err != nil {
			return fmt.Errorf("failed to import data into database: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("failed to resolve raw conection: %w", err)
	}
	return nil
}

func getRowsToBeInserted(csvStream io.Reader) ([][]interface{}, error) {
	csvReader := csv.NewReader(csvStream)
	csvReader.ReuseRecord = true // reuse slice to return the record line by line
	csvReader.FieldsPerRecord = -1
	rows := make([][]any, 1_000_000)
	// number of rows of the csv. If you don't know it beforehand
	// it's worthy to iterate in the stream to count it so you can pre-allocate the space you need
	isHeader := true
	row := make([]any, len(peopleColumns))
	count := 0
	for {
		r, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				return rows, nil
			}

			return nil, fmt.Errorf("failed to read csv: %w", err)
		}
		if isHeader {
			isHeader = false
			continue
		}

		for i, column := range r {
			row[i] = column
		}
		rows[count] = row
		count++
	}
}
