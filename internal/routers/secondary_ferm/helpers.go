package secondaryferm

import (
	"brewday/internal/recipe"
	"brewday/internal/tools"
	"brewday/internal/watcher"
	"errors"
	"time"

	"github.com/rs/zerolog/log"
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

// calculateSugar calculates the amount of sugar needed for a certain carbonation level
// It makes several calculations varying the water amount and stores all the results
// Values calculated are 0.1..0.5 liters of water (each 0.1)
func (r *SecondaryFermentationRouter) calculateSugar(id string, volume, carbonation, temperature, alcoholBefore float32, sugarType string) error {
	for i := 1; i <= 5; i++ {
		water := float32(i) / 10
		amount, alcohol := tools.SugarForCarbonation(volume, carbonation, temperature, alcoholBefore, water, sugarType)
		err := r.Store.AddSugarResult(id, &recipe.PrimingSugarResult{
			Water:   water,
			Amount:  amount,
			Alcohol: alcohol,
		})
		if err != nil {
			return err
		}
	}
	return nil
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

// checkWatchers will check it watchers were set for a given recipe.
// If they were not, it will fetch the notification dates from the store and set them up again
// This method helps notifications be persistent in case of restarts.
// It should be called in handlers after the initial watcher setup (where a watcher set up is assumed)
func (r *SecondaryFermentationRouter) checkWatchers(id string) error {
	reset := false
	if r.watchersSet == nil {
		reset = true
	} else {
		val, ok := r.watchersSet[id]
		if !ok {
			// In this case, the map exists but the recipe is not there
			// It could hae been that it got wiped, then restored by other recipe
			reset = true
		} else if !val {
			return errors.New("attempting to get watchers of recipe that does not expect them yet")
		}
	}
	if reset { // If the map is nil its been wiped, set up the watchers again
		re, err := r.Store.Retrieve(id)
		if err != nil {
			return err
		}
		dates, err := r.Store.RetrieveDates(id, "secondary_ferm_notification")
		if err != nil {
			return err
		}
		date := *dates[0]
		var logMessage, notMessage string
		if time.Until(date) < 0 {
			logMessage = "Sending expired secondary_ferm notification for recipe " + id
			notMessage = "Expired Secondary Fermentation Notification. You should have put in the fridge on " + date.Format("2006-01-02")
		} else {
			logMessage = "secondary fermentation notification triggered"
			notMessage = "Time to put bottles in the fridge"
		}
		watcher.NewWatcher(*dates[0], func() error {
			log.Info().Msgf(logMessage)
			return r.sendNotification(notMessage, "Secondary Fermentation "+re.Name, nil)
		}).Start()
		r.addWatchersSet(id)
	}
	return nil
}

func (r *SecondaryFermentationRouter) addWatchersSet(id string) {
	if r.watchersSet == nil {
		r.watchersSet = make(map[string]bool)
	}
	r.watchersSet[id] = true
}
