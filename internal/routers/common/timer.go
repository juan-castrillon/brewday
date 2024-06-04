package common

import (
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// RecipeStore represents a component that stores recipes
type RecipeStore interface {
	// AddDate allows to store a date with a certain purpose. It can be used to store notification dates, or timers
	AddDate(id string, date *time.Time, name string) error
	// RetrieveDates allows to retreive stored dates with its purpose (name).It can be used to store notification dates, or timers
	// It supports pattern in the name to retrieve multiple values
	RetrieveDates(id, namePattern string) ([]*time.Time, error)
	// AddBoolFlag allows to store a given flag that can be true or false in the store with a unique name
	AddBoolFlag(id, name string, flag bool) error
	// RetrieveBoolFlag gets a bool flag from the store given its name
	RetrieveBoolFlag(id, name string) (bool, error)
}

type RespGetTimestamp struct {
	EndTimestamp int64 `json:"end_timestamp"`
}

type RespGetRealDuration struct {
	RealDurationMinutes float32 `json:"real_duration_minutes,omitempty"`
}

type ReqPostStopTimer struct {
	StoppedTimestamp int64 `json:"stopped_timestamp"`
}

type Timer struct {
	Store RecipeStore
}

func NewTimer(store RecipeStore) *Timer {
	return &Timer{
		Store: store,
	}
}

func (t *Timer) getName(prefix, suffix, purpose string) string {
	var s string
	switch purpose {
	case "start":
		s = "started"
	case "stop":
		s = "stopped"
	case "end":
		s = "end"
	default:
		s = "unknown"
	}
	res := prefix + "_" + s
	if suffix != "" {
		res = res + "_" + suffix
	}
	return res
}

// GetBoolFlags returns whether the timer has started and has been stopped. Only the first suffix is used
func (t *Timer) GetBoolFlags(id string, prefix string, suffix ...string) (bool, bool, error) {
	singleSuffix := ""
	if len(suffix) > 0 {
		singleSuffix = suffix[0]
	}
	started, err := t.Store.RetrieveBoolFlag(id, t.getName(prefix, singleSuffix, "start"))
	if err != nil {
		return false, false, err
	}
	stopped, err := t.Store.RetrieveBoolFlag(id, t.getName(prefix, singleSuffix, "stop"))
	if err != nil {
		return false, false, err
	}
	return started, stopped, nil
}

// HandleStartTimer will respond with the correct json for the timer template to work. Only the first suffix is used
func (t *Timer) HandleStartTimer(c echo.Context, id string, duration time.Duration, prefix string, suffix ...string) error {
	singleSuffix := ""
	if len(suffix) > 0 {
		singleSuffix = suffix[0]
	}
	started, err := t.Store.RetrieveBoolFlag(id, t.getName(prefix, singleSuffix, "start"))
	if err != nil {
		return err
	}
	var stopTs time.Time
	if !started {
		err = t.Store.AddBoolFlag(id, t.getName(prefix, singleSuffix, "start"), true)
		if err != nil {
			return err
		}
		now := time.Now()
		err = t.Store.AddDate(id, &now, t.getName(prefix, singleSuffix, "start"))
		if err != nil {
			return err
		}
		stopTs = now.Add(duration)
		err = t.Store.AddDate(id, &stopTs, t.getName(prefix, singleSuffix, "end"))
		if err != nil {
			return err
		}
	} else {
		setDates, err := t.Store.RetrieveDates(id, t.getName(prefix, singleSuffix, "end"))
		if err != nil {
			return err
		}
		if len(setDates) == 0 {
			return errors.New("invalid empty date for " + prefix + " end")
		}
		stopTs = *setDates[0]
	}
	resp := &RespGetTimestamp{
		EndTimestamp: stopTs.Unix(),
	}
	return c.JSON(http.StatusOK, resp)
}

// HandleStopTimer will mark the timer as stopped. Only the first suffix is used
func (t *Timer) HandleStopTimer(c echo.Context, id string, prefix string, suffix ...string) error {
	singleSuffix := ""
	if len(suffix) > 0 {
		singleSuffix = suffix[0]
	}
	name := t.getName(prefix, singleSuffix, "stop")
	var req ReqPostStopTimer
	err := c.Bind(&req)
	if err != nil {
		return err
	}
	stopped, err := t.Store.RetrieveBoolFlag(id, name)
	if err != nil {
		return err
	}
	if !stopped {
		err = t.Store.AddBoolFlag(id, name, true)
		if err != nil {
			return err
		}
		st := time.Unix(req.StoppedTimestamp, 0)
		err = t.Store.AddDate(id, &st, name)
		if err != nil {
			return err
		}
	}
	return c.NoContent(http.StatusOK)
}

// HandleRealDuration will return the real duration to the timer template. Only the first suffix is used
func (t *Timer) HandleRealDuration(c echo.Context, id string, prefix string, suffix ...string) error {
	singleSuffix := ""
	if len(suffix) > 0 {
		singleSuffix = suffix[0]
	}
	startDates, err := t.Store.RetrieveDates(id, t.getName(prefix, singleSuffix, "start"))
	if err != nil {
		return err
	}
	if len(startDates) == 0 {
		return errors.New("invalid empty date for " + prefix + " start")
	}
	stoppedDates, err := t.Store.RetrieveDates(id, t.getName(prefix, singleSuffix, "stop"))
	if err != nil {
		return err
	}
	if len(stoppedDates) == 0 {
		return errors.New("invalid empty date for " + prefix + "  stopped")
	}
	start := *startDates[0]
	stopped := *stoppedDates[0]
	realDur := stopped.Sub(start)
	resp := &RespGetRealDuration{
		RealDurationMinutes: float32(realDur.Minutes()),
	}
	return c.JSON(http.StatusOK, resp)
}
