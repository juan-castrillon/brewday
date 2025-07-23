package sql

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// createTables will create:
// - the summaries table and initialize the foreign key constrain with the recipes table
// - the stats table that is not constrained by a recipe
func createTables(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS summaries (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		mash_temp REAL,
		mash_notes TEXT,
		mash_rasts TEXT,
		lautern_info TEXT,
		hopping_vol_bb REAL,
		hopping_vol_bb_notes TEXT, 
		hopping_hops TEXT,
		hopping_vol_ab REAL,
		hopping_vol_ab_notes TEXT,
		cooling_temp REAL,
		cooling_time REAL,
		cooling_notes TEXT,
		pre_ferm_vols TEXT,
		yeast_start_temp TEXT,
		yeast_start_notes TEXT,
		main_ferm_sgs TEXT,
		main_ferm_alcohol REAL,
		main_ferm_dry_hops TEXT,
		bottling_pre_bottle_volume REAL,
		bottling_carbonation REAL,
		bottling_sugar_amount REAL,
		bottling_sugar_type TEXT,
		bottling_temperature REAL,
		bottling_alcohol REAL,
		bottling_volume_bottled REAL,
		bottling_notes TEXT,
		sec_ferm_days INTEGER,
		sec_ferm_notes TEXT,
		stats_evaporation REAL,
		stats_effiency REAL,
		timeline TEXT,
		recipe_id INTEGER NOT NULL,
		FOREIGN KEY (recipe_id) 
			REFERENCES recipes (id)
				ON DELETE CASCADE
				ON UPDATE CASCADE
	)`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS ix_summaries ON summaries (recipe_id)`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS stats (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			recipe_title TEXT UNIQUE,
			evaporation REAL,
			efficiency REAL
	)`) // Eventually here we can create an index if we want to search for particular stats
	if err != nil {
		return err
	}
	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS ix_stats ON stats (recipe_title)`)
	return err
}
