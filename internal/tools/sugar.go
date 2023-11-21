package tools

import (
	"math"
)

// SugarForCarbonation calculates the amount of sugar needed for a certain carbonation level
// It assumes that the sugar will be made into a solution with water
// It supports different types of sugar:
// - glucose (Traubenzucker)
// - sucrose (Haushaltszucker)
// It returns the amount of sugar in grams and the estimated final alcohol content
func SugarForCarbonation(volume, carbonation, temperature, alcoholBefore, waterVolume float32, sugarType string) (float32, float32) {
	co2Present := co2InBeer(temperature)
	co2Needed := carbonation - float32(co2Present)
	unitAlcohol, unitCO2 := addedBySugar(sugarType)
	// Calculate the amount of sugar needed to reach the desired CO2 level
	sugar := co2Needed * volume / unitCO2
	// Calculate the amount of alcohol added by the sugar
	alcohol := sugar * unitAlcohol // this is grams of alcohol
	// Calculate final alcohol content based on added alcohol and alcohol present in beer and water in the sugar solution
	ethanolDensity := float32(789)                            // g/l
	alcoholBeforeLiters := alcoholBefore * volume / 100       // Liters of alcohol present in beer
	alcoholAddedLiters := alcohol / ethanolDensity            // Liters of alcohol added by sugar
	alcoholLiters := alcoholBeforeLiters + alcoholAddedLiters // Total liters of alcohol
	waterBefore := volume - alcoholBeforeLiters
	totalWater := waterBefore + waterVolume
	totalVolume := totalWater + alcoholLiters
	finalAlcohol := alcoholLiters / totalVolume * 100
	return sugar, finalAlcohol
}

// CO2InBeer calculates the amount of CO2 in beer at a certain temperature
// Its uses the Henry's law implementation in braureka.de
func co2InBeer(temp float32) float64 {
	return 10.13 * math.Exp(float64(-10.73797+2617.25/(temp+273.5)))
}

// addedBySugar calculates the grams of alcohol and co2 added by a gram of sugar in beer
// This is taken from lots of sources including:
// https://braureka.de/calculations/carbonation/ Uses 0.5 and 86,2% yield for glucose
// https://gradplato.com/kategorien/know-how/flaschengaerung-fuer-alle Use 0.5
// http://www.brsquared.org/wine/CalcInfo/HydSugAl.htm suggests 0.47
// https://homedistiller.org/wiki/htm/calcs/calcs_alcohol_yield.html) suggests 0.48 and 88%
// Brewer's Friend suggests 91.03% yield for glucose
// MMUM suggests 91.03% yield for glucose
// MaltMiller suggests 91% yield for glucose
// https://wiki.homebrewtalk.com/index.php/Fermentable_adjuncts#Sugar_Adjuncts suggest 90% yield for glucose
// MashCamp suggests 85.9% yield for glucose
// Fabier suggests 86.31% yield for glucose
// Northern Brewer suggests 91% yield for glucose
// Based on all this, this method uses 0.48 for the grams of alcohol added by a gram of sugar
// And a yield of 90% for glucose
func addedBySugar(sugar string) (float32, float32) {
	var alcoholGrams, co2Grams float32
	const sugarAlcohol = 0.48
	switch sugar {
	case "sucrose":
		alcoholGrams = sugarAlcohol
	case "glucose":
		alcoholGrams = sugarAlcohol * 0.9
	default:
		// Handle other types of sugar or provide an error message
		alcoholGrams = 0
	}

	co2Grams = alcoholGrams // Assuming 1:1 ratio for CO2 to ethanol

	return alcoholGrams, co2Grams
}

// CarbonationForSugar calculates the carbonation level for a certain amount of sugar
// It assumes that the sugar will be made into a solution with water
func CarbonationForSugar(volume, sugar, temperature float32, sugarType string) float32 {
	co2Present := co2InBeer(temperature)
	_, unitCO2 := addedBySugar(sugarType)
	// Calculate how much CO2 is added by the sugar
	co2Created := sugar * unitCO2 / volume // g\L

	return co2Created + float32(co2Present)
}
