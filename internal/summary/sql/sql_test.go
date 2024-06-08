package sql

import (
	"brewday/internal/summary"
	"database/sql"
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func provisionDB(t *testing.T, db *sql.DB, recipes []string) {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS recipes (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL
	)`)
	require.NoError(t, err)
	for _, r := range recipes {
		_, err := db.Exec(`INSERT INTO recipes (name) VALUES (?)`, r)
		require.NoError(t, err)
	}

}

func TestAddSummary(t *testing.T) {
	require := require.New(t)
	fileName := strings.ToLower(strings.TrimSpace(t.Name())) + ".sqlite"
	db, err := sql.Open("sqlite3", "file:"+fileName+"?_foreign_keys=true")
	require.NoError(err)
	provisionDB(t, db, []string{"recipe1", "recipe2", "recipe3", "recipe4"})
	store, err := NewSummaryPersistentStore(db)
	require.NoError(err)
	defer os.Remove(fileName)
	getSt, err := db.Prepare(`SELECT title FROM summaries WHERE recipe_id = ?`)
	require.NoError(err)
	testCases := []struct {
		Name     string
		Title    string
		RecipeID string
		Error    bool
	}{
		{
			Name:     "Normal case",
			Title:    "title1",
			RecipeID: "1",
			Error:    false,
		},
		{
			Name:     "Empty title",
			RecipeID: "2",
			Title:    "",
			Error:    false,
		},
		{
			Name:     "SQL Injection in RecipeID",
			RecipeID: "4'; DROP TABLE summaries; --",
			Title:    "title4",
			Error:    true,
		},
		{
			Name:     "SQL Injection in Title",
			RecipeID: "3",
			Title:    "5'; DROP TABLE summaries; --",
			Error:    false,
		},
		{
			Name:     "Special Characters in Title",
			RecipeID: "4",
			Title:    "title7$#%@^&* ()",
			Error:    false,
		},
		{
			Name:     "Empty recipeID",
			RecipeID: "",
			Title:    "my-title",
			Error:    true,
		},
		{
			Name:     "Non-existing recipeID",
			RecipeID: "5",
			Title:    "no-title",
			Error:    true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			err := store.AddSummary(tc.RecipeID, tc.Title)
			if tc.Error {
				require.Error(err)
			} else {
				require.NoError(err)
				var title string
				err = getSt.QueryRow(tc.RecipeID).Scan(&title)
				require.NoError(err)
				require.Equal(tc.Title, title)
			}
		})
	}
}

func TestDeleteSummary(t *testing.T) {
	require := require.New(t)
	fileName := strings.ToLower(strings.TrimSpace(t.Name())) + ".sqlite"
	db, err := sql.Open("sqlite3", "file:"+fileName+"?_foreign_keys=true")
	require.NoError(err)
	provisionDB(t, db, []string{"recipe1", "recipe2", "recipe3", "recipe4"})
	store, err := NewSummaryPersistentStore(db)
	require.NoError(err)
	defer os.Remove(fileName)
	getSt, err := db.Prepare(`SELECT title FROM summaries WHERE recipe_id = ?`)
	require.NoError(err)
	testCases := []struct {
		Name     string
		RecipeID string
		ToAdd    []struct {
			id    string
			title string
		}
		Error bool
	}{
		{
			Name:     "Normal",
			RecipeID: "1",
			ToAdd: []struct {
				id    string
				title string
			}{
				{id: "1", title: "title1"},
				{id: "2", title: "title2"},
			},
			Error: false,
		},
		{
			Name:     "SQL Injection in recipe id",
			RecipeID: "4'; DROP TABLE summaries; --",
			ToAdd: []struct {
				id    string
				title string
			}{
				{id: "4", title: "title4"},
			},
			Error: false,
		},
		{
			Name:     "Non existing recipe id",
			RecipeID: "5",
			Error:    false,
		},
		{
			Name:     "Empty recipe id",
			RecipeID: "",
			Error:    true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			for _, summ := range tc.ToAdd {
				err := store.AddSummary(summ.id, summ.title)
				require.NoError(err)
			}
			err = store.DeleteSummary(tc.RecipeID)
			if tc.Error {
				require.Error(err)
			} else {
				require.NoError(err)
				var title string
				err = getSt.QueryRow(tc.RecipeID).Scan(&title)
				require.ErrorIs(err, sql.ErrNoRows)
			}
		})
	}
}

func TestAddMashTemp(t *testing.T) {
	require := require.New(t)
	fileName := strings.ToLower(strings.TrimSpace(t.Name())) + ".sqlite"
	db, err := sql.Open("sqlite3", "file:"+fileName+"?_foreign_keys=true")
	require.NoError(err)
	provisionDB(t, db, []string{"recipe1", "recipe2", "recipe3", "recipe4"})
	store, err := NewSummaryPersistentStore(db)
	require.NoError(err)
	defer os.Remove(fileName)
	for i := 1; i <= 3; i++ {
		num := strconv.Itoa(i)
		require.NoError(store.AddSummary(num, "t"+num))
	}

	getSt, err := db.Prepare(`SELECT mash_temp, mash_notes FROM summaries WHERE recipe_id = ?`)
	require.NoError(err)
	testCases := []struct {
		Name     string
		Temp     float32
		Notes    string
		RecipeID string
		SkipRead bool
		Error    bool
	}{
		{
			Name:     "Normal case",
			Temp:     55.5,
			Notes:    "hello",
			RecipeID: "1",
			Error:    false,
		},
		{
			Name:     "no notes",
			Temp:     56.5,
			Notes:    "",
			RecipeID: "2",
			Error:    false,
		},
		{
			Name:     "SQL Injection in notes",
			Temp:     55.5,
			Notes:    "5'; DROP TABLE summaries; --",
			RecipeID: "3",
			Error:    false,
		},
		{
			Name:     "Non existing recipe id",
			Temp:     55.5,
			Notes:    "oe",
			RecipeID: "5",
			Error:    false,
			SkipRead: true,
		},
		{
			Name:     "summary not created",
			Temp:     55.5,
			Notes:    "oe",
			RecipeID: "4",
			Error:    false,
			SkipRead: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			err = store.AddMashTemp(tc.RecipeID, tc.Temp, tc.Notes)
			if tc.Error {
				require.Error(err)
			} else {
				require.NoError(err)
				if !tc.SkipRead {
					var temp float64
					var notes string
					err = getSt.QueryRow(tc.RecipeID).Scan(&temp, &notes)
					require.NoError(err)
					require.Equal(tc.Notes, notes)
					require.InDelta(tc.Temp, temp, 0.001)
				}
			}
		})
	}
}

func TestAddRast(t *testing.T) {
	require := require.New(t)
	fileName := strings.ToLower(strings.TrimSpace(t.Name())) + ".sqlite"
	db, err := sql.Open("sqlite3", "file:"+fileName+"?_foreign_keys=true")
	require.NoError(err)
	provisionDB(t, db, []string{"recipe1", "recipe2", "recipe3", "recipe4"})
	store, err := NewSummaryPersistentStore(db)
	require.NoError(err)
	defer os.Remove(fileName)
	for i := 1; i <= 3; i++ {
		num := strconv.Itoa(i)
		require.NoError(store.AddSummary(num, "t"+num))
	}

	getSt, err := db.Prepare(`SELECT mash_rasts FROM summaries WHERE recipe_id = ?`)
	require.NoError(err)
	testCases := []struct {
		Name     string
		Rasts    []*summary.MashRastInfo
		RecipeID string
		SkipRead bool
		Error    bool
	}{
		{
			Name: "Normal case",
			Rasts: []*summary.MashRastInfo{
				{Temperature: 55.5, Time: 60, Notes: "notes1"},
			},
			RecipeID: "1",
			Error:    false,
		},
		{
			Name: "Multiple rasts",
			Rasts: []*summary.MashRastInfo{
				{Temperature: 55.5, Time: 60, Notes: "notes1"},
				{Temperature: 70.5, Time: 30, Notes: "notes2"},
			},
			RecipeID: "2",
			Error:    false,
		},
		{
			Name: "SQL Injection in notes",
			Rasts: []*summary.MashRastInfo{
				{Temperature: 55.5, Time: 60, Notes: "5'; DROP TABLE summaries; --"},
			},
			RecipeID: "3",
			Error:    false,
		},
		{
			Name: "SQL Injection in recipe_id",
			Rasts: []*summary.MashRastInfo{
				{Temperature: 55.5, Time: 60, Notes: "notes1"},
			},
			RecipeID: "5'; DROP TABLE summaries; --",
			Error:    true,
		},
		{
			Name: "Non existing recipe_id",
			Rasts: []*summary.MashRastInfo{
				{Temperature: 55.5, Time: 60, Notes: "notes1"},
			},
			RecipeID: "10",
			Error:    true,
		},
		{
			Name: "Empty recipe_id",
			Rasts: []*summary.MashRastInfo{
				{Temperature: 55.5, Time: 60, Notes: "notes1"},
			},
			RecipeID: "",
			Error:    true,
		},
		{
			Name: "Summary not created",
			Rasts: []*summary.MashRastInfo{
				{Temperature: 55.5, Time: 60, Notes: "notes1"},
			},
			RecipeID: "4",
			Error:    true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			for _, rast := range tc.Rasts {
				err = store.AddRast(tc.RecipeID, rast.Temperature, rast.Time, rast.Notes)
				if tc.Error {
					require.Error(err)
				}
			}
			if tc.Error {
				require.Error(err)
			} else {
				require.NoError(err)
				if !tc.SkipRead {
					var rasts string
					err = getSt.QueryRow(tc.RecipeID).Scan(&rasts)
					require.NoError(err)
					expected, err := json.Marshal(tc.Rasts)
					require.NoError(err)
					require.Equal(string(expected), rasts)
				}
			}
		})
	}
}

func TestAddLauternNotes(t *testing.T) {
	require := require.New(t)
	fileName := strings.ToLower(strings.TrimSpace(t.Name())) + ".sqlite"
	db, err := sql.Open("sqlite3", "file:"+fileName+"?_foreign_keys=true")
	require.NoError(err)
	provisionDB(t, db, []string{"recipe1", "recipe2", "recipe3", "recipe4"})
	store, err := NewSummaryPersistentStore(db)
	require.NoError(err)
	defer os.Remove(fileName)
	for i := 1; i <= 3; i++ {
		num := strconv.Itoa(i)
		require.NoError(store.AddSummary(num, "t"+num))
	}
	getSt, err := db.Prepare(`SELECT lautern_info FROM summaries WHERE recipe_id = ?`)
	require.NoError(err)
	testCases := []struct {
		Name     string
		RecipeID string
		Notes    string
		SkipRead bool
		Error    bool
	}{
		{
			Name:     "Normal case",
			RecipeID: "1",
			Notes:    "notes1",
			Error:    false,
		},
		{
			Name:     "Empty RecipeID",
			RecipeID: "",
			Notes:    "Some notes",
			Error:    true,
		},
		{
			Name:     "Empty Notes",
			RecipeID: "2",
			Notes:    "",
			Error:    false,
		},
		{
			Name:     "SQL Injection in RecipeID",
			RecipeID: "3; DROP TABLE summaries;",
			Notes:    "Some notes",
			Error:    false,
			SkipRead: true,
		},
		{
			Name:     "SQL Injection in Notes",
			RecipeID: "3",
			Notes:    "Some notes; DROP TABLE summaries;",
			Error:    false,
		},
		{
			Name:     "Non-Existing RecipeID",
			RecipeID: "999",
			Notes:    "Some notes",
			Error:    false,
			SkipRead: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			err = store.AddLauternNotes(tc.RecipeID, tc.Notes)
			if tc.Error {
				require.Error(err)
			} else {
				require.NoError(err)
				if !tc.SkipRead {
					var notes string
					require.NoError(getSt.QueryRow(tc.RecipeID).Scan(&notes))
					require.Equal(tc.Notes, notes)
				}
			}
		})
	}
}

func TestAddHopping(t *testing.T) {
	require := require.New(t)
	fileName := strings.ToLower(strings.TrimSpace(t.Name())) + ".sqlite"
	db, err := sql.Open("sqlite3", "file:"+fileName+"?_foreign_keys=true")
	require.NoError(err)
	provisionDB(t, db, []string{"recipe1", "recipe2", "recipe3", "recipe4"})
	store, err := NewSummaryPersistentStore(db)
	require.NoError(err)
	defer os.Remove(fileName)
	for i := 1; i <= 3; i++ {
		num := strconv.Itoa(i)
		require.NoError(store.AddSummary(num, "t"+num))
	}
	getSt, err := db.Prepare(`SELECT hopping_hops FROM summaries WHERE recipe_id = ?`)
	require.NoError(err)
	testCases := []struct {
		Name     string
		RecipeID string
		Hops     []*summary.HopInfo
		Error    bool
	}{
		{
			Name:     "Valid Input",
			RecipeID: "1",
			Hops: []*summary.HopInfo{
				{
					Name:     "Hop1",
					Grams:    10,
					Alpha:    5,
					Time:     60,
					TimeUnit: "minutes",
					Notes:    "Some notes",
				},
			},
			Error: false,
		},
		{
			Name:     "Multiple hops",
			RecipeID: "2",
			Hops: []*summary.HopInfo{
				{Name: "hop1", Grams: 10, Alpha: 3.2, Time: 50, TimeUnit: "minutes", Notes: "notes 1"},
				{Name: "hop2", Grams: 20, Alpha: 5.2, Time: 70, TimeUnit: "minutes", Notes: "notes 2"},
			},
			Error: false,
		},
		{
			Name:     "Empty RecipeID",
			RecipeID: "",
			Hops: []*summary.HopInfo{
				{
					Name:     "Hop1",
					Grams:    10,
					Alpha:    5,
					Time:     60,
					TimeUnit: "minutes",
					Notes:    "Some notes",
				},
			},
			Error: true,
		},
		{
			Name:     "SQL Injection in RecipeID",
			RecipeID: "123; DROP TABLE summaries;",
			Hops: []*summary.HopInfo{
				{
					Name:     "Hop1",
					Grams:    10,
					Alpha:    5,
					Time:     60,
					TimeUnit: "minutes",
					Notes:    "Some notes",
				},
			},
			Error: true,
		},
		{
			Name:     "Non-Existing RecipeID",
			RecipeID: "999",
			Hops: []*summary.HopInfo{
				{
					Name:     "Hop1",
					Grams:    10,
					Alpha:    5,
					Time:     60,
					TimeUnit: "minutes",
					Notes:    "Some notes",
				},
			},
			Error: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			for _, hop := range tc.Hops {
				err = store.AddHopping(tc.RecipeID, hop.Name, hop.Grams, hop.Alpha, hop.Time, hop.Notes)
				if tc.Error {
					require.Error(err)
				}
			}
			if tc.Error {
				require.Error(err)
			} else {
				require.NoError(err)
				var hops string
				require.NoError(getSt.QueryRow(tc.RecipeID).Scan(&hops))
				expected, err := json.Marshal(tc.Hops)
				require.NoError(err)
				require.Equal(string(expected), hops)
			}
		})
	}
}

func TestAddVolumeBeforeBoil(t *testing.T) {
	require := require.New(t)
	fileName := strings.ToLower(strings.TrimSpace(t.Name())) + ".sqlite"
	db, err := sql.Open("sqlite3", "file:"+fileName+"?_foreign_keys=true")
	require.NoError(err)
	provisionDB(t, db, []string{"recipe1", "recipe2", "recipe3", "recipe4"})
	store, err := NewSummaryPersistentStore(db)
	require.NoError(err)
	defer os.Remove(fileName)
	for i := 1; i <= 3; i++ {
		num := strconv.Itoa(i)
		require.NoError(store.AddSummary(num, "t"+num))
	}

	getSt, err := db.Prepare(`SELECT hopping_vol_bb, hopping_vol_bb_notes FROM summaries WHERE recipe_id = ?`)
	require.NoError(err)
	testCases := []struct {
		Name     string
		RecipeID string
		Volume   float32
		Notes    string
		SkipRead bool
		Error    bool
	}{
		{
			Name:     "Valid Inputs",
			RecipeID: "1",
			Volume:   10.0,
			Notes:    "Some notes",
			Error:    false,
		},
		{
			Name:     "Empty RecipeID",
			RecipeID: "",
			Volume:   10.0,
			Notes:    "Some notes",
			Error:    true,
		},
		{
			Name:     "SQL Injection in RecipeID",
			RecipeID: "3; DROP TABLE summaries;",
			Volume:   10.0,
			Notes:    "Some notes",
			Error:    false,
			SkipRead: true,
		},
		{
			Name:     "SQL Injection in Notes",
			RecipeID: "3",
			Volume:   10.0,
			Notes:    "Some notes; DROP TABLE summaries;",
			Error:    false,
		},
		{
			Name:     "Non-Existing RecipeID",
			RecipeID: "999",
			Volume:   10.0,
			Notes:    "Some notes",
			Error:    false,
			SkipRead: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			err = store.AddVolumeBeforeBoil(tc.RecipeID, tc.Volume, tc.Notes)
			if tc.Error {
				require.Error(err)
			} else {
				require.NoError(err)
				if !tc.SkipRead {
					var vol float32
					var notes string
					require.NoError(getSt.QueryRow(tc.RecipeID).Scan(&vol, &notes))
					require.Equal(tc.Notes, notes)
					require.InDelta(tc.Volume, vol, 0.001)
				}
			}
		})
	}
}

func TestAddVolumeAfterBoil(t *testing.T) {
	require := require.New(t)
	fileName := strings.ToLower(strings.TrimSpace(t.Name())) + ".sqlite"
	db, err := sql.Open("sqlite3", "file:"+fileName+"?_foreign_keys=true")
	require.NoError(err)
	provisionDB(t, db, []string{"recipe1", "recipe2", "recipe3", "recipe4"})
	store, err := NewSummaryPersistentStore(db)
	require.NoError(err)
	defer os.Remove(fileName)
	for i := 1; i <= 3; i++ {
		num := strconv.Itoa(i)
		require.NoError(store.AddSummary(num, "t"+num))
	}

	getSt, err := db.Prepare(`SELECT hopping_vol_ab, hopping_vol_ab_notes FROM summaries WHERE recipe_id = ?`)
	require.NoError(err)
	testCases := []struct {
		Name     string
		RecipeID string
		Volume   float32
		Notes    string
		SkipRead bool
		Error    bool
	}{
		{
			Name:     "Valid Inputs",
			RecipeID: "1",
			Volume:   10.0,
			Notes:    "Some notes",
			Error:    false,
		},
		{
			Name:     "Empty RecipeID",
			RecipeID: "",
			Volume:   10.0,
			Notes:    "Some notes",
			Error:    true,
		},
		{
			Name:     "SQL Injection in RecipeID",
			RecipeID: "3; DROP TABLE summaries;",
			Volume:   10.0,
			Notes:    "Some notes",
			Error:    false,
			SkipRead: true,
		},
		{
			Name:     "SQL Injection in Notes",
			RecipeID: "3",
			Volume:   10.0,
			Notes:    "Some notes; DROP TABLE summaries;",
			Error:    false,
		},
		{
			Name:     "Non-Existing RecipeID",
			RecipeID: "999",
			Volume:   10.0,
			Notes:    "Some notes",
			Error:    false,
			SkipRead: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			err = store.AddVolumeAfterBoil(tc.RecipeID, tc.Volume, tc.Notes)
			if tc.Error {
				require.Error(err)
			} else {
				require.NoError(err)
				if !tc.SkipRead {
					var vol float32
					var notes string
					require.NoError(getSt.QueryRow(tc.RecipeID).Scan(&vol, &notes))
					require.Equal(tc.Notes, notes)
					require.InDelta(tc.Volume, vol, 0.001)
				}
			}
		})
	}
}

func TestAddCooling(t *testing.T) {
	require := require.New(t)
	fileName := strings.ToLower(strings.TrimSpace(t.Name())) + ".sqlite"
	db, err := sql.Open("sqlite3", "file:"+fileName+"?_foreign_keys=true")
	require.NoError(err)
	provisionDB(t, db, []string{"recipe1", "recipe2", "recipe3", "recipe4"})
	store, err := NewSummaryPersistentStore(db)
	require.NoError(err)
	defer os.Remove(fileName)
	for i := 1; i <= 3; i++ {
		num := strconv.Itoa(i)
		require.NoError(store.AddSummary(num, "t"+num))
	}
	getSt, err := db.Prepare(`SELECT cooling_temp, cooling_time, cooling_notes FROM summaries WHERE recipe_id = ?`)
	require.NoError(err)
	testCases := []struct {
		Name     string
		RecipeID string
		Temp     float32
		Time     float32
		Notes    string
		SkipRead bool
		Error    bool
	}{
		{
			Name:     "Valid Inputs",
			RecipeID: "1",
			Temp:     20.0,
			Time:     10.0,
			Notes:    "Some notes",
			Error:    false,
		}, {
			Name:     "Empty RecipeID",
			RecipeID: "",
			Temp:     20.0,
			Time:     10.0,
			Notes:    "Some notes",
			Error:    true,
		},
		{
			Name:     "SQL Injection in RecipeID",
			RecipeID: "123; DROP TABLE summaries;",
			Temp:     20.0,
			Time:     10.0,
			Notes:    "Some notes",
			Error:    false,
			SkipRead: true,
		},
		{
			Name:     "SQL Injection in Notes",
			RecipeID: "2",
			Temp:     20.0,
			Time:     10.0,
			Notes:    "Some notes; DROP TABLE summaries;",
			Error:    false,
		},
		{
			Name:     "Non-Existing RecipeID",
			RecipeID: "999",
			Temp:     20.0,
			Time:     10.0,
			Notes:    "Some notes",
			Error:    false,
			SkipRead: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			err = store.AddCooling(tc.RecipeID, tc.Temp, tc.Time, tc.Notes)
			if tc.Error {
				require.Error(err)
			} else {
				require.NoError(err)
				if !tc.SkipRead {
					var temp, time float32
					var notes string
					require.NoError(getSt.QueryRow(tc.RecipeID).Scan(&temp, &time, &notes))
					require.Equal(tc.Notes, notes)
					require.InDelta(tc.Temp, temp, 0.001)
					require.InDelta(tc.Time, time, 0.001)
				}
			}
		})
	}
}

func TestAddPreFermentationVolume(t *testing.T) {
	require := require.New(t)
	fileName := strings.ToLower(strings.TrimSpace(t.Name())) + ".sqlite"
	db, err := sql.Open("sqlite3", "file:"+fileName+"?_foreign_keys=true")
	require.NoError(err)
	provisionDB(t, db, []string{"recipe1", "recipe2", "recipe3", "recipe4"})
	store, err := NewSummaryPersistentStore(db)
	require.NoError(err)
	defer os.Remove(fileName)
	for i := 1; i <= 3; i++ {
		num := strconv.Itoa(i)
		require.NoError(store.AddSummary(num, "t"+num))
	}
	getSt, err := db.Prepare(`SELECT pre_ferm_vols FROM summaries WHERE recipe_id = ?`)
	require.NoError(err)
	testCases := []struct {
		Name     string
		RecipeID string
		Vols     []*summary.PreFermentationInfo
		Error    bool
	}{
		{
			Name:     "Normal case",
			RecipeID: "1",
			Vols: []*summary.PreFermentationInfo{
				{Volume: 10, SG: 1.054, Notes: "notes 1"},
			},
			Error: false,
		},
		{
			Name:     "Multiple vols",
			RecipeID: "2",
			Vols: []*summary.PreFermentationInfo{
				{Volume: 10, SG: 1.054, Notes: "notes 1"},
				{Volume: 12, SG: 1.067, Notes: "notes 2"},
			},
			Error: false,
		},
		{
			Name:     "SQL Injection in notes",
			RecipeID: "3",
			Vols: []*summary.PreFermentationInfo{
				{Volume: 10, SG: 1.054, Notes: "3; DROP TABLE summaries;"},
			},
			Error: false,
		},
		{
			Name:     "SQL Injection in recipe_id",
			RecipeID: "2; DROP TABLE summaries;",
			Vols: []*summary.PreFermentationInfo{
				{Volume: 10, SG: 1.054, Notes: "notes 1"},
			},
			Error: true,
		},
		{
			Name:     "Non existing recipe_id",
			RecipeID: "15",
			Vols: []*summary.PreFermentationInfo{
				{Volume: 10, SG: 1.054, Notes: "notes 1"},
			},
			Error: true,
		},
		{
			Name:     "Empty recipe_id",
			RecipeID: "",
			Vols: []*summary.PreFermentationInfo{
				{Volume: 10, SG: 1.054, Notes: "notes 1"},
			},
			Error: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			for _, vol := range tc.Vols {
				err = store.AddPreFermentationVolume(tc.RecipeID, vol.Volume, vol.SG, vol.Notes)
				if tc.Error {
					require.Error(err)
				}
			}
			if tc.Error {
				require.Error(err)
			} else {
				require.NoError(err)
				var vols string
				require.NoError(getSt.QueryRow(tc.RecipeID).Scan(&vols))
				expected, err := json.Marshal(tc.Vols)
				require.NoError(err)
				require.Equal(string(expected), vols)
			}
		})
	}
}

func TestAddYeastStart(t *testing.T) {
	require := require.New(t)
	fileName := strings.ToLower(strings.TrimSpace(t.Name())) + ".sqlite"
	db, err := sql.Open("sqlite3", "file:"+fileName+"?_foreign_keys=true")
	require.NoError(err)
	provisionDB(t, db, []string{"recipe1", "recipe2", "recipe3", "recipe4"})
	store, err := NewSummaryPersistentStore(db)
	require.NoError(err)
	defer os.Remove(fileName)
	for i := 1; i <= 3; i++ {
		num := strconv.Itoa(i)
		require.NoError(store.AddSummary(num, "t"+num))
	}
	getSt, err := db.Prepare(`SELECT yeast_start_temp, yeast_start_notes FROM summaries WHERE recipe_id = ?`)
	require.NoError(err)
	testCases := []struct {
		Name     string
		RecipeID string
		Temp     string
		Notes    string
		SkipRead bool
		Error    bool
	}{
		{
			Name:     "Valid Inputs",
			RecipeID: "1",
			Temp:     "20.0",
			Notes:    "Some notes",
			Error:    false,
		}, {
			Name:     "Empty RecipeID",
			RecipeID: "",
			Temp:     "20",
			Notes:    "Some notes",
			Error:    true,
		},
		{
			Name:     "SQL Injection in RecipeID",
			RecipeID: "3; DROP TABLE summaries;",
			Temp:     "20",
			Notes:    "Some notes",
			Error:    false,
			SkipRead: true,
		},
		{
			Name:     "SQL Injection in Notes",
			RecipeID: "2",
			Temp:     "20",
			Notes:    "Some notes; DROP TABLE summaries;",
			Error:    false,
		},
		{
			Name:     "SQL Injection in temp",
			RecipeID: "3",
			Temp:     "3; DROP TABLE summaries;",
			Notes:    "oe",
			Error:    false,
		},
		{
			Name:     "Non-Existing RecipeID",
			RecipeID: "999",
			Temp:     "20",
			Notes:    "Some notes",
			Error:    false,
			SkipRead: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			err = store.AddYeastStart(tc.RecipeID, tc.Temp, tc.Notes)
			if tc.Error {
				require.Error(err)
			} else {
				require.NoError(err)
				if !tc.SkipRead {
					var temp, notes string
					require.NoError(getSt.QueryRow(tc.RecipeID).Scan(&temp, &notes))
					require.Equal(tc.Notes, notes)
					require.Equal(tc.Temp, temp)
				}
			}
		})
	}
}

func TestAddMainFermentationAlcohol(t *testing.T) {
	require := require.New(t)
	fileName := strings.ToLower(strings.TrimSpace(t.Name())) + ".sqlite"
	db, err := sql.Open("sqlite3", "file:"+fileName+"?_foreign_keys=true")
	require.NoError(err)
	provisionDB(t, db, []string{"recipe1", "recipe2", "recipe3", "recipe4"})
	store, err := NewSummaryPersistentStore(db)
	require.NoError(err)
	defer os.Remove(fileName)
	for i := 1; i <= 3; i++ {
		num := strconv.Itoa(i)
		require.NoError(store.AddSummary(num, "t"+num))
	}
	getSt, err := db.Prepare(`SELECT main_ferm_alcohol FROM summaries WHERE recipe_id = ?`)
	require.NoError(err)
	testCases := []struct {
		Name     string
		RecipeID string
		Alcohol  float32
		SkipRead bool
		Error    bool
	}{
		{
			Name:     "Valid Inputs",
			RecipeID: "1",
			Alcohol:  5.5,
			Error:    false,
		}, {
			Name:     "Empty RecipeID",
			RecipeID: "",
			Alcohol:  5.5,
			Error:    true,
		},
		{
			Name:     "SQL Injection in RecipeID",
			RecipeID: "123; DROP TABLE summaries;",
			Alcohol:  5.5,
			Error:    false,
			SkipRead: true,
		},
		{
			Name:     "Non-Existing RecipeID",
			RecipeID: "999",
			Alcohol:  5.5,
			Error:    false,
			SkipRead: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			err = store.AddMainFermentationAlcohol(tc.RecipeID, tc.Alcohol)
			if tc.Error {
				require.Error(err)
			} else {
				require.NoError(err)
				if !tc.SkipRead {
					var alcohol float32
					require.NoError(getSt.QueryRow(tc.RecipeID).Scan(&alcohol))
					require.InDelta(tc.Alcohol, alcohol, 0.001)
				}
			}
		})
	}
}

func TestAddPreBottlingVolume(t *testing.T) {
	require := require.New(t)
	fileName := strings.ToLower(strings.TrimSpace(t.Name())) + ".sqlite"
	db, err := sql.Open("sqlite3", "file:"+fileName+"?_foreign_keys=true")
	require.NoError(err)
	provisionDB(t, db, []string{"recipe1", "recipe2", "recipe3", "recipe4"})
	store, err := NewSummaryPersistentStore(db)
	require.NoError(err)
	defer os.Remove(fileName)
	for i := 1; i <= 3; i++ {
		num := strconv.Itoa(i)
		require.NoError(store.AddSummary(num, "t"+num))
	}
	getSt, err := db.Prepare(`SELECT bottling_pre_bottle_volume FROM summaries WHERE recipe_id = ?`)
	require.NoError(err)
	testCases := []struct {
		Name     string
		RecipeID string
		Volume   float32
		SkipRead bool
		Error    bool
	}{
		{
			Name:     "Valid Inputs",
			RecipeID: "1",
			Volume:   10.3,
			Error:    false,
		}, {
			Name:     "Empty RecipeID",
			RecipeID: "",
			Volume:   10.3,
			Error:    true,
		},
		{
			Name:     "SQL Injection in RecipeID",
			RecipeID: "123; DROP TABLE summaries;",
			Volume:   10.3,
			Error:    false,
			SkipRead: true,
		},
		{
			Name:     "Non-Existing RecipeID",
			RecipeID: "999",
			Volume:   10.3,
			Error:    false,
			SkipRead: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			err = store.AddPreBottlingVolume(tc.RecipeID, tc.Volume)
			if tc.Error {
				require.Error(err)
			} else {
				require.NoError(err)
				if !tc.SkipRead {
					var vol float32
					require.NoError(getSt.QueryRow(tc.RecipeID).Scan(&vol))
					require.InDelta(tc.Volume, vol, 0.001)
				}
			}
		})
	}
}

func TestAddBottling(t *testing.T) {
	require := require.New(t)
	fileName := strings.ToLower(strings.TrimSpace(t.Name())) + ".sqlite"
	db, err := sql.Open("sqlite3", "file:"+fileName+"?_foreign_keys=true")
	require.NoError(err)
	provisionDB(t, db, []string{"recipe1", "recipe2", "recipe3", "recipe4"})
	store, err := NewSummaryPersistentStore(db)
	require.NoError(err)
	defer os.Remove(fileName)
	for i := 1; i <= 3; i++ {
		num := strconv.Itoa(i)
		require.NoError(store.AddSummary(num, "t"+num))
	}
	getSt, err := db.Prepare(`SELECT bottling_carbonation  ,
	bottling_sugar_amount  ,
	bottling_sugar_type  ,
	bottling_temperature  ,
	bottling_alcohol  ,
	bottling_volume_bottled  ,
	bottling_notes FROM summaries WHERE recipe_id == ?`)
	require.NoError(err)
	testCases := []struct {
		Name        string
		RecipeID    string
		Carbonation float32
		Sugar       float32
		SugarType   string
		Temp        float32
		Alcohol     float32
		Volume      float32
		Notes       string
		SkipRead    bool
		Error       bool
	}{
		{
			Name:        "Valid Inputs",
			RecipeID:    "1",
			Temp:        20.0,
			Carbonation: 5.5,
			Sugar:       100,
			SugarType:   "glucose",
			Alcohol:     5.69,
			Volume:      10.3,
			Notes:       "Some notes",
			Error:       false,
		}, {
			Name:        "Empty RecipeID",
			RecipeID:    "",
			Temp:        20.0,
			Carbonation: 5.5,
			Sugar:       100,
			SugarType:   "glucose",
			Alcohol:     5.69,
			Volume:      10.3,
			Notes:       "Some notes",
			Error:       true,
		},
		{
			Name:        "SQL Injection in RecipeID",
			RecipeID:    "2; DROP TABLE summaries;",
			Temp:        20.0,
			Carbonation: 5.5,
			Sugar:       100,
			SugarType:   "glucose",
			Alcohol:     5.69,
			Volume:      10.3,
			Notes:       "Some notes",
			Error:       false,
			SkipRead:    true,
		},
		{
			Name:        "SQL Injection in Notes",
			RecipeID:    "2",
			Temp:        20.0,
			Carbonation: 5.5,
			Sugar:       100,
			SugarType:   "glucose",
			Alcohol:     5.69,
			Volume:      10.3,
			Notes:       "Some notes; DROP TABLE summaries;",
			Error:       false,
		},
		{
			Name:        "SQL Injection in sugar type",
			SugarType:   "glucose; DROP TABLE summaries;",
			RecipeID:    "3",
			Temp:        20.0,
			Carbonation: 5.5,
			Sugar:       100,
			Alcohol:     5.69,
			Volume:      10.3,
			Notes:       "Some notes",
			Error:       false,
		},
		{
			Name:        "Non-Existing RecipeID",
			RecipeID:    "999",
			Temp:        20.0,
			Carbonation: 5.5,
			Sugar:       100,
			SugarType:   "glucose",
			Alcohol:     5.69,
			Volume:      10.3,
			Notes:       "Some notes",
			Error:       false,
			SkipRead:    true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			err = store.AddBottling(tc.RecipeID, tc.Carbonation, tc.Alcohol, tc.Sugar, tc.Temp, tc.Volume, tc.SugarType, tc.Notes)
			if tc.Error {
				require.Error(err)
			} else {
				require.NoError(err)
				if !tc.SkipRead {
					var temp, carbonation, alcohol, vol, sugar float32
					var st, notes string
					require.NoError(getSt.QueryRow(tc.RecipeID).Scan(&carbonation, &sugar, &st, &temp, &alcohol, &vol, &notes))
					require.Equal(tc.Notes, notes)
					require.Equal(tc.SugarType, st)
					require.InDelta(tc.Carbonation, carbonation, 0.001)
					require.InDelta(tc.Temp, temp, 0.001)
					require.InDelta(tc.Sugar, sugar, 0.001)
					require.InDelta(tc.Alcohol, alcohol, 0.001)
					require.InDelta(tc.Volume, vol, 0.001)
				}
			}
		})
	}
}

func TestAddSummarySecondary(t *testing.T) {
	require := require.New(t)
	fileName := strings.ToLower(strings.TrimSpace(t.Name())) + ".sqlite"
	db, err := sql.Open("sqlite3", "file:"+fileName+"?_foreign_keys=true")
	require.NoError(err)
	provisionDB(t, db, []string{"recipe1", "recipe2", "recipe3", "recipe4"})
	store, err := NewSummaryPersistentStore(db)
	require.NoError(err)
	defer os.Remove(fileName)
	for i := 1; i <= 3; i++ {
		num := strconv.Itoa(i)
		require.NoError(store.AddSummary(num, "t"+num))
	}
	getSt, err := db.Prepare(`SELECT sec_ferm_days, sec_ferm_notes FROM summaries WHERE recipe_id = ?`)
	require.NoError(err)
	testCases := []struct {
		Name     string
		RecipeID string
		Days     int
		Notes    string
		SkipRead bool
		Error    bool
	}{
		{
			Name:     "Valid Inputs",
			RecipeID: "1",
			Days:     5,
			Notes:    "Some notes",
			Error:    false,
		}, {
			Name:     "Empty RecipeID",
			RecipeID: "",
			Days:     5,
			Notes:    "Some notes",
			Error:    true,
		},
		{
			Name:     "SQL Injection in RecipeID",
			RecipeID: "123; DROP TABLE summaries;",
			Days:     5,
			Notes:    "Some notes",
			Error:    false,
			SkipRead: true,
		},
		{
			Name:     "SQL Injection in Notes",
			RecipeID: "2",
			Days:     5,
			Notes:    "Some notes; DROP TABLE summaries;",
			Error:    false,
		},
		{
			Name:     "Non-Existing RecipeID",
			RecipeID: "999",
			Days:     5,
			Notes:    "Some notes",
			Error:    false,
			SkipRead: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			err = store.AddSummarySecondary(tc.RecipeID, tc.Days, tc.Notes)
			if tc.Error {
				require.Error(err)
			} else {
				require.NoError(err)
				if !tc.SkipRead {
					var days int
					var notes string
					require.NoError(getSt.QueryRow(tc.RecipeID).Scan(&days, &notes))
					require.Equal(tc.Notes, notes)
					require.Equal(tc.Days, days)
				}
			}
		})
	}
}

func TestAddEvaporation(t *testing.T) {
	require := require.New(t)
	fileName := strings.ToLower(strings.TrimSpace(t.Name())) + ".sqlite"
	db, err := sql.Open("sqlite3", "file:"+fileName+"?_foreign_keys=true")
	require.NoError(err)
	provisionDB(t, db, []string{"recipe1", "recipe2", "recipe3", "recipe4"})
	store, err := NewSummaryPersistentStore(db)
	require.NoError(err)
	defer os.Remove(fileName)
	for i := 1; i <= 3; i++ {
		num := strconv.Itoa(i)
		require.NoError(store.AddSummary(num, "t"+num))
	}
	getSt, err := db.Prepare(`SELECT stats_evaporation FROM summaries WHERE recipe_id = ?`)
	require.NoError(err)
	testCases := []struct {
		Name        string
		RecipeID    string
		Evaporation float32
		SkipRead    bool
		Error       bool
	}{
		{
			Name:        "Valid Inputs",
			RecipeID:    "1",
			Evaporation: 16.66,
			Error:       false,
		}, {
			Name:        "Empty RecipeID",
			RecipeID:    "",
			Evaporation: 16.66,
			Error:       true,
		},
		{
			Name:        "SQL Injection in RecipeID",
			RecipeID:    "123; DROP TABLE summaries;",
			Evaporation: 16.66,
			Error:       false,
			SkipRead:    true,
		},
		{
			Name:        "Non-Existing RecipeID",
			RecipeID:    "999",
			Evaporation: 16.66,
			Error:       false,
			SkipRead:    true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			err = store.AddEvaporation(tc.RecipeID, tc.Evaporation)
			if tc.Error {
				require.Error(err)
			} else {
				require.NoError(err)
				if !tc.SkipRead {
					var evap float32
					require.NoError(getSt.QueryRow(tc.RecipeID).Scan(&evap))
					require.InDelta(tc.Evaporation, evap, 0.001)
				}
			}
		})
	}
}

func TestAddEfficency(t *testing.T) {
	require := require.New(t)
	fileName := strings.ToLower(strings.TrimSpace(t.Name())) + ".sqlite"
	db, err := sql.Open("sqlite3", "file:"+fileName+"?_foreign_keys=true")
	require.NoError(err)
	provisionDB(t, db, []string{"recipe1", "recipe2", "recipe3", "recipe4"})
	store, err := NewSummaryPersistentStore(db)
	require.NoError(err)
	defer os.Remove(fileName)
	for i := 1; i <= 3; i++ {
		num := strconv.Itoa(i)
		require.NoError(store.AddSummary(num, "t"+num))
	}
	getSt, err := db.Prepare(`SELECT stats_effiency FROM summaries WHERE recipe_id = ?`)
	require.NoError(err)
	testCases := []struct {
		Name      string
		RecipeID  string
		Efficency float32
		SkipRead  bool
		Error     bool
	}{
		{
			Name:      "Valid Inputs",
			RecipeID:  "1",
			Efficency: 16.66,
			Error:     false,
		}, {
			Name:      "Empty RecipeID",
			RecipeID:  "",
			Efficency: 16.66,
			Error:     true,
		},
		{
			Name:      "SQL Injection in RecipeID",
			RecipeID:  "123; DROP TABLE summaries;",
			Efficency: 16.66,
			Error:     false,
			SkipRead:  true,
		},
		{
			Name:      "Non-Existing RecipeID",
			RecipeID:  "999",
			Efficency: 16.66,
			Error:     false,
			SkipRead:  true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			err = store.AddEfficiency(tc.RecipeID, tc.Efficency)
			if tc.Error {
				require.Error(err)
			} else {
				require.NoError(err)
				if !tc.SkipRead {
					var eff float32
					require.NoError(getSt.QueryRow(tc.RecipeID).Scan(&eff))
					require.InDelta(tc.Efficency, eff, 0.001)
				}
			}
		})
	}
}

func TestAddMainFermentationSGMeasurement(t *testing.T) {
	require := require.New(t)
	fileName := strings.ToLower(strings.TrimSpace(t.Name())) + ".sqlite"
	db, err := sql.Open("sqlite3", "file:"+fileName+"?_foreign_keys=true")
	require.NoError(err)
	provisionDB(t, db, []string{"recipe1", "recipe2", "recipe3", "recipe4"})
	store, err := NewSummaryPersistentStore(db)
	require.NoError(err)
	defer os.Remove(fileName)
	for i := 1; i <= 3; i++ {
		num := strconv.Itoa(i)
		require.NoError(store.AddSummary(num, "t"+num))
	}

	getSt, err := db.Prepare(`SELECT main_ferm_sgs FROM summaries WHERE recipe_id = ?`)
	require.NoError(err)
	testCases := []struct {
		Name     string
		SGs      []*summary.SGMeasurement
		RecipeID string
		SkipRead bool
		Error    bool
	}{
		{
			Name: "Normal case",
			SGs: []*summary.SGMeasurement{
				{SG: 1.063, Date: "2013-02-21 13:45:43", Final: false, Notes: "some notes"},
			},
			RecipeID: "1",
			Error:    false,
		},
		{
			Name: "Multiple sgs",
			SGs: []*summary.SGMeasurement{
				{SG: 1.063, Date: "2013-02-21 13:45:43", Final: false, Notes: "some notes"},
				{SG: 1.010, Date: "2013-02-23 13:45:43", Final: true, Notes: "some notes2"},
			},
			RecipeID: "2",
			Error:    false,
		},
		{
			Name: "SQL Injection in notes",
			SGs: []*summary.SGMeasurement{
				{SG: 1.063, Date: "2013-02-21 13:45:43", Final: false, Notes: "5'; DROP TABLE summaries; --"},
			},
			RecipeID: "3",
			Error:    false,
		},
		{
			Name: "SQL Injection in recipe_id",
			SGs: []*summary.SGMeasurement{
				{SG: 1.063, Date: "2013-02-21 13:45:43", Final: false, Notes: "some notes"},
			},
			RecipeID: "5'; DROP TABLE summaries; --",
			Error:    true,
		},
		{
			Name: "Non existing recipe_id",
			SGs: []*summary.SGMeasurement{
				{SG: 1.063, Date: "2013-02-21 13:45:43", Final: false, Notes: "some notes"},
			},
			RecipeID: "10",
			Error:    true,
		},
		{
			Name: "Empty recipe_id",
			SGs: []*summary.SGMeasurement{
				{SG: 1.063, Date: "2013-02-21 13:45:43", Final: false, Notes: "some notes"},
			},
			RecipeID: "",
			Error:    true,
		},
		{
			Name: "Summary not created",
			SGs: []*summary.SGMeasurement{
				{SG: 1.063, Date: "2013-02-21 13:45:43", Final: false, Notes: "some notes"},
			},
			RecipeID: "4",
			Error:    true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			for _, sg := range tc.SGs {
				err = store.AddMainFermentationSGMeasurement(tc.RecipeID, sg.Date, sg.SG, sg.Final, sg.Notes)
				if tc.Error {
					require.Error(err)
				}
			}
			if tc.Error {
				require.Error(err)
			} else {
				require.NoError(err)
				if !tc.SkipRead {
					var sgs string
					err = getSt.QueryRow(tc.RecipeID).Scan(&sgs)
					require.NoError(err)
					expected, err := json.Marshal(tc.SGs)
					require.NoError(err)
					require.Equal(string(expected), sgs)
				}
			}
		})
	}
}

func TestAddDryHopStart(t *testing.T) {
	require := require.New(t)
	fileName := strings.ToLower(strings.TrimSpace(t.Name())) + ".sqlite"
	db, err := sql.Open("sqlite3", "file:"+fileName+"?_foreign_keys=true")
	require.NoError(err)
	provisionDB(t, db, []string{"recipe1", "recipe2", "recipe3", "recipe4"})
	store, err := NewSummaryPersistentStore(db)
	require.NoError(err)
	defer os.Remove(fileName)
	for i := 1; i <= 3; i++ {
		num := strconv.Itoa(i)
		require.NoError(store.AddSummary(num, "t"+num))
	}
	getSt, err := db.Prepare(`SELECT main_ferm_dry_hops FROM summaries WHERE recipe_id = ?`)
	require.NoError(err)
	testCases := []struct {
		Name     string
		RecipeID string
		Hops     []*summary.HopInfo
		Error    bool
	}{
		{
			Name:     "Valid Input",
			RecipeID: "1",
			Hops: []*summary.HopInfo{
				{
					Name:  "Hop1",
					Grams: 10,
					Alpha: 5,
					Notes: "Some notes",
				},
			},
			Error: false,
		},
		{
			Name:     "Multiple hops",
			RecipeID: "2",
			Hops: []*summary.HopInfo{
				{Name: "hop1", Grams: 10, Alpha: 3.2, Time: 50, TimeUnit: "days", Notes: "notes 1"},
				{Name: "hop2", Grams: 20, Alpha: 5.2, Time: 70, TimeUnit: "days", Notes: "notes 2"},
			},
			Error: false,
		},
		{
			Name:     "Empty RecipeID",
			RecipeID: "",
			Hops: []*summary.HopInfo{
				{
					Name:  "Hop1",
					Grams: 10,
					Alpha: 5,
					Notes: "Some notes",
				},
			},
			Error: true,
		},
		{
			Name:     "SQL Injection in RecipeID",
			RecipeID: "123; DROP TABLE summaries;",
			Hops: []*summary.HopInfo{
				{
					Name:  "Hop1",
					Grams: 10,
					Alpha: 5,
					Notes: "Some notes",
				},
			},
			Error: true,
		},
		{
			Name:     "Non-Existing RecipeID",
			RecipeID: "999",
			Hops: []*summary.HopInfo{
				{
					Name:  "Hop1",
					Grams: 10,
					Alpha: 5,
					Notes: "Some notes",
				},
			},
			Error: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			for _, hop := range tc.Hops {
				err = store.AddDryHopStart(tc.RecipeID, hop.Name, hop.Grams, hop.Alpha, hop.Notes)
				if tc.Error {
					require.Error(err)
				}
			}
			if tc.Error {
				require.Error(err)
			} else {
				require.NoError(err)
				var hops string
				require.NoError(getSt.QueryRow(tc.RecipeID).Scan(&hops))
				expected, err := json.Marshal(tc.Hops)
				require.NoError(err)
				require.Equal(string(expected), hops)
			}
		})
	}
}

func TestAddDryHopEnd(t *testing.T) {
	require := require.New(t)
	fileName := strings.ToLower(strings.TrimSpace(t.Name())) + ".sqlite"
	db, err := sql.Open("sqlite3", "file:"+fileName+"?_foreign_keys=true")
	require.NoError(err)
	provisionDB(t, db, []string{"recipe1", "recipe2", "recipe3", "recipe4"})
	store, err := NewSummaryPersistentStore(db)
	require.NoError(err)
	defer os.Remove(fileName)
	for i := 1; i <= 3; i++ {
		num := strconv.Itoa(i)
		require.NoError(store.AddSummary(num, "t"+num))
	}
	getSt, err := db.Prepare(`SELECT main_ferm_dry_hops FROM summaries WHERE recipe_id = ?`)
	require.NoError(err)
	testCases := []struct {
		Name            string
		RecipeID        string
		Hops            []*summary.HopInfo
		SkipDryHopStart bool
		Error           bool
	}{
		{
			Name:     "Normal case",
			RecipeID: "1",
			Hops: []*summary.HopInfo{
				{Name: "hop1", Grams: 10, Alpha: 3.2, Time: 50, TimeUnit: "hours", Notes: "notes 1"},
			},
			SkipDryHopStart: false,
			Error:           false,
		},
		{
			Name:     "Start was not called first",
			RecipeID: "2",
			Hops: []*summary.HopInfo{
				{Name: "hop1", Grams: 10, Alpha: 3.2, Time: 50, TimeUnit: "hours", Notes: "notes 1"},
			},
			SkipDryHopStart: true,
			Error:           true,
		},
		{
			Name:     "Multiple hops",
			RecipeID: "3",
			Hops: []*summary.HopInfo{
				{Name: "hop_1", Grams: 10, Alpha: 3.2, Time: 50, TimeUnit: "hours", Notes: "notes 1"},
				{Name: "hop 1", Grams: 10, Alpha: 3.2, Time: 60, TimeUnit: "hours", Notes: "notes 1"},
			},
			SkipDryHopStart: false,
			Error:           false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {

			for _, h := range tc.Hops {
				if !tc.SkipDryHopStart {
					err = store.AddDryHopStart(tc.RecipeID, h.Name, h.Grams, h.Alpha, h.Notes)
					require.NoError(err)
				}
				err = store.AddDryHopEnd(tc.RecipeID, h.Name, h.Time)
				if tc.Error {
					require.Error(err)

				} else {
					require.NoError(err)
				}
			}
			if !tc.Error {
				var hops string
				require.NoError(getSt.QueryRow(tc.RecipeID).Scan(&hops))
				expected, err := json.Marshal(tc.Hops)
				require.NoError(err)
				require.Equal(string(expected), hops)
			}

		})
	}
}

func TestGetSummary(t *testing.T) {
	require := require.New(t)
	fileName := strings.ToLower(strings.TrimSpace(t.Name())) + ".sqlite"
	db, err := sql.Open("sqlite3", "file:"+fileName+"?_foreign_keys=true")
	require.NoError(err)
	provisionDB(t, db, []string{"recipe1", "recipe2", "recipe3", "recipe4"})
	store, err := NewSummaryPersistentStore(db)
	require.NoError(err)
	defer os.Remove(fileName)
	testCases := []struct {
		Name     string
		RecipeID string
		Summ     *summary.Summary
		Error    bool
	}{
		{
			Name:     "Valid input",
			RecipeID: "1",
			Error:    false,
			Summ: &summary.Summary{
				Title: "summary 1",
				MashingInfo: &summary.MashingInfo{
					MashingTemperature: 57.5,
					MashingNotes:       "notes 1",
					RastInfos: []*summary.MashRastInfo{
						{Temperature: 63, Time: 30, Notes: "notes 2"},
						{Temperature: 72, Time: 45, Notes: "notes 3"},
					},
				},
				LauternInfo: "lautern",
				HoppingInfo: &summary.HoppingInfo{
					VolBeforeBoil: &summary.VolMeasurement{Volume: 12.5, Notes: "notes 4"},
					VolAfterBoil:  &summary.VolMeasurement{Volume: 9.4, Notes: "notes 5"},
					HopInfos: []*summary.HopInfo{
						{Name: "Amarillo", Grams: 20, Alpha: 6.8, Time: 60, TimeUnit: "minutes", Notes: "notes 6"},
						{Name: "Galaxy", Grams: 40, Alpha: 5.8, Time: 10, TimeUnit: "minutes", Notes: "notes 7"},
					},
				},
				CoolingInfo: &summary.CoolingInfo{Temperature: 21, Time: 58.789, Notes: "notes 8"},
				PreFermentationInfos: []*summary.PreFermentationInfo{
					{Volume: 7.5, SG: 1.098, Notes: "notes 9"},
					{Volume: 12, SG: 1.54, Notes: "notes 10"},
				},
				YeastInfo: &summary.YeastInfo{
					Temperature: "18-20", Notes: "notes 11",
				},
				MainFermentationInfo: &summary.MainFermentationInfo{
					SGs: []*summary.SGMeasurement{
						{SG: 1.013, Date: "2023-02-23 15:45:54", Final: false, Notes: "notes 12"},
						{SG: 1.011, Date: "2023-02-24 15:45:54", Final: true, Notes: "notes 13"},
					},
					Alcohol: 5.89,
					DryHopInfo: []*summary.HopInfo{
						{Name: "Amarillo", Grams: 20, Alpha: 6.8, Time: 4, TimeUnit: "days", Notes: "notes 14"},
						{Name: "Galaxy", Grams: 40, Alpha: 5.8, Time: 3, TimeUnit: "days", Notes: "notes 15"},
					},
				},
				BottlingInfo: &summary.BottlingInfo{
					PreBottleVolume: 12, Carbonation: 5.5, SugarAmount: 100, SugarType: "glucose",
					Temperature: 19, Alcohol: 6.5, VolumeBottled: 11, Notes: "notes 16",
				},
				SecondaryFermentationInfo: &summary.SecondaryFermentationInfo{
					Days: 5, Notes: "notes 17",
				},
				Statistics: &summary.Statistics{
					Evaporation: 16.66,
					Efficiency:  58.32,
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			require.NoError(storeSummary(tc.RecipeID, tc.Summ, store))
			summ, err := store.GetSummary(tc.RecipeID)
			if tc.Error {
				require.Error(err)
			} else {
				require.NoError(err)
				require.Equal(tc.Summ, summ)
			}
		})
	}
}

func storeSummary(id string, summ *summary.Summary, store *SummaryPersistentStore) error {
	err := store.AddSummary(id, summ.Title)
	if err != nil {
		return err
	}
	err = store.AddMashTemp(id, summ.MashingInfo.MashingTemperature, summ.MashingInfo.MashingNotes)
	if err != nil {
		return err
	}
	for _, rast := range summ.MashingInfo.RastInfos {
		err = store.AddRast(id, rast.Temperature, rast.Time, rast.Notes)
		if err != nil {
			return err
		}
	}
	err = store.AddLauternNotes(id, summ.LauternInfo)
	if err != nil {
		return err
	}
	for _, hop := range summ.HoppingInfo.HopInfos {
		err = store.AddHopping(id, hop.Name, hop.Grams, hop.Alpha, hop.Time, hop.Notes)
		if err != nil {
			return err
		}
	}
	err = store.AddVolumeBeforeBoil(id, summ.HoppingInfo.VolBeforeBoil.Volume, summ.HoppingInfo.VolBeforeBoil.Notes)
	if err != nil {
		return err
	}
	err = store.AddVolumeAfterBoil(id, summ.HoppingInfo.VolAfterBoil.Volume, summ.HoppingInfo.VolAfterBoil.Notes)
	if err != nil {
		return err
	}
	err = store.AddCooling(id, summ.CoolingInfo.Temperature, summ.CoolingInfo.Time, summ.CoolingInfo.Notes)
	if err != nil {
		return err
	}
	for _, preferm := range summ.PreFermentationInfos {
		err = store.AddPreFermentationVolume(id, preferm.Volume, preferm.SG, preferm.Notes)
		if err != nil {
			return err
		}
	}
	err = store.AddYeastStart(id, summ.YeastInfo.Temperature, summ.YeastInfo.Notes)
	if err != nil {
		return err
	}
	for _, sg := range summ.MainFermentationInfo.SGs {
		err = store.AddMainFermentationSGMeasurement(id, sg.Date, sg.SG, sg.Final, sg.Notes)
		if err != nil {
			return err
		}
	}
	err = store.AddMainFermentationAlcohol(id, summ.MainFermentationInfo.Alcohol)
	if err != nil {
		return err
	}
	for _, dh := range summ.MainFermentationInfo.DryHopInfo {
		err = store.AddDryHopStart(id, dh.Name, dh.Grams, dh.Alpha, dh.Notes)
		if err != nil {
			return err
		}
		err = store.AddDryHopEnd(id, dh.Name, dh.Time)
		if err != nil {
			return err
		}
	}
	err = store.AddPreBottlingVolume(id, summ.BottlingInfo.PreBottleVolume)
	if err != nil {
		return err
	}
	err = store.AddBottling(id, summ.BottlingInfo.Carbonation, summ.BottlingInfo.Alcohol, summ.BottlingInfo.SugarAmount,
		summ.BottlingInfo.Temperature, summ.BottlingInfo.VolumeBottled, summ.BottlingInfo.SugarType, summ.BottlingInfo.Notes)
	if err != nil {
		return err
	}
	err = store.AddSummarySecondary(id, summ.SecondaryFermentationInfo.Days, summ.SecondaryFermentationInfo.Notes)
	if err != nil {
		return err
	}
	err = store.AddEfficiency(id, summ.Statistics.Efficiency)
	if err != nil {
		return err
	}
	err = store.AddEvaporation(id, summ.Statistics.Evaporation)
	if err != nil {
		return err
	}
	return nil
}
