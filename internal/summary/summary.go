package summary

type Summary struct {
	Title string `json:"title,omitempty"`
	//GenerationDate is automatically populated by the printer when creating the summary file
	GenerationDate            string                     `json:"-"`
	MashingInfo               *MashingInfo               `json:"mashing_info,omitempty"`
	LauternInfo               string                     `json:"lautern_info,omitempty"`
	HoppingInfo               *HoppingInfo               `json:"hopping_info,omitempty"`
	CoolingInfo               *CoolingInfo               `json:"cooling_info,omitempty"`
	PreFermentationInfos      []*PreFermentationInfo     `json:"pre_fermentation_infos,omitempty"`
	YeastInfo                 *YeastInfo                 `json:"yeast_info,omitempty"`
	BottlingInfo              *BottlingInfo              `json:"bottling_info,omitempty"`
	MainFermentationInfo      *MainFermentationInfo      `json:"main_fermentation_info,omitempty"`
	SecondaryFermentationInfo *SecondaryFermentationInfo `json:"secondary_fermentation_info,omitempty"`
	Statistics                *Statistics                `json:"statistics,omitempty"`
	Timeline                  []string                   `json:"timeline,omitempty"`
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
	Volume float32 `json:"volume,omitempty"`
	Notes  string  `json:"notes,omitempty"`
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
