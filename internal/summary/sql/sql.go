package sql

import (
	"brewday/internal/summary"
	"database/sql"
	"encoding/json"
	"errors"

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

// AddSummary adds a summary for the given recipe id with the given title
func (s *SummaryRecorderPersistentStore) AddSummary(recipeID, title string) error {
	_, err := s.dbClient.Exec(`INSERT INTO summaries (title, recipe_id) VALUES (?, ?)`, title, recipeID)
	return err
}

// DeleteSummary deletes the summary for the given recipe id
func (s *SummaryRecorderPersistentStore) DeleteSummary(recipeID string) error {
	if recipeID == "" {
		return errors.New("invalid empty recipe id for deleting summary")
	}
	_, err := s.dbClient.Exec(`DELETE FROM summaries WHERE recipe_id == ?`, recipeID)
	return err
}

// AddMashTemp adds a mash temperature to the summary and notes related to it
func (s *SummaryRecorderPersistentStore) AddMashTemp(id string, temp float64, notes string) error {
	if id == "" {
		return errors.New("invalid empty recipe id")
	}
	_, err := s.dbClient.Exec(`UPDATE summaries SET mash_temp = ? , mash_notes = ? WHERE recipe_id == ?`, temp, notes, id)
	return err
}

func (s *SummaryRecorderPersistentStore) AddRast(id string, temp float64, duration float64, notes string) error {
	if id == "" {
		return errors.New("invalid empty recipe id")
	}
	newRast := summary.MashRastInfo{
		Temperature: temp,
		Time:        duration,
		Notes:       notes,
	}
	newRastBytes, err := json.Marshal(newRast)
	if err != nil {
		return err
	}
	return s.addToMarshalledArray(id, "mash_rasts", string(newRastBytes))
}

func (s *SummaryRecorderPersistentStore) AddLauternNotes(id, notes string) error {
	if id == "" {
		return errors.New("invalid empty recipe id")
	}
	_, err := s.dbClient.Exec(`UPDATE summaries SET lautern_info = ? WHERE recipe_id == ?`, notes, id)
	return err
}

func (s *SummaryRecorderPersistentStore) AddHopping(id string, name string, amount float32, alpha float32, duration float32, notes string) error {
	if id == "" {
		return errors.New("invalid empty recipe id")
	}
	newHop := summary.HopInfo{
		Name:     name,
		Grams:    amount,
		Alpha:    alpha,
		Time:     duration,
		TimeUnit: "minutes",
		Notes:    notes,
	}
	newHopBytes, err := json.Marshal(newHop)
	if err != nil {
		return err
	}
	return s.addToMarshalledArray(id, "hopping_hops", string(newHopBytes))
}

func (s *SummaryRecorderPersistentStore) AddVolumeBeforeBoil(id string, amount float32, notes string) error {
	if id == "" {
		return errors.New("invalid empty recipe id")
	}
	_, err := s.dbClient.Exec(`UPDATE summaries SET hopping_vol_bb = ? , hopping_vol_bb_notes = ? WHERE recipe_id == ?`, amount, notes, id)
	return err
}

func (s *SummaryRecorderPersistentStore) AddVolumeAfterBoil(id string, amount float32, notes string) error {
	if id == "" {
		return errors.New("invalid empty recipe id")
	}
	_, err := s.dbClient.Exec(`UPDATE summaries SET hopping_vol_ab = ? , hopping_vol_ab_notes = ? WHERE recipe_id == ?`, amount, notes, id)
	return err
}

func (s *SummaryRecorderPersistentStore) AddCooling(id string, finalTemp, coolingTime float32, notes string) error {
	if id == "" {
		return errors.New("invalid empty recipe id")
	}
	panic("Implement me!")
}

func (s *SummaryRecorderPersistentStore) AddPreFermentationVolume(id string, volume float32, sg float32, notes string) error {
	if id == "" {
		return errors.New("invalid empty recipe id")
	}
	panic("Implement me!")
}

func (s *SummaryRecorderPersistentStore) AddYeastStart(id string, temperature, notes string) error {
	if id == "" {
		return errors.New("invalid empty recipe id")
	}
	panic("Implement me!")
}

func (s *SummaryRecorderPersistentStore) AddMainFermentationSGMeasurement(id string, date string, gravity float32, final bool, notes string) error {
	if id == "" {
		return errors.New("invalid empty recipe id")
	}
	panic("Implement me!")
}

func (s *SummaryRecorderPersistentStore) AddMainFermentationAlcohol(id string, alcohol float32) error {
	if id == "" {
		return errors.New("invalid empty recipe id")
	}
	panic("Implement me!")
}

func (s *SummaryRecorderPersistentStore) AddMainFermentationDryHop(id string, name string, amount, alpha, duration float32, notes string) error {
	if id == "" {
		return errors.New("invalid empty recipe id")
	}
	panic("Implement me!")
}

func (s *SummaryRecorderPersistentStore) AddPreBottlingVolume(id string, volume float32) error {
	if id == "" {
		return errors.New("invalid empty recipe id")
	}
	panic("Implement me!")
}

func (s *SummaryRecorderPersistentStore) AddBottling(id string, carbonation, alcohol, sugar, temp, vol float32, sugarType, notes string) error {
	if id == "" {
		return errors.New("invalid empty recipe id")
	}
	panic("Implement me!")
}

func (s *SummaryRecorderPersistentStore) AddSummarySecondary(id string, days int, notes string) error {
	if id == "" {
		return errors.New("invalid empty recipe id")
	}
	panic("Implement me!")
}

func (s *SummaryRecorderPersistentStore) AddEvaporation(id string, amount float32) error {
	if id == "" {
		return errors.New("invalid empty recipe id")
	}
	panic("Implement me!")
}

func (s *SummaryRecorderPersistentStore) AddEfficiency(id string, efficiencyPercentage float32) error {
	if id == "" {
		return errors.New("invalid empty recipe id")
	}
	panic("Implement me!")
}

func (s *SummaryRecorderPersistentStore) AddTimeline(id string, timeline []string) error {
	if id == "" {
		return errors.New("invalid empty recipe id")
	}
	panic("Implement me!")
}

func (s *SummaryRecorderPersistentStore) GetSummary(id string) (*summary.Summary, error) {
	if id == "" {
		return nil, errors.New("invalid empty recipe id")
	}
	panic("Implement me!")
}