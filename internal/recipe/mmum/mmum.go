package mmum

import (
	"brewday/internal/recipe"
	"brewday/internal/tools"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// MMUMParser is a RecipeParser implementation that parses recipes in Maische Malz und Mehr format (json)
type MMUMParser struct{}

// MMUMRecipe represents a recipe in Maische Malz und Mehr format
type MMUMRecipe struct {
	Name                string  `json:"Name"`
	Style               string  `json:"Sorte"`
	Volume              int     `json:"Ausschlagswuerze"`
	InitialPlato        float32 `json:"Stammwuerze"`
	IBU                 float32 `json:"Bittere"`
	EBC                 string  `json:"Farbe"`
	Alcohol             float32 `json:"Alkohol"`
	MashType            string  `json:"Maischform"`
	MainWater           float32 `json:"Infusion_Hauptguss"`
	ExtraWater          float32 `json:"Nachguss"`
	Malt1Name           string  `json:"Malz1"`
	Malt1Amount         float32 `json:"Malz1_Menge"`
	Malt1Unit           string  `json:"Malz1_Einheit"`
	Malt2Name           string  `json:"Malz2"`
	Malt2Amount         float32 `json:"Malz2_Menge"`
	Malt2Unit           string  `json:"Malz2_Einheit"`
	Malt3Name           string  `json:"Malz3"`
	Malt3Amount         float32 `json:"Malz3_Menge"`
	Malt3Unit           string  `json:"Malz3_Einheit"`
	Malt4Name           string  `json:"Malz4"`
	Malt4Amount         float32 `json:"Malz4_Menge"`
	Malt4Unit           string  `json:"Malz4_Einheit"`
	Malt5Name           string  `json:"Malz5"`
	Malt5Amount         float32 `json:"Malz5_Menge"`
	Malt5Unit           string  `json:"Malz5_Einheit"`
	Malt6Name           string  `json:"Malz6"`
	Malt6Amount         float32 `json:"Malz6_Menge"`
	Malt6Unit           string  `json:"Malz6_Einheit"`
	Malt7Name           string  `json:"Malz7"`
	Malt7Amount         float32 `json:"Malz7_Menge"`
	Malt7Unit           string  `json:"Malz7_Einheit"`
	MashTemp            float32 `json:"Infusion_Einmaischtemperatur"`
	MashOutTemp         string  `json:"Abmaischtemperatur"`
	MashRast1Temp       string  `json:"Infusion_Rasttemperatur1"`
	MashRast1Time       string  `json:"Infusion_Rastzeit1"`
	MashRast2Temp       string  `json:"Infusion_Rasttemperatur2"`
	MashRast2Time       string  `json:"Infusion_Rastzeit2"`
	MashRast3Temp       string  `json:"Infusion_Rasttemperatur3"`
	MashRast3Time       string  `json:"Infusion_Rastzeit3"`
	MashRast4Temp       string  `json:"Infusion_Rasttemperatur4"`
	MashRast4Time       string  `json:"Infusion_Rastzeit4"`
	MashRast5Temp       string  `json:"Infusion_Rasttemperatur5"`
	MashRast5Time       string  `json:"Infusion_Rastzeit5"`
	MashRast6Temp       string  `json:"Infusion_Rasttemperatur6"`
	MashRast6Time       string  `json:"Infusion_Rastzeit6"`
	MashRast7Temp       string  `json:"Infusion_Rasttemperatur7"`
	MashRast7Time       string  `json:"Infusion_Rastzeit7"`
	CookingTime         float32 `json:"Kochzeit_Wuerze"`
	HopBefore1Name      string  `json:"Hopfen_VWH_1_Sorte"`
	HopBefore1Amount    float32 `json:"Hopfen_VWH_1_Menge"`
	HopBefore1Alpha     float32 `json:"Hopfen_VWH_1_alpha"`
	HopBefore2Name      string  `json:"Hopfen_VWH_2_Sorte"`
	HopBefore2Amount    float32 `json:"Hopfen_VWH_2_Menge"`
	HopBefore2Alpha     float32 `json:"Hopfen_VWH_2_alpha"`
	HopBefore3Name      string  `json:"Hopfen_VWH_3_Sorte"`
	HopBefore3Amount    float32 `json:"Hopfen_VWH_3_Menge"`
	HopBefore3Alpha     float32 `json:"Hopfen_VWH_3_alpha"`
	Hop1Name            string  `json:"Hopfen_1_Sorte"`
	Hop1Amount          float32 `json:"Hopfen_1_Menge"`
	Hop1Alpha           float32 `json:"Hopfen_1_alpha"`
	Hop1Time            float32 `json:"Hopfen_1_Kochzeit"`
	Hop2Name            string  `json:"Hopfen_2_Sorte"`
	Hop2Amount          float32 `json:"Hopfen_2_Menge"`
	Hop2Alpha           float32 `json:"Hopfen_2_alpha"`
	Hop2Time            float32 `json:"Hopfen_2_Kochzeit"`
	Hop3Name            string  `json:"Hopfen_3_Sorte"`
	Hop3Amount          float32 `json:"Hopfen_3_Menge"`
	Hop3Alpha           float32 `json:"Hopfen_3_alpha"`
	Hop3Time            float32 `json:"Hopfen_3_Kochzeit"`
	Hop4Name            string  `json:"Hopfen_4_Sorte"`
	Hop4Amount          float32 `json:"Hopfen_4_Menge"`
	Hop4Alpha           float32 `json:"Hopfen_4_alpha"`
	Hop4Time            float32 `json:"Hopfen_4_Kochzeit"`
	Hop5Name            string  `json:"Hopfen_5_Sorte"`
	Hop5Amount          float32 `json:"Hopfen_5_Menge"`
	Hop5Alpha           float32 `json:"Hopfen_5_alpha"`
	Hop5Time            float32 `json:"Hopfen_5_Kochzeit"`
	Hop6Name            string  `json:"Hopfen_6_Sorte"`
	Hop6Amount          float32 `json:"Hopfen_6_Menge"`
	Hop6Alpha           float32 `json:"Hopfen_6_alpha"`
	Hop6Time            float32 `json:"Hopfen_6_Kochzeit"`
	Hop7Name            string  `json:"Hopfen_7_Sorte"`
	Hop7Amount          float32 `json:"Hopfen_7_Menge"`
	Hop7Alpha           float32 `json:"Hopfen_7_alpha"`
	Hop7Time            float32 `json:"Hopfen_7_Kochzeit"`
	OtherSpice1Name     string  `json:"WeitereZutat_Wuerze_1_Name"`
	OtherSpice1Amount   float32 `json:"WeitereZutat_Wuerze_1_Menge"`
	OtherSpice1Unit     string  `json:"WeitereZutat_Wuerze_1_Einheit"`
	OtherSpice1Time     float32 `json:"WeitereZutat_Wuerze_1_Kochzeit"`
	OtherSpice2Name     string  `json:"WeitereZutat_Wuerze_2_Name"`
	OtherSpice2Amount   float32 `json:"WeitereZutat_Wuerze_2_Menge"`
	OtherSpice2Unit     string  `json:"WeitereZutat_Wuerze_2_Einheit"`
	OtherSpice2Time     float32 `json:"WeitereZutat_Wuerze_2_Kochzeit"`
	OtherSpice3Name     string  `json:"WeitereZutat_Wuerze_3_Name"`
	OtherSpice3Amount   float32 `json:"WeitereZutat_Wuerze_3_Menge"`
	OtherSpice3Unit     string  `json:"WeitereZutat_Wuerze_3_Einheit"`
	OtherSpice3Time     float32 `json:"WeitereZutat_Wuerze_3_Kochzeit"`
	OtherSpice4Name     string  `json:"WeitereZutat_Wuerze_4_Name"`
	OtherSpice4Amount   float32 `json:"WeitereZutat_Wuerze_4_Menge"`
	OtherSpice4Unit     string  `json:"WeitereZutat_Wuerze_4_Einheit"`
	OtherSpice4Time     float32 `json:"WeitereZutat_Wuerze_4_Kochzeit"`
	OtherSpice5Name     string  `json:"WeitereZutat_Wuerze_5_Name"`
	OtherSpice5Amount   float32 `json:"WeitereZutat_Wuerze_5_Menge"`
	OtherSpice5Unit     string  `json:"WeitereZutat_Wuerze_5_Einheit"`
	OtherSpice5Time     float32 `json:"WeitereZutat_Wuerze_5_Kochzeit"`
	DryHop1Name         string  `json:"Stopfhopfen_1_Sorte"`
	DryHop1Amount       float32 `json:"Stopfhopfen_1_Menge"`
	DryHop2Name         string  `json:"Stopfhopfen_2_Sorte"`
	DryHop2Amount       float32 `json:"Stopfhopfen_2_Menge"`
	DryHop3Name         string  `json:"Stopfhopfen_3_Sorte"`
	DryHop3Amount       float32 `json:"Stopfhopfen_3_Menge"`
	OtherFerment1Name   string  `json:"WeitereZutat_Gaerung_1_Name"`
	OtherFerment1Amount float32 `json:"WeitereZutat_Gaerung_1_Menge"`
	OtherFerment1Unit   string  `json:"WeitereZutat_Gaerung_1_Einheit"`
	OtherFerment2Name   string  `json:"WeitereZutat_Gaerung_2_Name"`
	OtherFerment2Amount float32 `json:"WeitereZutat_Gaerung_2_Menge"`
	OtherFerment2Unit   string  `json:"WeitereZutat_Gaerung_2_Einheit"`
	OtherFerment3Name   string  `json:"WeitereZutat_Gaerung_3_Name"`
	OtherFerment3Amount float32 `json:"WeitereZutat_Gaerung_3_Menge"`
	OtherFerment3Unit   string  `json:"WeitereZutat_Gaerung_3_Einheit"`
	Yeast               string  `json:"Hefe"`
	FermentationTemp    string  `json:"Gaertemperatur"`
	Carbonation         string  `json:"Karbonisierung"`
	Notes               string  `json:"Anmerkung_Autor"`
}

// Parse parses a recipe from a string
func (p *MMUMParser) Parse(recipe string) (*recipe.Recipe, error) {
	var r MMUMRecipe
	err := json.Unmarshal([]byte(recipe), &r)
	if err != nil {
		return nil, err
	}
	return mmumRecipeToRecipe(&r)
}

// mmumRecipeToRecipe converts a MMUMRecipe to a recipe.Recipe
func mmumRecipeToRecipe(r *MMUMRecipe) (*recipe.Recipe, error) {
	color, err := strconv.ParseFloat(r.EBC, 32)
	if err != nil {
		return nil, err
	}
	mashInst, err := getMashInstructions(r)
	if err != nil {
		return nil, err
	}
	hopInst, err := getHopInstructions(r)
	if err != nil {
		return nil, err
	}
	fermInst, err := getFermentationInstructions(r)
	if err != nil {
		return nil, err
	}
	return &recipe.Recipe{
		Name:         r.Name,
		Style:        r.Style,
		BatchSize:    float32(r.Volume),
		InitialSG:    tools.PlatoToSG(r.InitialPlato),
		Bitterness:   r.IBU,
		ColorEBC:     float32(color),
		Mashing:      *mashInst,
		Hopping:      *hopInst,
		Fermentation: *fermInst,
	}, nil
}

// getMashInstructions returns the mash instructions for a MMUMRecipe
func getMashInstructions(r *MMUMRecipe) (*recipe.MashInstructions, error) {
	var malts []recipe.Malt
	var rasts []recipe.Rast
	v := reflect.ValueOf(r).Elem()
	for i := 1; i <= 7; i++ {
		m := recipe.Malt{}
		nameValue := v.FieldByName(fmt.Sprintf("Malt%dName", i)).String()
		if nameValue == "" {
			continue
		}
		m.Name = strings.TrimSpace(nameValue)
		amountValue := v.FieldByName(fmt.Sprintf("Malt%dAmount", i)).Float()
		unitValue := v.FieldByName(fmt.Sprintf("Malt%dUnit", i)).String()
		if unitValue == "kg" {
			m.Amount = float32(amountValue * 1000)
		} else {
			m.Amount = float32(amountValue)
		}
		malts = append(malts, m)
	}
	for i := 1; i <= 7; i++ {
		r := recipe.Rast{}
		tempValue := v.FieldByName(fmt.Sprintf("MashRast%dTemp", i))
		if tempValue.IsZero() {
			continue
		}
		tempValueFloat, err := strconv.ParseFloat(tempValue.String(), 32)
		if err != nil {
			return nil, err
		}
		r.Temperature = float32(tempValueFloat)
		timeValue := v.FieldByName(fmt.Sprintf("MashRast%dTime", i)).String()
		timeValueFloat, err := strconv.ParseFloat(timeValue, 32)
		if err != nil {
			return nil, err
		}
		r.Duration = float32(timeValueFloat)
		rasts = append(rasts, r)
	}
	mashOutValue, err := strconv.ParseFloat(r.MashOutTemp, 32)
	if err != nil {
		return nil, err
	}
	return &recipe.MashInstructions{
		Malts:              malts,
		MainWaterVolume:    r.MainWater,
		MashTemperature:    r.MashTemp,
		MashOutTemperature: float32(mashOutValue),
		Rasts:              rasts,
	}, nil
}

// getHopInstructions returns the hop instructions for a MMUMRecipe
func getHopInstructions(r *MMUMRecipe) (*recipe.HopInstructions, error) {
	v := reflect.ValueOf(r).Elem()
	var hops []recipe.Hops
	var additions []recipe.AdditionalIngredient
	for i := 1; i <= 3; i++ {
		h := recipe.Hops{}
		h.DryHop = false
		nameValue := v.FieldByName(fmt.Sprintf("HopBefore%dName", i)).String()
		if nameValue == "" {
			continue
		}
		h.Name = strings.TrimSpace(nameValue) + " (VW)"
		amountValue := v.FieldByName(fmt.Sprintf("HopBefore%dAmount", i)).Float()
		alphaValue := v.FieldByName(fmt.Sprintf("HopBefore%dAlpha", i)).Float()
		h.Amount = float32(amountValue)
		h.Alpha = float32(alphaValue)
		h.Duration = r.CookingTime
		hops = append(hops, h)
	}
	for i := 1; i <= 7; i++ {
		h := recipe.Hops{}
		h.DryHop = false
		nameValue := v.FieldByName(fmt.Sprintf("Hop%dName", i)).String()
		if nameValue == "" {
			continue
		}
		h.Name = strings.TrimSpace(nameValue)
		amountValue := v.FieldByName(fmt.Sprintf("Hop%dAmount", i)).Float()
		alphaValue := v.FieldByName(fmt.Sprintf("Hop%dAlpha", i)).Float()
		durationValue := v.FieldByName(fmt.Sprintf("Hop%dTime", i)).Float()
		h.Amount = float32(amountValue)
		h.Alpha = float32(alphaValue)
		h.Duration = float32(durationValue)
		hops = append(hops, h)
	}
	for i := 1; i <= 3; i++ {
		h := recipe.Hops{}
		h.DryHop = true
		nameValue := v.FieldByName(fmt.Sprintf("DryHop%dName", i)).String()
		if nameValue == "" {
			continue
		}
		h.Name = strings.TrimSpace(nameValue)
		amountValue := v.FieldByName(fmt.Sprintf("DryHop%dAmount", i)).Float()
		h.Amount = float32(amountValue)
		hops = append(hops, h)
	}
	for i := 1; i <= 5; i++ {
		a := recipe.AdditionalIngredient{}
		nameValue := v.FieldByName(fmt.Sprintf("OtherSpice%dName", i)).String()
		if nameValue == "" {
			continue
		}
		a.Name = strings.TrimSpace(nameValue)
		amountValue := v.FieldByName(fmt.Sprintf("OtherSpice%dAmount", i)).Float()
		unitValue := v.FieldByName(fmt.Sprintf("OtherSpice%dUnit", i)).String()
		if unitValue == "kg" {
			a.Amount = float32(amountValue * 1000)
		} else {
			a.Amount = float32(amountValue)
		}
		timeValue := v.FieldByName(fmt.Sprintf("OtherSpice%dTime", i)).Float()
		a.Duration = float32(timeValue)
		additions = append(additions, a)
	}
	return &recipe.HopInstructions{
		Hops:                  hops,
		AdditionalIngredients: additions,
	}, nil
}

// getFermentationInstructions returns the fermentation instructions for a MMUMRecipe
func getFermentationInstructions(r *MMUMRecipe) (*recipe.FermentationInstructions, error) {
	carbonationValue, err := strconv.ParseFloat(r.Carbonation, 32)
	if err != nil {
		return nil, err
	}
	v := reflect.ValueOf(r).Elem()
	var additions []recipe.AdditionalIngredient
	for i := 1; i <= 3; i++ {
		a := recipe.AdditionalIngredient{}
		nameValue := v.FieldByName(fmt.Sprintf("OtherFerment%dName", i)).String()
		if nameValue == "" {
			continue
		}
		a.Name = strings.TrimSpace(nameValue)
		amountValue := v.FieldByName(fmt.Sprintf("OtherFerment%dAmount", i)).Float()
		unitValue := v.FieldByName(fmt.Sprintf("OtherFerment%dUnit", i)).String()
		if unitValue == "kg" {
			a.Amount = float32(amountValue * 1000)
		} else {
			a.Amount = float32(amountValue)
		}
		additions = append(additions, a)
	}
	return &recipe.FermentationInstructions{
		Yeast: recipe.Yeast{
			Name: r.Yeast,
		},
		Temperature:           r.FermentationTemp,
		AdditionalIngredients: additions,
		Carbonation:           float32(carbonationValue),
	}, nil
}
