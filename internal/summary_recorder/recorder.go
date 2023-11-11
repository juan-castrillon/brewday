package summaryrecorder

type SummaryRecorder interface {
	// AddMashTemp adds a mash temperature to the summary and notes related to it
	AddMashTemp(temp float64, notes string)
	// AddRast adds a rast to the summary and notes related to it
	AddRast(temp float64, duration float64, notes string)
	// // AddLauternNotes adds lautern notes to the summary
	AddLaunternNotes(notes string)
	// AddHopping adds a hopping to the summary and notes related to it
	AddHopping(name string, amount float32, alpha float32, duration float32, notes string)
	// AddMeasuredVolume adds a measured volume to the summary
	AddMeasuredVolume(name string, amount float32, notes string)
	// AddEvaporation adds an evaporation to the summary
	AddEvaporation(amount float32)
	// AddCooling adds a cooling to the summary and notes related to it
	AddCooling(finalTemp, coolingTime float32, notes string)
	// AddSummaryPreFermentation adds a summary of the pre fermentation
	AddSummaryPreFermentation(volume float32, sg float32, notes string)
	// AddEfficiency adds the efficiency (sudhausausbeute) to the summary
	AddEfficiency(efficiencyPercentage float32)
	// AddYeastStart adds the yeast start to the summary
	AddYeastStart(temperature, notes string)
	// AddSGMeasurement adds a SG measurement to the summary
	AddSGMeasurement(date string, gravity float32, final bool, notes string)
	// AddSummaryDryHop adds a summary of the dry hop
	AddSummaryDryHop(name string, amount float32)
	// Close closes the summary recorder
	Close()
	// GetSummary returns the summary
	GetSummary() string
	// GetExtension returns the extention of the summary
	GetExtension() string
	// AddTimeline adds a timeline to the summary
	AddTimeline(timeline []string)
	// AddAlcoholMainFermentation adds the alcohol after the main fermentation to the summary
	AddAlcoholMainFermentation(alcohol float32)
}
