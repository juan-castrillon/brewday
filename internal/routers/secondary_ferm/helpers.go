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

// addSummaryFinishedTime adds te finished time to the summary
func (r *SecondaryFermentationRouter) addSummaryFinishedTime(id string, t time.Time) error {
	if r.SummaryStore != nil {
		return r.SummaryStore.AddFinishedTime(id, t)
	}
	return nil
}

// addStats saves Statistics from the finished recipe for future reference
func (r *SecondaryFermentationRouter) addStats(id string) error {
	if r.StatsStore != nil {
		return r.StatsStore.AddStats(id)
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
			log.Info().Msgf("%s", logMessage)
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
