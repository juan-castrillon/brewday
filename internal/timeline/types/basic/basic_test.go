package basic

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestAddEventSimple(t *testing.T) {
	tl := NewBasicTimeline()
	tl.AddEvent("test")
	require.Equal(t, 1, len(tl.events))
	require.Equal(t, "test", tl.events[0].Message)
}

func TestAddEventSequential(t *testing.T) {
	tl := NewBasicTimeline()
	tl.AddEvent("test1")
	tl.AddEvent("test2")
	require.Equal(t, 2, len(tl.events))
	require.Equal(t, "test1", tl.events[0].Message)
	require.Equal(t, "test2", tl.events[1].Message)
}

func TestAddEventParallel(t *testing.T) {
	tl := NewBasicTimeline()
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		i := i
		go func(index int) {
			defer wg.Done()
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
			tl.AddEvent(fmt.Sprintf("test%d", index))
		}(i)
	}
	wg.Wait()
	require.Equal(t, 10, len(tl.events))
	for i := 0; i < 10; i++ {
		wanted := fmt.Sprintf("test%d", i)
		found := false
		for _, event := range tl.events {
			if event.Message == wanted {
				found = true
				break
			}
		}
		require.True(t, found)
	}
}

func TestGetTimeline(t *testing.T) {
	tl := NewBasicTimeline()
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		i := i
		go func(index int) {
			defer wg.Done()
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			tl.AddEvent(fmt.Sprintf("test%d", index))
		}(i)
	}
	wg.Wait()
	require.Equal(t, 10, len(tl.events))
	// Now check that calling GetTimeline() returns the events sorted by timestamp
	result := tl.GetTimeline()
	require.Equal(t, 10, len(result))
	dates := make([]time.Time, len(result))
	for i := 0; i < len(result); i++ {
		sp := strings.Split(result[i], " ")
		rawDate := sp[0]
		date, err := time.Parse(time.RFC3339, rawDate)
		require.NoError(t, err)
		dates[i] = date
	}
	for i := 0; i < len(dates)-1; i++ {
		require.True(t, dates[i].Before(dates[i+1]))
	}
}
