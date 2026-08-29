package brewday

import (
	"brewday/internal/recipe"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	require := require.New(t)
	basePath := "../../../test/recipe/brewday/"
	testCases := []struct {
		Name     string
		FileName string
		Error    bool
		Expected *recipe.Recipe
	}{
		{
			Name:     "All",
			FileName: "All.json",
			Error:    false,
			Expected: &recipe.Recipe{
				Name:       "Recipe with all",
				Style:      "Brewday IPA",
				BatchSize:  10.5,
				InitialSG:  1.048,
				Bitterness: 28,
				ColorEBC:   10,
				Mashing: recipe.MashInstructions{
					Malts: []recipe.Malt{
						{Name: "Pale Ale", Amount: 1540},
						{Name: "CARAHELL", Amount: 500.5},
					},
					MainWaterVolume:    7.3,
					Nachguss:           3.7,
					MashTemperature:    55,
					MashOutTemperature: 78.1,
					Rasts: []recipe.Rast{
						{Temperature: 63, Duration: 45},
						{Temperature: 72, Duration: 23.5},
					},
				},
				Hopping: recipe.HopInstructions{
					TotalCookingTime: 90,
					Hops: []recipe.Hops{
						{Name: "Vic Secret", Alpha: 8.9, Amount: 10, Duration: 0, DryHop: false, Vorderwuerze: true},
						{Name: "Saaz", Alpha: 4, Amount: 20, Duration: 10, DryHop: false, Vorderwuerze: false},
						{Name: "Vic Secret", Alpha: 8.9, Amount: 20, Duration: 0, DryHop: true, Vorderwuerze: false},
					},
					AdditionalIngredients: []recipe.AdditionalIngredient{
						{Name: "Orange Peel", Amount: 25.4, Duration: 10},
					},
				},
				Fermentation: recipe.FermentationInstructions{
					Yeast:       recipe.Yeast{Name: "Safale US-05", Amount: 11},
					Temperature: "21",
					AdditionalIngredients: []recipe.AdditionalIngredient{
						{Name: "Roses", Amount: 30, Duration: 10},
					},
					Carbonation: 5.5,
				},
				BrewingSystem: "system1",
			},
		},
		{
			Name:     "No batch size",
			FileName: "no_batch.json",
			Error:    true,
			Expected: nil,
		},
		{
			Name:     "No initial gravity",
			FileName: "no_ig.json",
			Error:    true,
			Expected: nil,
		},
		{
			Name:     "No malts",
			FileName: "no_malts.json",
			Error:    true,
			Expected: nil,
		},
		{
			Name:     "No rasts",
			FileName: "no_rasts.json",
			Error:    true,
			Expected: nil,
		},
		{
			Name:     "No Carbonation",
			FileName: "no_carb.json",
			Error:    true,
			Expected: nil,
		},
		{
			Name:     "No bs",
			FileName: "no_bs.json",
			Error:    false,
			Expected: &recipe.Recipe{
				Name:       "Recipe with all",
				Style:      "Brewday IPA",
				BatchSize:  10.5,
				InitialSG:  1.048,
				Bitterness: 28,
				ColorEBC:   10,
				Mashing: recipe.MashInstructions{
					Malts: []recipe.Malt{
						{Name: "Pale Ale", Amount: 1540},
						{Name: "CARAHELL", Amount: 500.5},
					},
					MainWaterVolume:    7.3,
					Nachguss:           3.7,
					MashTemperature:    55,
					MashOutTemperature: 78.1,
					Rasts: []recipe.Rast{
						{Temperature: 63, Duration: 45},
						{Temperature: 72, Duration: 23.5},
					},
				},
				Hopping: recipe.HopInstructions{
					TotalCookingTime: 90,
					Hops: []recipe.Hops{
						{Name: "Vic Secret", Alpha: 8.9, Amount: 10, Duration: 0, DryHop: false, Vorderwuerze: true},
						{Name: "Saaz", Alpha: 4, Amount: 20, Duration: 10, DryHop: false, Vorderwuerze: false},
						{Name: "Vic Secret", Alpha: 8.9, Amount: 20, Duration: 0, DryHop: true, Vorderwuerze: false},
					},
					AdditionalIngredients: []recipe.AdditionalIngredient{
						{Name: "Orange Peel", Amount: 25.4, Duration: 10},
					},
				},
				Fermentation: recipe.FermentationInstructions{
					Yeast:       recipe.Yeast{Name: "Safale US-05", Amount: 11},
					Temperature: "21",
					AdditionalIngredients: []recipe.AdditionalIngredient{
						{Name: "Roses", Amount: 30, Duration: 10},
					},
					Carbonation: 5.5,
				},
				BrewingSystem: "",
			},
		},
	}
	parser := BrewdayParser{}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			bytes, err := os.ReadFile(basePath + tc.FileName)
			require.NoError(err)
			actual, err := parser.Parse(string(bytes))
			if tc.Error {
				require.Error(err)
			} else {
				require.NoError(err)
				require.Equal(tc.Expected, actual)
			}
		})
	}
}

func TestNotZero(t *testing.T) {
	require := require.New(t)
	testCases := []struct {
		Name  string
		Value any
		Error bool
	}{
		{Name: "Empty string", Value: "", Error: true},
		{Name: "Non Empty string", Value: "oe", Error: false},
		{Name: "Zero", Value: 0, Error: true},
		{Name: "Not zero int", Value: 2, Error: false},
		{Name: "Zero with decimal", Value: 0.0, Error: true},
		{Name: "Non zero decimal", Value: 3.2, Error: false},
		{Name: "Almost zero", Value: 0.0000001, Error: false},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			err := notZero(tc.Value)
			if tc.Error {
				require.Error(err)
			} else {
				require.NoError(err)
			}
		})
	}
}
