package tools

// CalculateEfficiencyPlato returns the efficiency of a brew
// It takes the measured gravity and measured volume, as well as the total amount of malt
// The measured gravity is in plato, the measured volume in liters and the total malt in grams
// It returns the efficiency in percent
func CalculateEfficiencyPlato(measuredGravity, measuredVolume, totalMalt float32) float32 {
	specificGravity := (measuredGravity / (258.6 - ((measuredGravity / 258.2) * 227.1))) + 1
	return 1000 * measuredVolume * specificGravity * 0.96 * measuredGravity / totalMalt
}

// CalculateEfficiencySG returns the efficiency of a brew
// It takes the measured gravity and measured volume, as well as the total amount of malt
// The measured gravity is in SG, the measured volume in liters and the total malt in grams
// It returns the efficiency in percent
func CalculateEfficiencySG(measuredGravity, measuredVolume, totalMalt float32) float32 {
	gPlato := SGToPlato(measuredGravity)
	return CalculateEfficiencyPlato(gPlato, measuredVolume, totalMalt)
}
