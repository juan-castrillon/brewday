package fermentation

import (
	"brewday/internal/recipe"
	"brewday/internal/routers/common"
	"brewday/internal/tools"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type FermentationRouter struct {
	TL           Timeline
	SummaryStore SummaryRecorderStore
	Store        RecipeStore
}

// addTimelineEvent adds an event to the timeline
func (r *FermentationRouter) addTimelineEvent(message string) {
	if r.TL != nil {
		r.TL.AddEvent(message)
	}
}

// addSummaryPreFermentation adds a pre fermentation summary
func (r *FermentationRouter) addSummaryPreFermentation(id string, volume, sg float32, notes string) error {
	if r.SummaryStore != nil {
		return r.SummaryStore.AddSummaryPreFermentation(id, volume, sg, notes)
	}
	return nil
}

// addSummaryEfficiency adds an efficiency summary
func (r *FermentationRouter) addSummaryEfficiency(id string, efficiencyPercentage float32) error {
	if r.SummaryStore != nil {
		return r.SummaryStore.AddEfficiency(id, efficiencyPercentage)
	}
	return nil
}

// addSummaryYeastStart adds a yeast start summary
func (r *FermentationRouter) addSummaryYeastStart(id string, temperature, notes string) error {
	if r.SummaryStore != nil {
		return r.SummaryStore.AddYeastStart(id, temperature, notes)
	}
	return nil
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
	r.addTimelineEvent("Started Pre Fermentation")
	re.SetStatus(recipe.RecipeStatusPreFermentation, "measure")
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
	re, err := r.Store.Retrieve(id)
	if err != nil {
		return err
	}
	var req ReqPostPreFermentation
	err = c.Bind(&req)
	if err != nil {
		return err
	}
	err = r.addSummaryPreFermentation(id, req.Volume, req.SG, req.Notes)
	if err != nil {
		log.Error().Str("id", id).Err(err).Msg("could not add pre fermentation summary")
	}
	volumeDiff := req.Volume - (re.BatchSize + 1) // +1 for the 1l of yeast
	sgDiff := re.InitialSG - req.SG
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
	re, err := r.Store.Retrieve(id)
	if err != nil {
		return err
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
	currentSG := re.InitialSG - float32(sgDiff)
	currentVol := re.BatchSize + float32(volumeDiff) + 1
	if sgDiff < 0.0 {
		toAdd, finalVol := tools.WaterForGravity(currentSG, re.InitialSG, currentVol)
		options = append(options, WaterOption{
			ToAdd:        toAdd,
			FinalVolume:  finalVol,
			FinalSG:      re.InitialSG,
			FinalSGPlato: tools.SGToPlato(re.InitialSG),
		})
		if volumeDiff < 0.0 {
			targetVol := re.BatchSize + 1
			toAdd, finalSG := tools.WaterForVolume(currentVol, targetVol, currentSG)
			options = append(options, WaterOption{
				ToAdd:        toAdd,
				FinalVolume:  targetVol,
				FinalSG:      finalSG,
				FinalSGPlato: tools.SGToPlato(finalSG),
			})
		}
	}
	re.SetStatus(recipe.RecipeStatusPreFermentation, "water", volumeDiff, sgDiff)
	return c.Render(http.StatusOK, "fermentation_pre_water.html", map[string]interface{}{
		"Title":         "Pre Fermentation Water",
		"RecipeID":      id,
		"RecipeVolume":  re.BatchSize + 1,
		"RecipeSG":      re.InitialSG,
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
	re, err := r.Store.Retrieve(id)
	if err != nil {
		return err
	}
	var req ReqPostPreFermentationWater
	err = c.Bind(&req)
	if err != nil {
		return err
	}
	r.addTimelineEvent("Finished Adding Water")
	err = r.addSummaryPreFermentation(id, req.FinalVolume, req.FinalSG, req.Notes)
	if err != nil {
		log.Error().Str("id", id).Err(err).Msg("could not add pre fermentation summary")
	}
	eff := tools.CalculateEfficiencySG(req.FinalSG, req.FinalVolume, re.Mashing.GetTotalMaltWeight())
	err = r.addSummaryEfficiency(id, eff)
	if err != nil {
		log.Error().Str("id", id).Err(err).Msg("could not add efficiency to summary")
	}
	r.addTimelineEvent("Finished Pre Fermentation")
	return c.Redirect(http.StatusFound, c.Echo().Reverse("getFermentation", id))
}

// getFermentationHandler returns the handler for the fermentation page
func (r *FermentationRouter) getFermentationHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	re, err := r.Store.Retrieve(id)
	if err != nil {
		return err
	}
	r.addTimelineEvent("Started Fermentation")
	re.SetStatus(recipe.RecipeStatusFermenting, "start")
	return c.Render(http.StatusOK, "fermentation.html", map[string]interface{}{
		"Title":       "Fermentation",
		"RecipeID":    id,
		"Yeast":       re.Fermentation.Yeast,
		"Temperature": re.Fermentation.Temperature,
	})
}

// postFermentationHandler handles the post request for the fermentation page
func (r *FermentationRouter) postFermentationHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	var req ReqPostFermentation
	err := c.Bind(&req)
	if err != nil {
		return err
	}
	r.addTimelineEvent("Inserted Yeast")
	err = r.addSummaryYeastStart(id, req.Temperature, req.Notes)
	if err != nil {
		log.Error().Str("id", id).Err(err).Msg("could not add yeast start to summary")
	}
	return c.Redirect(http.StatusFound, c.Echo().Reverse("getEndFermentation", id))
}

// getEndFermentationHandler handles the get request for the end fermentation page
func (r *FermentationRouter) getEndFermentationHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	re, err := r.Store.Retrieve(id)
	if err != nil {
		return err
	}
	var hops []recipe.Hops
	for _, h := range re.Hopping.Hops {
		if h.DryHop {
			hops = append(hops, h)
		}
	}
	r.addTimelineEvent("Finished Day")
	re.SetStatus(recipe.RecipeStatusFinished)
	return c.Render(http.StatusOK, "finished_day.html", map[string]interface{}{
		"Title":     "End Fermentation",
		"RecipeID":  id,
		"Subtitle":  "Congratulations, you've finished the brew day!",
		"Hops":      hops,
		"Additions": re.Fermentation.AdditionalIngredients,
	})
}
