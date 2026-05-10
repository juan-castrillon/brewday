package brewday

import (
	"brewday/internal/recipe"
	"encoding/json"
	"errors"
)

// BrewdayParser is a RecipeParser implementation for brewday's own
// recipe format
type BrewdayParser struct{}

func (p *BrewdayParser) Parse(rec string) (*recipe.Recipe, error) {
	var r recipe.Recipe
	err := json.Unmarshal([]byte(rec), &r)
	if err != nil {
		return nil, err
	}
	err = p.validate(&r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (p *BrewdayParser) validate(r *recipe.Recipe) error {
	if err := notZero(r.BatchSize); err != nil {
		return errors.New("Invalid batch size, must be defined")
	}
	if err := notZero(r.InitialSG); err != nil {
		return errors.New("Invalid Initial gravity, must be defined")
	}
	if len(r.Mashing.Malts) == 0 {
		return errors.New("At least one malt is expected")
	}
	if len(r.Mashing.Rasts) == 0 {
		return errors.New("At least one rast is needed")
	}
	if err := notZero(r.Fermentation.Carbonation); err != nil {
		return errors.New("Carbonation must be defined")
	}
	return nil
}

func notZero(field any) error {
	isZero := false
	assertionError := errors.New("Failed to convert type")
	switch v := field.(type) {
	case string:
		isZero = v == ""
	case int:
		isZero = v == 0
	case float32:
		isZero = v == float32(0)
	case float64:
		isZero = v == float64(0)
	default:
		return assertionError
	}
	if isZero {
		return errors.New("Value is zero")
	}
	return nil
}
