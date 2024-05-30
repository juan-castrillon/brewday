package memory

import (
	"brewday/internal/recipe"
	"encoding/hex"
	"errors"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Date struct {
	date *time.Time
	name string
}

// MemoryStore is a store that stores data in memory
type MemoryStore struct {
	lock      sync.Mutex
	recipes   map[string]*recipe.Recipe
	datesLock sync.Mutex
	dates     map[string][]*Date
}

// NewMemoryStore creates a new MemoryStore
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		recipes: make(map[string]*recipe.Recipe),
	}
}

// CreateID creates a new identifier based on a recipe name
// It encodes the recipe name using hexadecimal encoding
func (s *MemoryStore) CreateID(recipeName string) string {
	dst := make([]byte, hex.EncodedLen(len(recipeName)))
	hex.Encode(dst, []byte(recipeName))
	return string(dst)
}

// Store stores a recipe and returns an identifier that can be used to retrieve it
func (s *MemoryStore) Store(recipe *recipe.Recipe) (string, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	id := s.CreateID(recipe.Name)
	s.recipes[id] = recipe
	recipe.ID = id
	recipe.InitResults()
	return id, nil
}

// Retrieve retrieves a recipe based on an identifier
func (s *MemoryStore) Retrieve(id string) (*recipe.Recipe, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	re, ok := s.recipes[id]
	if !ok {
		return nil, errors.New("recipe not found")
	}
	return re, nil
}

// List lists all the recipes
func (s *MemoryStore) List() ([]*recipe.Recipe, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	var recipes []*recipe.Recipe
	for _, re := range s.recipes {
		recipes = append(recipes, re)
	}
	return recipes, nil
}

// Delete deletes a recipe
func (s *MemoryStore) Delete(id string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	_, ok := s.recipes[id]
	if !ok {
		return errors.New("recipe not found")
	}
	delete(s.recipes, id)
	return nil
}

// UpdateStatus updates the status of a recipe in the store
func (s *MemoryStore) UpdateStatus(id string, status recipe.RecipeStatus, statusParams ...string) error {
	r, err := s.Retrieve(id)
	if err != nil {
		return err
	}
	r.SetStatus(status, statusParams...)
	return nil
}

// UpdateResult updates a certain result of a recipe
func (s *MemoryStore) UpdateResult(id string, resultType recipe.ResultType, value float32) error {
	r, err := s.Retrieve(id)
	if err != nil {
		return err
	}
	switch resultType {
	case recipe.ResultHotWortVolume:
		r.SetHotWortVolume(value)
	case recipe.ResultOriginalGravity:
		r.SetOriginalGravity(value)
	case recipe.ResultFinalGravity:
		r.SetFinalGravity(value)
	case recipe.ResultAlcohol:
		r.SetAlcohol(value)
	case recipe.ResultMainFermentationVolume:
		r.SetMainFermentationVolume(value)
	case recipe.ResultVolumeBeforeBoil:
		r.SetVolumeBeforeBoil(value)
	default:
		return errors.New("invalid result not present in struct: " + strconv.Itoa(int(resultType)))
	}
	return nil
}

// RetrieveResult gets a certain result value from a recipe
func (s *MemoryStore) RetrieveResult(id string, resultType recipe.ResultType) (float32, error) {
	res, err := s.RetrieveResults(id)
	if err != nil {
		return 0, err
	}
	switch resultType {
	case recipe.ResultHotWortVolume:
		return res.HotWortVolume, nil
	case recipe.ResultOriginalGravity:
		return res.OriginalGravity, nil
	case recipe.ResultFinalGravity:
		return res.FinalGravity, nil
	case recipe.ResultAlcohol:
		return res.Alcohol, nil
	case recipe.ResultMainFermentationVolume:
		return res.MainFermentationVolume, nil
	case recipe.ResultVolumeBeforeBoil:
		return res.VolumeBeforeBoil, nil
	default:
		return 0, errors.New("invalid result not present in struct: " + strconv.Itoa(int(resultType)))
	}
}

// RetrieveResults gets the results from a certain recipe
func (s *MemoryStore) RetrieveResults(id string) (*recipe.RecipeResults, error) {
	r, err := s.Retrieve(id)
	if err != nil {
		return nil, err
	}
	res := r.GetResults()
	return &res, nil
}

// AddMainFermSG adds a new specific gravity measurement to a given recipe
func (s *MemoryStore) AddMainFermSG(id string, m *recipe.SGMeasurement) error {
	r, err := s.Retrieve(id)
	if err != nil {
		return err
	}
	r.SetSGMeasurement(m)
	return nil
}

// RetrieveMainFermSGs returns all measured sgs for a recipe
func (s *MemoryStore) RetrieveMainFermSGs(id string) ([]*recipe.SGMeasurement, error) {
	r, err := s.Retrieve(id)
	if err != nil {
		return nil, err
	}
	return r.GetSGMeasurements(), nil
}

// AddDate allows to store a date with a certain purpose. It can be used to store notification dates, or timers
func (s *MemoryStore) AddDate(id string, date *time.Time, name string) error {
	s.datesLock.Lock()
	defer s.datesLock.Unlock()
	if s.dates == nil {
		s.dates = make(map[string][]*Date)
	}
	_, ok := s.dates[id]
	if !ok {
		s.dates[id] = make([]*Date, 0)
	}
	d := &Date{
		date: date,
		name: name,
	}
	s.dates[id] = append(s.dates[id], d)
	return nil
}

// RetrieveDates allows to retreive stored dates with its purpose (name).It can be used to store notification dates, or timers
// It supports pattern in the name to retrieve multiple values
// The pattern is searched in the names with strings.Contains
func (s *MemoryStore) RetrieveDates(id, namePattern string) ([]*time.Time, error) {
	s.datesLock.Lock()
	defer s.datesLock.Unlock()
	results := make([]*time.Time, 0)
	for _, d := range s.dates[id] {
		if strings.Contains(d.name, namePattern) {
			results = append(results, d.date)
		}
	}
	return results, nil
}

// AddSugarResult adds a new priming sugar result to a given recipe
func (s *MemoryStore) AddSugarResult(id string, result *recipe.PrimingSugarResult) error {
	r, err := s.Retrieve(id)
	if err != nil {
		return err
	}
	r.SetPrimingSugarResult(result)
	return nil
}

// RetrieveSugarResults returns all sugar results for a recipe
func (s *MemoryStore) RetrieveSugarResults(id string) ([]*recipe.PrimingSugarResult, error) {
	r, err := s.Retrieve(id)
	if err != nil {
		return nil, err
	}
	return r.GetPrimingSugarResults(), nil
}
