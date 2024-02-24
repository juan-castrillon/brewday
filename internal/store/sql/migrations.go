package sql

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func createTable(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS recipes (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		style TEXT,
		batch_size_l REAL,
		initial_sg REAL,
		ibu REAL,
		ebc REAL,
		status INTEGER NOT NULL, 
		status_args TEXT,
		mash_malts TEXT,
		mash_main_water REAL,
		mash_nachguss REAL,
		mash_temp REAL,
		mash_out_temp REAL,
		mash_rasts TEXT,
		hop_cooking_time REAL,
		hop_hops TEXT,
		hop_additional TEXT,
		ferm_yeast TEXT,
		ferm_temp TEXT,
		ferm_additional TEXT,
		ferm_carbonation REAL
	)`)
	return err
}

func createResultsTable(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS recipe_results (
		id INTEGER NOT NULL PRIMARY KEY,
		hot_wort_vol REAL,
		original_sg REAL,
		final_sg REAL,
		alcohol REAL,
		main_ferm_vol REAL,
		recipe_id INTEGER NOT NULL,
		FOREIGN KEY (recipe_id) 
			REFERENCES recipes (id)
				ON DELETE CASCADE
				ON UPDATE CASCADE
	)`)
	return err
}