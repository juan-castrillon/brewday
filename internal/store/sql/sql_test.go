package sql

import (
	"brewday/internal/recipe"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStoreAndRetrieve(t *testing.T) {
	require := require.New(t)
	fileName := strings.ToLower(strings.TrimSpace(t.Name()))
	store, err := NewPersistentStore(fileName + ".sqlite")
	require.NoError(err)
	defer os.Remove(fileName)
	testRecipe := &recipe.Recipe{
		Name:       "Huell Saison",
		Style:      "Saison",
		BatchSize:  10,
		InitialSG:  1.0534434,
		Bitterness: 31.3,
		ColorEBC:   13,
		Mashing: recipe.MashInstructions{
			Malts: []recipe.Malt{
				{Name: "Pilsner Malz", Amount: 1900},
				{Name: "Weizenmalz hell", Amount: 350},
				{Name: "Haferflocken", Amount: 100},
				{Name: "Karamellmalz Amber", Amount: 150},
			},
			MainWaterVolume:    7.5,
			Nachguss:           6.8,
			MashTemperature:    57,
			MashOutTemperature: 78,
			Rasts: []recipe.Rast{
				{Temperature: 57, Duration: 10},
				{Temperature: 64, Duration: 45},
				{Temperature: 72, Duration: 25},
				{Temperature: 78, Duration: 5},
			},
		},
		Hopping: recipe.HopInstructions{
			TotalCookingTime: 90,
			Hops: []recipe.Hops{
				{Name: "Huell Melon (VW)", Alpha: 7.2, Amount: 10, Duration: 90, DryHop: false, Vorderwuerze: true},
				{Name: "Huell Melon", Alpha: 7.2, Amount: 3, Duration: 60, DryHop: false, Vorderwuerze: false},
				{Name: "Huell Melon", Alpha: 7.2, Amount: 2, Duration: 10, DryHop: false, Vorderwuerze: false},
				{Name: "Huell Melon", Alpha: 7.2, Amount: 10, Duration: 0, DryHop: false, Vorderwuerze: false},
				{Name: "Huell Melon", Amount: 35, Duration: 0, DryHop: true, Vorderwuerze: false},
			},
			AdditionalIngredients: nil,
		},
		Fermentation: recipe.FermentationInstructions{
			Yeast: recipe.Yeast{
				Name: "LalBrew Belle Saison",
			},
			Temperature:           "21",
			AdditionalIngredients: nil,
			Carbonation:           5,
		},
	}
	testCases := []struct {
		Name         string
		Modify       bool
		Status       recipe.RecipeStatus
		StatusParams []string
	}{
		{Name: "Original", Modify: false},
		{Name: "StatusNoParams", Modify: true, Status: recipe.RecipeStatusLautering, StatusParams: []string{}},
		{Name: "StatusParams", Modify: true, Status: recipe.RecipeStatusPreFermentation, StatusParams: []string{
			"water", "15.324", "0.032",
		}},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			r := testRecipe
			if tc.Modify {
				r.SetStatus(tc.Status, tc.StatusParams...)
			}
			id, err := store.Store(testRecipe)
			require.NoError(err)
			actual, err := store.Retrieve(id)
			require.NoError(err)
			r.ID = id
			require.Equal(r, actual)
			actualStatus, actualParams := actual.GetStatus()
			expectedStatus, expectedParams := r.GetStatus()
			require.Equal(expectedStatus, actualStatus)
			require.ElementsMatch(expectedParams, actualParams)
		})
	}
}
