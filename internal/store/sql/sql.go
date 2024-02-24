package sql

import (
	"brewday/internal/recipe"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type PersistentStore struct {
	dbClient        *sql.DB
	insertStatement *sql.Stmt
	// updateStatement   *sql.Stmt
	// retrieveStatement *sql.Stmt
	// deleteStatement   *sql.Stmt
}

type MarshallResult struct {
	StatusParams string
	MashingMalts string
	MashingRasts string
	HopHops      string
	HopAdd       string
	FermAdd      string
	Yeast        string
}

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

func NewPersistentStore(path string) (*PersistentStore, error) {
	db, err := sql.Open("sqlite3", "file:"+path+"?_foreign_keys=true")
	if err != nil {
		return nil, err
	}
	err = createTable(db)
	if err != nil {
		return nil, err
	}
	err = createResultsTable(db)
	if err != nil {
		return nil, err
	}
	is, err := db.Prepare(`INSERT INTO recipes 
	(
		name, style, batch_size_l, initial_sg, ibu, ebc, status, status_args,
		mash_malts, mash_main_water, mash_nachguss, mash_temp, mash_out_temp, mash_rasts,
		hop_cooking_time, hop_hops, hop_additional,
		ferm_yeast, ferm_temp, ferm_additional, ferm_carbonation
	) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return nil, err
	}
	// us, err := db.Prepare("UPDATE my_objects SET name = ? WHERE id == ?")
	// if err != nil {
	// 	return nil, err
	// }
	// rs, err := db.Prepare("SELECT id, name, count, simple_name FROM my_objects WHERE id == ?")
	// if err != nil {
	// 	return nil, err
	// }
	// ds, err := db.Prepare("DELETE FROM my_objects WHERE id == ?")
	// if err != nil {
	// 	return nil, err
	// }

	return &PersistentStore{
		dbClient:        db,
		insertStatement: is,
		// updateStatement:   us,
		// retrieveStatement: rs,
		// deleteStatement:   ds,
	}, nil
}

// Store stores a recipe and returns an identifier that can be used to retrieve it
func (s *PersistentStore) Store(r *recipe.Recipe) (string, error) {
	status, _ := r.GetStatus()
	marshalled := s.marshallArrays(r)
	res, err := s.insertStatement.Exec(
		r.Name, r.Style, r.BatchSize, r.InitialSG, r.Bitterness, r.ColorEBC, status, marshalled.StatusParams,
		marshalled.MashingMalts, r.Mashing.MainWaterVolume, r.Mashing.Nachguss, r.Mashing.MashTemperature, r.Mashing.MashOutTemperature, marshalled.MashingRasts,
		r.Hopping.TotalCookingTime, marshalled.HopHops, marshalled.HopAdd,
		marshalled.Yeast, r.Fermentation.Temperature, marshalled.FermAdd, r.Fermentation.Carbonation,
	)
	if err != nil {
		return "", err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(id, 10), nil
}

// Retrieve retrieves a recipe based on an identifier
func (s *PersistentStore) Retrieve(id string) (*recipe.Recipe, error) {
	return nil, nil
}

// List lists all the recipes
func (s *PersistentStore) List() ([]*recipe.Recipe, error) {
	return nil, nil
}

// Delete deletes a recipe based on an identifier
func (s *PersistentStore) Delete(id string) error {
	return nil
}

func (s *PersistentStore) marshallArrays(r *recipe.Recipe) *MarshallResult {
	_, statusParams := r.GetStatus()
	var paramsString, malzString, rastString, hopsString, hopsAddString, fermAddString strings.Builder
	for _, p := range statusParams {
		paramsString.WriteString("[" + p + "]")
	}
	for _, m := range r.Mashing.Malts {
		malzString.WriteString(fmt.Sprintf("[%s-%.2f]", m.Name, m.Amount))
	}
	for _, rast := range r.Mashing.Rasts {
		rastString.WriteString(fmt.Sprintf("[%.2f-%.2f]", rast.Temperature, rast.Duration))
	}
	for _, h := range r.Hopping.Hops {
		hopsString.WriteString(fmt.Sprintf("[%s-%.2f-%.2f-%2.f-%t-%t]", h.Name, h.Alpha, h.Amount, h.Duration, h.DryHop, h.Vorderwuerze))
	}
	for _, ha := range r.Hopping.AdditionalIngredients {
		hopsAddString.WriteString(fmt.Sprintf("[%s-%.2f-%.2f]", ha.Name, ha.Amount, ha.Duration))
	}
	for _, fa := range r.Fermentation.AdditionalIngredients {
		fermAddString.WriteString(fmt.Sprintf("[%s-%.2f-%.2f]", fa.Name, fa.Amount, fa.Duration))
	}
	return &MarshallResult{
		StatusParams: paramsString.String(),
		MashingMalts: malzString.String(),
		MashingRasts: rastString.String(),
		HopHops:      hopsString.String(),
		HopAdd:       hopsAddString.String(),
		FermAdd:      fermAddString.String(),
		Yeast:        fmt.Sprintf("[%s-%.2f]", r.Fermentation.Yeast.Name, r.Fermentation.Yeast.Amount),
	}
}

// type MySimpleObject struct {
// 	AName string
// }

// type MyObject struct {
// 	ID             int
// 	Name           string
// 	Count          int
// 	MySimpleObject MySimpleObject
// }

// func (s *PersistentStore) create() error {
// 	const q = `CREATE TABLE IF NOT EXISTS my_objects (
// 		id INTEGER NOT NULL PRIMARY KEY,
// 		name TEXT,
// 		count INTEGER,
// 		simple_name TEXT
// 	)`
// 	_, err := s.dbClient.Exec(q)
// 	return err
// }

// func (s *PersistentStore) testSave(o *MyObject) error {
// 	err := s.create()
// 	if err != nil {
// 		return err
// 	}
// 	// stmt, _ := db.Prepare("INSERT INTO people (id, first_name, last_name, email, ip_address) VALUES (?, ?, ?, ?, ?)")
// 	stmt, err := s.dbClient.Prepare("INSERT INTO my_objects (id, name, count, simple_name) VALUES (?, ?, ?, ?)")
// 	if err != nil {
// 		return err
// 	}
// 	defer stmt.Close()
// 	_, err = stmt.Exec(o.ID, o.Name, o.Count, o.MySimpleObject.AName)
// 	return err
// }

// func (s *PersistentStore) testRetrieve(id int) (*MyObject, error) {
// 	err := s.create()
// 	if err != nil {
// 		return nil, err
// 	}
// 	stmt, err := s.dbClient.Prepare("SELECT id, name, count, simple_name FROM my_objects WHERE id == ?")
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer stmt.Close()
// 	//result := &MyObject{MySimpleObject: MySimpleObject{}}
// 	var result MyObject
// 	err = stmt.QueryRow(id).Scan(&result.ID, &result.Name, &result.Count, &result.MySimpleObject.AName)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &result, nil
// }

// func (s *PersistentStore) testUpdate(id int, name string) error {
// 	err := s.create()
// 	if err != nil {
// 		return err
// 	}
// 	stmt, err := s.dbClient.Prepare("UPDATE my_objects SET name = ? WHERE id == ?")
// 	if err != nil {
// 		return err
// 	}
// 	defer stmt.Close()
// 	_, err = stmt.Exec(name, id)
// 	return err
// }

// func (s *PersistentStore) testDelete(id int) error {
// 	err := s.create()
// 	if err != nil {
// 		return err
// 	}
// 	stmt, err := s.dbClient.Prepare("DELETE FROM my_objects WHERE id == ?")
// 	if err != nil {
// 		return err
// 	}
// 	defer stmt.Close()
// 	_, err = stmt.Exec(id)
// 	return err
// }
