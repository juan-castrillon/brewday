package sql

import (
	"brewday/internal/recipe"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type PersistentStore struct {
	dbClient             *sql.DB
	retrieveStatement    *sql.Stmt
	updateStateStatement *sql.Stmt
}

func NewPersistentStore(db *sql.DB) (*PersistentStore, error) {
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
		(hot_wort_vol, original_sg, final_sg, alcohol, main_ferm_vol, vol_bb, recipe_id)
		VALUES ( ?, ?, ?, ?, ?, ?, ?)
	`, initialResults.HotWortVolume, initialResults.OriginalGravity,
		initialResults.FinalGravity, initialResults.Alcohol,
		initialResults.MainFermentationVolume, initialResults.VolumeBeforeBoil, idString)
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

func (s *PersistentStore) columnNameFromType(resultType recipe.ResultType) string {
	switch resultType {
	case recipe.ResultHotWortVolume:
		return "hot_wort_vol"
	case recipe.ResultOriginalGravity:
		return "original_sg"
	case recipe.ResultFinalGravity:
		return "final_sg"
	case recipe.ResultAlcohol:
		return "alcohol"
	case recipe.ResultMainFermentationVolume:
		return "main_ferm_vol"
	case recipe.ResultVolumeBeforeBoil:
		return "vol_bb"
	default:
		return ""
	}
}

// UpdateResult updates a certain result of a recipe
func (s *PersistentStore) UpdateResult(id string, resultType recipe.ResultType, value float32) error {
	columnName := s.columnNameFromType(resultType)
	if columnName == "" {
		return errors.New("invalid result not present in schema: " + strconv.Itoa(int(resultType)))
	}
	_, err := s.dbClient.Exec("UPDATE recipe_results SET "+columnName+" = ? WHERE recipe_id == ?", value, id)
	return err
}

// RetrieveResult gets a certain result value from a recipe
func (s *PersistentStore) RetrieveResult(id string, resultType recipe.ResultType) (float32, error) {
	columnName := s.columnNameFromType(resultType)
	if columnName == "" {
		return 0, errors.New("invalid result not present in schema: " + strconv.Itoa(int(resultType)))
	}
	var val float32
	err := s.dbClient.QueryRow(fmt.Sprintf(`SELECT %s FROM recipe_results WHERE recipe_id == ?`, columnName), id).Scan(&val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

// RetrieveResults gets all the results from a certain recipe
func (s *PersistentStore) RetrieveResults(id string) (*recipe.RecipeResults, error) {
	var actual recipe.RecipeResults
	err := s.dbClient.QueryRow(`
		SELECT hot_wort_vol, original_sg, final_sg, alcohol, main_ferm_vol, vol_bb
		FROM recipe_results WHERE recipe_id == ?`, id).Scan(
		&actual.HotWortVolume, &actual.OriginalGravity,
		&actual.FinalGravity, &actual.Alcohol, &actual.MainFermentationVolume,
		&actual.VolumeBeforeBoil,
	)
	if err != nil {
		return nil, err
	}
	return &actual, nil
}

// AddMainFermSG adds a new specific gravity measurement to a given recipe
func (s *PersistentStore) AddMainFermSG(id string, m *recipe.SGMeasurement) error {
	_, err := s.dbClient.Exec(`INSERT INTO main_ferm_sgs (sg, date, recipe_id) VALUES (?, ?, ?)`, m.Value, m.Date, id)
	return err
}

// RetrieveMainFermSGs returns all measured sgs for a recipe
func (s *PersistentStore) RetrieveMainFermSGs(id string) ([]*recipe.SGMeasurement, error) {
	rows, err := s.dbClient.Query(`SELECT sg, date FROM main_ferm_sgs WHERE recipe_id == ? ORDER BY id ASC`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	results := make([]*recipe.SGMeasurement, 0)
	for rows.Next() {
		var m recipe.SGMeasurement
		err = rows.Scan(&m.Value, &m.Date)
		if err != nil {
			return nil, err
		}
		results = append(results, &m)
	}
	return results, nil
}

// AddDate allows to store a date with a certain purpose. It can be used to store notification dates, or timers
func (s *PersistentStore) AddDate(id string, date *time.Time, name string) error {
	dateString := date.Format(time.RFC3339)
	_, err := s.dbClient.Exec(`INSERT INTO dates (date, name, recipe_id) VALUES (?, ?, ?)`, dateString, name, id)
	return err
}

// RetrieveDates allows to retreive stored dates with its purpose (name).It can be used to store notification dates, or timers
// It supports pattern in the name to retrieve multiple values
func (s *PersistentStore) RetrieveDates(id, namePattern string) ([]*time.Time, error) {
	sanitizedPattern := strings.ReplaceAll(strings.ReplaceAll(namePattern, "_", "!_"), "%", "!%") + "%"
	rows, err := s.dbClient.Query(`SELECT date FROM dates WHERE recipe_id == ? AND name LIKE ? ESCAPE '!'`, id, sanitizedPattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	results := make([]*time.Time, 0)
	for rows.Next() {
		var date string
		err = rows.Scan(&date)
		if err != nil {
			return nil, err
		}
		t, err := time.Parse(time.RFC3339, date)
		if err != nil {
			return nil, err
		}
		results = append(results, &t)
	}
	return results, nil
}

// AddSugarResult adds a new priming sugar result to a given recipe
func (s *PersistentStore) AddSugarResult(id string, r *recipe.PrimingSugarResult) error {
	_, err := s.dbClient.Exec(`INSERT INTO sugar_results (water, sugar, alcohol, recipe_id) VALUES (?, ?, ?, ?)`, r.Water, r.Amount, r.Alcohol, id)
	return err
}

// RetrieveSugarResults returns all sugar results for a recipe
func (s *PersistentStore) RetrieveSugarResults(id string) ([]*recipe.PrimingSugarResult, error) {
	rows, err := s.dbClient.Query(`SELECT water, sugar, alcohol FROM sugar_results WHERE recipe_id == ? ORDER BY id ASC`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	results := make([]*recipe.PrimingSugarResult, 0)
	for rows.Next() {
		var m recipe.PrimingSugarResult
		err = rows.Scan(&m.Water, &m.Amount, &m.Alcohol)
		if err != nil {
			return nil, err
		}
		results = append(results, &m)
	}
	return results, nil
}

// AddBoolFlag allows to store a given flag that can be true or false in the store with a unique name
func (s *PersistentStore) AddBoolFlag(id, name string, flag bool) error {
	var flagID int
	err := s.dbClient.QueryRow(`SELECT id FROM bool_flags WHERE recipe_id == ? AND name == ?`, id, name).Scan(&flagID)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		_, err := s.dbClient.Exec(`INSERT INTO bool_flags (name, value, recipe_id) VALUES (?, ?, ?)`, name, flag, id)
		return err
	}
	_, err = s.dbClient.Exec(`UPDATE bool_flags SET value = ? WHERE recipe_id == ? AND name == ?`, flag, id, name)
	return err
}

// RetrieveBoolFlag gets a bool flag from the store given its name
func (s *PersistentStore) RetrieveBoolFlag(id, name string) (bool, error) {
	var value bool
	err := s.dbClient.QueryRow(`SELECT value FROM bool_flags WHERE recipe_id == ? AND name == ?`, id, name).Scan(&value)
	if err != nil {
		if err != sql.ErrNoRows {
			return false, err
		}
		return false, nil
	}
	return value, nil
}
