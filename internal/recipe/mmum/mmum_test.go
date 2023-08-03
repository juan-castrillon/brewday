package mmum

import (
	"brewday/internal/recipe"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetMashInstructions(t *testing.T) {
	require := require.New(t)
	basePath := "../../../test/recipe/mmum/"
	type testCase struct {
		Name     string
		FileName string
		Expected recipe.MashInstructions
	}
	testCases := []testCase{
		{
			Name:     "Hula Hula IPA",
			FileName: "Hula_Hula_IPA.json",
			Expected: recipe.MashInstructions{
				Malts: []recipe.Malt{
					{Name: "Golden Promise PA", Amount: 5600},
					{Name: "Barke Pilsner", Amount: 5000},
					{Name: "Haferflocken", Amount: 500},
					{Name: "Gerstenflocken", Amount: 500},
					{Name: "Carapils", Amount: 500},
					{Name: "Sauermalz", Amount: 300},
					{Name: "Cara Red", Amount: 300},
				},
				MainWaterVolume:    41,
				MashTemperature:    69,
				Nachguss:           19,
				MashOutTemperature: 77,
				Rasts: []recipe.Rast{
					{Temperature: 67.5, Duration: 45},
					{Temperature: 72, Duration: 15},
				},
			},
		},
		{
			Name:     "Callippo Mango",
			FileName: "Callippo_Mango.json",
			Expected: recipe.MashInstructions{
				Malts: []recipe.Malt{
					{Name: "Pilsner extra hell", Amount: 3800},
					{Name: "Weizenmalz", Amount: 1200},
					{Name: "Carahell", Amount: 400},
					{Name: "Carapils", Amount: 400},
					{Name: "Invertzucker No3", Amount: 300},
				},
				MainWaterVolume:    22,
				MashTemperature:    65,
				Nachguss:           18,
				MashOutTemperature: 72,
				Rasts: []recipe.Rast{
					{Temperature: 65, Duration: 40},
					{Temperature: 45, Duration: 20},
					{Temperature: 65, Duration: 40},
					{Temperature: 72, Duration: 20},
				},
			},
		},
		{
			Name:     "4S Saison",
			FileName: "4S_Saison.json",
			Expected: recipe.MashInstructions{
				Malts: []recipe.Malt{
					{Name: "Wiener", Amount: 1400},
					{Name: "Pilsner", Amount: 1000},
					{Name: "Weizen roh", Amount: 700},
					{Name: "Special W Weyermann", Amount: 70},
				},
				MainWaterVolume:    14,
				MashTemperature:    50,
				Nachguss:           10,
				MashOutTemperature: 73,
				Rasts: []recipe.Rast{
					{Temperature: 67, Duration: 60},
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			bytes, err := os.ReadFile(basePath + tc.FileName)
			require.NoError(err)
			var r MMUMRecipe
			err = json.Unmarshal(bytes, &r)
			require.NoError(err)
			actual, err := getMashInstructions(&r)
			require.NoError(err)
			require.Equal(tc.Expected, *actual)
		})
	}
}

func TestGetHopInstructions(t *testing.T) {
	require := require.New(t)
	basePath := "../../../test/recipe/mmum/"
	type testCase struct {
		Name     string
		FileName string
		Expected recipe.HopInstructions
	}
	testCases := []testCase{
		{
			Name:     "Hula Hula IPA",
			FileName: "Hula_Hula_IPA.json",
			Expected: recipe.HopInstructions{
				TotalCookingTime: 75,
				Hops: []recipe.Hops{
					{Name: "Simcoe (VW)", Amount: 34, Alpha: 12.5, Duration: 75, DryHop: false, Vorderwuerze: true},
					{Name: "Simcoe", Amount: 180, Alpha: 12.5, Duration: 0, DryHop: false},
					{Name: "Simcoe", Amount: 75, Alpha: 0, Duration: 0, DryHop: true},
					{Name: "Citra", Amount: 100, Alpha: 0, Duration: 0, DryHop: true},
					{Name: "Mosaic", Amount: 100, Alpha: 0, Duration: 0, DryHop: true},
				},
				AdditionalIngredients: nil,
			},
		},
		{
			Name:     "Callippo Mango",
			FileName: "Callippo_Mango.json",
			Expected: recipe.HopInstructions{
				TotalCookingTime: 60,
				Hops: []recipe.Hops{
					{Name: "Callista Nectar", Amount: 20, Alpha: 6.5, Duration: 60, DryHop: false},
					{Name: "Callista Nectar", Amount: 30, Alpha: 6.5, Duration: 0, DryHop: false},
					{Name: "Callista Nectar", Amount: 50, Alpha: 0, Duration: 0, DryHop: true},
				},
				AdditionalIngredients: nil,
			},
		},
		{
			Name:     "4S Saison",
			FileName: "4S_Saison.json",
			Expected: recipe.HopInstructions{
				TotalCookingTime: 60,
				Hops: []recipe.Hops{
					{Name: "Saphir", Amount: 40, Alpha: 4.3, Duration: 30, DryHop: false},
					{Name: "Styrian Celeia", Amount: 25, Alpha: 3.4, Duration: 5, DryHop: false},
					{Name: "Sorachi Ace", Amount: 20, Alpha: 9, Duration: 0, DryHop: false},
					{Name: "Simcoe", Amount: 60, Alpha: 12.9, Duration: 0, DryHop: false},
				},
				AdditionalIngredients: []recipe.AdditionalIngredient{
					{Name: "Demerara Zucker", Amount: 360, Duration: 10},
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			bytes, err := os.ReadFile(basePath + tc.FileName)
			require.NoError(err)
			var r MMUMRecipe
			err = json.Unmarshal(bytes, &r)
			require.NoError(err)
			actual, err := getHopInstructions(&r)
			require.NoError(err)
			require.Equal(tc.Expected, *actual)
		})
	}
}

func TestGetFermentationInstructions(t *testing.T) {
	require := require.New(t)
	basePath := "../../../test/recipe/mmum/"
	type testCase struct {
		Name     string
		FileName string
		Expected recipe.FermentationInstructions
	}
	testCases := []testCase{
		{
			Name:     "Hula Hula IPA",
			FileName: "Hula_Hula_IPA.json",
			Expected: recipe.FermentationInstructions{
				Yeast:       recipe.Yeast{Name: "WY 1007"},
				Temperature: "18-20",
				AdditionalIngredients: []recipe.AdditionalIngredient{
					{Name: "Cryo Citra", Amount: 60, Duration: 0},
					{Name: "Cryo Simcoe", Amount: 60, Duration: 0},
					{Name: "Motueka", Amount: 40, Duration: 0},
				},
				Carbonation: 5.5,
			},
		},
		{
			Name:     "Callippo Mango",
			FileName: "Callippo_Mango.json",
			Expected: recipe.FermentationInstructions{
				Yeast:       recipe.Yeast{Name: "Philly Sour"},
				Temperature: "24",
				AdditionalIngredients: []recipe.AdditionalIngredient{
					{Name: "TK Mango", Amount: 3000, Duration: 0},
				},
				Carbonation: 4.5,
			},
		},
		{
			Name:     "4S Saison",
			FileName: "4S_Saison.json",
			Expected: recipe.FermentationInstructions{
				Yeast:                 recipe.Yeast{Name: "Mangrove Jacks M29 French Saison"},
				Temperature:           "21",
				AdditionalIngredients: nil,
				Carbonation:           6,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			bytes, err := os.ReadFile(basePath + tc.FileName)
			require.NoError(err)
			var r MMUMRecipe
			err = json.Unmarshal(bytes, &r)
			require.NoError(err)
			actual, err := getFermentationInstructions(&r)
			require.NoError(err)
			require.Equal(tc.Expected, *actual)
		})
	}
}
