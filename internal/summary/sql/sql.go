package sql

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type SummaryRecorderPersistentStore struct {
	dbClient *sql.DB
}

func NewSummaryRecorderPersistentStore(db *sql.DB) (*SummaryRecorderPersistentStore, error) {
	err := createTable(db)
	if err != nil {
		return nil, err
	}
	return &SummaryRecorderPersistentStore{
		dbClient: db,
	}, nil
}
