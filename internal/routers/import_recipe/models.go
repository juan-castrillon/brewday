package import_recipe

import "brewday/internal/recipe"

// RecipeParser represents a component that parses recipes
type RecipeParser interface {
	// Parse parses a recipe from a string
	Parse(recipe string) (*recipe.Recipe, error)
}

// RecipeStore represents a component that stores recipes
type RecipeStore interface {
	// Store stores a recipe and returns an identifier that can be used to retrieve it
	Store(recipe *recipe.Recipe) (string, error)
}
