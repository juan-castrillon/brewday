package hopping

import (
	"brewday/internal/recipe"
	"brewday/internal/routers/common"
	"brewday/internal/tools"
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type HoppingRouter struct {
	Store           RecipeStore
	TLStore         TimelineStore
	SummaryStore    SummaryRecorderStore
	ingredientCache map[string]ingredientList
	initialVolCache map[string]float32
}

// addTimelineEvent adds an event to the timeline
func (r *HoppingRouter) addTimelineEvent(id, message string) error {
	if r.TLStore != nil {
		return r.TLStore.AddEvent(id, message)
	}
	return nil
}

// addSummaryHopping adds a hopping to the summary and notes related to it
func (r *HoppingRouter) addSummaryHopping(id string, name string, amount float32, alpha float32, duration float32, notes string) error {
	if r.SummaryStore != nil {
		return r.SummaryStore.AddHopping(id, name, amount, alpha, duration, notes)
	}
	return nil
}

// addSummaryMeasuredVolume adds a measured volume to the summary
func (r *HoppingRouter) addSummaryMeasuredVolume(id string, name string, amount float32, notes string) error {
	if r.SummaryStore != nil {
		return r.SummaryStore.AddMeasuredVolume(id, name, amount, notes)
	}
	return nil
}

// addSummaryEvaporation adds an evaporation to the summary
func (r *HoppingRouter) addSummaryEvaporation(id string, amount float32) error {
	if r.SummaryStore != nil {
		return r.SummaryStore.AddEvaporation(id, amount)
	}
	return nil
}

// storeIngredients stores the ingredients in the router
// It initializes the ingredients map if it is nil
func (r *HoppingRouter) storeIngredients(id string, re *recipe.Recipe) {
	if r.ingredientCache == nil {
		r.ingredientCache = make(map[string]ingredientList)
	}
	r.ingredientCache[id] = organizeIngredients(re)
}

// getIngredients returns the ingredients for the given recipe from the cache
// If the ingredients are not in the cache, it calculates them and stores them in the cache
func (r *HoppingRouter) getIngredients(id string, re *recipe.Recipe) ingredientList {
	ings, ok := r.ingredientCache[id]
	if !ok {
		r.storeIngredients(id, re)
		ings = r.ingredientCache[id]
	}
	return ings
}

// storeInitialVolume stores the initial volume in the router
// It initializes the initial volume map if it is nil
func (r *HoppingRouter) storeInitialVolume(id string, vol float32) {
	if r.initialVolCache == nil {
		r.initialVolCache = make(map[string]float32)
	}
	r.initialVolCache[id] = vol
}

// getInitialVolume returns the initial volume for the given recipe from the cache
// If the initial volume is not in the cache, it returns an error
func (r *HoppingRouter) getInitialVolume(id string) (float32, error) {
	vol, ok := r.initialVolCache[id]
	if !ok {
		return 0, errors.New("initial volume not found for recipe")
	}
	return vol, nil
}

// RegisterRoutes registers the routes for the hopping router
func (r *HoppingRouter) RegisterRoutes(root *echo.Echo, parent *echo.Group) {
	hopping := parent.Group("/hopping")
	hopping.GET("/start/:recipe_id", r.getStartHoppingHandler).Name = "getStartHopping"
	hopping.POST("/start/:recipe_id", r.postStartHoppingHandler).Name = "postStartHopping"
	hopping.GET("/boil/:recipe_id", r.getBoilingHandler).Name = "getBoiling"
	hopping.GET("/end/:recipe_id", r.getEndHoppingHandler).Name = "getEndHopping"
	hopping.POST("/end/:recipe_id", r.postEndHoppingHandler).Name = "postEndHopping"
	hopping.GET("/hop/:recipe_id/:ingr_num", r.getHoppingHandler).Name = "getHopping"
	hopping.POST("/hop/:recipe_id/:ingr_num", r.postHoppingHandler).Name = "postHopping"
}

// getStartHoppingHandler returns the handler for the start hopping route
func (r *HoppingRouter) getStartHoppingHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	err := r.addTimelineEvent(id, "Started Hopping")
	if err != nil {
		log.Error().Str("id", id).Err(err).Msg("could not add timeline event")
	}
	re, err := r.Store.Retrieve(id)
	if err != nil {
		return err
	}
	re.SetStatus(recipe.RecipeStatusBoiling, "initialVol")
	return c.Render(http.StatusOK, "hopping_start.html", map[string]interface{}{
		"Title":    "Hopping " + re.Name,
		"Subtitle": "1. Measure volume before boiling",
		"RecipeID": id,
	})
}

// postStartHoppingHandler returns the handler for the start hopping route
func (r *HoppingRouter) postStartHoppingHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	err := r.addTimelineEvent(id, "Start heating up")
	if err != nil {
		log.Error().Str("id", id).Err(err).Msg("could not add timeline event")
	}
	var req ReqPostStartHopping
	err = c.Bind(&req)
	if err != nil {
		return err
	}
	err = r.addSummaryMeasuredVolume(id, "Measured volume before boiling", req.InitialVolume, req.Notes)
	if err != nil {
		log.Error().Str("id", id).Err(err).Msg("could not add measured volume to summary")
	}
	r.storeInitialVolume(id, req.InitialVolume)
	return c.Redirect(http.StatusFound, c.Echo().Reverse("getBoiling", id))
}

// getBoilingHandler returns the handler for the start hopping route
func (r *HoppingRouter) getBoilingHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	re, err := r.Store.Retrieve(id)
	if err != nil {
		return err
	}
	ings := r.getIngredients(id, re)
	re.SetStatus(recipe.RecipeStatusBoiling, "beforeBoil")
	return c.Render(http.StatusOK, "hopping_boiling.html", map[string]interface{}{
		"Title":       "Hopping " + re.Name,
		"Subtitle":    "2. Boil",
		"RecipeID":    id,
		"Ingredients": ings,
	})
}

// getHoppingHandler returns the handler for the hopping route
func (r *HoppingRouter) getHoppingHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	re, err := r.Store.Retrieve(id)
	if err != nil {
		return err
	}
	ingrNumStr := c.Param("ingr_num")
	if ingrNumStr == "" {
		return errors.New("no hop number provided")
	}
	ingrNum, err := strconv.Atoi(ingrNumStr)
	if err != nil {
		return err
	}
	ings := r.getIngredients(id, re)
	var cookingTime float32
	if ingrNum > len(ings) {
		return c.Redirect(http.StatusFound, c.Echo().Reverse("getEndHopping", id))
	}
	if ingrNum == len(ings) {
		if ings[ingrNum-1].Duration != 0 {
			re.SetStatus(recipe.RecipeStatusBoiling, "lastBoil", ingrNum)
			return c.Render(http.StatusOK, "hopping_last_boil.html", map[string]interface{}{
				"Title":       "Hopping " + re.Name,
				"Subtitle":    "3. Add hops",
				"RecipeID":    id,
				"BoilingTime": ings[ingrNum-1].Duration,
				"IngrNum":     ingrNum,
			})
		} else {
			return c.Redirect(http.StatusFound, c.Echo().Reverse("getEndHopping", id))
		}
	}
	if ingrNum == 0 {
		cookingTime = re.Hopping.TotalCookingTime
	} else {
		cookingTime = ings[ingrNum-1].Duration
	}
	ingredient := ings[ingrNum]
	re.SetStatus(recipe.RecipeStatusBoiling, "hop", ingrNum)
	return c.Render(http.StatusOK, "hopping_hop.html", map[string]interface{}{
		"Title":            "Hopping " + re.Name,
		"Subtitle":         "3. Add hops",
		"RecipeID":         id,
		"Ingredient":       ingredient,
		"IngrNum":          ingrNum,
		"TotalCookingTime": cookingTime,
	})
}

// postHoppingHandler returns the handler for the hopping route
func (r *HoppingRouter) postHoppingHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	re, err := r.Store.Retrieve(id)
	if err != nil {
		return err
	}
	ingrNumStr := c.Param("ingr_num")
	if ingrNumStr == "" {
		return errors.New("no hop number provided")
	}
	ingrNum, err := strconv.Atoi(ingrNumStr)
	if err != nil {
		return err
	}
	ings := r.getIngredients(id, re)
	if ingrNum < 0 || ingrNum > len(ings) {
		return errors.New("invalid hop number")
	} else if ingrNum <= len(ings) {
		var req ReqPostHopping
		err = c.Bind(&req)
		if err != nil {
			return err
		}
		if ingrNum == len(ings) {
			err = r.addTimelineEvent(id, "Finished hopping boiling time")
			if err != nil {
				log.Error().Err(err).Str("id", id).Msg("could not add timeline event")
			}
		} else {
			ingredient := ings[ingrNum]
			err = r.addTimelineEvent(id, "Added "+ingredient.Name)
			if err != nil {
				log.Error().Err(err).Str("id", id).Msg("could not add timeline event")
			}
			err = r.addSummaryHopping(id, ingredient.Name, req.RealAmount, req.RealAlpha, req.RealDuration, "")
			if err != nil {
				log.Error().Str("id", id).Err(err).Msg("could not add hopping to summary")
			}
		}
	}
	return c.Redirect(http.StatusFound, c.Echo().Reverse("getHopping", id, ingrNum+1))
}

// getEndHoppingHandler returns the handler for the end hopping route
func (r *HoppingRouter) getEndHoppingHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	re, err := r.Store.Retrieve(id)
	if err != nil {
		return err
	}
	_, err = r.getInitialVolume(id)
	if err != nil {
		return err
	}
	re.SetStatus(recipe.RecipeStatusBoiling, "finalVol")
	err = r.addTimelineEvent(id, "Finished Hopping")
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("could not add timeline event")
	}
	return c.Render(http.StatusOK, "hopping_end.html", map[string]interface{}{
		"Title":    "Hopping " + re.Name,
		"Subtitle": "4. Measure volume after boiling",
		"RecipeID": id,
	})
}

// postEndHoppingHandler returns the handler for the end hopping route
func (r *HoppingRouter) postEndHoppingHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return common.ErrNoRecipeIDProvided
	}
	re, err := r.Store.Retrieve(id)
	if err != nil {
		return err
	}
	err = r.addTimelineEvent(id, "Boil finished")
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("could not add timeline event")
	}
	var req ReqPostEndHopping
	err = c.Bind(&req)
	if err != nil {
		return err
	}
	initialVol, err := r.getInitialVolume(id)
	if err != nil {
		return err
	}
	err = r.addSummaryMeasuredVolume(id, "Measured volume after boiling", req.FinalVolume, req.Notes)
	if err != nil {
		log.Error().Str("id", id).Err(err).Msg("could not add measured volume to summary")
	}
	re.SetHotWortVolume(req.FinalVolume)
	evap := tools.CalculateEvaporation(initialVol, req.FinalVolume, re.Hopping.TotalCookingTime)
	err = r.addSummaryEvaporation(id, evap)
	if err != nil {
		log.Error().Str("id", id).Err(err).Msg("could not add evaporation to summary")
	}
	return c.Redirect(http.StatusFound, c.Echo().Reverse("getCooling", id))
}
