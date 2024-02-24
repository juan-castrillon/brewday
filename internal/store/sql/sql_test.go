package sql

import (
	"brewday/internal/recipe"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStore(t *testing.T) {
	require := require.New(t)
	fileName := strings.ToLower(strings.TrimSpace(t.Name()))
	store, err := NewPersistentStore(fileName)
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
			var name, style, status_args, mash_malts, mash_rasts, hop_hops, hop_add, ferm_yeast, ferm_temp, ferm_add string
			var batch_size_l, initial_sg, ibu, ebc, mash_main_water, mash_nachguss, mash_temp, mash_out_temp, hop_cooking, ferm_carbonation float32
			var status recipe.RecipeStatus
			err = store.dbClient.QueryRow(`SELECT 
				name, style, batch_size_l, initial_sg, ibu, ebc, status, status_args,
				mash_malts, mash_main_water, mash_nachguss, mash_temp, mash_out_temp, mash_rasts,
				hop_cooking_time, hop_hops, hop_additional,
				ferm_yeast, ferm_temp, ferm_additional, ferm_carbonation
			FROM recipes WHERE id == ?`, id).Scan(&name, &style, &batch_size_l, &initial_sg, &ibu, &ebc, &status, &status_args,
				&mash_malts, &mash_main_water, &mash_nachguss, &mash_temp, &mash_out_temp, &mash_rasts,
				&hop_cooking, &hop_hops, &hop_add,
				&ferm_yeast, &ferm_temp, &ferm_add, &ferm_carbonation)
			require.NoError(err)
			require.Equal(r.Name, name)
			require.Equal(r.Style, style)
			require.Equal(r.BatchSize, batch_size_l)
			require.Equal(r.InitialSG, initial_sg)
			require.Equal(r.Bitterness, ibu)
			require.Equal(r.ColorEBC, ebc)
			require.Equal(r.Mashing.MainWaterVolume, mash_main_water)
			require.Equal(r.Mashing.Nachguss, mash_nachguss)
			require.Equal(r.Mashing.MashTemperature, mash_temp)
			require.Equal(r.Mashing.MashOutTemperature, mash_out_temp)
			require.Equal(r.Hopping.TotalCookingTime, hop_cooking)
			require.Equal(r.Fermentation.Temperature, ferm_temp)
			require.Equal(r.Fermentation.Carbonation, ferm_carbonation)
			actualStatus, _ := r.GetStatus()
			require.Equal(actualStatus, status)
			actualArrays := store.marshallArrays(r)
			require.Equal(actualArrays.StatusParams, status_args)
			require.Equal(actualArrays.MashingRasts, mash_rasts)
			require.Equal(actualArrays.MashingMalts, mash_malts)
			require.Equal(actualArrays.HopHops, hop_hops)
			require.Equal(actualArrays.HopAdd, hop_add)
			require.Equal(actualArrays.Yeast, ferm_yeast)
			require.Equal(actualArrays.FermAdd, ferm_add)
		})
	}
}
