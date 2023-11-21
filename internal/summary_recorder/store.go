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

// getSummaryRecorder returns the summary recorder for the given recipe id
func (s *SummaryRecorderStore) getSummaryRecorder(recipeID string) (SummaryRecorder, error) {
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

// DeleteSummaryRecorder deletes the summary recorder for the given recipe id
func (s *SummaryRecorderStore) DeleteSummaryRecorder(recipeID string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.summaryRecorders, recipeID)
}

// AddMashTemp adds a mash temperature to the summary and notes related to it
func (s *SummaryRecorderStore) AddMashTemp(id string, temp float64, notes string) error {
	rec, err := s.getSummaryRecorder(id)
	if err != nil {
		return err
	}
	rec.AddMashTemp(temp, notes)
	return nil
}

// AddRast adds a rast to the summary and notes related to it
func (s *SummaryRecorderStore) AddRast(id string, temp float64, duration float64, notes string) error {
	rec, err := s.getSummaryRecorder(id)
	if err != nil {
		return err
	}
	rec.AddRast(temp, duration, notes)
	return nil
}

// AddLauternNotes adds lautern notes to the summary
func (s *SummaryRecorderStore) AddLaunternNotes(id, notes string) error {
	rec, err := s.getSummaryRecorder(id)
	if err != nil {
		return err
	}
	rec.AddLaunternNotes(notes)
	return nil
}

// AddHopping adds a hopping to the summary and notes related to it
func (s *SummaryRecorderStore) AddHopping(id string, name string, amount float32, alpha float32, duration float32, notes string) error {
	rec, err := s.getSummaryRecorder(id)
	if err != nil {
		return err
	}
	rec.AddHopping(name, amount, alpha, duration, notes)
	return nil
}

// AddMeasuredVolume adds a measured volume to the summary
func (s *SummaryRecorderStore) AddMeasuredVolume(id string, name string, amount float32, notes string) error {
	rec, err := s.getSummaryRecorder(id)
	if err != nil {
		return err
	}
	rec.AddMeasuredVolume(name, amount, notes)
	return nil
}

// AddEvaporation adds an evaporation to the summary
func (s *SummaryRecorderStore) AddEvaporation(id string, amount float32) error {
	rec, err := s.getSummaryRecorder(id)
	if err != nil {
		return err
	}
	rec.AddEvaporation(amount)
	return nil
}

// AddCooling adds a cooling to the summary and notes related to it
func (s *SummaryRecorderStore) AddCooling(id string, finalTemp, coolingTime float32, notes string) error {
	rec, err := s.getSummaryRecorder(id)
	if err != nil {
		return err
	}
	rec.AddCooling(finalTemp, coolingTime, notes)
	return nil
}

// AddSummaryPreFermentation adds a summary of the pre fermentation
func (s *SummaryRecorderStore) AddSummaryPreFermentation(id string, volume float32, sg float32, notes string) error {
	rec, err := s.getSummaryRecorder(id)
	if err != nil {
		return err
	}
	rec.AddSummaryPreFermentation(volume, sg, notes)
	return nil
}

// AddEfficiency adds the efficiency (sudhausausbeute) to the summary
func (s *SummaryRecorderStore) AddEfficiency(id string, efficiencyPercentage float32) error {
	rec, err := s.getSummaryRecorder(id)
	if err != nil {
		return err
	}
	rec.AddEfficiency(efficiencyPercentage)
	return nil
}

// AddYeastStart adds the yeast start to the summary
func (s *SummaryRecorderStore) AddYeastStart(id string, temperature, notes string) error {
	rec, err := s.getSummaryRecorder(id)
	if err != nil {
		return err
	}
	rec.AddYeastStart(temperature, notes)
	return nil
}

// AddSGMeasurement adds a SG measurement to the summary
func (s *SummaryRecorderStore) AddSGMeasurement(id string, date string, gravity float32, final bool, notes string) error {
	rec, err := s.getSummaryRecorder(id)
	if err != nil {
		return err
	}
	rec.AddSGMeasurement(date, gravity, final, notes)
	return nil
}

// AddAlcoholMainFermentation adds the alcohol after the main fermentation to the summary
func (s *SummaryRecorderStore) AddAlcoholMainFermentation(id string, alcohol float32) error {
	rec, err := s.getSummaryRecorder(id)
	if err != nil {
		return err
	}
	rec.AddAlcoholMainFermentation(alcohol)
	return nil
}

// AddSummaryDryHop adds a summary of the dry hop
func (s *SummaryRecorderStore) AddSummaryDryHop(id string, name string, amount float32) error {
	rec, err := s.getSummaryRecorder(id)
	if err != nil {
		return err
	}
	rec.AddSummaryDryHop(name, amount)
	return nil
}

// AddSummaryPreBottle adds a summary of the pre bottling
func (s *SummaryRecorderStore) AddSummaryPreBottle(id string, volume float32) error {
	rec, err := s.getSummaryRecorder(id)
	if err != nil {
		return err
	}
	rec.AddSummaryPreBottle(volume)
	return nil
}

// AddSummaryBottle adds a summary of the bottling
func (s *SummaryRecorderStore) AddSummaryBottle(id string, carbonation, alcohol, sugar, temp, vol float32, sugarType, notes string) error {
	rec, err := s.getSummaryRecorder(id)
	if err != nil {
		return err
	}
	rec.AddSummaryBottle(carbonation, alcohol, sugar, temp, vol, sugarType, notes)
	return nil
}

// AddSummarySecondary adds a summary of the secondary fermentation
func (s *SummaryRecorderStore) AddSummarySecondary(id string, days int, notes string) error {
	rec, err := s.getSummaryRecorder(id)
	if err != nil {
		return err
	}
	rec.AddSummarySecondary(days, notes)
	return nil
}

// AddTimeline adds a timeline to the summary
func (s *SummaryRecorderStore) AddTimeline(id string, timeline []string) error {
	rec, err := s.getSummaryRecorder(id)
	if err != nil {
		return err
	}
	rec.AddTimeline(timeline)
	return nil
}

// GetSummary returns the summary
func (s *SummaryRecorderStore) GetSummary(id string) (string, error) {
	rec, err := s.getSummaryRecorder(id)
	if err != nil {
		return "", err
	}
	return rec.GetSummary(), nil
}

// GetExtension returns the extension of the summary
func (s *SummaryRecorderStore) GetExtension(id string) (string, error) {
	rec, err := s.getSummaryRecorder(id)
	if err != nil {
		return "", err
	}
	return rec.GetExtension(), nil
}

// Close closes the summary recorder
func (s *SummaryRecorderStore) Close(id string) error {
	rec, err := s.getSummaryRecorder(id)
	if err != nil {
		return err
	}
	rec.Close()
	return nil
}
