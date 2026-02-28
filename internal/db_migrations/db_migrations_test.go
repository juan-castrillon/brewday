package dbmigrations

import (
	"database/sql"
	"embed"
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"
)

//go:embed test_data/migrations/*.sql
var testMigrationsFS embed.FS

//go:embed test_data/nothing/*
var notMigrationsFS embed.FS

func TestRunMigrations(t *testing.T) {
	emptyFS := fstest.MapFS{}
	require := require.New(t)
	testCases := []struct {
		Name    string
		FS      []fs.FS
		Path    string
		Error   bool
		Tables  []string
		Indexes []string
	}{
		{
			Name:  "Empty FS",
			Error: true,
			FS:    []fs.FS{emptyFS},
		},
		{
			Name:  "FS with no migrations",
			Error: true,
			FS:    []fs.FS{notMigrationsFS},
			Path:  "test_data/nothing",
		},
		{
			Name:  "Multiple FS",
			Error: true,
			FS:    []fs.FS{notMigrationsFS, testMigrationsFS},
		},
		{
			Name:    "Valid FS",
			Error:   false,
			FS:      []fs.FS{testMigrationsFS},
			Path:    "test_data/migrations",
			Tables:  []string{"dummy", "indexed"},
			Indexes: []string{"ix_event"},
		},
		{
			Name:    "Prod FS",
			Error:   false,
			FS:      []fs.FS{},
			Path:    "migrations",
			Tables:  []string{"bool_flags", "dates", "main_ferm_sgs", "recipe_results", "recipes", "stats", "sugar_results", "summaries", "timelines"},
			Indexes: []string{"ix_bool_flags", "ix_dates", "ix_main_ferm_sgs", "ix_stats", "ix_sugar_results", "ix_summaries", "ix_timelines"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			db, err := sql.Open("sqlite3", ":memory:")
			require.NoError(err)
			err = RunMigrations(db, tc.Path, tc.FS...)
			if tc.Error {
				require.Error(err)
			} else {
				require.NoError(err)
				var name string
				for _, tableName := range tc.Tables {
					require.NoError(db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name=?", tableName).Scan(&name))
				}
				for _, indexName := range tc.Indexes {
					require.NoError(db.QueryRow("SELECT name FROM sqlite_master WHERE type='index' AND name=?", indexName).Scan(&name))
				}
			}
		})
	}
}

func TestSeedDataMigrations(t *testing.T) {
	require := require.New(t)
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(err)
	err = RunMigrations(db, "migrations")
	testCases := []struct {
		Name        string
		Query       string
		QueryValues []any
		Expected    string
	}{
		{
			Name:        "Test a recipe in stats",
			Query:       "SELECT efficiency FROM stats WHERE recipe_title == ?", // ('U2Fib3IgVHJvcGljYWw=', 1719180000, 16.4, 72) -- Sabor Tropical
			QueryValues: []any{"U2Fib3IgVHJvcGljYWw="},
			Expected:    "72",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			var result string
			require.NoError(db.QueryRow(tc.Query, tc.QueryValues...).Scan(&result))
			require.Equal(tc.Expected, result)
		})
	}

}
