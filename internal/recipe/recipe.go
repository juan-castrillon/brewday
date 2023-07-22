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
	// Nachguss is the nachguss volume in liters
	Nachguss float32
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

// HopInstructions is the struct for the hopping instructions
// It contains the hops and the additional ingredients (like spices)
type HopInstructions struct {
	// Hops is the list of hops to use
	Hops []Hops
	// AdditionalIngredients is the list of additional ingredients to use in the boil
	AdditionalIngredients []AdditionalIngredient
}

// Hops is the struct for a hop
// It contains the name, the alpha acid percentage, the amount in grams, the duration in minutes and if it is a dry hop
type Hops struct {
	// Name of the hop
	Name string
	// Alpha is the alpha acid percentage
	Alpha float32
	// Amount in grams
	Amount float32
	// Duration of cooking in minutes
	Duration float32
	// DryHop is true if this hop is for dry hopping
	DryHop bool
}

// AdditionalIngredient is the struct for an additional ingredient
// It contains the name, the amount in grams and the duration in minutes
// It can represent spices, sugar, fruits, ...
// It is used in the boil or in the fermentation
type AdditionalIngredient struct {
	// Name of the additional ingredient
	Name string
	// Amount in grams
	Amount float32
	// Duration in minutes
	Duration float32
}

// FermentationInstructions is the struct for the fermentation instructions
// It contains the yeast, the temperature, the additional ingredients and the carbonation in g/l
type FermentationInstructions struct {
	// Yeast is the yeast to use
	Yeast Yeast
	// Temperature is the fermentation temperature in °C
	// It can be a range of temperature (e.g. 18-22°C)
	Temperature string
	// AdditionalIngredients is the list of additional ingredients to use in the fermentation
	AdditionalIngredients []AdditionalIngredient
	// Carbonation is the carbonation in g/l
	Carbonation float32
}

// Yeast is the struct for a yeast
// It contains the name and the amount in grams
type Yeast struct {
	// Name of the yeast
	Name string
	// Amount in grams
	Amount float32
}