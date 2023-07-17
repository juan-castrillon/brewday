package recipe

// Recipe is the main struct for a recipe
type Recipe struct {
	// Name is the name of the recipe
	Name string
	// Style of the beer
	Style string
	// BatchSize is the size of the batch in liters
	BatchSize float32
	// InitialSG is the initial specific gravity (in SG)
	InitialSG float32
	// Bitterness is the bitterness in IBU
	Bitterness float32
	// ColorEBC is the color in EBC
	ColorEBC float32
	// Mashing is the mashing instructions
	Mashing MashInstructions
	// Hopping is the hopping instructions
	Hopping HopInstructions
	// Fermentation is the fermentation instructions
	Fermentation FermentationInstructions
}

// MashInstructions is the struct for the mashing instructions
// It contains the malts, the main water volume, the mash temperature, the mash out temperature and the rasts
type MashInstructions struct {
	// List of malts to use
	Malts []Malt
	// MainWaterVolume is the main water volume in liters
	MainWaterVolume float32
	// MashTemperature is the mash temperature in °C
	MashTemperature float32
	// MashOutTemperature is the mash out temperature in °C
	MashOutTemperature float32
	// Rasts is the list of rasts to perform
	Rasts []Rast
}

// Rast is the struct for a rast
// It represent maintaining a temperature for a given duration
type Rast struct {
	// Temperature is the temperature in °C
	Temperature float32
	// Duration is the duration in minutes
	Duration float32
}

// Malt is the struct for a malt
// It contains the name and the amount in grams
type Malt struct {
	// Name of the malt
	Name string
	// Amount in grams
	Amount float32
}

type HopInstructions struct {
	Hops                  []Hops
	AdditionalIngredients []AdditionalIngredient
}

type Hops struct {
	Name     string
	Alpha    float32
	Amount   float32
	Duration float32
	DryHop   bool
}

type AdditionalIngredient struct {
	Name     string
	Amount   float32
	Duration float32
}

type FermentationInstructions struct {
	Yeast                 Yeast
	Temperature           float32
	AdditionalIngredients []AdditionalIngredient
	Carbonation           float32
}

type Yeast struct {
	Name   string
	Amount float32
}
