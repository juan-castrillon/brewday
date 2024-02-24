package recipe

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetTotalMaltWeight(t *testing.T) {
	require := require.New(t)
	type testCase struct {
		Name     string
		Mash     MashInstructions
		Expected float32
	}
	testCases := []testCase{
		{
			Name: "Single malt",
			Mash: MashInstructions{
				Malts: []Malt{
					{Name: "Carahell", Amount: 5600},
				},
			},
			Expected: 5600,
		},
		{
			Name: "Multiple malts",
			Mash: MashInstructions{
				Malts: []Malt{
					{Name: "Carahell", Amount: 5600},
					{Name: "Pilsner", Amount: 5000},
					{Name: "Haferflocken", Amount: 500},
				},
			},
			Expected: 11100,
		},
		{
			Name: "No malts",
			Mash: MashInstructions{
				Malts: []Malt{},
			},
			Expected: 0,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			require.Equal(tc.Expected, tc.Mash.GetTotalMaltWeight())
		})
	}

}

func TestGetStatus(t *testing.T) {
	require := require.New(t)
	type testCase struct {
		Name           string
		Recipe         *Recipe
		ExpectedStatus RecipeStatus
		ExpectedParams []string
	}
	testCases := []testCase{
		{
			Name:           "No params",
			Recipe:         &Recipe{status: RecipeStatusCreated, statusParams: []string{}},
			ExpectedStatus: RecipeStatusCreated,
			ExpectedParams: []string{},
		},
		{
			Name:           "Nil params",
			Recipe:         &Recipe{status: RecipeStatusCreated, statusParams: nil},
			ExpectedStatus: RecipeStatusCreated,
			ExpectedParams: nil,
		},
		{
			Name:           "With one param",
			Recipe:         &Recipe{status: RecipeStatusCreated, statusParams: []string{"test"}},
			ExpectedStatus: RecipeStatusCreated,
			ExpectedParams: []string{"test"},
		},
		{
			Name:           "With multiple params",
			Recipe:         &Recipe{status: RecipeStatusCreated, statusParams: []string{"test", "1", "2.0"}},
			ExpectedStatus: RecipeStatusCreated,
			ExpectedParams: []string{"test", "1", "2.0"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			status, params := tc.Recipe.GetStatus()
			require.Equal(tc.ExpectedStatus, status)
			require.Equal(tc.ExpectedParams, params)
		})
	}
}

func TestSetStatus(t *testing.T) {
	require := require.New(t)
	type testCase struct {
		Name   string
		Recipe *Recipe
		Status RecipeStatus
		Params []string
	}
	testCases := []testCase{
		{
			Name:   "No params to params",
			Recipe: &Recipe{status: RecipeStatusCreated, statusParams: []string{}},
			Status: RecipeStatusMashing,
			Params: []string{"test", "1", "2.0"},
		},
		{
			Name:   "No params to nil",
			Recipe: &Recipe{status: RecipeStatusCreated, statusParams: []string{}},
			Status: RecipeStatusMashing,
			Params: nil,
		},
		{
			Name:   "Params to no params",
			Recipe: &Recipe{status: RecipeStatusCreated, statusParams: []string{"test", "1", "2.0"}},
			Status: RecipeStatusMashing,
			Params: []string{},
		},
		{
			Name:   "Params to nil",
			Recipe: &Recipe{status: RecipeStatusCreated, statusParams: []string{"test", "1", "2.0"}},
			Status: RecipeStatusMashing,
			Params: nil,
		},
		{
			Name:   "Nil to no params",
			Recipe: &Recipe{status: RecipeStatusCreated, statusParams: nil},
			Status: RecipeStatusMashing,
			Params: []string{},
		},
		{
			Name:   "Nil to params",
			Recipe: &Recipe{status: RecipeStatusCreated, statusParams: nil},
			Status: RecipeStatusMashing,
			Params: []string{"test", "1", "2.0"},
		},
		{
			Name:   "Params to params",
			Recipe: &Recipe{status: RecipeStatusCreated, statusParams: []string{"test", "1", "2.0"}},
			Status: RecipeStatusMashing,
			Params: []string{"test", "1", "2.0"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			tc.Recipe.SetStatus(tc.Status, tc.Params...)
			require.Equal(tc.Status, tc.Recipe.status)
			require.Equal(tc.Params, tc.Recipe.statusParams)
		})
	}
}

func TestGetStatusString(t *testing.T) {
	require := require.New(t)
	type testCase struct {
		Name     string
		Recipe   *Recipe
		Expected string
	}
	testCases := []testCase{
		{
			Name:     "Created",
			Recipe:   &Recipe{status: RecipeStatusCreated},
			Expected: "Created",
		},
		{
			Name:     "Mashing",
			Recipe:   &Recipe{status: RecipeStatusMashing},
			Expected: "Mashing",
		},
		{
			Name:     "Lautering",
			Recipe:   &Recipe{status: RecipeStatusLautering},
			Expected: "Lautering",
		},
		{
			Name:     "Boiling",
			Recipe:   &Recipe{status: RecipeStatusBoiling},
			Expected: "Boiling",
		},
		{
			Name:     "Cooling",
			Recipe:   &Recipe{status: RecipeStatusCooling},
			Expected: "Cooling",
		},
		{
			Name:     "Pre-fermentation",
			Recipe:   &Recipe{status: RecipeStatusPreFermentation},
			Expected: "Pre-fermentation",
		},
		{
			Name:     "Fermenting",
			Recipe:   &Recipe{status: RecipeStatusFermenting},
			Expected: "Fermenting",
		},
		{
			Name:     "Bottled",
			Recipe:   &Recipe{status: RecipeStatusBottled},
			Expected: "Bottled",
		},
		{
			Name:     "Fridge",
			Recipe:   &Recipe{status: RecipeStatusFridge},
			Expected: "Fridge",
		},
		{
			Name:     "Finished",
			Recipe:   &Recipe{status: RecipeStatusFinished},
			Expected: "Finished",
		},
		{
			Name:     "Unknown",
			Recipe:   &Recipe{status: RecipeStatusUnknown},
			Expected: "Unknown",
		},
		{
			Name:     "Status not set",
			Recipe:   &Recipe{},
			Expected: "Unknown",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			status := tc.Recipe.GetStatusString()
			require.Equal(tc.Expected, status)
		})
	}
}
