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
		vol_bb REAL,
		recipe_id INTEGER NOT NULL,
		FOREIGN KEY (recipe_id) 
			REFERENCES recipes (id)
				ON DELETE CASCADE
				ON UPDATE CASCADE
	)`)
	return err
}

func createSGsTable(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS main_ferm_sgs (
		id INTEGER NOT NULL PRIMARY KEY,
		sg REAL,
		date TEXT,
		recipe_id INTEGER NOT NULL,
		FOREIGN KEY (recipe_id) 
			REFERENCES recipes (id)
				ON DELETE CASCADE
				ON UPDATE CASCADE
	)`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS ix_main_ferm_sgs ON main_ferm_sgs (recipe_id, id)`)
	return err
}

func createTimeTable(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS dates (
		id INTEGER NOT NULL PRIMARY KEY,
		date TEXT,
		name TEXT,
		recipe_id INTEGER NOT NULL,
		FOREIGN KEY (recipe_id) 
			REFERENCES recipes (id)
				ON DELETE CASCADE
				ON UPDATE CASCADE
	)`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS ix_dates ON dates (recipe_id)`)
	return err
}

func createSugarResultsTable(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS sugar_results (
		id INTEGER NOT NULL PRIMARY KEY,
		water REAL,
		sugar REAL,
		alcohol REAL,
		recipe_id INTEGER NOT NULL,
		FOREIGN KEY (recipe_id) 
			REFERENCES recipes (id)
				ON DELETE CASCADE
				ON UPDATE CASCADE
	)`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS ix_sugar_results ON sugar_results (recipe_id, id)`)
	return err
}

func createBoolFlagsTable(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS bool_flags (
		id INTEGER NOT NULL PRIMARY KEY,
		value INTEGER,
		name TEXT,
		recipe_id INTEGER NOT NULL,
		FOREIGN KEY (recipe_id) 
			REFERENCES recipes (id)
				ON DELETE CASCADE
				ON UPDATE CASCADE
	)`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS ix_bool_flags ON dates (recipe_id, name)`)
	return err
}
