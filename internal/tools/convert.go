package tools

// SGToPlato converts a specific gravity to a plato value
func SGToPlato(sg float32) float32 {
	if sg > 1.000 {
		return (135.997 * sg * sg * sg) - (630.272 * sg * sg) + (1111.14 * sg) - 616.868
	} else {
		return 0.0
	}
}

// PlatoToSG converts a plato value to a specific gravity
func PlatoToSG(plato float32) float32 {
	return 1 + (plato / (258.6 - ((plato / 258.2) * 227.1)))
}

// EBCtoSRM converts EBC to SRM
func EBCtoSRM(ebc float32) float32 {
	return ebc / 1.97
}

// SRMtoEBC converts SRM to EBC
func SRMtoEBC(srm float32) float32 {
	return srm * 1.97
}
