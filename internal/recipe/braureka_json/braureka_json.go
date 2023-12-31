package braureka_json

import (
	"brewday/internal/recipe"
	"brewday/internal/tools"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// BraurekaJSONParser is a RecipeParser implementation that parses recipes in MMuM format generated by braureka (json)
// For some reason, the json format is not consistent with MMuM format, so we need to do some extra worK
type BraurekaJSONParser struct{}

// BraurekaJSON represents a recipe in Braureka JSON format
type BraurekaJSONRecipe struct {
	Name                string  `json:"Name"`
	Style               string  `json:"Sorte"`
	Volume              string  `json:"Ausschlagswuerze"`
	InitialPlato        string  `json:"Stammwuerze"`
	IBU                 string  `json:"Bittere"`
	EBC                 string  `json:"Farbe"`
	Alcohol             string  `json:"Alkohol"`
	MashType            string  `json:"Maischform"`
	MainWater           string  `json:"Infusion_Hauptguss"`
	ExtraWater          string  `json:"Nachguss"`
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
	MashTemp            string  `json:"Infusion_Einmaischtemperatur"`
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
	CookingTime         string  `json:"Kochzeit_Wuerze"`
	HopBefore1Name      string  `json:"Hopfen_VWH_1_Sorte"`
	HopBefore1Amount    string  `json:"Hopfen_VWH_1_Menge"`
	HopBefore1Alpha     string  `json:"Hopfen_VWH_1_alpha"`
	HopBefore2Name      string  `json:"Hopfen_VWH_2_Sorte"`
	HopBefore2Amount    string  `json:"Hopfen_VWH_2_Menge"`
	HopBefore2Alpha     string  `json:"Hopfen_VWH_2_alpha"`
	HopBefore3Name      string  `json:"Hopfen_VWH_3_Sorte"`
	HopBefore3Amount    string  `json:"Hopfen_VWH_3_Menge"`
	HopBefore3Alpha     string  `json:"Hopfen_VWH_3_alpha"`
	Hop1Name            string  `json:"Hopfen_1_Sorte"`
	Hop1Amount          string  `json:"Hopfen_1_Menge"`
	Hop1Alpha           string  `json:"Hopfen_1_alpha"`
	Hop1Time            string  `json:"Hopfen_1_Kochzeit"`
	Hop2Name            string  `json:"Hopfen_2_Sorte"`
	Hop2Amount          string  `json:"Hopfen_2_Menge"`
	Hop2Alpha           string  `json:"Hopfen_2_alpha"`
	Hop2Time            string  `json:"Hopfen_2_Kochzeit"`
	Hop3Name            string  `json:"Hopfen_3_Sorte"`
	Hop3Amount          string  `json:"Hopfen_3_Menge"`
	Hop3Alpha           string  `json:"Hopfen_3_alpha"`
	Hop3Time            string  `json:"Hopfen_3_Kochzeit"`
	Hop4Name            string  `json:"Hopfen_4_Sorte"`
	Hop4Amount          string  `json:"Hopfen_4_Menge"`
	Hop4Alpha           string  `json:"Hopfen_4_alpha"`
	Hop4Time            string  `json:"Hopfen_4_Kochzeit"`
	Hop5Name            string  `json:"Hopfen_5_Sorte"`
	Hop5Amount          string  `json:"Hopfen_5_Menge"`
	Hop5Alpha           string  `json:"Hopfen_5_alpha"`
	Hop5Time            string  `json:"Hopfen_5_Kochzeit"`
	Hop6Name            string  `json:"Hopfen_6_Sorte"`
	Hop6Amount          string  `json:"Hopfen_6_Menge"`
	Hop6Alpha           string  `json:"Hopfen_6_alpha"`
	Hop6Time            string  `json:"Hopfen_6_Kochzeit"`
	Hop7Name            string  `json:"Hopfen_7_Sorte"`
	Hop7Amount          string  `json:"Hopfen_7_Menge"`
	Hop7Alpha           string  `json:"Hopfen_7_alpha"`
	Hop7Time            string  `json:"Hopfen_7_Kochzeit"`
	OtherSpice1Name     string  `json:"WeitereZutat_Wuerze_1_Name"`
	OtherSpice1Amount   string  `json:"WeitereZutat_Wuerze_1_Menge"`
	OtherSpice1Unit     string  `json:"WeitereZutat_Wuerze_1_Einheit"`
	OtherSpice1Time     string  `json:"WeitereZutat_Wuerze_1_Kochzeit"`
	OtherSpice2Name     string  `json:"WeitereZutat_Wuerze_2_Name"`
	OtherSpice2Amount   string  `json:"WeitereZutat_Wuerze_2_Menge"`
	OtherSpice2Unit     string  `json:"WeitereZutat_Wuerze_2_Einheit"`
	OtherSpice2Time     string  `json:"WeitereZutat_Wuerze_2_Kochzeit"`
	OtherSpice3Name     string  `json:"WeitereZutat_Wuerze_3_Name"`
	OtherSpice3Amount   string  `json:"WeitereZutat_Wuerze_3_Menge"`
	OtherSpice3Unit     string  `json:"WeitereZutat_Wuerze_3_Einheit"`
	OtherSpice3Time     string  `json:"WeitereZutat_Wuerze_3_Kochzeit"`
	OtherSpice4Name     string  `json:"WeitereZutat_Wuerze_4_Name"`
	OtherSpice4Amount   string  `json:"WeitereZutat_Wuerze_4_Menge"`
	OtherSpice4Unit     string  `json:"WeitereZutat_Wuerze_4_Einheit"`
	OtherSpice4Time     string  `json:"WeitereZutat_Wuerze_4_Kochzeit"`
	OtherSpice5Name     string  `json:"WeitereZutat_Wuerze_5_Name"`
	OtherSpice5Amount   string  `json:"WeitereZutat_Wuerze_5_Menge"`
	OtherSpice5Unit     string  `json:"WeitereZutat_Wuerze_5_Einheit"`
	OtherSpice5Time     string  `json:"WeitereZutat_Wuerze_5_Kochzeit"`
	DryHop1Name         string  `json:"Stopfhopfen_1_Sorte"`
	DryHop1Amount       string  `json:"Stopfhopfen_1_Menge"`
	DryHop2Name         string  `json:"Stopfhopfen_2_Sorte"`
	DryHop2Amount       string  `json:"Stopfhopfen_2_Menge"`
	DryHop3Name         string  `json:"Stopfhopfen_3_Sorte"`
	DryHop3Amount       string  `json:"Stopfhopfen_3_Menge"`
	OtherFerment1Name   string  `json:"WeitereZutat_Gaerung_1_Name"`
	OtherFerment1Amount string  `json:"WeitereZutat_Gaerung_1_Menge"`
	OtherFerment1Unit   string  `json:"WeitereZutat_Gaerung_1_Einheit"`
	OtherFerment2Name   string  `json:"WeitereZutat_Gaerung_2_Name"`
	OtherFerment2Amount string  `json:"WeitereZutat_Gaerung_2_Menge"`
	OtherFerment2Unit   string  `json:"WeitereZutat_Gaerung_2_Einheit"`
	OtherFerment3Name   string  `json:"WeitereZutat_Gaerung_3_Name"`
	OtherFerment3Amount string  `json:"WeitereZutat_Gaerung_3_Menge"`
	OtherFerment3Unit   string  `json:"WeitereZutat_Gaerung_3_Einheit"`
	Yeast               string  `json:"Hefe"`
	FermentationTemp    string  `json:"Gaertemperatur"`
	Carbonation         string  `json:"Karbonisierung"`
	Notes               string  `json:"Anmerkung_Autor"`
}

// Parse parses a recipe from a string
func (p *BraurekaJSONParser) Parse(recipe string) (*recipe.Recipe, error) {
	var r BraurekaJSONRecipe
	err := json.Unmarshal([]byte(recipe), &r)
	if err != nil {
		return nil, err
	}
	return bJsonRecipeToRecipe(&r)
}

// stringToFloat converts a string to a float32
func stringToFloat(s string) (float32, error) {
	if s == "" {
		return 0, nil
	}
	f, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return 0, err
	}
	return float32(f), nil
}

// bJsonRecipeToRecipe converts a BraurekaJSONRecipe to a recipe.Recipe
func bJsonRecipeToRecipe(r *BraurekaJSONRecipe) (*recipe.Recipe, error) {
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
	vol, err := stringToFloat(r.Volume)
	if err != nil {
		return nil, err
	}
	initialPlato, err := stringToFloat(r.InitialPlato)
	if err != nil {
		return nil, err
	}
	ibu, err := stringToFloat(r.IBU)
	if err != nil {
		return nil, err
	}
	return &recipe.Recipe{
		Name:         r.Name,
		Style:        r.Style,
		BatchSize:    vol,
		InitialSG:    tools.PlatoToSG(initialPlato),
		Bitterness:   ibu,
		ColorEBC:     float32(color),
		Mashing:      *mashInst,
		Hopping:      *hopInst,
		Fermentation: *fermInst,
	}, nil
}

// getMashInstructions returns the mash instructions for a BraurekaJSONRecipe
func getMashInstructions(r *BraurekaJSONRecipe) (*recipe.MashInstructions, error) {
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
	mainWater, err := stringToFloat(r.MainWater)
	if err != nil {
		return nil, err
	}
	extraWater, err := stringToFloat(r.ExtraWater)
	if err != nil {
		return nil, err
	}
	mashTemp, err := stringToFloat(r.MashTemp)
	if err != nil {
		return nil, err
	}
	return &recipe.MashInstructions{
		Malts:              malts,
		MainWaterVolume:    mainWater,
		Nachguss:           extraWater,
		MashTemperature:    mashTemp,
		MashOutTemperature: float32(mashOutValue),
		Rasts:              rasts,
	}, nil
}

// getHopInstructions returns the hop instructions for a BraurekaJSONRecipe
func getHopInstructions(r *BraurekaJSONRecipe) (*recipe.HopInstructions, error) {
	v := reflect.ValueOf(r).Elem()
	var hops []recipe.Hops
	var additions []recipe.AdditionalIngredient
	for i := 1; i <= 3; i++ {
		h := recipe.Hops{}
		h.DryHop = false
		h.Vorderwuerze = true
		nameValue := v.FieldByName(fmt.Sprintf("HopBefore%dName", i)).String()
		if nameValue == "" {
			continue
		}
		h.Name = strings.TrimSpace(nameValue) + " (VW)"
		amountValue := v.FieldByName(fmt.Sprintf("HopBefore%dAmount", i)).String()
		alphaValue := v.FieldByName(fmt.Sprintf("HopBefore%dAlpha", i)).String()
		amountValueFloat, err := stringToFloat(amountValue)
		if err != nil {
			return nil, err
		}
		alphaValueFloat, err := stringToFloat(alphaValue)
		if err != nil {
			return nil, err
		}
		h.Amount = amountValueFloat
		h.Alpha = alphaValueFloat
		cookingTime, err := stringToFloat(r.CookingTime)
		if err != nil {
			return nil, err
		}
		h.Duration = cookingTime
		hops = append(hops, h)
	}
	for i := 1; i <= 7; i++ {
		h := recipe.Hops{}
		h.DryHop = false
		h.Vorderwuerze = false
		nameValue := v.FieldByName(fmt.Sprintf("Hop%dName", i)).String()
		if nameValue == "" {
			continue
		}
		h.Name = strings.TrimSpace(nameValue)
		amountValue := v.FieldByName(fmt.Sprintf("Hop%dAmount", i)).String()
		alphaValue := v.FieldByName(fmt.Sprintf("Hop%dAlpha", i)).String()
		durationValue := v.FieldByName(fmt.Sprintf("Hop%dTime", i)).String()
		amountValueFloat, err := stringToFloat(amountValue)
		if err != nil {
			return nil, err
		}
		alphaValueFloat, err := stringToFloat(alphaValue)
		if err != nil {
			return nil, err
		}
		if durationValue != "Whirlpool" {
			durationValueFloat, err := stringToFloat(durationValue)
			if err != nil {
				return nil, err
			}
			h.Duration = durationValueFloat
		} else {
			h.Duration = 0
		}
		h.Amount = amountValueFloat
		h.Alpha = alphaValueFloat

		hops = append(hops, h)
	}
	for i := 1; i <= 3; i++ {
		h := recipe.Hops{}
		h.DryHop = true
		h.Vorderwuerze = false
		nameValue := v.FieldByName(fmt.Sprintf("DryHop%dName", i)).String()
		if nameValue == "" {
			continue
		}
		h.Name = strings.TrimSpace(nameValue)
		amountValue := v.FieldByName(fmt.Sprintf("DryHop%dAmount", i)).String()
		amountValueFloat, err := stringToFloat(amountValue)
		if err != nil {
			return nil, err
		}
		h.Amount = amountValueFloat
		hops = append(hops, h)
	}
	for i := 1; i <= 5; i++ {
		a := recipe.AdditionalIngredient{}
		nameValue := v.FieldByName(fmt.Sprintf("OtherSpice%dName", i)).String()
		if nameValue == "" {
			continue
		}
		a.Name = strings.TrimSpace(nameValue)
		amountValue := v.FieldByName(fmt.Sprintf("OtherSpice%dAmount", i)).String()
		amountValueFloat, err := stringToFloat(amountValue)
		if err != nil {
			return nil, err
		}
		unitValue := v.FieldByName(fmt.Sprintf("OtherSpice%dUnit", i)).String()
		if unitValue == "kg" {
			a.Amount = amountValueFloat * 1000
		} else {
			a.Amount = amountValueFloat
		}
		timeValue := v.FieldByName(fmt.Sprintf("OtherSpice%dTime", i)).String()
		timeValueFloat, err := stringToFloat(timeValue)
		if err != nil {
			return nil, err
		}
		a.Duration = timeValueFloat
		additions = append(additions, a)
	}
	cookingTime, err := stringToFloat(r.CookingTime)
	if err != nil {
		return nil, err
	}
	return &recipe.HopInstructions{
		TotalCookingTime:      cookingTime,
		Hops:                  hops,
		AdditionalIngredients: additions,
	}, nil
}

// getFermentationInstructions returns the fermentation instructions for a BraurekaJSONRecipe
func getFermentationInstructions(r *BraurekaJSONRecipe) (*recipe.FermentationInstructions, error) {
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
		amountValue := v.FieldByName(fmt.Sprintf("OtherFerment%dAmount", i)).String()
		amountValueFloat, err := stringToFloat(amountValue)
		if err != nil {
			return nil, err
		}
		unitValue := v.FieldByName(fmt.Sprintf("OtherFerment%dUnit", i)).String()
		if unitValue == "kg" {
			a.Amount = amountValueFloat * 1000
		} else {
			a.Amount = amountValueFloat
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
