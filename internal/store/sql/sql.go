package sql

import (
	"brewday/internal/recipe"
	"database/sql"
	"errors"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type PersistentStore struct {
	dbClient             *sql.DB
	retrieveStatement    *sql.Stmt
	updateStateStatement *sql.Stmt
}

func NewPersistentStore(db *sql.DB) (*PersistentStore, error) {
	err := createTable(db)
	if err != nil {
		return nil, err
	}
	err = createResultsTable(db)
	if err != nil {
		return nil, err
	}
	rs, err := db.Prepare(`SELECT 
		name, style, batch_size_l, initial_sg, ibu, ebc, status, status_args,
		mash_malts, mash_main_water, mash_nachguss, mash_temp, mash_out_temp, mash_rasts,
		hop_cooking_time, hop_hops, hop_additional,
		ferm_yeast, ferm_temp, ferm_additional, ferm_carbonation
	FROM recipes WHERE id == ?`)
	if err != nil {
		return nil, err
	}
	uss, err := db.Prepare("UPDATE recipes SET status = ? , status_args = ? WHERE id == ?")
	if err != nil {
		return nil, err
	}
	return &PersistentStore{
		dbClient:             db,
		retrieveStatement:    rs,
		updateStateStatement: uss,
	}, nil
}

// Close closes the underlying connections to the database. It must always be called to avoid leaks
func (s *PersistentStore) Close() error {
	err := s.retrieveStatement.Close()
	if err != nil {
		return err
	}
	err = s.updateStateStatement.Close()
	if err != nil {
		return err
	}
	return nil
}

// Store stores a recipe and returns an identifier that can be used to retrieve it
func (s *PersistentStore) Store(r *recipe.Recipe) (string, error) {
	status, _ := r.GetStatus()
	marshalled, err := s.marshalStructs(r)
	if err != nil {
		return "", err
	}
	res, err := s.dbClient.Exec(`INSERT INTO recipes 
	(
		name, style, batch_size_l, initial_sg, ibu, ebc, status, status_args,
		mash_malts, mash_main_water, mash_nachguss, mash_temp, mash_out_temp, mash_rasts,
		hop_cooking_time, hop_hops, hop_additional,
		ferm_yeast, ferm_temp, ferm_additional, ferm_carbonation
	) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, r.Name, r.Style, r.BatchSize, r.InitialSG, r.Bitterness, r.ColorEBC, status, marshalled.StatusParams,
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
	idString := strconv.FormatInt(id, 10)
	r.InitResults()
	initialResults := r.GetResults()
	_, err = s.dbClient.Exec(`INSERT INTO recipe_results 
		(hot_wort_vol, original_sg, final_sg, alcohol, main_ferm_vol, recipe_id)
		VALUES ( ?, ?, ?, ?, ?, ?)
	`, initialResults.HotWortVolume, initialResults.OriginalGravity,
		initialResults.FinalGravity, initialResults.Alcohol,
		initialResults.MainFermentationVolume, idString)
	if err != nil {
		return "", err
	}
	return idString, nil
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
	// Just getting what I need for now to display, if future calls to list require more, they are to be added here
	rows, err := s.dbClient.Query("SELECT id, name, style, status FROM recipes")
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
	_, err := s.dbClient.Exec("DELETE FROM recipes WHERE id == ?", id)
	return err
}

// UpdateStatus updates the status of a recipe in the store
func (s *PersistentStore) UpdateStatus(id string, status recipe.RecipeStatus, statusParams ...string) error {
	statusArgs, err := s.marshalStatusParams(statusParams...)
	if err != nil {
		return err
	}
	_, err = s.updateStateStatement.Exec(status, statusArgs, id)
	return err
}

// UpdateResults updates a certain result of a recipe
func (s *PersistentStore) UpdateResults(id string, resultType recipe.ResultType, value float32) error {
	var columnName string
	switch resultType {
	case recipe.ResultHotWortVolume:
		columnName = "hot_wort_vol"
	case recipe.ResultOriginalGravity:
		columnName = "original_sg"
	case recipe.ResultFinalGravity:
		columnName = "final_sg"
	case recipe.ResultAlcohol:
		columnName = "alcohol"
	case recipe.ResultMainFermentationVolume:
		columnName = "main_ferm_vol"
	default:
		return errors.New("invalid result not present in schema: " + strconv.Itoa(int(resultType)))
	}
	_, err := s.dbClient.Exec("UPDATE recipe_results SET "+columnName+" = ? WHERE recipe_id == ?", value, id)
	return err
}

// RetrieveResults gets the results from a certain recipe
func (s *PersistentStore) RetrieveResults(id string) (*recipe.RecipeResults, error) {
	var actual recipe.RecipeResults
	err := s.dbClient.QueryRow(`
		SELECT hot_wort_vol, original_sg, final_sg, alcohol, main_ferm_vol
		FROM recipe_results WHERE recipe_id == ?`, id).Scan(
		&actual.HotWortVolume, &actual.OriginalGravity,
		&actual.FinalGravity, &actual.Alcohol, &actual.MainFermentationVolume,
	)
	if err != nil {
		return nil, err
	}
	return &actual, nil
}
