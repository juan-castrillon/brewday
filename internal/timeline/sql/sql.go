package sql

import (
	"database/sql"
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type TimelinePersistentStore struct {
	dbClient        *sql.DB
	insertStatement *sql.Stmt
}

// NewTimelinePersistentStore creates a new TimelineStore
func NewTimelinePersistentStore(db *sql.DB) (*TimelinePersistentStore, error) {
	err := createTable(db)
	if err != nil {
		return nil, err
	}
	is, err := db.Prepare(`INSERT INTO timelines (event, timestamp_unix, recipe_id) VALUES (?, ?, ?)`)
	if err != nil {
		return nil, err
	}
	return &TimelinePersistentStore{
		dbClient:        db,
		insertStatement: is,
	}, nil
}

// AddEvent adds an event to the timeline
func (s *TimelinePersistentStore) AddEvent(id, message string) error {
	if message == "" {
		return errors.New("invalid empty event for timeline")
	}
	if id == "" {
		return errors.New("invalid empty recipe id for adding event")
	}
	_, err := s.insertStatement.Exec(message, time.Now().Unix(), id)
	return err
}

// GetTimeline returns a timeline of events
func (s *TimelinePersistentStore) GetTimeline(id string) ([]string, error) {
	if id == "" {
		return nil, errors.New("invalid empty recipe id for getting timeline")
	}
	rows, err := s.dbClient.Query(`SELECT timestamp_unix, event FROM timelines WHERE recipe_id = ? ORDER BY timestamp_unix ASC`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []string{}
	for rows.Next() {
		var event string
		var ts int64
		err = rows.Scan(&ts, &event)
		if err != nil {
			return nil, err
		}
		result = append(result, time.Unix(ts, 0).Format(time.RFC3339Nano)+" "+event)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

// AddTimeline adds a timeline to the store
func (s *TimelinePersistentStore) AddTimeline(recipeID string) error {
	if recipeID == "" {
		return errors.New("invalid empty recipe id for addming timeline")
	}
	return s.AddEvent(recipeID, "Initialized Recipe")
}

// DeleteTimeline deletes the timeline for the given recipe id
func (s *TimelinePersistentStore) DeleteTimeline(recipeID string) error {
	if recipeID == "" {
		return errors.New("invalid empty recipe id for deleting timeline")
	}
	_, err := s.dbClient.Exec("DELETE FROM timelines WHERE recipe_id == ?", recipeID)
	return err
}
