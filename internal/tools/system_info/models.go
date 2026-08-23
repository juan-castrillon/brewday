package systeminfo

// SystemProperties represent a brewing system. SImilar to config/BrewingSystemConfig it might differ in the future thus the decouple
type SystemProperties struct {
	LD        float32
	UD        float32
	Power     int
	MaxVol    float32
	MaxHeight float32
}
