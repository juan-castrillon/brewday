package watcher

import (
	"sync"
	"time"
)

// Watcher is a component that waits for a certain amount of time and then executes a callback
type Watcher struct {
	// The callback to execute
	callback func() error
	// The date to execute the callback
	executionDate time.Time
	// The channel to receive the stop signal
	stopChan chan bool
	// The mutex to protect the error
	errMu sync.Mutex
	// The error returned by the callback function
	err error
}

// NewWatcher creates a new watcher
func NewWatcher(date time.Time, callback func() error) *Watcher {
	return &Watcher{
		callback:      callback,
		executionDate: date,
		stopChan:      make(chan bool, 1),
	}
}

// Start starts the watcher
// It calculates the time to wait and then waits for that time
func (w *Watcher) Start() {
	go func() {
		select {
		case <-w.stopChan:
			return
		case <-time.After(time.Until(w.executionDate)):
			err := w.callback()
			w.errMu.Lock()
			w.err = err
			w.errMu.Unlock()
		}
	}()
}

// Stop stops the watcher
func (w *Watcher) Stop() {
	w.stopChan <- true
}

// MissingTime returns the time until the callback is executed
func (w *Watcher) MissingTime() time.Duration {
	missing := time.Until(w.executionDate)
	if missing < 0 {
		return 0
	}
	return missing
}

// Err returns the error returned by the callback
func (w *Watcher) Err() error {
	w.errMu.Lock()
	defer w.errMu.Unlock()
	return w.err
}
