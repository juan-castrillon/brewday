package braureka_json

import (
	"brewday/internal/recipe"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBJsonRecipeToRecipe(t *testing.T) {
	require := require.New(t)
	basePath := "../../../test/recipe/braureka_json/"
	type testCase struct {
		Name     string
		FileName string
		Expected recipe.Recipe
	}
	testCases := []testCase{
		{
			Name:     "Ginger_Wit",
			FileName: "Ginger_Wit.json",
			Expected: recipe.Recipe{
				Name:       "Ginger Wit",
				Style:      "Witbier",
				BatchSize:  23,
				InitialSG:  1.0471187,
				Bitterness: 17.7,
				ColorEBC:   7,
				Mashing: recipe.MashInstructions{
					Malts: []recipe.Malt{
						{Name: "Pilsner Malz", Amount: 2434},
						{Name: "Weizenmalz hell", Amount: 1850},
						{Name: "Haferflocken", Amount: 235},
						{Name: "Pilsner Malz extra hell", Amount: 135},
					},
					MainWaterVolume:    15,
					Nachguss:           15.2,
					MashTemperature:    42,
					MashOutTemperature: 72,
					Rasts: []recipe.Rast{
						{Temperature: 42, Duration: 1},
						{Temperature: 50, Duration: 20},
						{Temperature: 55, Duration: 20},
						{Temperature: 62, Duration: 30},
						{Temperature: 72, Duration: 30},
					},
				},
				Hopping: recipe.HopInstructions{
					TotalCookingTime: 90,
					Hops: []recipe.Hops{
						{Name: "Cascade US (VW)", Alpha: 7.1, Amount: 25, Duration: 90, DryHop: false, Vorderwuerze: true},
					},
					AdditionalIngredients: []recipe.AdditionalIngredient{
						{Name: "Orangenschale", Amount: 15, Duration: 10},
						{Name: "Koriander", Amount: 10, Duration: 10},
						{Name: "Ingwer", Amount: 50, Duration: 10},
					},
				},
				Fermentation: recipe.FermentationInstructions{
					Yeast: recipe.Yeast{
						Name:   "Gozdawa Classic Belgian Witbier CBW",
						Amount: 0,
					},
					Temperature:           "22",
					AdditionalIngredients: nil,
					Carbonation:           5,
				},
			},
		},
		{
			Name:     "Hermano_Juan",
			FileName: "Hermano_Juan.json",
			Expected: recipe.Recipe{
				Name:       "Hermano Juan",
				Style:      "Belgisches Dubbel",
				BatchSize:  11,
				InitialSG:  1.0667312,
				Bitterness: 23.7,
				ColorEBC:   41,
				Mashing: recipe.MashInstructions{
					Malts: []recipe.Malt{
						{Name: "Wiener Malz", Amount: 3500},
						{Name: "Weizenmalz dunkel", Amount: 330},
						{Name: "Karamellmalz Aroma", Amount: 145},
					},
					MainWaterVolume:    11.9,
					Nachguss:           4.5,
					MashTemperature:    57,
					MashOutTemperature: 78,
					Rasts: []recipe.Rast{
						{Temperature: 57, Duration: 5},
						{Temperature: 63, Duration: 40},
						{Temperature: 72, Duration: 30},
						{Temperature: 78, Duration: 1},
					},
				},
				Hopping: recipe.HopInstructions{
					TotalCookingTime: 90,
					Hops: []recipe.Hops{
						{Name: "Perle", Alpha: 8.1, Amount: 14, Duration: 70, DryHop: false, Vorderwuerze: false},
					},
					AdditionalIngredients: []recipe.AdditionalIngredient{
						{Name: "Kandiszucker", Amount: 150, Duration: 15},
					},
				},
				Fermentation: recipe.FermentationInstructions{
					Yeast: recipe.Yeast{
						Name: "Fermentis SafAle BE-256",
					},
					Temperature:           "21",
					AdditionalIngredients: nil,
					Carbonation:           5.5,
				},
			},
		},
		{
			Name:     "Huell_Saison",
			FileName: "Huell_Saison.json",
			Expected: recipe.Recipe{
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
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			bytes, err := os.ReadFile(basePath + tc.FileName)
			require.NoError(err)
			var r BraurekaJSONRecipe
			err = json.Unmarshal(bytes, &r)
			require.NoError(err)
			actual, err := bJsonRecipeToRecipe(&r)
			require.NoError(err)
			tc.Expected.SetStatus(recipe.RecipeStatusCreated)
			require.Equal(tc.Expected, *actual)
		})
	}
}
