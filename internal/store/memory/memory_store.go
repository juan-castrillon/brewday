package memory

import (
	"brewday/internal/recipe"
	"encoding/hex"
	"errors"
)

// MemoryStore is a store that stores data in memory
type MemoryStore struct {
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
	id := s.CreateID(recipe.Name)
	s.recipes[id] = recipe
	return id, nil
}

// Retrieve retrieves a recipe based on an identifier
func (s *MemoryStore) Retrieve(id string) (*recipe.Recipe, error) {
	re, ok := s.recipes[id]
	if !ok {
		return nil, errors.New("recipe not found")
	}
	return re, nil
}
