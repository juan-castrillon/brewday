package sql

import (
	"brewday/internal/summary"
	"brewday/internal/tools"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type SummaryPersistentStore struct {
	dbClient *sql.DB
}

func NewSummaryPersistentStore(db *sql.DB) (*SummaryPersistentStore, error) {
	return &SummaryPersistentStore{
		dbClient: db,
	}, nil
}

// AddSummary adds a summary for the given recipe id with the given title, it also inits an entry in the stats table
func (s *SummaryPersistentStore) AddSummary(recipeID, title string) error {
	_, err := s.dbClient.Exec(`INSERT INTO summaries (title, recipe_id) VALUES (?, ?)`, title, recipeID)
	if err != nil {
		return err
	}
	_, err = s.dbClient.Exec(`INSERT INTO stats (recipe_title) VALUES (?)`, tools.B64Encode(title))
	return err
}

// DeleteSummary deletes the summary for the given recipe id
func (s *SummaryPersistentStore) DeleteSummary(recipeID string) error {
	if recipeID == "" {
		return errors.New("invalid empty recipe id for deleting summary")
	}
	_, err := s.dbClient.Exec(`DELETE FROM summaries WHERE recipe_id == ?`, recipeID)
	return err
}

// AddMashTemp adds a mash temperature to the summary and notes related to it
func (s *SummaryPersistentStore) AddMashTemp(id string, temp float32, notes string) error {
	if id == "" {
		return errors.New("invalid empty recipe id")
	}
	_, err := s.dbClient.Exec(`UPDATE summaries SET mash_temp = ? , mash_notes = ? WHERE recipe_id == ?`, temp, notes, id)
	return err
}

// AddRast adds a rast to the summary and notes related to it
func (s *SummaryPersistentStore) AddRast(id string, temp float32, duration float32, notes string) error {
	if id == "" {
		return errors.New("invalid empty recipe id")
	}
	newRast := summary.MashRastInfo{
		Temperature: temp,
		Time:        duration,
		Notes:       notes,
	}
	newRastBytes, err := json.Marshal(newRast)
	if err != nil {
		return err
	}
	return s.addToMarshalledArray(id, "mash_rasts", string(newRastBytes))
}

// AddLauternNotes adds lautern notes to the summary
func (s *SummaryPersistentStore) AddLauternNotes(id, notes string) error {
	if id == "" {
		return errors.New("invalid empty recipe id")
	}
	_, err := s.dbClient.Exec(`UPDATE summaries SET lautern_info = ? WHERE recipe_id == ?`, notes, id)
	return err
}

// AddHopping adds a hopping to the summary and notes related to it
func (s *SummaryPersistentStore) AddHopping(id string, name string, amount float32, alpha float32, duration float32, notes string) error {
	if id == "" {
		return errors.New("invalid empty recipe id")
	}
	newHop := summary.HopInfo{
		Name:     name,
		Grams:    amount,
		Alpha:    alpha,
		Time:     duration,
		TimeUnit: "minutes",
		Notes:    notes,
	}
	newHopBytes, err := json.Marshal(newHop)
	if err != nil {
		return err
	}
	return s.addToMarshalledArray(id, "hopping_hops", string(newHopBytes))
}

// AddVolumeAfterBoil adds the measured volume before boiling the wort to the summary
func (s *SummaryPersistentStore) AddVolumeBeforeBoil(id string, amount float32, notes string) error {
	if id == "" {
		return errors.New("invalid empty recipe id")
	}
	_, err := s.dbClient.Exec(`UPDATE summaries SET hopping_vol_bb = ? , hopping_vol_bb_notes = ? WHERE recipe_id == ?`, amount, notes, id)
	return err
}

// AddVolumeBeforeBoil adds the measured volume after boiling the wort to the summary
func (s *SummaryPersistentStore) AddVolumeAfterBoil(id string, amount float32, notes string) error {
	if id == "" {
		return errors.New("invalid empty recipe id")
	}
	_, err := s.dbClient.Exec(`UPDATE summaries SET hopping_vol_ab = ? , hopping_vol_ab_notes = ? WHERE recipe_id == ?`, amount, notes, id)
	return err
}

// AddCooling adds a cooling to the summary and notes related to it
func (s *SummaryPersistentStore) AddCooling(id string, finalTemp, coolingTime float32, notes string) error {
	if id == "" {
		return errors.New("invalid empty recipe id")
	}
	_, err := s.dbClient.Exec(`UPDATE summaries SET cooling_temp = ? , cooling_time = ? , cooling_notes = ? WHERE recipe_id == ?`, finalTemp, coolingTime, notes, id)
	return err
}

// AddPreFermentationVolume adds a summary of the pre fermentation
func (s *SummaryPersistentStore) AddPreFermentationVolume(id string, volume float32, sg float32, notes string) error {
	if id == "" {
		return errors.New("invalid empty recipe id")
	}
	newVol := summary.PreFermentationInfo{
		Volume: volume,
		SG:     sg,
		Notes:  notes,
	}
	newVolBytes, err := json.Marshal(newVol)
	if err != nil {
		return err
	}
	return s.addToMarshalledArray(id, "pre_ferm_vols", string(newVolBytes))
}

// AddYeastStart adds the yeast start to the summary
func (s *SummaryPersistentStore) AddYeastStart(id string, temperature, notes string) error {
	if id == "" {
		return errors.New("invalid empty recipe id")
	}
	_, err := s.dbClient.Exec(`UPDATE summaries SET yeast_start_temp = ? , yeast_start_notes = ? WHERE recipe_id == ?`, temperature, notes, id)
	return err
}

// AddMainFermentationSGMeasurement adds a SG measurement to the summary
func (s *SummaryPersistentStore) AddMainFermentationSGMeasurement(id string, date string, gravity float32, final bool, notes string) error {
	if id == "" {
		return errors.New("invalid empty recipe id")
	}
	newSG := summary.SGMeasurement{
		SG:    gravity,
		Date:  date,
		Final: final,
		Notes: notes,
	}
	newSGBytes, err := json.Marshal(newSG)
	if err != nil {
		return err
	}
	return s.addToMarshalledArray(id, "main_ferm_sgs", string(newSGBytes))
}

// AddMainFermentationAlcohol adds the alcohol after the main fermentation to the summary
func (s *SummaryPersistentStore) AddMainFermentationAlcohol(id string, alcohol float32) error {
	if id == "" {
		return errors.New("invalid empty recipe id")
	}
	_, err := s.dbClient.Exec(`UPDATE summaries SET main_ferm_alcohol = ? WHERE recipe_id == ?`, alcohol, id)
	return err
}

func (s *SummaryPersistentStore) AddDryHopStart(id string, name string, amount, alpha float32, notes string) error {
	if id == "" {
		return errors.New("invalid empty recipe id")
	}
	newHop := summary.HopInfo{
		Name:  name,
		Grams: amount,
		Alpha: alpha,
		Notes: notes,
	}
	newHopBytes, err := json.Marshal(newHop)
	if err != nil {
		return err
	}
	return s.addToMarshalledArray(id, "main_ferm_dry_hops", string(newHopBytes))
}
func (s *SummaryPersistentStore) AddDryHopEnd(id string, name string, durationHours float32) error {
	var current sql.NullString
	err := s.dbClient.QueryRow(`SELECT main_ferm_dry_hops FROM summaries WHERE recipe_id == ?`, id).Scan(&current)
	if err != nil {
		return err
	}
	var slice []*summary.HopInfo
	err = json.Unmarshal([]byte(s.sliceFromNullString(current)), &slice)
	if err != nil {
		return err
	}
	found := false
	for i, dh := range slice {
		if dh.Name == name {
			slice[i].Time = durationHours
			slice[i].TimeUnit = "hours"
			found = true
			break
		}
	}
	if !found {
		return errors.New("attempting to end a dry hop that is not started")
	}
	newBytes, err := json.Marshal(slice)
	if err != nil {
		return err
	}
	_, err = s.dbClient.Exec(`UPDATE summaries SET main_ferm_dry_hops = ? WHERE recipe_id == ?`, string(newBytes), id)
	return err
}

// AddPreBottlingVolume adds the volume before bottling
func (s *SummaryPersistentStore) AddPreBottlingVolume(id string, volume float32) error {
	if id == "" {
		return errors.New("invalid empty recipe id")
	}
	_, err := s.dbClient.Exec(`UPDATE summaries SET bottling_pre_bottle_volume = ? WHERE recipe_id == ?`, volume, id)
	return err
}

// AddBottling adds a summary of the bottling
func (s *SummaryPersistentStore) AddBottling(id string, carbonation, alcohol, sugar, temp, vol float32, sugarType, notes string) error {
	if id == "" {
		return errors.New("invalid empty recipe id")
	}
	_, err := s.dbClient.Exec(`UPDATE summaries SET
		bottling_carbonation = ? ,
		bottling_sugar_amount = ? ,
		bottling_sugar_type = ? ,
		bottling_temperature = ? ,
		bottling_alcohol = ? ,
		bottling_volume_bottled = ? ,
		bottling_notes = ?
	WHERE recipe_id == ?`, carbonation, sugar, sugarType, temp, alcohol, vol, notes, id)
	return err
}

// AddSummarySecondary adds a summary of the secondary fermentation
func (s *SummaryPersistentStore) AddSummarySecondary(id string, days int, notes string) error {
	if id == "" {
		return errors.New("invalid empty recipe id")
	}
	_, err := s.dbClient.Exec(`UPDATE summaries SET sec_ferm_days = ?, sec_ferm_notes = ? WHERE recipe_id == ?`, days, notes, id)
	return err
}

func (s *SummaryPersistentStore) getRecipeTitleB64(id string) (string, error) {
	if id == "" {
		return "", errors.New("invalid empty recipe id")
	}
	var title string
	err := s.dbClient.QueryRow(`SELECT title FROM summaries WHERE recipe_id == ?`, id).Scan(&title)
	if err != nil {
		return "", err
	}
	return tools.B64Encode(title), nil
}

// AddEvaporation adds an evaporation to the summary
func (s *SummaryPersistentStore) AddEvaporation(id string, amount float32) error {
	if id == "" {
		return errors.New("invalid empty recipe id")
	}
	t, err := s.getRecipeTitleB64(id)
	if err != nil {
		return err
	}
	_, err = s.dbClient.Exec(`UPDATE stats SET evaporation = ? WHERE recipe_title == ?`, amount, t)
	return err
}

// AddEfficiency adds the efficiency (sudhausausbeute) to the summary
func (s *SummaryPersistentStore) AddEfficiency(id string, efficiencyPercentage float32) error {
	if id == "" {
		return errors.New("invalid empty recipe id")
	}
	t, err := s.getRecipeTitleB64(id)
	if err != nil {
		return err
	}
	_, err = s.dbClient.Exec(`UPDATE stats SET efficiency = ? WHERE recipe_title == ?`, efficiencyPercentage, t)
	return err
}

// GetSummary returns the summary
func (s *SummaryPersistentStore) GetSummary(id string) (*summary.Summary, error) {
	if id == "" {
		return nil, errors.New("invalid empty recipe id")
	}
	var title string
	var mash_notes, mash_rasts, lautern_info, hopping_vol_bb_notes, hopping_hops, hopping_vol_ab_notes, cooling_notes, pre_ferm_vols, yeast_start_temp, yeast_start_notes, main_ferm_sgs, main_ferm_dry_hops, bottling_sugar_type, bottling_notes, sec_ferm_notes sql.NullString
	var mash_temp, hopping_vol_bb, hopping_vol_ab, cooling_temp, cooling_time, main_ferm_alcohol, bottling_pre_bottle_volume, bottling_carbonation, bottling_sugar_amount, bottling_temperature, bottling_alcohol, bottling_volume_bottled, evaporation, efficiency sql.NullFloat64
	var sec_ferm_days sql.NullInt32
	err := s.dbClient.QueryRow(
		`SELECT title, mash_temp, mash_notes, mash_rasts,
		lautern_info, hopping_vol_bb, hopping_vol_bb_notes, hopping_hops,
		hopping_vol_ab, hopping_vol_ab_notes, cooling_temp, cooling_time,
		cooling_notes, pre_ferm_vols, yeast_start_temp, yeast_start_notes,
		main_ferm_sgs, main_ferm_alcohol, main_ferm_dry_hops, bottling_pre_bottle_volume,
		bottling_carbonation, bottling_sugar_amount, bottling_sugar_type, bottling_temperature,
		bottling_alcohol, bottling_volume_bottled, bottling_notes, sec_ferm_days,
		sec_ferm_notes FROM summaries WHERE recipe_id == ?`, id).Scan(
		&title, &mash_temp, &mash_notes, &mash_rasts,
		&lautern_info, &hopping_vol_bb, &hopping_vol_bb_notes, &hopping_hops,
		&hopping_vol_ab, &hopping_vol_ab_notes, &cooling_temp, &cooling_time,
		&cooling_notes, &pre_ferm_vols, &yeast_start_temp, &yeast_start_notes,
		&main_ferm_sgs, &main_ferm_alcohol, &main_ferm_dry_hops, &bottling_pre_bottle_volume,
		&bottling_carbonation, &bottling_sugar_amount, &bottling_sugar_type, &bottling_temperature,
		&bottling_alcohol, &bottling_volume_bottled, &bottling_notes, &sec_ferm_days,
		&sec_ferm_notes,
	)
	if err != nil {
		return nil, err
	}
	err = s.dbClient.QueryRow(`SELECT evaporation, efficiency FROM stats WHERE recipe_title == ?`, tools.B64Encode(title)).Scan(&evaporation, &efficiency)
	if err != nil {
		return nil, err
	}
	var rastInfos []*summary.MashRastInfo
	err = json.Unmarshal([]byte(s.sliceFromNullString(mash_rasts)), &rastInfos)
	if err != nil {
		return nil, err
	}
	var hopInfos, dryHopInfos []*summary.HopInfo
	err = json.Unmarshal([]byte(s.sliceFromNullString(hopping_hops)), &hopInfos)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(s.sliceFromNullString(main_ferm_dry_hops)), &dryHopInfos)
	if err != nil {
		return nil, err
	}
	var preFermentationInfos []*summary.PreFermentationInfo
	err = json.Unmarshal([]byte(s.sliceFromNullString(pre_ferm_vols)), &preFermentationInfos)
	if err != nil {
		return nil, err
	}
	var sgs []*summary.SGMeasurement
	err = json.Unmarshal([]byte(s.sliceFromNullString(main_ferm_sgs)), &sgs)
	if err != nil {
		return nil, err
	}
	return &summary.Summary{
		Title: title,
		MashingInfo: &summary.MashingInfo{
			MashingTemperature: s.valueFromNullFloat(mash_temp),
			MashingNotes:       s.valueFromNullString(mash_notes),
			RastInfos:          rastInfos,
		},
		LauternInfo: s.valueFromNullString(lautern_info),
		HoppingInfo: &summary.HoppingInfo{
			VolBeforeBoil: &summary.VolMeasurement{
				Volume: s.valueFromNullFloat(hopping_vol_bb),
				Notes:  s.valueFromNullString(hopping_vol_bb_notes),
			},
			VolAfterBoil: &summary.VolMeasurement{
				Volume: s.valueFromNullFloat(hopping_vol_ab),
				Notes:  s.valueFromNullString(hopping_vol_ab_notes),
			},
			HopInfos: hopInfos,
		},
		CoolingInfo: &summary.CoolingInfo{
			Temperature: s.valueFromNullFloat(cooling_temp),
			Time:        s.valueFromNullFloat(cooling_time),
			Notes:       s.valueFromNullString(cooling_notes),
		},
		PreFermentationInfos: preFermentationInfos,
		YeastInfo: &summary.YeastInfo{
			Temperature: s.valueFromNullString(yeast_start_temp),
			Notes:       s.valueFromNullString(yeast_start_notes),
		},
		MainFermentationInfo: &summary.MainFermentationInfo{
			SGs:        sgs,
			DryHopInfo: dryHopInfos,
			Alcohol:    s.valueFromNullFloat(main_ferm_alcohol),
		},
		BottlingInfo: &summary.BottlingInfo{
			PreBottleVolume: s.valueFromNullFloat(bottling_pre_bottle_volume),
			Carbonation:     s.valueFromNullFloat(bottling_carbonation),
			SugarAmount:     s.valueFromNullFloat(bottling_sugar_amount),
			SugarType:       s.valueFromNullString(bottling_sugar_type),
			Temperature:     s.valueFromNullFloat(bottling_temperature),
			Alcohol:         s.valueFromNullFloat(bottling_alcohol),
			VolumeBottled:   s.valueFromNullFloat(bottling_volume_bottled),
			Notes:           s.valueFromNullString(bottling_notes),
		},
		SecondaryFermentationInfo: &summary.SecondaryFermentationInfo{
			Days:  s.valueFromNullInt(sec_ferm_days),
			Notes: s.valueFromNullString(sec_ferm_notes),
		},
		Statistics: &summary.Statistics{
			Evaporation: s.valueFromNullFloat(evaporation),
			Efficiency:  s.valueFromNullFloat(efficiency),
		},
	}, nil

}

// GetAllStats returns all the statistics
func (s *SummaryPersistentStore) GetAllStats() (map[string]*summary.Statistics, error) {
	rows, err := s.dbClient.Query(`SELECT recipe_title, evaporation, efficiency, finished_epoch FROM stats`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := make(map[string]*summary.Statistics)
	for rows.Next() {
		r := summary.Statistics{}
		var title string
		var epoch sql.NullInt64
		err = rows.Scan(&title, &r.Evaporation, &r.Efficiency, &epoch)
		if err != nil {
			return nil, err
		}
		r.FinishedTime = time.Unix(s.valueFromNullInt64(epoch), 0)
		titleDecoded, err := tools.B64Decode(title)
		if err != nil {
			return nil, err
		}
		res[titleDecoded] = &r
	}
	return res, nil
}

// AddFinishedTime adds the time when the recipe was done, mainly for statistics
func (s *SummaryPersistentStore) AddFinishedTime(id string, t time.Time) error {
	if id == "" {
		return errors.New("invalid empty recipe id")
	}
	title, err := s.getRecipeTitleB64(id)
	if err != nil {
		return err
	}
	_, err = s.dbClient.Exec(`UPDATE stats SET finished_epoch = ? WHERE recipe_title == ?`, t.Unix(), title)
	return err
}

func (s *SummaryPersistentStore) AddStatsExternal(recipeName string, stats *summary.Statistics) error {
	_, err := s.dbClient.Exec(`INSERT INTO stats (recipe_title, finished_epoch, evaporation, efficiency) VALUES (?, ?, ?, ?)`,
		tools.B64Encode(recipeName),
		stats.FinishedTime.Unix(),
		stats.Evaporation,
		stats.Efficiency,
	)
	return err
}
