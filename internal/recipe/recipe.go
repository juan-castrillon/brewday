package recipe

import "sync"

type RecipeStatus int
type ResultType int

const (
	RecipeStatusUnknown RecipeStatus = iota
	RecipeStatusCreated
	RecipeStatusMashing
	RecipeStatusLautering
	RecipeStatusBoiling
	RecipeStatusCooling
	RecipeStatusPreFermentation
	RecipeStatusFermenting
	RecipeStatusBottled
	RecipeStatusFridge
	RecipeStatusFinished
	ResultHotWortVolume ResultType = iota
	ResultOriginalGravity
	ResultFinalGravity
	ResultAlcohol
	ResultMainFermentationVolume
	ResultVolumeBeforeBoil
)

type SGMeasurement struct {
	Value float32 `json:"value,omitempty"`
	Date  string  `json:"date,omitempty"`
}

type RecipeResults struct {
	VolumeBeforeBoil       float32
	HotWortVolume          float32
	OriginalGravity        float32
	FinalGravity           float32
	Alcohol                float32
	MainFermentationVolume float32
}

// Recipe is the main struct for a recipe
type Recipe struct {
	// Name is the name of the recipe
	Name string `json:"Name"`
	// Style of the beer
	Style string `json:"Style"`
	// BatchSize is the size of the batch in liters
	BatchSize float32 `json:"BatchSizeL"`
	// InitialSG is the initial specific gravity (in SG)
	InitialSG float32 `json:"InitialSG"`
	// Bitterness is the bitterness in IBU
	Bitterness float32 `json:"IBU"`
	// ColorEBC is the color in EBC
	ColorEBC float32 `json:"ColorEBC"`
	// Mashing is the mashing instructions
	Mashing MashInstructions `json:"Mashing"`
	// Hopping is the hopping instructions
	Hopping HopInstructions `json:"Hopping"`
	// Fermentation is the fermentation instructions
	Fermentation FermentationInstructions `json:"Fermentation"`
	// status is the status of the recipe
	status RecipeStatus `json:"-"`
	// statusParams is the parameters for the status
	// This is particular for each status, and it can be empty or things like rast number, hop number, ...
	// It must always be a string to avoid having to use reflection when storing the recipe.
	statusParams []string `json:"-"`
	// statusLock is the lock for the status
	statusLock sync.Mutex `json:"-"`
	// ID is the identifier of the recipe. This is populated by the appropriate store and should not be set manually
	ID string `json:"ID"`
	// results is the results of the recipe. This is populated by the appropriate handlers and should not be set manually
	results RecipeResults `json:"-"`
	// mainFermSGs is the measured specific gravities during the main fermentation process. This is populated by the appropriate handlers and should not be set manually
	mainFermSGs []*SGMeasurement `json:"-"`
	// resultsLock is the lock for the results
	resultsLock sync.Mutex `json:"-"`
	// mainFermSGsLock is the lock for the mainFermSGs
	mainFermSGsLock sync.Mutex `json:"-"`
}

// MashInstructions is the struct for the mashing instructions
// It contains the malts, the main water volume, the mash temperature, the mash out temperature and the rasts
type MashInstructions struct {
	// List of malts to use
	Malts []Malt `json:"Malts"`
	// MainWaterVolume is the main water volume in liters
	MainWaterVolume float32 `json:"MainWaterVolume"`
	// Nachguss is the nachguss volume in liters
	Nachguss float32 `json:"Nachguss"`
	// MashTemperature is the mash temperature in °C
	MashTemperature float32 `json:"MashTemperature"`
	// MashOutTemperature is the mash out temperature in °C
	MashOutTemperature float32 `json:"MashOutTemperature"`
	// Rasts is the list of rasts to perform
	Rasts []Rast `json:"Rasts"`
}

// Rast is the struct for a rast
// It represent maintaining a temperature for a given duration
type Rast struct {
	// Temperature is the temperature in °C
	Temperature float32 `json:"Temperature"`
	// Duration is the duration in minutes
	Duration float32 `json:"Duration"`
}

// Malt is the struct for a malt
// It contains the name and the amount in grams
type Malt struct {
	// Name of the malt
	Name string `json:"Name"`
	// Amount in grams
	Amount float32 `json:"Amount"`
}

// HopInstructions is the struct for the hopping instructions
// It contains the hops and the additional ingredients (like spices)
type HopInstructions struct {
	// TotalCookingTime is the total cooking time in minutes
	TotalCookingTime float32 `json:"TotalCookingTime"`
	// Hops is the list of hops to use
	Hops []Hops `json:"Hops"`
	// AdditionalIngredients is the list of additional ingredients to use in the boil
	AdditionalIngredients []AdditionalIngredient `json:"HoppingAdditionalIngredients"`
}

// Hops is the struct for a hop
// It contains the name, the alpha acid percentage, the amount in grams, the duration in minutes and if it is a dry hop
type Hops struct {
	// Name of the hop
	Name string `json:"Name"`
	// Alpha is the alpha acid percentage
	Alpha float32 `json:"Alpha"`
	// Amount in grams
	Amount float32 `json:"Amount"`
	// Duration of cooking in minutes
	Duration float32 `json:"Duration"`
	// DryHop is true if this hop is for dry hopping
	DryHop bool `json:"DryHop"`
	// Vorderwürze is true if this hop is for vorderwürze hopping
	Vorderwuerze bool `json:"Vorderwuerze"`
}

// AdditionalIngredient is the struct for an additional ingredient
// It contains the name, the amount in grams and the duration in minutes
// It can represent spices, sugar, fruits, ...
// It is used in the boil or in the fermentation
type AdditionalIngredient struct {
	// Name of the additional ingredient
	Name string `json:"Name"`
	// Amount in grams
	Amount float32 `json:"Amount"`
	// Duration in minutes
	Duration float32 `json:"Duration"`
}

// FermentationInstructions is the struct for the fermentation instructions
// It contains the yeast, the temperature, the additional ingredients and the carbonation in g/l
type FermentationInstructions struct {
	// Yeast is the yeast to use
	Yeast Yeast `json:"Yeast"`
	// Temperature is the fermentation temperature in °C
	// It can be a range of temperature (e.g. 18-22°C)
	Temperature string `json:"Temperature"`
	// AdditionalIngredients is the list of additional ingredients to use in the fermentation
	AdditionalIngredients []AdditionalIngredient `json:"FermentationAdditionalIngredients"`
	// Carbonation is the carbonation in g/l
	Carbonation float32 `json:"Carbonation"`
}

// Yeast is the struct for a yeast
// It contains the name and the amount in grams
type Yeast struct {
	// Name of the yeast
	Name string `json:"Name"`
	// Amount in grams
	Amount float32 `json:"Amount"`
}

// GetTotalMaltWeight returns the total weight of the malts in g
func (mash MashInstructions) GetTotalMaltWeight() float32 {
	var total float32
	for _, malt := range mash.Malts {
		total += malt.Amount
	}
	return total
}

// GetStatus returns the status of the recipe
func (r *Recipe) GetStatus() (RecipeStatus, []string) {
	r.statusLock.Lock()
	defer r.statusLock.Unlock()
	return r.status, r.statusParams
}

// SetStatus sets the status of the recipe together with optional parameters
func (r *Recipe) SetStatus(status RecipeStatus, params ...string) {
	r.statusLock.Lock()
	defer r.statusLock.Unlock()
	r.status = status
	r.statusParams = params
}

// GetStatusString returns the status of the recipe as a string
func (r *Recipe) GetStatusString() string {
	status, _ := r.GetStatus()
	switch status {
	case RecipeStatusCreated:
		return "Created"
	case RecipeStatusMashing:
		return "Mashing"
	case RecipeStatusLautering:
		return "Lautering"
	case RecipeStatusBoiling:
		return "Boiling"
	case RecipeStatusCooling:
		return "Cooling"
	case RecipeStatusPreFermentation:
		return "Pre-fermentation"
	case RecipeStatusFermenting:
		return "Fermenting"
	case RecipeStatusBottled:
		return "Bottled"
	case RecipeStatusFridge:
		return "Fridge"
	case RecipeStatusFinished:
		return "Finished"
	default:
		return "Unknown"
	}
}

// InitResults initializes the results of the recipe
func (r *Recipe) InitResults() {
	r.resultsLock.Lock()
	defer r.resultsLock.Unlock()
	r.results = RecipeResults{}
}

// GetResults returns the results of the recipe
func (r *Recipe) GetResults() RecipeResults {
	r.resultsLock.Lock()
	defer r.resultsLock.Unlock()
	return r.results
}

// SetHotWortVolume sets the hot wort volume of the recipe
func (r *Recipe) SetHotWortVolume(volume float32) {
	r.resultsLock.Lock()
	defer r.resultsLock.Unlock()
	r.results.HotWortVolume = volume
}

// SetOriginalGravity sets the original gravity of the recipe
func (r *Recipe) SetOriginalGravity(og float32) {
	r.resultsLock.Lock()
	defer r.resultsLock.Unlock()
	r.results.OriginalGravity = og
}

// SetFinalGravity sets the final gravity of the recipe
func (r *Recipe) SetFinalGravity(fg float32) {
	r.resultsLock.Lock()
	defer r.resultsLock.Unlock()
	r.results.FinalGravity = fg
}

// SetAlcohol sets the alcohol of the recipe
func (r *Recipe) SetAlcohol(alcohol float32) {
	r.resultsLock.Lock()
	defer r.resultsLock.Unlock()
	r.results.Alcohol = alcohol
}

// SetMainFermentationVolume sets the main fermentation volume of the recipe
func (r *Recipe) SetMainFermentationVolume(volume float32) {
	r.resultsLock.Lock()
	defer r.resultsLock.Unlock()
	r.results.MainFermentationVolume = volume
}

// SetVolumeBeforeBoil sets the volume measured before boiling the wort
func (r *Recipe) SetVolumeBeforeBoil(volume float32) {
	r.resultsLock.Lock()
	defer r.resultsLock.Unlock()
	r.results.VolumeBeforeBoil = volume
}

// SetSGMeasurement adds an sg measurement to the results of the recipe
func (r *Recipe) SetSGMeasurement(measurement *SGMeasurement) {
	r.mainFermSGsLock.Lock()
	defer r.mainFermSGsLock.Unlock()
	if r.mainFermSGs == nil {
		r.mainFermSGs = make([]*SGMeasurement, 0)
	}
	r.mainFermSGs = append(r.mainFermSGs, measurement)
}

// GetSGMeasurements returns all main fermentation sg measurements for this recipe
func (r *Recipe) GetSGMeasurements() []*SGMeasurement {
	r.mainFermSGsLock.Lock()
	defer r.mainFermSGsLock.Unlock()
	return r.mainFermSGs
}
