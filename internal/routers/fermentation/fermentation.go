package fermentation

import (
	"brewday/internal/recipe"
	"brewday/internal/routers/common"
	"brewday/internal/tools"
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
		return common.ErrNoRecipeIDProvided
	}
	re, err := r.Store.Retrieve(id)
	if err != nil {
		return err
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
		return common.ErrNoRecipeIDProvided
	}
	if r.recipe == nil {
		return common.ErrNoRecipeLoaded
	}
	var req ReqPostPreFermentation
	err := c.Bind(&req)
	if err != nil {
		return err
	}
	r.addSummaryPreFermentation(req.Volume, req.SG, req.Notes)
	volumeDiff := req.Volume - (r.recipe.BatchSize + 1) // +1 for the 1l of yeast
	sgDiff := r.recipe.InitialSG - req.SG
	if volumeDiff >= 0 && sgDiff >= 0 {
		eff := tools.CalculateEfficiencySG(req.SG, req.Volume, r.recipe.Mashing.GetTotalMaltWeight())
		r.addSummaryEfficiency(eff)
		r.addTimelineEvent("Finished Pre Fermentation")
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
		return common.ErrNoRecipeIDProvided
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
		return common.ErrNoRecipeIDProvided
	}
	if r.recipe == nil {
		return common.ErrNoRecipeLoaded
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
	r.addTimelineEvent("Finished Pre Fermentation")
	return c.Redirect(http.StatusFound, c.Echo().Reverse("getFermentation", id))
}

// getFermentationHandler returns the handler for the fermentation page
func (r *FermentationRouter) getFermentationHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	if r.recipe == nil {
		return common.ErrNoRecipeLoaded
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
		return common.ErrNoRecipeIDProvided
	}
	if r.recipe == nil {
		return common.ErrNoRecipeLoaded
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
		return common.ErrNoRecipeIDProvided
	}
	if r.recipe == nil {
		return common.ErrNoRecipeLoaded
	}
	var hops []recipe.Hops
	for _, h := range r.recipe.Hopping.Hops {
		if h.DryHop {
			hops = append(hops, h)
		}
	}
	r.addTimelineEvent("Finished Day")
	return c.Render(http.StatusOK, "finished_day.html", map[string]interface{}{
		"Title":     "End Fermentation",
		"RecipeID":  id,
		"Subtitle":  "Congratulations, you've finished the brew day!",
		"Hops":      hops,
		"Additions": r.recipe.Fermentation.AdditionalIngredients,
	})
}
