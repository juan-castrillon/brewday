package sql

import (
	"brewday/internal/recipe"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var testRecipe = recipe.Recipe{
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

func TestStoreAndRetrieve(t *testing.T) {
	require := require.New(t)
	fileName := strings.ToLower(strings.TrimSpace(t.Name())) + ".sqlite"
	store, err := NewPersistentStore(fileName)
	require.NoError(err)
	defer os.Remove(fileName)
	testCases := []struct {
		Name         string
		Modify       bool
		SkipStore    bool
		SkipStoreID  string
		Status       recipe.RecipeStatus
		StatusParams []string
		Error        bool
	}{
		{Name: "Original", Modify: false},
		{Name: "StatusNoParams", Modify: true, Status: recipe.RecipeStatusLautering, StatusParams: []string{}},
		{Name: "StatusParams", Modify: true, Status: recipe.RecipeStatusPreFermentation, StatusParams: []string{
			"water", "15.324", "0.032",
		}},
		{Name: "Read non existent recipe", Modify: false, SkipStore: true, SkipStoreID: "10", Error: true},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			r := testRecipe
			if tc.Modify {
				r.SetStatus(tc.Status, tc.StatusParams...)
			}
			var id string
			if !tc.SkipStore {
				id, err = store.Store(&r)
				require.NoError(err)
			} else {
				id = tc.SkipStoreID
			}
			actual, err := store.Retrieve(id)
			if tc.Error {
				require.Error(err)
			} else {
				require.NoError(err)
				r.ID = id
				require.Equal(&r, actual)
				actualStatus, actualParams := actual.GetStatus()
				expectedStatus, expectedParams := r.GetStatus()
				require.Equal(expectedStatus, actualStatus)
				require.ElementsMatch(expectedParams, actualParams)
			}
		})
	}
}

func TestList(t *testing.T) {
	require := require.New(t)
	testCases := []struct {
		Name           string
		RecipesToStore int
		Expected       []*recipe.Recipe
		Status         []recipe.RecipeStatus
		Error          bool
	}{
		{
			Name:           "Single recipe",
			RecipesToStore: 1,
			Expected: []*recipe.Recipe{
				{Name: "Huell Saison0", Style: "Saison"},
			},
			Status: []recipe.RecipeStatus{
				recipe.RecipeStatusBoiling,
			},
			Error: false,
		},
		{
			Name:           "No recipe",
			RecipesToStore: 0,
			Expected:       []*recipe.Recipe{},
			Status:         []recipe.RecipeStatus{},
			Error:          false,
		},
		{
			Name:           "Multiple recipes",
			RecipesToStore: 3,
			Expected: []*recipe.Recipe{
				{Name: "Huell Saison0", Style: "Saison"},
				{Name: "Huell Saison1", Style: "Saison"},
				{Name: "Huell Saison2", Style: "Saison"},
			},
			Status: []recipe.RecipeStatus{
				recipe.RecipeStatusBoiling,
				recipe.RecipeStatusLautering,
				recipe.RecipeStatusFermenting,
			},
			Error: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			fileName := strings.ToLower(strings.TrimSpace("testlist_"+tc.Name)) + ".sqlite"
			store, err := NewPersistentStore(fileName)
			require.NoError(err)
			defer os.Remove(fileName)
			for i := 0; i < tc.RecipesToStore; i++ {
				r := testRecipe
				r.Name = fmt.Sprintf("%s%d", r.Name, i)
				r.SetStatus(tc.Status[i])
				id, err := store.Store(&r)
				require.NoError(err)
				tc.Expected[i].ID = id
				tc.Expected[i].SetStatus(tc.Status[i])
			}
			actual, err := store.List()
			if tc.Error {
				require.Error(err)
			} else {
				require.NoError(err)
				require.Equal(tc.RecipesToStore, len(actual))
				for i := 0; i < tc.RecipesToStore; i++ {
					require.Equal(tc.Expected[i].Name, actual[i].Name)
					require.Equal(tc.Expected[i].Style, actual[i].Style)
					actualStatus, _ := actual[i].GetStatus()
					require.Equal(tc.Status[i], actualStatus)
				}
			}
		})
	}
}

func TestDelete(t *testing.T) {
	require := require.New(t)
	testCases := []struct {
		Name           string
		RecipesToStore int
		DeleteID       string
		Expected       []*recipe.Recipe
		Error          bool
	}{
		{
			Name:           "Single recipe",
			RecipesToStore: 1,
			DeleteID:       "1",
			Expected:       []*recipe.Recipe{},
			Error:          false,
		},
		{
			Name:           "Delete first from two",
			RecipesToStore: 2,
			DeleteID:       "1",
			Expected:       []*recipe.Recipe{{Name: "Huell Saison1", Style: "Saison"}},
			Error:          false,
		},
		{
			Name:           "Delete second from two",
			RecipesToStore: 2,
			DeleteID:       "2",
			Expected:       []*recipe.Recipe{{Name: "Huell Saison0", Style: "Saison"}},
			Error:          false,
		},
		{
			Name:           "Delete middle from three",
			RecipesToStore: 3,
			DeleteID:       "2",
			Expected: []*recipe.Recipe{
				{Name: "Huell Saison0", Style: "Saison"},
				{Name: "Huell Saison2", Style: "Saison"},
			},
			Error: false,
		},
		{
			Name:           "Delete from empty db",
			RecipesToStore: 0,
			DeleteID:       "1",
			Expected:       []*recipe.Recipe{},
			Error:          false,
		},
		{
			Name:           "Delete non existent ID",
			RecipesToStore: 1,
			DeleteID:       "3",
			Expected: []*recipe.Recipe{
				{Name: "Huell Saison0", Style: "Saison"},
			},
			Error: false,
		},
		{
			Name:           "Delete non valid ID",
			RecipesToStore: 1,
			DeleteID:       "invalid",
			Expected: []*recipe.Recipe{
				{Name: "Huell Saison0", Style: "Saison"},
			},
			Error: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			fileName := strings.ToLower(strings.TrimSpace("testdelete_"+tc.Name)) + ".sqlite"
			store, err := NewPersistentStore(fileName)
			require.NoError(err)
			defer os.Remove(fileName)
			for i := 0; i < tc.RecipesToStore; i++ {
				r := testRecipe
				r.Name = fmt.Sprintf("%s%d", r.Name, i)
				_, err := store.Store(&r)
				require.NoError(err)
			}
			err = store.Delete(tc.DeleteID)
			if tc.Error {
				require.Error(err)
			} else {
				require.NoError(err)
				actual, err := store.List()
				require.NoError(err)
				require.Equal(len(tc.Expected), len(actual))
				for i := 0; i < len(actual); i++ {
					require.Equal(tc.Expected[i].Name, actual[i].Name)
					require.Equal(tc.Expected[i].Style, actual[i].Style)
				}
			}
		})
	}
}

func TestUpdateStatus(t *testing.T) {
	require := require.New(t)
	fileName := strings.ToLower(strings.TrimSpace(t.Name())) + ".sqlite"
	store, err := NewPersistentStore(fileName)
	require.NoError(err)
	id, err := store.Store(&testRecipe)
	require.NoError(err)
	defer os.Remove(fileName)
	testCases := []struct {
		Name         string
		Status       recipe.RecipeStatus
		StatusParams []string
	}{
		{
			Name:         "Status with no params",
			Status:       recipe.RecipeStatusLautering,
			StatusParams: []string{},
		},
		{
			Name:         "Status with params",
			Status:       recipe.RecipeStatusPreFermentation,
			StatusParams: []string{"water", "15.324", "0.032"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			err = store.UpdateStatus(id, tc.Status, tc.StatusParams...)
			require.NoError(err)
			actual, err := store.Retrieve(id)
			require.NoError(err)
			actualStatus, actualParams := actual.GetStatus()
			require.Equal(testRecipe.Name, actual.Name)
			require.Equal(tc.Status, actualStatus)
			require.ElementsMatch(tc.StatusParams, actualParams)
		})
	}
}

func TestUpdateResults(t *testing.T) {
	require := require.New(t)
	testCases := []struct {
		Name        string
		RecipeID    string
		UpdateType  recipe.ResultType
		UpdateValue float32
		Expected    *recipe.RecipeResults
		Error       bool
	}{
		{
			Name:        "Update hot wort volume",
			RecipeID:    "test",
			UpdateType:  recipe.ResultHotWortVolume,
			UpdateValue: 10.5,
			Expected: &recipe.RecipeResults{
				HotWortVolume: 10.5,
			},
			Error: false,
		},
		{
			Name:        "Update original sg",
			RecipeID:    "test",
			UpdateType:  recipe.ResultOriginalGravity,
			UpdateValue: 1.068,
			Expected: &recipe.RecipeResults{
				OriginalGravity: 1.068,
			},
			Error: false,
		},
		{
			Name:        "Update final sg",
			RecipeID:    "test",
			UpdateType:  recipe.ResultFinalGravity,
			UpdateValue: 1.059,
			Expected: &recipe.RecipeResults{
				FinalGravity: 1.059,
			},
			Error: false,
		},
		{
			Name:        "Update alcohol",
			RecipeID:    "test",
			UpdateType:  recipe.ResultAlcohol,
			UpdateValue: 5.5,
			Expected: &recipe.RecipeResults{
				Alcohol: 5.5,
			},
			Error: false,
		},
		{
			Name:        "Update main volume",
			RecipeID:    "test",
			UpdateType:  recipe.ResultMainFermentationVolume,
			UpdateValue: 10.3,
			Expected: &recipe.RecipeResults{
				MainFermentationVolume: 10.3,
			},
			Error: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			fileName := strings.ToLower(strings.TrimSpace("testupresult_"+tc.Name)) + ".sqlite"
			store, err := NewPersistentStore(fileName)
			require.NoError(err)
			id, err := store.Store(&testRecipe)
			require.NoError(err)
			defer os.Remove(fileName)
			if tc.RecipeID == "test" {
				tc.RecipeID = id
			}
			err = store.UpdateResults(tc.RecipeID, tc.UpdateType, tc.UpdateValue)
			if tc.Error {
				require.Error(err)
			} else {
				require.NoError(err)
				actual, err := store.RetrieveResults(tc.RecipeID)
				require.NoError(err)
				require.Equal(tc.Expected, actual)
			}
		})
	}
}
