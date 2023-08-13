package lautern

import (
	"brewday/internal/recipe"
	"errors"

	"github.com/labstack/echo/v4"
)

type LauternRouter struct {
	TL      Timeline
	Summary SummaryRecorder
	Store   RecipeStore
}

// RegisterRoutes adds routes to the web server
func (r *LauternRouter) RegisterRoutes(root *echo.Echo, parent *echo.Group) {
	lautern := parent.Group("/lautern")
	lautern.GET("/:recipe_id", r.getLauternHandler).Name = "getLautern"
	lautern.POST("/:recipe_id", r.postLauternHandler).Name = "postLautern"
}

// addTimelineEvent adds an event to the timeline
func (r *LauternRouter) addTimelineEvent(message string) {
	if r.TL != nil {
		r.TL.AddEvent(message)
	}
}

// addSummaryLauternNotes adds lautern notes to the summary
func (r *LauternRouter) addSummaryLauternNotes(notes string) {
	if r.Summary != nil {
		r.Summary.AddLaunternNotes(notes)
	}
}

// getLauternHandler is the handler for the lautern page
func (r *LauternRouter) getLauternHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return errors.New("no recipe id provided")
	}
	r.addTimelineEvent("Started Läutern")
	// re, err := r.Store.Retrieve(id)
	// if err != nil {
	// 	return err
	// }
	// http://localhost:8080/mash/start/48756c612048756c6120495041
	re := &recipe.Recipe{
		Name:       "Hula Hula IPA",
		Style:      "IPA",
		BatchSize:  40,
		InitialSG:  1.073,
		Bitterness: 25,
		ColorEBC:   11,
		Mashing: recipe.MashInstructions{
			Malts: []recipe.Malt{
				{Name: "Golden Promise PA", Amount: 5600},
				{Name: "Barke Pilsner", Amount: 5000},
				{Name: "Haferflocken", Amount: 500},
				{Name: "Gerstenflocken", Amount: 500},
				{Name: "Carapils", Amount: 500},
				{Name: "Sauermalz", Amount: 300},
				{Name: "Cara Red", Amount: 300},
			},
			MainWaterVolume:    41,
			MashTemperature:    69,
			MashOutTemperature: 77,
			Rasts: []recipe.Rast{
				{Temperature: 67.5, Duration: 45},
				{Temperature: 72, Duration: 15},
			},
		},
		Hopping: recipe.HopInstructions{
			Hops: []recipe.Hops{
				{Name: "Simcoe (VW)", Amount: 34, Alpha: 12.5, Duration: 75, DryHop: false, Vorderwuerze: true},
				{Name: "Simcoe", Amount: 180, Alpha: 12.5, Duration: 0, DryHop: false},
				{Name: "Simcoe", Amount: 75, Alpha: 0, Duration: 0, DryHop: true},
				{Name: "Citra", Amount: 100, Alpha: 0, Duration: 0, DryHop: true},
				{Name: "Mosaic", Amount: 100, Alpha: 0, Duration: 0, DryHop: true},
			},
			AdditionalIngredients: nil,
		},
		Fermentation: recipe.FermentationInstructions{
			Yeast:       recipe.Yeast{Name: "WY 1007"},
			Temperature: "18-20",
			AdditionalIngredients: []recipe.AdditionalIngredient{
				{Name: "Cryo Citra", Amount: 60, Duration: 0},
				{Name: "Cryo Simcoe", Amount: 60, Duration: 0},
				{Name: "Motueka", Amount: 40, Duration: 0},
			},
			Carbonation: 5.5,
		},
	}
	return c.Render(200, "lautern.html", map[string]interface{}{
		"Title":       "Mash " + re.Name,
		"Subtitle":    "Läutern",
		"RecipeID":    id,
		"MashOutTemp": re.Mashing.MashOutTemperature,
		"Hops":        re.Hopping.Hops,
		"RestTime":    1,
	})
}

// postLauternHandler is the handler for the lautern page
func (r *LauternRouter) postLauternHandler(c echo.Context) error {
	var req ReqPostLautern
	err := c.Bind(&req)
	if err != nil {
		return err
	}
	id := c.Param("recipe_id")
	if id == "" {
		return errors.New("no recipe id provided")
	}
	r.addSummaryLauternNotes(req.Notes)
	return c.Redirect(302, c.Echo().Reverse("getStartHopping", id))
}
