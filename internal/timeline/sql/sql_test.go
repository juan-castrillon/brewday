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
		name TEXT NOT NULL,
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
	store, err := NewTimelinePersistentStore(db)
	require.NoError(err)
	defer os.Remove(fileName)
	provisionDB(t, db, []string{"recipe1", "recipe2", "recipe3"})
	getSt, err := db.Prepare(`SELECT timestamp, event FROM timestamps WHERE recipe_id = ?`)
	require.NoError(err)
	testCases := []struct {
		Name     string
		RecipeID string
		Event    string
	}{
		{
			Name:     "pe",
			RecipeID: "1",
			Event:    "my-event",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			current := time.Now()
			err := store.AddEvent(tc.RecipeID, tc.Event)
			require.NoError(err)
			var ts int64
			var event string
			err = getSt.QueryRow(tc.RecipeID).Scan(&ts, &event)
			require.NoError(err)
			require.Equal(tc.Event, event)
			require.WithinRange(time.Unix(ts, 0), current, current.Add(2*time.Second))
		})
	}
}

func TestGetTimeline(t *testing.T) {

}

func TestAddTimeline(t *testing.T) {

}

func TestDeleteTimeline(t *testing.T) {

}
