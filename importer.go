package main

import (
	"context"
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"io"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func importData(ctx context.Context, csvStream io.Reader, db *sql.DB) error {
	conn, err := db.Conn(ctx)
	if err != nil {
		return fmt.Errorf("failed to get Db connection fro sql.DB instance: %w", err)
	}

	if err = conn.Raw(func(driverConn any) error {
		pgxConn := driverConn.(*stdlib.Conn).Conn()
		_, err := pgxConn.CopyFrom(ctx, pgx.Identifier{"people"}, peopleColumns, newPeopleCopyFromSource(csvStream))
		if err != nil {
			return fmt.Errorf("failed to import data into database: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("failed to resolve raw conection: %w", err)
	}
	return nil
}

var peopleColumns = []string{
	"first_name",
	"last_name",
	"city",
}

func newPeopleCopyFromSource(csvStream io.Reader) *peopleCopyFromSource {
	csvReader := csv.NewReader(csvStream)
	csvReader.ReuseRecord = true // reuse slice to return the record line by line
	csvReader.FieldsPerRecord = -1

	return &peopleCopyFromSource{
		reader: csvReader,
		isBOF:  true, // first line is header
		record: make([]interface{}, len(peopleColumns)),
	}
}

type peopleCopyFromSource struct {
	reader        *csv.Reader
	err           error
	currentCsvRow []string
	record        []interface{}
	isEOF         bool
	isBOF         bool
}

func (pfs *peopleCopyFromSource) Values() ([]any, error) {
	if pfs.isEOF {
		return nil, nil
	}

	if pfs.err != nil {
		return nil, pfs.err
	}

	// the order of the elements of the record array, must match with
	// the order of the columns in passed into the copy method
	pfs.record[0] = pfs.currentCsvRow[0]
	pfs.record[1] = pfs.currentCsvRow[1]
	pfs.record[2] = pfs.currentCsvRow[2]
	return pfs.record, nil
}

func (pfs *peopleCopyFromSource) Next() bool {
	pfs.currentCsvRow, pfs.err = pfs.reader.Read()
	if pfs.err != nil {

		// when get to the end of the file return false and clean the error
		if errors.Is(pfs.err, io.EOF) {
			pfs.isEOF = true
			pfs.err = nil
		}
		return false
	}

	if pfs.isBOF {
		pfs.isBOF = false
		return pfs.Next()
	}

	return true
}

func (pfs *peopleCopyFromSource) Err() error {
	return pfs.err
}
