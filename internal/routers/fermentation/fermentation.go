package fermentation

import (
	"brewday/internal/recipe"
	"brewday/internal/tools"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type FermentationRouter struct {
	TL      Timeline
	Summary SummaryRecorder
	Store   RecipeStore
	recipe  *recipe.Recipe
}

// addTimelineEvent adds an event to the timeline
func (r *FermentationRouter) addTimelineEvent(message string) {
	if r.TL != nil {
		r.TL.AddEvent(message)
	}
}

// addSummaryPreFermentation adds a pre fermentation summary
func (r *FermentationRouter) addSummaryPreFermentation(volume, sg float32, notes string) {
	if r.Summary != nil {
		r.Summary.AddSummaryPreFermentation(volume, sg, notes)
	}
}

// addSummaryEfficiency adds an efficiency summary
func (r *FermentationRouter) addSummaryEfficiency(efficiencyPercentage float32) {
	if r.Summary != nil {
		r.Summary.AddEfficiency(efficiencyPercentage)
	}
}

// addSummaryYeastStart adds a yeast start summary
func (r *FermentationRouter) addSummaryYeastStart(temperature, notes string) {
	if r.Summary != nil {
		r.Summary.AddYeastStart(temperature, notes)
	}
}

// closeSummary closes the summary
func (r *FermentationRouter) closeSummary() {
	if r.Summary != nil {
		r.Summary.Close()
	}
}

// registerRoutes registers the routes for the fermentation router
func (r *FermentationRouter) RegisterRoutes(root *echo.Echo, parent *echo.Group) {
	fermentation := parent.Group("/fermentation")
	fermentation.GET("/:recipe_id", r.getFermentationHandler).Name = "getFermentation"
	fermentation.POST("/:recipe_id", r.postFermentationHandler).Name = "postFermentation"
	fermentation.GET("/pre/:recipe_id", r.getPreFermentationHandler).Name = "getPreFermentation"
	fermentation.POST("/pre/:recipe_id", r.postPreFermentationHandler).Name = "postPreFermentation"
	fermentation.GET("/pre/water/:recipe_id", r.getPreFermentationWaterHandler).Name = "getPreFermentationWater"
	fermentation.POST("/pre/water/:recipe_id", r.postPreFermentationWaterHandler).Name = "postPreFermentationWater"
	root.GET("/end/:recipe_id", r.getEndFermentationHandler).Name = "getEndFermentation"
}

// getPreFermentationHandler returns the handler for the pre fermentation page
func (r *FermentationRouter) getPreFermentationHandler(c echo.Context) error {
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
		},
		Hopping: recipe.HopInstructions{
			TotalCookingTime: 3,
			Hops: []recipe.Hops{
				{Name: "Saphir", Amount: 40, Alpha: 4.3, Duration: 2, DryHop: false},
				{Name: "Styrian Celeia", Amount: 25, Alpha: 3.4, Duration: 1, DryHop: false},
				{Name: "Sorachi Ace", Amount: 20, Alpha: 9, Duration: 0, DryHop: false},
				{Name: "Simcoe", Amount: 60, Alpha: 12.9, Duration: 0, DryHop: false},
			},
			AdditionalIngredients: []recipe.AdditionalIngredient{
				{Name: "Demerara Zucker", Amount: 360, Duration: 2},
			},
		},
		Fermentation: recipe.FermentationInstructions{
			Yeast: recipe.Yeast{
				Name:   "WY 1007",
				Amount: 11,
			},
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
	r.addTimelineEvent("Started Pre Fermentation")
	return c.Render(http.StatusOK, "fermentation_pre.html", map[string]interface{}{
		"Title":    "Pre Fermentation",
		"RecipeID": id,
	})
}

// postPreFermentationHandler handles the post request for the pre fermentation page
func (r *FermentationRouter) postPreFermentationHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return errors.New("no recipe id provided")
	}
	if r.recipe == nil {
		return errors.New("no recipe loaded")
	}
	var req ReqPostPreFermentation
	err := c.Bind(&req)
	if err != nil {
		return err
	}
	r.addTimelineEvent("Finished Pre Fermentation")
	r.addSummaryPreFermentation(req.Volume, req.SG, req.Notes)
	volumeDiff := req.Volume - (r.recipe.BatchSize + 1) // +1 for the 1l of yeast
	sgDiff := r.recipe.InitialSG - req.SG
	if volumeDiff >= 0 && sgDiff >= 0 {
		eff := tools.CalculateEfficiencySG(req.SG, req.Volume, r.recipe.Mashing.GetTotalMaltWeight())
		r.addSummaryEfficiency(eff)
		return c.Redirect(http.StatusFound, c.Echo().Reverse("getFermentation", id))
	}
	redirect := "getPreFermentationWater"
	queryParams := fmt.Sprintf("?volumeDiff=%f&sgDiff=%f", volumeDiff, sgDiff)
	return c.Redirect(http.StatusFound, c.Echo().Reverse(redirect, id)+queryParams)
}

// getPreFermentationWaterHandler returns the handler for the pre fermentation water page
func (r *FermentationRouter) getPreFermentationWaterHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return errors.New("no recipe id provided")
	}
	r.addTimelineEvent("Started Pre Fermentation Water")
	volumeDiffRaw := c.QueryParam("volumeDiff")
	sgDiffRaw := c.QueryParam("sgDiff")
	options := []WaterOption{}
	sgDiff, err := strconv.ParseFloat(sgDiffRaw, 32)
	if err != nil {
		return err
	}
	volumeDiff, err := strconv.ParseFloat(volumeDiffRaw, 32)
	if err != nil {
		return err
	}
	currentSG := r.recipe.InitialSG - float32(sgDiff)
	currentVol := r.recipe.BatchSize + float32(volumeDiff) + 1
	if sgDiff < 0.0 {
		toAdd, finalVol := tools.WaterForGravity(currentSG, r.recipe.InitialSG, currentVol)
		options = append(options, WaterOption{
			ToAdd:        toAdd,
			FinalVolume:  finalVol,
			FinalSG:      r.recipe.InitialSG,
			FinalSGPlato: tools.SGToPlato(r.recipe.InitialSG),
		})
		if volumeDiff < 0.0 {
			targetVol := r.recipe.BatchSize + 1
			toAdd, finalSG := tools.WaterForVolume(currentVol, targetVol, currentSG)
			options = append(options, WaterOption{
				ToAdd:        toAdd,
				FinalVolume:  targetVol,
				FinalSG:      finalSG,
				FinalSGPlato: tools.SGToPlato(finalSG),
			})
		}
	}
	return c.Render(http.StatusOK, "fermentation_pre_water.html", map[string]interface{}{
		"Title":         "Pre Fermentation Water",
		"RecipeID":      id,
		"RecipeVolume":  r.recipe.BatchSize + 1,
		"RecipeSG":      r.recipe.InitialSG,
		"CurrentSG":     currentSG,
		"CurrentVolume": currentVol,
		"Options":       options,
	})
}

// postPreFermentationWaterHandler handles the post request for the pre fermentation water page
func (r *FermentationRouter) postPreFermentationWaterHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return errors.New("no recipe id provided")
	}
	if r.recipe == nil {
		return errors.New("no recipe loaded")
	}
	var req ReqPostPreFermentationWater
	err := c.Bind(&req)
	if err != nil {
		return err
	}
	r.addTimelineEvent("Finished Adding Water")
	r.addSummaryPreFermentation(req.FinalVolume, req.FinalSG, req.Notes)
	eff := tools.CalculateEfficiencySG(req.FinalSG, req.FinalVolume, r.recipe.Mashing.GetTotalMaltWeight())
	r.addSummaryEfficiency(eff)
	return c.Redirect(http.StatusFound, c.Echo().Reverse("getFermentation", id))
}

// getFermentationHandler returns the handler for the fermentation page
func (r *FermentationRouter) getFermentationHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return errors.New("no recipe id provided")
	}
	if r.recipe == nil {
		return errors.New("no recipe loaded")
	}
	r.addTimelineEvent("Started Fermentation")
	return c.Render(http.StatusOK, "fermentation.html", map[string]interface{}{
		"Title":       "Fermentation",
		"RecipeID":    id,
		"Yeast":       r.recipe.Fermentation.Yeast,
		"Temperature": r.recipe.Fermentation.Temperature,
	})
}

// postFermentationHandler handles the post request for the fermentation page
func (r *FermentationRouter) postFermentationHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return errors.New("no recipe id provided")
	}
	if r.recipe == nil {
		return errors.New("no recipe loaded")
	}
	var req ReqPostFermentation
	err := c.Bind(&req)
	if err != nil {
		return err
	}
	r.addTimelineEvent("Inserted Yeast")
	r.addSummaryYeastStart(req.Temperature, req.Notes)
	return c.Redirect(http.StatusFound, c.Echo().Reverse("getEndFermentation", id))
}

// getEndFermentationHandler handles the get request for the end fermentation page
func (r *FermentationRouter) getEndFermentationHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return errors.New("no recipe id provided")
	}
	if r.recipe == nil {
		return errors.New("no recipe loaded")
	}
	var hops []recipe.Hops
	for _, h := range r.recipe.Hopping.Hops {
		if h.DryHop {
			hops = append(hops, h)
		}
	}
	r.addTimelineEvent("Finished Day")
	r.closeSummary()
	return c.Render(http.StatusOK, "finished_day.html", map[string]interface{}{
		"Title":     "End Fermentation",
		"RecipeID":  id,
		"Subtitle":  "Congratulations, you've finished the brew day!",
		"Hops":      hops,
		"Additions": r.recipe.Fermentation.AdditionalIngredients,
	})
}
