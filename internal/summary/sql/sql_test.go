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
	store, err := NewSummaryRecorderPersistentStore(db)
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

func Test(t *testing.T) {
	require := require.New(t)
	fileName := strings.ToLower(strings.TrimSpace(t.Name())) + ".sqlite"
	db, err := sql.Open("sqlite3", "file:"+fileName+"?_foreign_keys=true")
	require.NoError(err)
	provisionDB(t, db, []string{"recipe1", "recipe2", "recipe3", "recipe4"})
	store, err := NewSummaryRecorderPersistentStore(db)
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
	store, err := NewSummaryRecorderPersistentStore(db)
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
		Temp     float64
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
	store, err := NewSummaryRecorderPersistentStore(db)
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
	store, err := NewSummaryRecorderPersistentStore(db)
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
	store, err := NewSummaryRecorderPersistentStore(db)
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
	store, err := NewSummaryRecorderPersistentStore(db)
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
	store, err := NewSummaryRecorderPersistentStore(db)
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