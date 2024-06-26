package sql

import (
	"brewday/internal/recipe"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"database/sql"

	"github.com/stretchr/testify/require"

	_ "github.com/mattn/go-sqlite3"
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
	db, err := sql.Open("sqlite3", "file:"+fileName+"?_foreign_keys=true")
	require.NoError(err)
	store, err := NewPersistentStore(db)
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
			db, err := sql.Open("sqlite3", "file:"+fileName+"?_foreign_keys=true")
			require.NoError(err)
			store, err := NewPersistentStore(db)
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
			db, err := sql.Open("sqlite3", "file:"+fileName+"?_foreign_keys=true")
			require.NoError(err)
			store, err := NewPersistentStore(db)
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
	db, err := sql.Open("sqlite3", "file:"+fileName+"?_foreign_keys=true")
	require.NoError(err)
	store, err := NewPersistentStore(db)
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

func TestUpdateResult(t *testing.T) {
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
		{
			Name:        "Update vol before boil",
			RecipeID:    "test",
			UpdateType:  recipe.ResultVolumeBeforeBoil,
			UpdateValue: 12,
			Expected: &recipe.RecipeResults{
				VolumeBeforeBoil: 12,
			},
			Error: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			fileName := strings.ToLower(strings.TrimSpace("testupresult_"+tc.Name)) + ".sqlite"
			db, err := sql.Open("sqlite3", "file:"+fileName+"?_foreign_keys=true")
			require.NoError(err)
			store, err := NewPersistentStore(db)
			require.NoError(err)
			id, err := store.Store(&testRecipe)
			require.NoError(err)
			defer os.Remove(fileName)
			if tc.RecipeID == "test" {
				tc.RecipeID = id
			}
			err = store.UpdateResult(tc.RecipeID, tc.UpdateType, tc.UpdateValue)
			if tc.Error {
				require.Error(err)
			} else {
				require.NoError(err)
				actual, err := store.RetrieveResults(tc.RecipeID)
				require.NoError(err)
				require.Equal(tc.Expected, actual)
				val, err := store.RetrieveResult(tc.RecipeID, tc.UpdateType)
				require.NoError(err)
				require.InDelta(tc.UpdateValue, val, 0.0001)
			}
		})
	}
}

func TestUpdateSGs(t *testing.T) {
	require := require.New(t)
	testCases := []struct {
		Name  string
		ToAdd map[string][]*recipe.SGMeasurement
		Error bool
	}{
		{
			Name: "Single sg",
			ToAdd: map[string][]*recipe.SGMeasurement{
				"recipe1": {
					{Value: 1.013, Date: time.Now().Format("2006-01-02")},
				},
			},
			Error: false,
		},
		{
			Name: "Multiple sg",
			ToAdd: map[string][]*recipe.SGMeasurement{
				"recipe1": {
					{Value: 1.013, Date: time.Now().Format("2006-01-02")},
					{Value: 1.011, Date: time.Now().Add(10 * time.Second).Format("2006-01-02")},
				},
			},
			Error: false,
		},
		{
			Name: "Multiple sg multiple recipes",
			ToAdd: map[string][]*recipe.SGMeasurement{
				"recipe1": {
					{Value: 1.013, Date: time.Now().Format("2006-01-02")},
					{Value: 1.011, Date: time.Now().Add(10 * time.Second).Format("2006-01-02")},
				},
				"recipe2": {
					{Value: 1.013, Date: time.Now().Format("2006-01-02")},
					{Value: 1.011, Date: time.Now().Add(10 * time.Second).Format("2006-01-02")},
				},
			},
			Error: false,
		},
		{
			Name: "Order is respected",
			ToAdd: map[string][]*recipe.SGMeasurement{
				"recipe1": {
					{Value: 1.013, Date: time.Now().Format("2006-01-02")},
					{Value: 1.011, Date: time.Now().Add(10 * time.Second).Format("2006-01-02")},
					{Value: 1.011, Date: time.Now().Add(-10 * time.Second).Format("2006-01-02")},
				},
			},
			Error: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			fileName := strings.ToLower(strings.TrimSpace("testupresult_"+tc.Name)) + ".sqlite"
			db, err := sql.Open("sqlite3", "file:"+fileName+"?_foreign_keys=true")
			require.NoError(err)
			store, err := NewPersistentStore(db)
			require.NoError(err)
			defer os.Remove(fileName)
			for recipeTitle, sgs := range tc.ToAdd {
				id, err := store.Store(&recipe.Recipe{Name: recipeTitle})
				require.NoError(err)
				for _, sg := range sgs {
					err = store.AddMainFermSG(id, sg)
					if tc.Error {
						require.Error(err)
					} else {
						require.NoError(err)
					}
				}
				if !tc.Error {
					realSGs, err := store.RetrieveMainFermSGs(id)
					require.NoError(err)
					for i := 0; i < len(tc.ToAdd[recipeTitle]); i++ {
						require.Equal(tc.ToAdd[recipeTitle][i], realSGs[i])
					}
				}
			}
		})
	}
}

func TestDates(t *testing.T) {
	require := require.New(t)
	t1 := time.Now().Add(1 * time.Second)
	t2 := time.Now().Add(2 * time.Second)
	t3 := time.Now().Add(3 * time.Second)
	t4 := time.Now().Add(4 * time.Second)
	testCases := []struct {
		Name  string
		ToAdd map[string][]time.Time
		Error bool
	}{
		{
			Name: "Single date",
			ToAdd: map[string][]time.Time{
				"recipe_1": {t1},
			},
			Error: false,
		},
		{
			Name: "Multiple dates",
			ToAdd: map[string][]time.Time{
				"recipe_1": {t1, t2},
			},
			Error: false,
		},
		{
			Name: "Multiple dates and recipes",
			ToAdd: map[string][]time.Time{
				"recipe_1": {t1, t2},
				"recipe_2": {t3, t4},
			},
			Error: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			fileName := strings.ToLower(strings.TrimSpace("testupresult_"+tc.Name)) + ".sqlite"
			db, err := sql.Open("sqlite3", "file:"+fileName+"?_foreign_keys=true")
			require.NoError(err)
			store, err := NewPersistentStore(db)
			require.NoError(err)
			defer os.Remove(fileName)
			for recipeTitle, dates := range tc.ToAdd {
				id, err := store.Store(&recipe.Recipe{Name: recipeTitle})
				require.NoError(err)
				for i, date := range dates {
					err = store.AddDate(id, &date, fmt.Sprintf("%s_%d", recipeTitle, i))
					if tc.Error {
						require.Error(err)
					} else {
						require.NoError(err)
					}
				}
				if !tc.Error {
					realDates, err := store.RetrieveDates(id, recipeTitle)
					require.NotEmpty(realDates)
					require.NoError(err)
					realDatesString := make([]string, len(realDates))
					for i, d := range realDates {
						realDatesString[i] = d.Format(time.RFC3339)
					}
					datesString := make([]string, len(dates))
					for i, d := range dates {
						datesString[i] = d.Format(time.RFC3339)
					}
					require.ElementsMatch(datesString, realDatesString)
				}
			}
		})
	}
}

func TestUpdateSugarResults(t *testing.T) {
	require := require.New(t)
	testCases := []struct {
		Name  string
		ToAdd map[string][]*recipe.PrimingSugarResult
		Error bool
	}{
		{
			Name: "Single sugar result",
			ToAdd: map[string][]*recipe.PrimingSugarResult{
				"recipe1": {
					{Water: 10, Amount: 20.5, Alcohol: 5.5},
				},
			},
			Error: false,
		},
		{
			Name: "Multiple sugar result",
			ToAdd: map[string][]*recipe.PrimingSugarResult{
				"recipe1": {
					{Water: 10, Amount: 20.5, Alcohol: 5.5},
					{Water: 20, Amount: 20.5, Alcohol: 4.5},
				},
			},
			Error: false,
		},
		{
			Name: "Multiple results multiple recipes",
			ToAdd: map[string][]*recipe.PrimingSugarResult{
				"recipe1": {
					{Water: 10, Amount: 20.5, Alcohol: 5.5},
					{Water: 20, Amount: 20.5, Alcohol: 4.5},
				},
				"recipe2": {
					{Water: 30, Amount: 10.4, Alcohol: 9.5},
					{Water: 40, Amount: 10.4, Alcohol: 8.5},
				},
			},
			Error: false,
		},
		{
			Name: "Order is respected",
			ToAdd: map[string][]*recipe.PrimingSugarResult{
				"recipe1": {
					{Water: 10, Amount: 20.5, Alcohol: 5.5},
					{Water: 20, Amount: 20.5, Alcohol: 4.5},
					{Water: 5, Amount: 20.5, Alcohol: 6.5},
				},
			},
			Error: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			fileName := strings.ToLower(strings.TrimSpace("testupresult_"+tc.Name)) + ".sqlite"
			db, err := sql.Open("sqlite3", "file:"+fileName+"?_foreign_keys=true")
			require.NoError(err)
			store, err := NewPersistentStore(db)
			require.NoError(err)
			defer os.Remove(fileName)
			for recipeTitle, results := range tc.ToAdd {
				id, err := store.Store(&recipe.Recipe{Name: recipeTitle})
				require.NoError(err)
				for _, res := range results {
					err = store.AddSugarResult(id, res)
					if tc.Error {
						require.Error(err)
					} else {
						require.NoError(err)
					}
				}
				if !tc.Error {
					realResults, err := store.RetrieveSugarResults(id)
					require.NoError(err)
					for i := 0; i < len(tc.ToAdd[recipeTitle]); i++ {
						require.Equal(tc.ToAdd[recipeTitle][i], realResults[i])
					}
				}
			}
		})
	}
}

func TestBoolFlags(t *testing.T) {
	require := require.New(t)
	type boolFlag struct {
		Value bool
		Name  string
	}
	testCases := []struct {
		Name  string
		ToAdd map[string][]*boolFlag
		Error bool
	}{
		{
			Name: "Store single false",
			ToAdd: map[string][]*boolFlag{
				"recipe1": {
					{Value: false, Name: "name_1"},
				},
			},
			Error: false,
		},
		{
			Name: "Store single true",
			ToAdd: map[string][]*boolFlag{
				"recipe2": {
					{Value: true, Name: "name_2"},
				},
			},
			Error: false,
		},
		{
			Name: "Same name different recipe",
			ToAdd: map[string][]*boolFlag{
				"recipe3": {
					{Value: false, Name: "name_3"},
				},
				"recipe4": {
					{Value: true, Name: "name_3"},
				},
			},
			Error: false,
		},
		{
			Name: "Overwrite",
			ToAdd: map[string][]*boolFlag{
				"recipe5": {
					{Value: false, Name: "name_4"},
					{Value: true, Name: "name_4"},
				},
			},
			Error: false,
		},
		{
			Name: "SQL Injection in name",
			ToAdd: map[string][]*boolFlag{
				"recipe_6": {
					{Value: true, Name: "5'; DROP TABLE bool_flags; --"},
				},
			},
			Error: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			fileName := strings.ToLower(strings.TrimSpace("testupresult_"+tc.Name)) + ".sqlite"
			db, err := sql.Open("sqlite3", "file:"+fileName+"?_foreign_keys=true")
			require.NoError(err)
			store, err := NewPersistentStore(db)
			require.NoError(err)
			defer os.Remove(fileName)
			for recipeTitle, v := range tc.ToAdd {
				id, err := store.Store(&recipe.Recipe{Name: recipeTitle})
				require.NoError(err)
				for _, b := range v {
					err = store.AddBoolFlag(id, b.Name, b.Value)
					if tc.Error {
						require.Error(err)
					} else {
						require.NoError(err)
						actual, err := store.RetrieveBoolFlag(id, b.Name)
						require.NoError(err)
						require.Equal(b.Value, actual)
					}
				}
			}

		})
	}
}
