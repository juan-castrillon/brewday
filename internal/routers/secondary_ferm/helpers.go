package secondaryferm

import (
	"brewday/internal/watcher"
	"errors"
)

// addTimelineEvent adds a timeline event if the timeline store is available
func (r *SecondaryFermentationRouter) addTimelineEvent(id, message string) error {
	if r.TLStore != nil {
		return r.TLStore.AddEvent(id, message)
	}
	return nil
}

// sendNotification sends a notification if the notifier is available
func (r *SecondaryFermentationRouter) sendNotification(message, title string, opts map[string]interface{}) error {
	if r.Notifier != nil {
		return r.Notifier.Send(message, title, opts)
	}
	return nil
}

// addDryHop adds a dry hop to a recipe
func (r *SecondaryFermentationRouter) addDryHop(id string, dryHop *DryHop) {
	r.dryHopsLock.Lock()
	defer r.dryHopsLock.Unlock()
	if r.DryHops == nil {
		r.DryHops = make(map[string]DryHopMap)
	}
	_, ok := r.DryHops[id]
	if !ok {
		r.DryHops[id] = make(DryHopMap)
	}
	r.DryHops[id][dryHop.id] = dryHop
}

// getDryHops retrieves the dry hops for a recipe
func (r *SecondaryFermentationRouter) getDryHops(id string) (DryHopMap, error) {
	r.dryHopsLock.Lock()
	defer r.dryHopsLock.Unlock()
	list, ok := r.DryHops[id]
	if !ok {
		return nil, errors.New("no dry hops found for recipe")
	}
	return list, nil
}

// addDryHopNotification adds a notification for a certain dry hop in a recipe
func (r *SecondaryFermentationRouter) addDryHopNotification(id, dryHopID string, watcher *watcher.Watcher) error {
	r.hopWatchersLock.Lock()
	defer r.hopWatchersLock.Unlock()
	if r.HopWatchers == nil {
		r.HopWatchers = make(map[string]DryHopNotification)
	}
	_, ok := r.HopWatchers[id]
	if !ok {
		r.HopWatchers[id] = make(DryHopNotification)
	}
	r.HopWatchers[id][dryHopID] = watcher
	dh, err := r.getDryHops(id)
	if err != nil {
		return err
	}
	dh[dryHopID].NotificationSet = true
	return nil
}

// getDryHopNotification retrieves a notification for a certain dry hop in a recipe
func (r *SecondaryFermentationRouter) getDryHopNotifications(id string) (DryHopNotification, error) {
	r.hopWatchersLock.Lock()
	defer r.hopWatchersLock.Unlock()
	list, ok := r.HopWatchers[id]
	if !ok {
		return nil, errors.New("no dry hop notifications found for recipe")
	}
	return list, nil
}

// addSummaryDryHop adds a pre fermentation summary
func (r *SecondaryFermentationRouter) addSummaryDryHop(id string, name string, amount float32) error {
	if r.SummaryStore != nil {
		return r.SummaryStore.AddSummaryDryHop(id, name, amount)
	}
	return nil
}
