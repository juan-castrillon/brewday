package sql

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func createTable(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS timelines (
		id INTEGER NOT NULL PRIMARY KEY,
		event TEXT NOT NULL,
		timestamp_unix INTEGER NOT NULL,
		recipe_id INTEGER NOT NULL,
		FOREIGN KEY (recipe_id) 
			REFERENCES recipes (id)
				ON DELETE CASCADE
				ON UPDATE CASCADE
	)`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS ix_timelines ON timelines (recipe_id, timestamp_unix)`)
	return err
}
