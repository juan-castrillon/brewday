package summaryrecorder

import (
	"brewday/internal/summary_recorder/markdown"
	"errors"
	"sync"
)

// SummaryRecorderStore represents a component that stores summary recorders
// The recipe id is used as key
type SummaryRecorderStore struct {
	lock             sync.Mutex
	summaryRecorders map[string]SummaryRecorder
}

// NewSummaryRecorderStore creates a new SummaryRecorderStore
func NewSummaryRecorderStore() *SummaryRecorderStore {
	return &SummaryRecorderStore{
		summaryRecorders: make(map[string]SummaryRecorder),
	}
}

// SummaryRecorderFactory is a factory for SummaryRecorder
func SummaryRecorderFactory(recorderType string) SummaryRecorder {
	switch recorderType {
	case "markdown":
		return markdown.NewMarkdownSummaryRecorder()
	default:
		return markdown.NewMarkdownSummaryRecorder()
	}
}

// GetSummaryRecorder returns the summary recorder for the given recipe id
func (s *SummaryRecorderStore) GetSummaryRecorder(recipeID string) (SummaryRecorder, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	sr, ok := s.summaryRecorders[recipeID]
	if !ok {
		return nil, errors.New("no summary recorder found for recipe id " + recipeID)
	}
	return sr, nil
}

// AddSummaryRecorder adds a summary recorder for the given recipe id
func (s *SummaryRecorderStore) AddSummaryRecorder(recipeID string, recorderType string) {
	summaryRecorder := SummaryRecorderFactory(recorderType)
	s.lock.Lock()
	defer s.lock.Unlock()
	s.summaryRecorders[recipeID] = summaryRecorder
}

// AddMashTemp adds a mash temperature to the summary and notes related to it
func (s *SummaryRecorderStore) AddMashTemp(id string, temp float64, notes string) error {
	rec, err := s.GetSummaryRecorder(id)
	if err != nil {
		return err
	}
	rec.AddMashTemp(temp, notes)
	return nil
}

// AddRast adds a rast to the summary and notes related to it
func (s *SummaryRecorderStore) AddRast(id string, temp float64, duration float64, notes string) error {
	rec, err := s.GetSummaryRecorder(id)
	if err != nil {
		return err
	}
	rec.AddRast(temp, duration, notes)
	return nil
}

// AddLauternNotes adds lautern notes to the summary
func (s *SummaryRecorderStore) AddLaunternNotes(id, notes string) error {
	rec, err := s.GetSummaryRecorder(id)
	if err != nil {
		return err
	}
	rec.AddLaunternNotes(notes)
	return nil
}

// AddHopping adds a hopping to the summary and notes related to it
func (s *SummaryRecorderStore) AddHopping(id string, name string, amount float32, alpha float32, duration float32, notes string) error {
	rec, err := s.GetSummaryRecorder(id)
	if err != nil {
		return err
	}
	rec.AddHopping(name, amount, alpha, duration, notes)
	return nil
}

// AddMeasuredVolume adds a measured volume to the summary
func (s *SummaryRecorderStore) AddMeasuredVolume(id string, name string, amount float32, notes string) error {
	rec, err := s.GetSummaryRecorder(id)
	if err != nil {
		return err
	}
	rec.AddMeasuredVolume(name, amount, notes)
	return nil
}

// AddEvaporation adds an evaporation to the summary
func (s *SummaryRecorderStore) AddEvaporation(id string, amount float32) error {
	rec, err := s.GetSummaryRecorder(id)
	if err != nil {
		return err
	}
	rec.AddEvaporation(amount)
	return nil
}

// AddCooling adds a cooling to the summary and notes related to it
func (s *SummaryRecorderStore) AddCooling(id string, finalTemp, coolingTime float32, notes string) error {
	rec, err := s.GetSummaryRecorder(id)
	if err != nil {
		return err
	}
	rec.AddCooling(finalTemp, coolingTime, notes)
	return nil
}