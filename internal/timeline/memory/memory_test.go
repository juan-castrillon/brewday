package memory

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDeleteTimeline(t *testing.T) {
	require := require.New(t)
	testCases := []struct {
		Name  string
		Tls   map[string]*BasicTimeline
		ToDel []string
		Error bool
	}{
		{
			Name:  "Simple timeline",
			Error: false,
			Tls: map[string]*BasicTimeline{
				"id1": &BasicTimeline{},
			},
			ToDel: []string{"id1"},
		},
		{
			Name:  "Two stats",
			Error: false,
			Tls: map[string]*BasicTimeline{
				"id1": &BasicTimeline{},
				"id2": &BasicTimeline{},
			},
			ToDel: []string{"id1"},
		},
		{
			Name:  "Two stats delete ne",
			Error: true,
			Tls: map[string]*BasicTimeline{
				"id1": &BasicTimeline{},
				"id2": &BasicTimeline{},
			},
			ToDel: []string{"id3"},
		},
		{
			Name:  "Delete from empty",
			Error: true,
			Tls:   map[string]*BasicTimeline{},
			ToDel: []string{"id1"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			s := NewTimelineMemoryStore()
			for id, _ := range tc.Tls {
				require.NoError(s.AddTimeline(id))
			}
			initialLen := len(s.timelines)
			for _, id := range tc.ToDel {
				err := s.DeleteTimeline(id)
				if tc.Error {
					require.Error(err)
				} else {
					require.NoError(err)
					_, ok := s.timelines[id]
					require.False(ok)
				}
			}
			if !tc.Error {
				require.Equal(initialLen-len(tc.ToDel), len(s.timelines))
			}
		})
	}
}
