package summary

type Summary struct {
	Title string
	//GenerationDate is automatically populated by the printer when creating the summary file
	GenerationDate            string
	MashingInfo               *MashingInfo
	LauternInfo               string
	HoppingInfo               *HoppingInfo
	CoolingInfo               *CoolingInfo
	PreFermentationInfos      []*PreFermentationInfo
	YeastInfo                 *YeastInfo
	BottlingInfo              *BottlingInfo
	MainFermentationInfo      *MainFermentationInfo
	SecondaryFermentationInfo *SecondaryFermentationInfo
	Statistics                *Statistics
	Timeline                  []string
}

type MashingInfo struct {
	MashingTemperature float64
	MashingNotes       string
	RastInfos          []*MashRastInfo
}

type MashRastInfo struct {
	Temperature float64 `json:"temperature,omitempty"`
	Time        float64 `json:"time,omitempty"`
	Notes       string  `json:"notes,omitempty"`
}

type HoppingInfo struct {
	VolBeforeBoil *VolMeasurement
	HopInfos      []*HopInfo
	VolAfterBoil  *VolMeasurement
}

type VolMeasurement struct {
	Volume float32
	Notes  string
}

type HopInfo struct {
	Name     string  `json:"name,omitempty"`
	Grams    float32 `json:"grams,omitempty"`
	Alpha    float32 `json:"alpha,omitempty"`
	Time     float32 `json:"time,omitempty"`
	TimeUnit string  `json:"time_unit,omitempty"`
	Notes    string  `json:"notes,omitempty"`
}

type CoolingInfo struct {
	Temperature float32
	Time        float32
	Notes       string
}

type PreFermentationInfo struct {
	Volume float32
	SG     float32
	Notes  string
}

type YeastInfo struct {
	Temperature string //String to allow for ranges like 18-20
	Notes       string
}

type MainFermentationInfo struct {
	SGs        []*SGMeasurement
	Alcohol    float32
	DryHopInfo []*HopInfo
}

type SGMeasurement struct {
	SG    float32
	Date  string
	Final bool
	Notes string
}

type BottlingInfo struct {
	PreBottleVolume float32
	Carbonation     float32
	SugarAmount     float32
	SugarType       string
	Temperature     float32
	Alcohol         float32
	VolumeBottled   float32
	Notes           string
}

type SecondaryFermentationInfo struct {
	Days  int
	Notes string
}

type Statistics struct {
	Evaporation float32
	Efficiency  float32
}

func NewSummary() *Summary {
	return &Summary{}
}
