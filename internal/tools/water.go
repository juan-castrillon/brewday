package tools

// WaterForVolume returns the amount of water that needs to be added to a volume to reach a target volume
// It also returns the approximate gravity of the final solution
// Formula is based on Ray Daniels' "Designing Great Beers"
// It uses gravity in SG units and calculates gravity points, then applies the formula
func WaterForVolume(current, target, currentSG float32) (float32, float32) {
	toAdd := target - current
	finalSG := ((current / target) * (currentSG - 1)) + 1
	return toAdd, finalSG
}

// WaterForGravity returns the amount of water that needs to be added to a volume to reach a target gravity
// It also returns the volume of the final solution
func WaterForGravity(current, target, currentVolume float32) (float32, float32) {
	targetVol := currentVolume * (current - 1) / (target - 1)
	toAdd := targetVol - currentVolume
	return toAdd, targetVol
}
