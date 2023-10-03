package recipes

import "brewday/internal/recipe"

// RecipeStore represents a component that stores recipes
type RecipeStore interface {
	// List lists all the recipes
	List() ([]*recipe.Recipe, error)
}
