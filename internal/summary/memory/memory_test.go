package memory

import (
	"brewday/internal/summary"
	"brewday/internal/tools"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDeleteStats(t *testing.T) {
	require := require.New(t)
	testCases := []struct {
		Name  string
		Stats map[string]*summary.Statistics
		ToDel []string
		Error bool
	}{
		{
			Name:  "Simple stat",
			Error: false,
			Stats: map[string]*summary.Statistics{
				"title1": &summary.Statistics{Evaporation: 50, Efficiency: 50},
			},
			ToDel: []string{"title1"},
		},
		{
			Name:  "Two stats",
			Error: false,
			Stats: map[string]*summary.Statistics{
				"title1": &summary.Statistics{Evaporation: 50, Efficiency: 50},
				"title2": &summary.Statistics{Evaporation: 50, Efficiency: 50},
			},
			ToDel: []string{"title1"},
		},
		{
			Name:  "Two stats delete ne",
			Error: true,
			Stats: map[string]*summary.Statistics{
				"title1": &summary.Statistics{Evaporation: 50, Efficiency: 50},
				"title2": &summary.Statistics{Evaporation: 50, Efficiency: 50},
			},
			ToDel: []string{"title3"},
		},
		{
			Name:  "Delete from empty",
			Error: true,
			Stats: map[string]*summary.Statistics{},
			ToDel: []string{"title1"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			s := NewSummaryMemoryStore()
			for title, st := range tc.Stats {
				s.stats[tools.B64Encode(title)] = st
			}
			initialLen := len(s.stats)
			for _, t := range tc.ToDel {
				err := s.DeleteStats(t)
				if tc.Error {
					require.Error(err)
				} else {
					require.NoError(err)
					_, ok := s.stats[tools.B64Encode(t)]
					require.False(ok)
				}
			}
			if !tc.Error {
				require.Equal(initialLen-len(tc.ToDel), len(s.stats))
			}
		})
	}
}

func TestDeleteSummary(t *testing.T) {
	require := require.New(t)
	testCases := []struct {
		Name         string
		InitialSum   []*summary.Summary
		InitialStats []*summary.Statistics
		ToDel        []string
		Error        bool
	}{
		{
			Name:         "Single summary",
			InitialSum:   []*summary.Summary{{Title: "Summary One"}},
			InitialStats: []*summary.Statistics{{Evaporation: 50.2, Efficiency: 70}},
			ToDel:        []string{"0"},
			Error:        false,
		},
		{
			Name:       "Two summaries",
			InitialSum: []*summary.Summary{{Title: "Summary One"}, {Title: "Summary Two"}},
			InitialStats: []*summary.Statistics{
				{Evaporation: 50.2, Efficiency: 70},
				{Evaporation: 20, Efficiency: 72},
			},
			ToDel: []string{"0"},
			Error: false,
		},
		{
			Name:       "Two summaries delete ne",
			InitialSum: []*summary.Summary{{Title: "Summary One"}, {Title: "Summary Two"}},
			InitialStats: []*summary.Statistics{
				{Evaporation: 50.2, Efficiency: 70},
				{Evaporation: 20, Efficiency: 72},
			},
			ToDel: []string{"2"},
			Error: true,
		},
		{
			Name:         "Delete from empty ",
			InitialSum:   []*summary.Summary{},
			InitialStats: []*summary.Statistics{},
			ToDel:        []string{"0"},
			Error:        true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			s := NewSummaryMemoryStore()
			titles := []string{}
			for i, sum := range tc.InitialSum {
				s.summaries[strconv.Itoa(i)] = sum
				s.stats[tools.B64Encode(sum.Title)] = tc.InitialStats[i]
				titles = append(titles, tools.B64Encode(sum.Title))
			}
			initialLen := len(s.summaries)
			for _, id := range tc.ToDel {
				err := s.DeleteSummary(id)
				if tc.Error {
					require.Error(err)
				} else {
					require.NoError(err)
					_, ok := s.summaries[id]
					require.False(ok)
					index, err := strconv.Atoi(id)
					require.NoError(err)
					_, ok = s.stats[titles[index]]
					require.True(ok)
				}
			}
			if !tc.Error {
				require.Equal(initialLen-len(tc.ToDel), len(s.summaries))
				require.Equal(initialLen, len(s.stats))
			}
		})
	}
}
