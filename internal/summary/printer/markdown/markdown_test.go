package markdown

import (
	"brewday/internal/summary"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

// func Test(t *testing.T) {
// 	//filetext, err := os.ReadFile("./md.tmpl")
// 	_, err := os.ReadFile("./md.tmpl")

// 	require.NoError(t, err)
// 	//te := template.Must(template.New("t1").Parse(string(filetext)))
// 	te := template.Must(template.ParseFiles("./md.tmpl"))
// 	summ := &summary.Summary{
// 		Title:          "My title",
// 		GenerationDate: "oe",
// 	}
// 	err = te.Execute(os.Stdout, summ)
// 	// err = te.ExecuteTemplate(os.Stdout, "md.tmpl", summ)
// 	require.NoError(t, err)
// }

func TestPrint(t *testing.T) {
	require := require.New(t)
	testCases := []struct {
		Name     string
		Summ     *summary.Summary
		Expected string
		Error    bool
	}{
		{
			Name: "",
			Summ: &summary.Summary{
				Title:          "My Title",
				GenerationDate: "date",
				MashingInfo: &summary.MashingInfo{
					MashingTemperature: 57,
					MashingNotes:       "notes",
					RastInfos: []*summary.MashRastInfo{
						{
							Temperature: 63,
							Time:        30,
							Notes:       "notes1",
						},
						{
							Temperature: 72,
							Time:        40,
							Notes:       "notes2",
						},
						{
							Temperature: 78,
							Time:        15,
							Notes:       "notes3",
						},
					},
				},
				LauternInfo: "lautern info 1\nlautern info 2",
				HoppingInfo: &summary.HoppingInfo{
					VolBeforeBoil: &summary.VolMeasurement{Volume: 13.10, Notes: "notes4"},
					HopInfos: []*summary.HopInfo{
						{Name: "Karamellsirup", Grams: 450.00, Alpha: 0.00, Time: 60.00, TimeUnit: "minutes", Notes: "notes5"},
						{Name: "Hallertauer Tradition", Grams: 15.00, Alpha: 5.50, Time: 60.00, TimeUnit: "minutes", Notes: "notes6"},
						{Name: "Saazer", Grams: 16.00, Alpha: 3.60, Time: 20.00, TimeUnit: "minutes", Notes: "notes7"},
						{Name: "Koriander", Grams: 4.00, Alpha: 0.00, Time: 10.00, TimeUnit: "minutes", Notes: "notes8"},
					},
					VolAfterBoil: &summary.VolMeasurement{Volume: 8.70, Notes: "notes9"},
				},
				CoolingInfo: &summary.CoolingInfo{
					Temperature: 20.00,
					Time:        20.00,
					Notes:       "notes10",
				},
				PreFermentationInfos: []*summary.PreFermentationInfo{
					{Volume: 7.70, SG: 1.098, Notes: "notes11"},
					{Volume: 12.90, SG: 1.054, Notes: "notes12"},
				},
				YeastInfo: &summary.YeastInfo{
					Temperature: "20.00",
					Notes:       "notes13",
				},
				MainFermentationInfo: &summary.MainFermentationInfo{
					SGs: []*summary.SGMeasurement{
						{SG: 1.013, Date: "2024-03-22", Final: false, Notes: "notes14"},
						{SG: 1.011, Date: "2024-03-23", Final: false, Notes: "notes15"},
						{SG: 1.011, Date: "2024-03-24", Final: true, Notes: "notes16"},
					},
					Alcohol: 5.64,
					DryHopInfo: []*summary.HopInfo{
						{Name: "Citra", Grams: 3.00, Alpha: 13.70, Time: 5.00, TimeUnit: "days", Notes: "notes17"},
						{Name: "Solero", Grams: 4.00, Alpha: 3.70, Time: 2.00, TimeUnit: "days", Notes: "notes18"},
					},
				},
				BottlingInfo: &summary.BottlingInfo{
					PreBottleVolume: 8.80,
					Carbonation:     5.52,
					SugarAmount:     70.00,
					SugarType:       "glucose",
					Temperature:     20.00,
					Alcohol:         5.86,
					VolumeBottled:   10.00,
					Notes:           "notes19",
				},
				SecondaryFermentationInfo: &summary.SecondaryFermentationInfo{Days: 5, Notes: "notes20"},
				Statistics: &summary.Statistics{
					Evaporation: 25.19,
					Efficiency:  49.11,
				},
				Timeline: []string{
					"2024-02-14T07:39:20.732108778Z@Started mashing",
					"2024-02-14T07:39:28.454239024Z@Finished Einmaischen",
					"2024-02-14T07:39:31.412846385Z@Started Rast 0",
				},
			},
			Expected: `# My Title

The following summary was generated on %s


## Mash

- **Mashing temperature**: 57.00°C (notes)
- **Rast**: 63.00°C for 30.00 minutes (notes1)
- **Rast**: 72.00°C for 40.00 minutes (notes2)
- **Rast**: 78.00°C for 15.00 minutes (notes3)


## Lautern

lautern info 1
lautern info 2


## Hopping

- **Measured volume - Measured volume before boiling**: 13.10L (notes4)
- **Karamellsirup**: 450.00g (0.00%% alpha) [60.00 minutes] (notes5)
- **Hallertauer Tradition**: 15.00g (5.50%% alpha) [60.00 minutes] (notes6)
- **Saazer**: 16.00g (3.60%% alpha) [20.00 minutes] (notes7)
- **Koriander**: 4.00g (0.00%% alpha) [10.00 minutes] (notes8)
- **Measured volume - Measured volume after boiling**: 8.70L (notes9)


## Cooling

Reached 20.00°C in 20.00 minutes (notes10)


## Pre-fermentation

- Measured 7.70L with SG: 1.098 (notes11)
- Measured 12.90L with SG: 1.054 (notes12)


## Yeast start

- **Temperature**: 20.00°C

notes13


## Main fermentation

Date | SG | Final | Notes
---|---|---|---
2024-03-22 | 1.013 | No | notes14
2024-03-23 | 1.011 | No | notes15
2024-03-24 | 1.011 | Yes | notes16

- Alcohol after main fermentation: 5.64%%

### Dry Hopping

- **Citra**: 3.00g (13.70%% alpha) [5.00 days] (notes17)
- **Solero**: 4.00g (3.70%% alpha) [2.00 days] (notes18)


## Bottling

- **Volume in tank**: 8.80L

- **Sugar**: 70.00 g (glucose)
- **Temperature**: 20.00°C

- **Carbonation**: 5.52 g/L
- **Final Alcohol**: 5.86%% vol
- **Volume Bottled**: 10.00L

notes19


## Secondary fermentation

- **Days**: 5

notes20


## Calculations

- **Evaporation**: 25.19%%/h
- **Efficiency**: 49.11%%


## Timeline

Timestamp | Event
--- | ---
2024-02-14T07:39:20.732108778Z | Started mashing
2024-02-14T07:39:28.454239024Z | Finished Einmaischen
2024-02-14T07:39:31.412846385Z | Started Rast 0`,
			Error: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			p := &MarkdownPrinter{}
			actual, err := p.Print(tc.Summ)
			if tc.Error {
				require.Error(err)
			} else {
				require.NoError(err)
				exp := fmt.Sprintf(tc.Expected, tc.Summ.GenerationDate)
				require.Equal(exp, actual)
			}
		})
	}
}
