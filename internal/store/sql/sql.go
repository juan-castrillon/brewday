package sql

import (
	"brewday/internal/recipe"
	"database/sql"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type PersistentStore struct {
	dbClient        *sql.DB
	insertStatement *sql.Stmt
	// updateStatement   *sql.Stmt
	retrieveStatement *sql.Stmt
	listStatement     *sql.Stmt
	// deleteStatement   *sql.Stmt
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
	rs, err := db.Prepare(`SELECT 
		name, style, batch_size_l, initial_sg, ibu, ebc, status, status_args,
		mash_malts, mash_main_water, mash_nachguss, mash_temp, mash_out_temp, mash_rasts,
		hop_cooking_time, hop_hops, hop_additional,
		ferm_yeast, ferm_temp, ferm_additional, ferm_carbonation
	FROM recipes WHERE id == ?`)
	if err != nil {
		return nil, err
	}
	// Just getting what I need for now to display, if future calls to list require more, they are to be added here
	ls, err := db.Prepare("SELECT id, name, style, status FROM recipes")
	if err != nil {
		return nil, err
	}
	// ds, err := db.Prepare("DELETE FROM my_objects WHERE id == ?")
	// if err != nil {
	// 	return nil, err
	// }

	return &PersistentStore{
		dbClient:        db,
		insertStatement: is,
		// updateStatement:   us,
		retrieveStatement: rs,
		listStatement:     ls,
		// deleteStatement:   ds,
	}, nil
}

// Store stores a recipe and returns an identifier that can be used to retrieve it
func (s *PersistentStore) Store(r *recipe.Recipe) (string, error) {
	status, _ := r.GetStatus()
	marshalled, err := s.marshalStructs(r)
	if err != nil {
		return "", err
	}
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
	var name, style, fermTemp string
	var batchSizeL, initialSg, ibu, ebc, mashMainWater, mashNachguss, mashTemp, mashOutTemp, hopCooking, fermCarbonation float32
	var status recipe.RecipeStatus
	toUnmarshall := &MarshalResult{}
	err := s.retrieveStatement.QueryRow(id).Scan(&name, &style, &batchSizeL, &initialSg, &ibu, &ebc, &status, &toUnmarshall.StatusParams,
		&toUnmarshall.MashingMalts, &mashMainWater, &mashNachguss, &mashTemp, &mashOutTemp, &toUnmarshall.MashingRasts,
		&hopCooking, &toUnmarshall.HopHops, &toUnmarshall.HopAdd,
		&toUnmarshall.Yeast, &fermTemp, &toUnmarshall.FermAdd, &fermCarbonation)
	if err != nil {
		return nil, err
	}
	unmarshaled, err := s.unmarshalStructs(toUnmarshall)
	if err != nil {
		return nil, err
	}
	r := &recipe.Recipe{
		ID:         id,
		Name:       name,
		Style:      style,
		BatchSize:  batchSizeL,
		InitialSG:  initialSg,
		Bitterness: ibu,
		ColorEBC:   ebc,
		Mashing: recipe.MashInstructions{
			Malts:              unmarshaled.MashingMalts,
			MainWaterVolume:    mashMainWater,
			Nachguss:           mashNachguss,
			MashTemperature:    mashTemp,
			MashOutTemperature: mashOutTemp,
			Rasts:              unmarshaled.MashingRasts,
		},
		Hopping: recipe.HopInstructions{
			TotalCookingTime:      hopCooking,
			Hops:                  unmarshaled.HopHops,
			AdditionalIngredients: unmarshaled.HopAdd,
		},
		Fermentation: recipe.FermentationInstructions{
			Yeast:                 unmarshaled.Yeast,
			Temperature:           fermTemp,
			AdditionalIngredients: unmarshaled.FermAdd,
			Carbonation:           fermCarbonation,
		},
	}
	r.SetStatus(status, unmarshaled.StatusParams...)
	return r, nil
}

// List lists all the recipes
func (s *PersistentStore) List() ([]*recipe.Recipe, error) {
	rows, err := s.listStatement.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []*recipe.Recipe{}
	for rows.Next() {
		var status recipe.RecipeStatus
		var id, name, style string
		err = rows.Scan(&id, &name, &style, &status)
		if err != nil {
			return nil, err
		}
		r := &recipe.Recipe{
			ID:    id,
			Name:  name,
			Style: style,
		}
		r.SetStatus(status)
		result = append(result, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, err
}

// Delete deletes a recipe based on an identifier
func (s *PersistentStore) Delete(id string) error {
	return nil
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
