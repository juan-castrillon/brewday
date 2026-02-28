package stats

import (
	"brewday/internal/summary"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type mockStore struct {
	store map[string]*summary.Statistics
}

func (s *mockStore) GetAllStats() (map[string]*summary.Statistics, error) {
	return s.store, nil
}

func (s *mockStore) AddStatsExternal(recipeName string, stats *summary.Statistics) error {
	s.store[recipeName] = stats
	return nil
}

func ptrFloat32(num float32) *float32 {
	return &num
}

func TestGetStats(t *testing.T) {
	require := require.New(t)
	type testCase struct {
		Name     string
		Store    map[string]*summary.Statistics
		Expected []StatEntry
		Error    bool
	}
	testCases := []testCase{
		{
			Name: "One summaries",
			Store: map[string]*summary.Statistics{
				"Test1": {Evaporation: 20.5, Efficiency: 62.3, FinishedTime: time.Unix(150, 0)},
			},
			Expected: []StatEntry{
				{RecipeName: "Test1", Evaporation: ptrFloat32(20.5), Efficiency: ptrFloat32(62.3), FinishedTimeString: "1970-01-01", FinishedTimeEpoch: 150},
			},
			Error: false,
		},
		{
			Name: "Two summaries",
			Store: map[string]*summary.Statistics{
				"Test1": {Evaporation: 20.5, Efficiency: 62.3, FinishedTime: time.Unix(150, 0)},
				"Test2": {Evaporation: 16.333, Efficiency: 72.84, FinishedTime: time.Unix(150000, 0)},
			},
			Expected: []StatEntry{
				{RecipeName: "Test1", Evaporation: ptrFloat32(20.5), Efficiency: ptrFloat32(62.3), FinishedTimeString: "1970-01-01", FinishedTimeEpoch: 150},
				{RecipeName: "Test2", Evaporation: ptrFloat32(16.333), Efficiency: ptrFloat32(72.84), FinishedTimeString: "1970-01-02", FinishedTimeEpoch: 150000},
			},
			Error: false,
		},
		{
			Name:     "No Summaries",
			Store:    map[string]*summary.Statistics{},
			Expected: []StatEntry{},
			Error:    false,
		},
		{
			Name:     "Nil store",
			Store:    nil,
			Expected: []StatEntry{},
			Error:    false, // To be fixed when fixing interfaces nil comparison
		},
		{
			Name: "Summary with Efficiency 0",
			Store: map[string]*summary.Statistics{
				"Test1": {Evaporation: 20.5, Efficiency: 0, FinishedTime: time.Unix(150, 0)},
			},
			Expected: []StatEntry{
				{RecipeName: "Test1", Evaporation: ptrFloat32(20.5), Efficiency: nil, FinishedTimeString: "1970-01-01", FinishedTimeEpoch: 150},
			},
			Error: false,
		},
		{
			Name: "Summary with Evaporation 0",
			Store: map[string]*summary.Statistics{
				"Test1": {Evaporation: 0, Efficiency: 62.3, FinishedTime: time.Unix(150, 0)},
			},
			Expected: []StatEntry{
				{RecipeName: "Test1", Evaporation: nil, Efficiency: ptrFloat32(62.3), FinishedTimeString: "1970-01-01", FinishedTimeEpoch: 150},
			},
			Error: false,
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
				require.ElementsMatch(tc.Expected, res)
			}
		})
	}
}

func TestAddStats(t *testing.T) {
	require := require.New(t)
	type testCase struct {
		Name     string
		Store    map[string]*summary.Statistics
		ToAdd    *ReqPostAddStat
		Expected map[string]*summary.Statistics
		Error    bool
	}
	testCases := []testCase{
		{
			Name:  "Add recipe, empty store",
			Store: map[string]*summary.Statistics{},
			ToAdd: &ReqPostAddStat{
				RecipeName:         "Test1",
				Evaporation:        60.2,
				Efficiency:         60.3,
				FinishedTimeString: "2025-12-25",
			},
			Expected: map[string]*summary.Statistics{
				"Test1": {
					Evaporation:  60.2,
					Efficiency:   60.3,
					FinishedTime: time.Date(2025, time.December, 25, 0, 0, 0, 0, time.UTC),
				},
			},
			Error: false,
		},
		{
			Name: "Add recipe, non-empty store",
			Store: map[string]*summary.Statistics{
				"Test1": {
					Evaporation:  60.2,
					Efficiency:   60.3,
					FinishedTime: time.Date(2025, time.December, 25, 0, 0, 0, 0, time.UTC),
				},
			},
			ToAdd: &ReqPostAddStat{
				RecipeName:         "Test2",
				Evaporation:        60.2,
				Efficiency:         60.3,
				FinishedTimeString: "2025-12-25",
			},
			Expected: map[string]*summary.Statistics{
				"Test1": {
					Evaporation:  60.2,
					Efficiency:   60.3,
					FinishedTime: time.Date(2025, time.December, 25, 0, 0, 0, 0, time.UTC),
				},
				"Test2": {
					Evaporation:  60.2,
					Efficiency:   60.3,
					FinishedTime: time.Date(2025, time.December, 25, 0, 0, 0, 0, time.UTC),
				},
			},
			Error: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			mockS := mockStore{store: tc.Store}
			router := StatsRouter{StatsStore: &mockS}
			err := router.addStats(tc.ToAdd)
			if tc.Error {
				require.Error(err)
			} else {
				require.NoError(err)
				require.Equal(tc.Expected, mockS.store)
			}
		})
	}
}
