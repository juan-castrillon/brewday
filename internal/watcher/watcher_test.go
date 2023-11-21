package watcher

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestWatcher(t *testing.T) {
	executionTime := time.Time{}
	var lock sync.Mutex
	tests := []struct {
		name                string
		timeToAdd           time.Duration
		callback            func() error
		expectedErr         error
		expectedMissingTime time.Duration
		timeToStop          time.Duration
		expectFinish        bool
	}{
		{
			name:      "callback is executed after waiting",
			timeToAdd: 200 * time.Millisecond,
			callback: func() error {
				lock.Lock()
				defer lock.Unlock()
				executionTime = time.Now()
				return nil
			},
			expectedErr:         nil,
			expectedMissingTime: 0,
			timeToStop:          0,
			expectFinish:        true,
		},
		{
			name:      "callback is not executed if stopped",
			timeToAdd: 200 * time.Millisecond,
			callback: func() error {
				lock.Lock()
				defer lock.Unlock()
				executionTime = time.Now()
				return nil
			},
			expectedErr:         nil,
			expectedMissingTime: 100 * time.Millisecond,
			timeToStop:          100 * time.Millisecond,
			expectFinish:        false,
		},
		{
			name:      "callback returns an error",
			timeToAdd: 200 * time.Millisecond,
			callback: func() error {
				return errors.New("error")
			},
			expectedErr:         errors.New("error"),
			expectedMissingTime: 0,
			timeToStop:          0,
			expectFinish:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executionDate := time.Now().Add(tt.timeToAdd)
			w := NewWatcher(executionDate, tt.callback)
			w.Start()
			if tt.timeToStop > 0 {
				time.Sleep(tt.timeToStop)
				w.Stop()
			} else {
				time.Sleep(tt.timeToAdd + 100*time.Millisecond)
			}
			require.WithinDuration(t, time.Now().Add(tt.expectedMissingTime), time.Now().Add(w.MissingTime()), 20*time.Millisecond)
			require.Equal(t, tt.expectedErr, w.Err())
			lock.Lock()
			defer lock.Unlock()
			if tt.expectFinish {
				require.WithinDuration(t, executionDate, executionTime, 20*time.Millisecond)
			} else {
				require.Equal(t, time.Time{}, executionTime)
			}
			executionTime = time.Time{}
		})
	}
}
