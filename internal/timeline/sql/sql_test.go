package sql

import (
	"database/sql"
	"os"
	"strings"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
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

func TestAddEvent(t *testing.T) {
	require := require.New(t)
	fileName := strings.ToLower(strings.TrimSpace(t.Name())) + ".sqlite"
	db, err := sql.Open("sqlite3", "file:"+fileName+"?_foreign_keys=true")
	require.NoError(err)
	provisionDB(t, db, []string{"recipe1", "recipe2", "recipe3", "recipe4"})
	store, err := NewTimelinePersistentStore(db)
	require.NoError(err)
	defer os.Remove(fileName)
	getSt, err := db.Prepare(`SELECT timestamp_unix, event FROM timelines WHERE recipe_id = ?`)
	require.NoError(err)
	testCases := []struct {
		Name     string
		RecipeID string
		Event    string
		Error    bool
	}{
		{
			Name:     "Successful add",
			RecipeID: "1",
			Event:    "my-event",
			Error:    false,
		},
		{
			Name:     "Empty event",
			RecipeID: "2",
			Event:    "",
			Error:    true,
		},
		{
			Name:     "SQL Injection in RecipeID",
			RecipeID: "4'; DROP TABLE timelines; --",
			Event:    "event4",
			Error:    true,
		},
		{
			Name:     "SQL Injection in Event",
			RecipeID: "3",
			Event:    "5'; DROP TABLE timelines; --",
			Error:    false,
		},
		{
			Name:     "Special Characters in Event",
			RecipeID: "4",
			Event:    "event7$#%@^&*()",
			Error:    false,
		},
		{
			Name:     "Empty recipeID",
			RecipeID: "",
			Event:    "my-event",
			Error:    true,
		},
		{
			Name:     "Non-existing recipeID",
			RecipeID: "5",
			Event:    "no-event",
			Error:    true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			current := time.Now()
			err := store.AddEvent(tc.RecipeID, tc.Event)
			if tc.Error {
				require.Error(err)
			} else {
				require.NoError(err)
				var ts int64
				var event string
				err = getSt.QueryRow(tc.RecipeID).Scan(&ts, &event)
				require.NoError(err)
				require.Equal(tc.Event, event)
				require.WithinDuration(current, time.Unix(ts, 0), 1*time.Second)
			}
		})
	}
}

func TestGetTimeline(t *testing.T) {
	require := require.New(t)
	fileName := strings.ToLower(strings.TrimSpace(t.Name())) + ".sqlite"
	db, err := sql.Open("sqlite3", "file:"+fileName+"?_foreign_keys=true")
	require.NoError(err)
	provisionDB(t, db, []string{"recipe1", "recipe2", "recipe3", "recipe4"})
	store, err := NewTimelinePersistentStore(db)
	require.NoError(err)
	defer os.Remove(fileName)
	testCases := []struct {
		Name     string
		RecipeID string
		Expected []string
		Insert   bool
		Error    bool
	}{
		{
			Name:     "Single event",
			RecipeID: "1",
			Expected: []string{"single"},
			Insert:   true,
			Error:    false,
		},
		{
			Name:     "SQL Injection in recipe ID",
			RecipeID: "4'; DROP TABLE timelines; --",
			Error:    false,
		},
		{
			Name:     "Normal timeline",
			RecipeID: "2",
			Expected: []string{"first", "second", "third"},
			Insert:   true,
			Error:    false,
		},
		{
			Name:     "Non-existing recipeID",
			RecipeID: "5",
			Expected: []string{},
			Insert:   false,
			Error:    false,
		},
		{
			Name:     "Empty recipeID",
			RecipeID: "",
			Error:    true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			var times []int64
			if tc.Insert {
				for _, e := range tc.Expected {
					t := time.Now().Unix()
					times = append(times, t)
					_, err := store.insertStatement.Exec(e, t, tc.RecipeID)
					require.NoError(err)
				}
			}
			tl, err := store.GetTimeline(tc.RecipeID)
			if tc.Error {
				require.Error(err)
			} else {
				require.NoError(err)
				require.Equal(len(tc.Expected), len(tl))
				for i, t := range times {
					require.Equal(time.Unix(t, 0).Format(time.RFC3339Nano)+" "+tc.Expected[i], tl[i])
				}

			}
		})
	}

}

func TestAddTimeline(t *testing.T) {
	require := require.New(t)
	fileName := strings.ToLower(strings.TrimSpace(t.Name())) + ".sqlite"
	db, err := sql.Open("sqlite3", "file:"+fileName+"?_foreign_keys=true")
	require.NoError(err)
	provisionDB(t, db, []string{"recipe1", "recipe2", "recipe3", "recipe4"})
	store, err := NewTimelinePersistentStore(db)
	require.NoError(err)
	defer os.Remove(fileName)
	getSt, err := db.Prepare(`SELECT event FROM timelines WHERE recipe_id = ?`)
	require.NoError(err)
	testCases := []struct {
		Name     string
		RecipeID string
		Error    bool
	}{
		{
			Name:     "SQL Injection in recipe ID",
			RecipeID: "4'; DROP TABLE timelines; --",
			Error:    true,
		},
		{
			Name:     "Normal",
			RecipeID: "1",
			Error:    false,
		},
		{
			Name:     "Non-existing recipe ID",
			RecipeID: "5",
			Error:    true,
		},
		{
			Name:     "Empty recipe ID",
			RecipeID: "",
			Error:    true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			err := store.AddTimeline(tc.RecipeID)
			if tc.Error {
				require.Error(err)
			} else {
				require.NoError(err)
				var event string
				err = getSt.QueryRow(tc.RecipeID).Scan(&event)
				require.NoError(err)
				require.Equal("Initialized Recipe", event)
			}
		})
	}

}

func TestDeleteTimeline(t *testing.T) {
	require := require.New(t)
	fileName := strings.ToLower(strings.TrimSpace(t.Name())) + ".sqlite"
	db, err := sql.Open("sqlite3", "file:"+fileName+"?_foreign_keys=true")
	require.NoError(err)
	provisionDB(t, db, []string{"recipe1", "recipe2", "recipe3", "recipe4"})
	store, err := NewTimelinePersistentStore(db)
	require.NoError(err)
	defer os.Remove(fileName)
	getSt, err := db.Prepare(`SELECT event FROM timelines WHERE recipe_id = ?`)
	require.NoError(err)
	testCases := []struct {
		Name     string
		RecipeID string
		ToAdd    []string
		Error    bool
	}{
		{
			Name:     "Single event timeline",
			RecipeID: "1",
			ToAdd:    []string{"event"},
			Error:    false,
		},
		{
			Name:     "SQL Injection in recipe ID",
			RecipeID: "4'; DROP TABLE timelines; --",
			Error:    false,
		},
		{
			Name:     "Multiple event timeline",
			RecipeID: "2",
			ToAdd:    []string{"event1", "event2", "event3"},
			Error:    false,
		},
		{
			Name:     "Empty timeline",
			RecipeID: "3",
			ToAdd:    []string{},
			Error:    false,
		},
		{
			Name:     "Non-existing recipe ID",
			RecipeID: "5",
			Error:    false,
		},
		{
			Name:     "Empty recipe ID",
			RecipeID: "",
			Error:    true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			for _, e := range tc.ToAdd {
				_, err := store.insertStatement.Exec(e, time.Now().Unix(), tc.RecipeID)
				require.NoError(err)
			}
			err := store.DeleteTimeline(tc.RecipeID)
			if tc.Error {
				require.Error(err)
			} else {
				require.NoError(err)
				var event string
				err = getSt.QueryRow(tc.RecipeID).Scan(&event)
				require.ErrorIs(err, sql.ErrNoRows)
			}
		})
	}
}
