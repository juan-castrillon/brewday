package mash

import (
	"brewday/internal/recipe"
	"errors"
	"fmt"
	"strconv"

	"github.com/labstack/echo/v4"
)

type MashRouter struct {
	Store   RecipeStore
	TL      Timeline
	Summary SummaryRecorder
	recipe  *recipe.Recipe
}

func (r *MashRouter) RegisterRoutes(root *echo.Echo, parent *echo.Group) {
	mash := parent.Group("/mash")
	mash.GET("/start/:recipe_id", r.getMashStartHandler).Name = "getMashStart"
	mash.POST("/rasts/:recipe_id/:rast_num", r.postRastsHandler).Name = "postRasts"
}

// addTimelineEvent adds an event to the timeline
func (r *MashRouter) addTimelineEvent(message string) {
	if r.TL != nil {
		r.TL.AddEvent(message)
	}
}

// addSummaryMashTemp adds a mash temperature to the summary and notes related to it
func (r *MashRouter) addSummaryMashTemp(temp float64, notes string) {
	if r.Summary != nil {
		r.Summary.AddMashTemp(temp, notes)
	}
}

// addSummaryRast adds a rast to the summary and notes related to it
func (r *MashRouter) addSummaryRast(temp float64, duration float64, notes string) {
	if r.Summary != nil {
		r.Summary.AddRast(temp, duration, notes)
	}
}

// getMashStartHandler is the handler for the mash start page
func (r *MashRouter) getMashStartHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return errors.New("no recipe id provided")
	}
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
			Nachguss: 5,
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
	r.recipe = re
	r.addTimelineEvent("Started mashing")
	return c.Render(200, "mash_start.html", map[string]interface{}{
		"Title":     "Mash " + re.Name,
		"MainWater": re.Mashing.MainWaterVolume,
		"MashTemp":  re.Mashing.MashTemperature,
		"RecipeID":  id,
	})
}

// postRastsHandler is the handler for the mash rasts page
func (r *MashRouter) postRastsHandler(c echo.Context) error {
	r.addTimelineEvent("Finished einmaischen")
	id := c.Param("recipe_id")
	if id == "" {
		return errors.New("no recipe id provided")
	}
	if r.recipe == nil {
		return errors.New("no recipe loaded")
	}
	rastNumStr := c.Param("rast_num")
	if rastNumStr == "" {
		return errors.New("no rast number provided")
	}
	rastNum, err := strconv.Atoi(rastNumStr)
	if err != nil {
		return err
	}
	var nextRastNum int
	switch rastNum {
	case 0:
		r.addTimelineEvent(fmt.Sprintf("Starting rast %d", rastNum))
		var req ReqPostFirstRast
		err := c.Bind(&req)
		if err != nil {
			return err
		}
		r.addSummaryMashTemp(req.RealMashTemperature, req.Notes)
		nextRastNum = 1
	case len(r.recipe.Mashing.Rasts):
		r.addTimelineEvent("Finished mashing")
		return c.Redirect(302, c.Echo().Reverse("getLautern", id))
	default:
		r.addTimelineEvent(fmt.Sprintf("Starting rast %d", rastNum))
		var req ReqPostRasts
		err := c.Bind(&req)
		if err != nil {
			return err
		}
		r.addSummaryRast(req.RealTemperature, req.RealDuration, req.Notes)
		nextRastNum = rastNum + 1
	}
	return c.Render(200, "mash_rasts.html", map[string]interface{}{
		"Title":       "Mash " + r.recipe.Name,
		"Rast":        r.recipe.Mashing.Rasts[rastNum],
		"NextRast":    nextRastNum,
		"MissingRast": len(r.recipe.Mashing.Rasts) - rastNum - 1,
		"Nachguss":    r.recipe.Mashing.Nachguss,
		"RecipeID":    id,
	})
}
