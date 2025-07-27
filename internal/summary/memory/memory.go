package memory

import (
	"brewday/internal/summary"
	"brewday/internal/tools"
	"errors"
	"sync"
	"time"
)

type SummaryMemoryStore struct {
	lock      sync.Mutex
	summaries map[string]*summary.Summary
	stats     map[string]*summary.Statistics // In here stats is a backup, in case the summary is deleted
}

func NewSummaryMemoryStore() *SummaryMemoryStore {
	return &SummaryMemoryStore{
		summaries: make(map[string]*summary.Summary),
		stats:     make(map[string]*summary.Statistics),
	}
}

// getSummary returns the summary for the given recipe id
func (s *SummaryMemoryStore) getSummary(recipeID string) (*summary.Summary, error) {
	sr, ok := s.summaries[recipeID]
	if !ok {
		return nil, errors.New("no summary recorder found for recipe id " + recipeID)
	}
	return sr, nil
}

// AddSummary adds a summary for the given recipe id with the given title
func (s *SummaryMemoryStore) AddSummary(recipeID, title string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	summ := summary.NewSummary()
	summ.Title = title
	s.summaries[recipeID] = summ
	s.stats[tools.B64Encode(title)] = summ.Statistics
	return nil
}

// DeleteSummary deletes the summary for the given recipe id
func (s *SummaryMemoryStore) DeleteSummary(recipeID string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.summaries, recipeID)
	return nil
}

// AddMashTemp adds a mash temperature to the summary and notes related to it
func (s *SummaryMemoryStore) AddMashTemp(id string, temp float32, notes string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	sum, err := s.getSummary(id)
	if err != nil {
		return err
	}
	if sum.MashingInfo == nil {
		sum.MashingInfo = &summary.MashingInfo{}
	}
	sum.MashingInfo.MashingTemperature = temp
	sum.MashingInfo.MashingNotes = notes
	return nil
}

// AddRast adds a rast to the summary and notes related to it
func (s *SummaryMemoryStore) AddRast(id string, temp float32, duration float32, notes string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	sum, err := s.getSummary(id)
	if err != nil {
		return err
	}
	if sum.MashingInfo == nil {
		sum.MashingInfo = &summary.MashingInfo{}
	}
	if sum.MashingInfo.RastInfos == nil {
		sum.MashingInfo.RastInfos = make([]*summary.MashRastInfo, 0)
	}
	sum.MashingInfo.RastInfos = append(sum.MashingInfo.RastInfos, &summary.MashRastInfo{
		Temperature: temp, Time: duration, Notes: notes,
	})
	return nil
}

// AddLauternNotes adds lautern notes to the summary
func (s *SummaryMemoryStore) AddLauternNotes(id, notes string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	sum, err := s.getSummary(id)
	if err != nil {
		return err
	}
	sum.LauternInfo = notes
	return nil
}

// AddHopping adds a hopping to the summary and notes related to it
func (s *SummaryMemoryStore) AddHopping(id string, name string, amount float32, alpha float32, duration float32, notes string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	sum, err := s.getSummary(id)
	if err != nil {
		return err
	}
	if sum.HoppingInfo == nil {
		sum.HoppingInfo = &summary.HoppingInfo{}
	}
	if sum.HoppingInfo.HopInfos == nil {
		sum.HoppingInfo.HopInfos = make([]*summary.HopInfo, 0)
	}
	sum.HoppingInfo.HopInfos = append(sum.HoppingInfo.HopInfos, &summary.HopInfo{
		Name:     name,
		Grams:    amount,
		Alpha:    alpha,
		Time:     duration,
		TimeUnit: "minutes",
		Notes:    notes,
	})
	return nil
}

// AddVolumeBeforeBoil adds the measured volume before boiling the wort to the summary
func (s *SummaryMemoryStore) AddVolumeBeforeBoil(id string, amount float32, notes string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	sum, err := s.getSummary(id)
	if err != nil {
		return err
	}
	if sum.HoppingInfo == nil {
		sum.HoppingInfo = &summary.HoppingInfo{}
	}
	if sum.HoppingInfo.VolBeforeBoil == nil {
		sum.HoppingInfo.VolBeforeBoil = &summary.VolMeasurement{}
	}
	sum.HoppingInfo.VolBeforeBoil.Volume = amount
	sum.HoppingInfo.VolBeforeBoil.Notes = notes
	return nil
}

// AddVolumeAfterBoil adds the measured volume after boiling the wort to the summary
func (s *SummaryMemoryStore) AddVolumeAfterBoil(id string, amount float32, notes string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	sum, err := s.getSummary(id)
	if err != nil {
		return err
	}
	if sum.HoppingInfo == nil {
		sum.HoppingInfo = &summary.HoppingInfo{}
	}
	if sum.HoppingInfo.VolAfterBoil == nil {
		sum.HoppingInfo.VolAfterBoil = &summary.VolMeasurement{}
	}
	sum.HoppingInfo.VolAfterBoil.Volume = amount
	sum.HoppingInfo.VolAfterBoil.Notes = notes
	return nil
}

// AddCooling adds a cooling to the summary and notes related to it
func (s *SummaryMemoryStore) AddCooling(id string, finalTemp, coolingTime float32, notes string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	sum, err := s.getSummary(id)
	if err != nil {
		return err
	}
	if sum.CoolingInfo == nil {
		sum.CoolingInfo = &summary.CoolingInfo{}
	}
	sum.CoolingInfo.Notes = notes
	sum.CoolingInfo.Temperature = finalTemp
	sum.CoolingInfo.Time = coolingTime
	return nil
}

// AddPreFermentationVolume adds a summary of the pre fermentation
func (s *SummaryMemoryStore) AddPreFermentationVolume(id string, volume float32, sg float32, notes string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	sum, err := s.getSummary(id)
	if err != nil {
		return err
	}
	if sum.PreFermentationInfos == nil {
		sum.PreFermentationInfos = make([]*summary.PreFermentationInfo, 0)
	}
	sum.PreFermentationInfos = append(sum.PreFermentationInfos, &summary.PreFermentationInfo{
		Volume: volume,
		SG:     sg,
		Notes:  notes,
	})
	return nil
}

// AddYeastStart adds the yeast start to the summary
func (s *SummaryMemoryStore) AddYeastStart(id string, temperature, notes string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	sum, err := s.getSummary(id)
	if err != nil {
		return err
	}
	if sum.YeastInfo == nil {
		sum.YeastInfo = &summary.YeastInfo{}
	}
	sum.YeastInfo.Notes = notes
	sum.YeastInfo.Temperature = temperature
	return nil
}

// AddMainFermentationSGMeasurement adds a SG measurement to the summary
func (s *SummaryMemoryStore) AddMainFermentationSGMeasurement(id string, date string, gravity float32, final bool, notes string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	sum, err := s.getSummary(id)
	if err != nil {
		return err
	}
	if sum.MainFermentationInfo == nil {
		sum.MainFermentationInfo = &summary.MainFermentationInfo{}
	}
	if sum.MainFermentationInfo.SGs == nil {
		sum.MainFermentationInfo.SGs = make([]*summary.SGMeasurement, 0)
	}
	sum.MainFermentationInfo.SGs = append(sum.MainFermentationInfo.SGs, &summary.SGMeasurement{
		SG:    gravity,
		Date:  date,
		Final: final,
		Notes: notes,
	})
	return nil
}

// AddMainFermentationAlcohol adds the alcohol after the main fermentation to the summary
func (s *SummaryMemoryStore) AddMainFermentationAlcohol(id string, alcohol float32) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	sum, err := s.getSummary(id)
	if err != nil {
		return err
	}
	if sum.MainFermentationInfo == nil {
		sum.MainFermentationInfo = &summary.MainFermentationInfo{}
	}
	sum.MainFermentationInfo.Alcohol = alcohol
	return nil
}

func (s *SummaryMemoryStore) AddDryHopStart(id string, name string, amount, alpha float32, notes string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	sum, err := s.getSummary(id)
	if err != nil {
		return err
	}
	if sum.MainFermentationInfo == nil {
		sum.MainFermentationInfo = &summary.MainFermentationInfo{}
	}
	if sum.MainFermentationInfo.DryHopInfo == nil {
		sum.MainFermentationInfo.DryHopInfo = make([]*summary.HopInfo, 0)
	}
	sum.MainFermentationInfo.DryHopInfo = append(sum.MainFermentationInfo.DryHopInfo, &summary.HopInfo{
		Name:  name,
		Grams: amount,
		Alpha: alpha,
		Notes: notes,
	})
	return nil
}
func (s *SummaryMemoryStore) AddDryHopEnd(id string, name string, durationHours float32) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	sum, err := s.getSummary(id)
	if err != nil {
		return err
	}
	if sum.MainFermentationInfo == nil || sum.MainFermentationInfo.DryHopInfo == nil {
		return errors.New("attempting to end a dry hop that is not started")
	}
	found := false
	for i, dh := range sum.MainFermentationInfo.DryHopInfo {
		if dh.Name == name {
			sum.MainFermentationInfo.DryHopInfo[i].Time = durationHours
			sum.MainFermentationInfo.DryHopInfo[i].TimeUnit = "hours"
			found = true
			break
		}
	}
	if !found {
		return errors.New("attempting to end a dry hop that is not started")
	}
	return nil
}

// AddPreBottlingVolume adds the volume before bottling
func (s *SummaryMemoryStore) AddPreBottlingVolume(id string, volume float32) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	sum, err := s.getSummary(id)
	if err != nil {
		return err
	}
	if sum.BottlingInfo == nil {
		sum.BottlingInfo = &summary.BottlingInfo{}
	}
	sum.BottlingInfo.PreBottleVolume = volume
	return nil
}

// AddBottling adds a summary of the bottling
func (s *SummaryMemoryStore) AddBottling(id string, carbonation, alcohol, sugar, temp, vol float32, sugarType, notes string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	sum, err := s.getSummary(id)
	if err != nil {
		return err
	}
	if sum.BottlingInfo == nil {
		sum.BottlingInfo = &summary.BottlingInfo{}
	}
	sum.BottlingInfo.Alcohol = alcohol
	sum.BottlingInfo.Carbonation = carbonation
	sum.BottlingInfo.SugarAmount = sugar
	sum.BottlingInfo.SugarType = sugarType
	sum.BottlingInfo.Temperature = temp
	sum.BottlingInfo.VolumeBottled = vol
	sum.BottlingInfo.Notes = notes
	return nil
}

// AddSummarySecondary adds a summary of the secondary fermentation
func (s *SummaryMemoryStore) AddSummarySecondary(id string, days int, notes string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	sum, err := s.getSummary(id)
	if err != nil {
		return err
	}
	if sum.SecondaryFermentationInfo == nil {
		sum.SecondaryFermentationInfo = &summary.SecondaryFermentationInfo{}
	}
	sum.SecondaryFermentationInfo.Days = days
	sum.SecondaryFermentationInfo.Notes = notes
	return nil
}

// AddEvaporation adds an evaporation to the summary
func (s *SummaryMemoryStore) AddEvaporation(id string, amount float32) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	sum, err := s.getSummary(id)
	if err != nil {
		return err
	}
	if sum.Statistics == nil {
		sum.Statistics = &summary.Statistics{}
	}
	sum.Statistics.Evaporation = amount
	s.stats[tools.B64Encode(sum.Title)] = sum.Statistics
	return nil
}

// AddEfficiency adds the efficiency (sudhausausbeute) to the summary
func (s *SummaryMemoryStore) AddEfficiency(id string, efficiencyPercentage float32) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	sum, err := s.getSummary(id)
	if err != nil {
		return err
	}
	if sum.Statistics == nil {
		sum.Statistics = &summary.Statistics{}
	}
	sum.Statistics.Efficiency = efficiencyPercentage
	s.stats[tools.B64Encode(sum.Title)] = sum.Statistics
	return nil
}

// GetSummary returns the summary
func (s *SummaryMemoryStore) GetSummary(id string) (*summary.Summary, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	sum, err := s.getSummary(id)
	if err != nil {
		return nil, err
	}
	return sum, nil
}

// GetAllStats returns all the statistics
func (s *SummaryMemoryStore) GetAllStats() (map[string]*summary.Statistics, error) {
	res := map[string]*summary.Statistics{}
	for nb64, v := range s.stats {
		decoded, err := tools.B64Decode(nb64)
		if err != nil {
			return nil, err
		}
		res[decoded] = v
	}
	return res, nil
}

// AddFinishedTime adds the time when the recipe was done, mainly for statistics
func (s *SummaryMemoryStore) AddFinishedTime(id string, t time.Time) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	sum, err := s.getSummary(id)
	if err != nil {
		return err
	}
	if sum.Statistics == nil {
		sum.Statistics = &summary.Statistics{}
	}
	sum.Statistics.FinishedTime = t
	s.stats[tools.B64Encode(sum.Title)] = sum.Statistics
	return nil
}

func (s *SummaryMemoryStore) AddStats(recipeName string, stats *summary.Statistics) error {
	s.stats[tools.B64Encode(recipeName)] = stats
	return nil
}
