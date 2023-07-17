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
				MashOutTemperature: 72,
				Rasts: []recipe.Rast{
					{Temperature: 65, Duration: 40},
					{Temperature: 45, Duration: 20},
					{Temperature: 65, Duration: 40},
					{Temperature: 72, Duration: 20},
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
