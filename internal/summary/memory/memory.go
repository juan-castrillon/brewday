package memory

import (
	"brewday/internal/summary/recorders"
	"errors"
	"sync"
)

type SummaryRecorderMemoryStore struct {
	lock      sync.Mutex
	recorders map[string]recorders.SummaryRecorder
}

func NewSummaryRecorderMemoryStore() *SummaryRecorderMemoryStore {
	return &SummaryRecorderMemoryStore{
		recorders: make(map[string]recorders.SummaryRecorder),
	}
}

// getRecorder returns the summary recorder for the given recipe id
func (s *SummaryRecorderMemoryStore) getRecorder(recipeID string) (recorders.SummaryRecorder, error) {
	sr, ok := s.recorders[recipeID]
	if !ok {
		return nil, errors.New("no summary recorder found for recipe id " + recipeID)
	}
	return sr, nil
}

// AddSummaryRecorder adds a summary recorder for the given recipe id
func (s *SummaryRecorderMemoryStore) AddSummaryRecorder(recipeID string, recorderType string) error {
	summaryRecorder, err := recorders.RecorderFactory(recorderType)
	if err != nil {
		return err
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	s.recorders[recipeID] = summaryRecorder
	return nil
}

// DeleteSummaryRecorder deletes the summary recorder for the given recipe id
func (s *SummaryRecorderMemoryStore) DeleteSummaryRecorder(recipeID string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.recorders, recipeID)
	return nil
}

// AddMashTemp adds a mash temperature to the summary and notes related to it
func (s *SummaryRecorderMemoryStore) AddMashTemp(id string, temp float64, notes string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	rec, err := s.getRecorder(id)
	if err != nil {
		return err
	}
	rec.AddMashTemp(temp, notes)
	return nil
}

// AddRast adds a rast to the summary and notes related to it
func (s *SummaryRecorderMemoryStore) AddRast(id string, temp float64, duration float64, notes string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	rec, err := s.getRecorder(id)
	if err != nil {
		return err
	}
	rec.AddRast(temp, duration, notes)
	return nil
}

// AddLauternNotes adds lautern notes to the summary
func (s *SummaryRecorderMemoryStore) AddLaunternNotes(id, notes string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	rec, err := s.getRecorder(id)
	if err != nil {
		return err
	}
	rec.AddLaunternNotes(notes)
	return nil
}

// AddHopping adds a hopping to the summary and notes related to it
func (s *SummaryRecorderMemoryStore) AddHopping(id string, name string, amount float32, alpha float32, duration float32, notes string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	rec, err := s.getRecorder(id)
	if err != nil {
		return err
	}
	rec.AddHopping(name, amount, alpha, duration, notes)
	return nil
}

// AddMeasuredVolume adds a measured volume to the summary
func (s *SummaryRecorderMemoryStore) AddMeasuredVolume(id string, name string, amount float32, notes string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	rec, err := s.getRecorder(id)
	if err != nil {
		return err
	}
	rec.AddMeasuredVolume(name, amount, notes)
	return nil
}

// AddEvaporation adds an evaporation to the summary
func (s *SummaryRecorderMemoryStore) AddEvaporation(id string, amount float32) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	rec, err := s.getRecorder(id)
	if err != nil {
		return err
	}
	rec.AddEvaporation(amount)
	return nil
}

// AddCooling adds a cooling to the summary and notes related to it
func (s *SummaryRecorderMemoryStore) AddCooling(id string, finalTemp, coolingTime float32, notes string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	rec, err := s.getRecorder(id)
	if err != nil {
		return err
	}
	rec.AddCooling(finalTemp, coolingTime, notes)
	return nil
}

// AddSummaryPreFermentation adds a summary of the pre fermentation
func (s *SummaryRecorderMemoryStore) AddSummaryPreFermentation(id string, volume float32, sg float32, notes string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	rec, err := s.getRecorder(id)
	if err != nil {
		return err
	}
	rec.AddSummaryPreFermentation(volume, sg, notes)
	return nil
}

// AddEfficiency adds the efficiency (sudhausausbeute) to the summary
func (s *SummaryRecorderMemoryStore) AddEfficiency(id string, efficiencyPercentage float32) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	rec, err := s.getRecorder(id)
	if err != nil {
		return err
	}
	rec.AddEfficiency(efficiencyPercentage)
	return nil
}

// AddYeastStart adds the yeast start to the summary
func (s *SummaryRecorderMemoryStore) AddYeastStart(id string, temperature, notes string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	rec, err := s.getRecorder(id)
	if err != nil {
		return err
	}
	rec.AddYeastStart(temperature, notes)
	return nil
}

// AddSGMeasurement adds a SG measurement to the summary
func (s *SummaryRecorderMemoryStore) AddSGMeasurement(id string, date string, gravity float32, final bool, notes string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	rec, err := s.getRecorder(id)
	if err != nil {
		return err
	}
	rec.AddSGMeasurement(date, gravity, final, notes)
	return nil
}

// AddAlcoholMainFermentation adds the alcohol after the main fermentation to the summary
func (s *SummaryRecorderMemoryStore) AddAlcoholMainFermentation(id string, alcohol float32) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	rec, err := s.getRecorder(id)
	if err != nil {
		return err
	}
	rec.AddAlcoholMainFermentation(alcohol)
	return nil
}

// AddSummaryDryHop adds a summary of the dry hop
func (s *SummaryRecorderMemoryStore) AddSummaryDryHop(id string, name string, amount float32) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	rec, err := s.getRecorder(id)
	if err != nil {
		return err
	}
	rec.AddSummaryDryHop(name, amount)
	return nil
}

// AddSummaryPreBottle adds a summary of the pre bottling
func (s *SummaryRecorderMemoryStore) AddSummaryPreBottle(id string, volume float32) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	rec, err := s.getRecorder(id)
	if err != nil {
		return err
	}
	rec.AddSummaryPreBottle(volume)
	return nil
}

// AddSummaryBottle adds a summary of the bottling
func (s *SummaryRecorderMemoryStore) AddSummaryBottle(id string, carbonation, alcohol, sugar, temp, vol float32, sugarType, notes string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	rec, err := s.getRecorder(id)
	if err != nil {
		return err
	}
	rec.AddSummaryBottle(carbonation, alcohol, sugar, temp, vol, sugarType, notes)
	return nil
}

// AddSummarySecondary adds a summary of the secondary fermentation
func (s *SummaryRecorderMemoryStore) AddSummarySecondary(id string, days int, notes string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	rec, err := s.getRecorder(id)
	if err != nil {
		return err
	}
	rec.AddSummarySecondary(days, notes)
	return nil
}

// AddTimeline adds a timeline to the summary
func (s *SummaryRecorderMemoryStore) AddTimeline(id string, timeline []string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	rec, err := s.getRecorder(id)
	if err != nil {
		return err
	}
	rec.AddTimeline(timeline)
	return nil
}

// GetSummary returns the summary
func (s *SummaryRecorderMemoryStore) GetSummary(id string) (string, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	rec, err := s.getRecorder(id)
	if err != nil {
		return "", err
	}
	return rec.GetSummary(), nil
}

// GetExtension returns the extension of the summary
func (s *SummaryRecorderMemoryStore) GetExtension(id string) (string, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	rec, err := s.getRecorder(id)
	if err != nil {
		return "", err
	}
	return rec.GetExtension(), nil
}

// Close closes the summary recorder
func (s *SummaryRecorderMemoryStore) Close(id string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	rec, err := s.getRecorder(id)
	if err != nil {
		return err
	}
	rec.Close()
	return nil
}
