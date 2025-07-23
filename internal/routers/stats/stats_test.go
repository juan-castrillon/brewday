package stats

import (
	"brewday/internal/summary"
	"testing"

	"github.com/stretchr/testify/require"
)

type mockStore struct {
	store map[string]*summary.Statistics
}

func (s *mockStore) GetAllStats() (map[string]*summary.Statistics, error) {
	return s.store, nil
}

func TestGetStats(t *testing.T) {
	require := require.New(t)
	type testCase struct {
		Name     string
		Store    map[string]*summary.Statistics
		Expected map[string]*summary.Statistics
		Error    bool
	}
	testCases := []testCase{
		{
			Name: "Two summaries",
			Store: map[string]*summary.Statistics{
				"1": {Evaporation: 20.5, Efficiency: 62.3},
				"2": {Evaporation: 16.333, Efficiency: 72.84},
			},
			Expected: map[string]*summary.Statistics{
				"1": {Evaporation: 20.5, Efficiency: 62.3},
				"2": {Evaporation: 16.333, Efficiency: 72.84},
			},
			Error: false,
		},
		{
			Name:     "No Summaries",
			Store:    map[string]*summary.Statistics{},
			Expected: map[string]*summary.Statistics{},
			Error:    false,
		},
		{
			Name:     "Nil store",
			Store:    nil,
			Expected: nil,
			Error:    false, // To be fixed when fixing interfaces nil comparison
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			mockS := mockStore{store: tc.Store}
			router := StatsRouter{StatsStore: &mockS}
			res, err := router.getStats()
			if tc.Error {
				require.Error(err)
			} else {
				require.NoError(err)
				require.Equal(tc.Expected, res)
			}
		})
	}
}
