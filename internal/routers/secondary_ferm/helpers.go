package secondaryferm

import (
	"brewday/internal/tools"
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

// addSugarResult adds a sugar result to a recipe
func (r *SecondaryFermentationRouter) addSugarResult(id string, result *SugarResult) {
	r.sugarResultsLock.Lock()
	defer r.sugarResultsLock.Unlock()
	if r.SugarResults == nil {
		r.SugarResults = make(map[string][]SugarResult)
	}
	_, ok := r.SugarResults[id]
	if !ok {
		r.SugarResults[id] = []SugarResult{}
	}
	r.SugarResults[id] = append(r.SugarResults[id], *result)
}

// getSugarResults retrieves the sugar results for a recipe
func (r *SecondaryFermentationRouter) getSugarResults(id string) ([]SugarResult, error) {
	r.sugarResultsLock.Lock()
	defer r.sugarResultsLock.Unlock()
	list, ok := r.SugarResults[id]
	if !ok {
		return nil, errors.New("no sugar results found for recipe")
	}
	return list, nil
}

// addSecondaryWatcher adds a watcher for the secondary fermentation
func (r *SecondaryFermentationRouter) addSecondaryWatcher(id string, watcher *watcher.Watcher) error {
	r.secondaryWatchersLock.Lock()
	defer r.secondaryWatchersLock.Unlock()
	if r.SecondaryWatchers == nil {
		r.SecondaryWatchers = make(map[string]SecondaryFermentationWatcher)
	}
	r.SecondaryWatchers[id] = SecondaryFermentationWatcher{
		watch: watcher,
	}
	return nil
}

// getSecondaryWatcher retrieves a watcher for the secondary fermentation
// if none is found, it returns nil
func (r *SecondaryFermentationRouter) getSecondaryWatcher(id string) *watcher.Watcher {
	r.secondaryWatchersLock.Lock()
	defer r.secondaryWatchersLock.Unlock()
	w, ok := r.SecondaryWatchers[id]
	if !ok {
		return nil
	}
	return w.watch
}

// calculateSugar calculates the amount of sugar needed for a certain carbonation level
// It makes several calculations varying the water amount and stores all the results
// Values calculated are 0.1..0.5 liters of water (each 0.1)
func (r *SecondaryFermentationRouter) calculateSugar(id string, volume, carbonation, temperature, alcoholBefore float32, sugarType string) {
	for i := 1; i <= 5; i++ {
		water := float32(i) / 10
		amount, alcohol := tools.SugarForCarbonation(volume, carbonation, temperature, alcoholBefore, water, sugarType)
		r.addSugarResult(id, &SugarResult{
			Water:   water,
			Amount:  amount,
			Alcohol: alcohol,
		})
	}
}

// addSummaryDryHop adds a pre fermentation summary
func (r *SecondaryFermentationRouter) addSummaryDryHop(id string, name string, amount, alpha, days float32, notes string) error {
	if r.SummaryStore != nil {
		return r.SummaryStore.AddMainFermentationDryHop(id, name, amount, alpha, days, notes)
	}
	return nil
}

// addSummaryBottle adds a summary of the bottling
func (r *SecondaryFermentationRouter) addSummaryBottle(id string, carbonation, alcohol, sugar, temp, vol float32, st, notes string) error {
	if r.SummaryStore != nil {
		return r.SummaryStore.AddBottling(id, carbonation, alcohol, sugar, temp, vol, st, notes)
	}
	return nil
}

// addSummaryPreBottle adds a summary of the pre bottling
func (r *SecondaryFermentationRouter) addSummaryPreBottle(id string, volume float32) error {
	if r.SummaryStore != nil {
		return r.SummaryStore.AddPreBottlingVolume(id, volume)
	}
	return nil
}

// addSummarySecondaryFermentation adds a summary of the secondary fermentation
func (r *SecondaryFermentationRouter) addSummarySecondaryFermentation(id string, days int, notes string) error {
	if r.SummaryStore != nil {
		return r.SummaryStore.AddSummarySecondary(id, days, notes)
	}
	return nil
}
