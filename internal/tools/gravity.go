package tools

// CorrectRefractometerAlcohol return the apparent gravity in SG correcting refractometer measurements for alcohol
// It applies Terril Linear and Novotny linear correlation models and picks the best one
// Input parameters are
// - OG: Original gravity (uncorrected) in SG
// - MG; Measured gravity (uncorrected) in SG
// - WCF: Wort Correction Factor for the given refractometer
// Based on procedure in https://www.braucampus.at/site/alkoholmessung-mit-dem-refraktometer/
func CorrectRefractometerAlcohol(originalGravity, measuredGravity, wcf float32) float32 {
	ogPlato := SGToPlato(originalGravity)
	mgPlato := SGToPlato(measuredGravity)
	ogCorrected := ogPlato / wcf
	mgCorrected := mgPlato / wcf
	//Terril Linear (result is in sg)
	tlResult := 1.0 - (0.000856829 * ogCorrected) + (0.00349412 * mgCorrected)
	// Novotny Linear (result is in sg)
	nvResult := (-0.002349 * ogCorrected) + (0.006276 * mgCorrected) + 1.0
	resultAvg := (tlResult + nvResult) / 2.0
	if resultAvg < 1.014 {
		// By the end, terril is better
		return tlResult
	}
	return nvResult
}

// CorrectGravityAlcohol returns the real extract in SG (TRE) given the apparent one
// Input parameters are
// - OG: Original gravity (uncorrected) in SG
// - MG; ApparentGravity gravity in SG
// - WCF: Wort Correction Factor for the given refractometer. If using hydrometer 1 can be used
// Based on Balling formula in https://www.maischemalzundmehr.de/index.php?inhaltmitte=toolsrefraktorechner
func CorrectGravityAlcohol(originalGravity, apparentGravity, wcf float32) float32 {
	ogPlato := SGToPlato(originalGravity)
	agPlato := SGToPlato(apparentGravity)
	ogCorrected := ogPlato / wcf
	resultPlato := (0.1808 * ogCorrected) + (0.8192 * agPlato)
	return PlatoToSG(resultPlato)
}
