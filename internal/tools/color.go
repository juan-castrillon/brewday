package tools

import (
	"math"
)

// SRMHex is a lookup table for SRM to hex color code conversion.
// The index is the SRM value, the value is the hex color code.
var SRMHex = []string{
	"#FFE699",
	"#FFD878",
	"#FFCA5A",
	"#FFBF42",
	"#FBB123",
	"#F8A600",
	"#F39C00",
	"#EA8F00",
	"#E58500",
	"#DE7C00",
	"#D77200",
	"#CF6900",
	"#CB6200",
	"#C35900",
	"#BB5100",
	"#B54C00",
	"#B04500",
	"#A63E00",
	"#A13700",
	"#9B3200",
	"#952D00",
	"#8E2900",
	"#882300",
	"#821E00",
	"#7B1A00",
	"#771900",
	"#701400",
	"#6A0E00",
	"#660D00",
	"#5E0B00",
	"#5A0A02",
	"#560A05",
	"#520907",
	"#4C0505",
	"#470606",
	"#420607",
	"#3D0708",
	"#370607",
	"#2D0607",
	"#1F0506",
	"#1A0404",
	"#160404",
	"#120303",
	"#0E0202",
	"#0A0202",
	"#080101",
	"#060101",
	"#040100",
	"#020100",
	"#000000",
}

// SRMToHex converts SRM value to hex color code using a lookup table.
func SRMToHex(srm float64) string {
	// Round the SRM value to the nearest integer.
	roundFloat := math.Round(srm)
	roundInt := int(roundFloat)
	return SRMHex[roundInt-1]
}

// EBCtoHex converts EBC value to hex color code using a simplified linear conversion.
func EBCtoHex(ebc float32) string {
	srm := EBCtoSRM(ebc)
	return SRMToHex(float64(srm))
}
