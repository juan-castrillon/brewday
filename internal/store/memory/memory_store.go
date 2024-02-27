package memory

import (
	"brewday/internal/recipe"
	"encoding/hex"
	"errors"
	"strconv"
	"sync"
)

// MemoryStore is a store that stores data in memory
type MemoryStore struct {
	lock    sync.Mutex
	recipes map[string]*recipe.Recipe
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

// UpdateResults updates a certain result of a recipe
func (s *MemoryStore) UpdateResults(id string, resultType recipe.ResultType, value float32) error {
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
	default:
		return errors.New("invalid result not present in struct: " + strconv.Itoa(int(resultType)))
	}
	return nil
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
