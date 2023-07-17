package recipe

// RecipeParser represents a component that parses recipes
type RecipeParser interface {
	// Parse parses a recipe from a string
	Parse(recipe string) (*Recipe, error)
}
