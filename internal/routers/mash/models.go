package mash

import "brewday/internal/recipe"

// RecipeStore represents a component that stores recipes
type RecipeStore interface {
	// Retrieve retrieves a recipe based on an identifier
	Retrieve(id string) (*recipe.Recipe, error)
}
