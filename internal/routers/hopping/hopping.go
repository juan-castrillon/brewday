package hopping

import (
	"brewday/internal/recipe"
	"brewday/internal/tools"
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type HoppingRouter struct {
	Store       RecipeStore
	TL          Timeline
	Summary     SummaryRecorder
	recipe      *recipe.Recipe
	ingredients ingredientList
	initialVol  float32
}

// addTimelineEvent adds an event to the timeline
func (r *HoppingRouter) addTimelineEvent(message string) {
	if r.TL != nil {
		r.TL.AddEvent(message)
	}
}

// addSummaryHopping adds a hopping to the summary and notes related to it
func (r *HoppingRouter) addSummaryHopping(name string, amount float32, alpha float32, notes string) {
	if r.Summary != nil {
		r.Summary.AddHopping(name, amount, alpha, notes)
	}
}

// addSummaryMeasuredVolume adds a measured volume to the summary
func (r *HoppingRouter) addSummaryMeasuredVolume(name string, amount float32, notes string) {
	if r.Summary != nil {
		r.Summary.AddMeasuredVolume(name, amount, notes)
	}
}

// addSummaryEvaporation adds an evaporation to the summary
func (r *HoppingRouter) addSummaryEvaporation(amount float32) {
	if r.Summary != nil {
		r.Summary.AddEvaporation(amount)
	}
}

// RegisterRoutes registers the routes for the hopping router
func (r *HoppingRouter) RegisterRoutes(root *echo.Echo, parent *echo.Group) {
	hopping := parent.Group("/hopping")
	hopping.GET("/start/:recipe_id", r.getStartHoppingHandler).Name = "getStartHopping"
	hopping.POST("/start/:recipe_id", r.postStartHoppingHandler).Name = "postStartHopping"
	hopping.GET("/end/:recipe_id", r.getEndHoppingHandler).Name = "getEndHopping"
	hopping.POST("/end/:recipe_id", r.postEndHoppingHandler).Name = "postEndHopping"
	hopping.GET("/hop/:recipe_id/:ingr_num", r.getHoppingHandler).Name = "getHopping"
	hopping.POST("/hop/:recipe_id/:ingr_num", r.postHoppingHandler).Name = "postHopping"
}

// getStartHoppingHandler returns the handler for the start hopping route
func (r *HoppingRouter) getStartHoppingHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return errors.New("no recipe id provided")
	}
	r.addTimelineEvent("Started Hopping")
	re, err := r.Store.Retrieve(id)
	if err != nil {
		return err
	}
	r.recipe = re
	r.ingredients = organizeIngredients(re)
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
		return errors.New("no recipe id provided")
	}
	if r.recipe == nil {
		return errors.New("no recipe loaded")
	}
	r.addTimelineEvent("Start heating up")
	var req ReqPostStartHopping
	err := c.Bind(&req)
	if err != nil {
		return err
	}
	r.addSummaryMeasuredVolume("Measured volume before boiling", req.InitialVolume, req.Notes)
	r.initialVol = req.InitialVolume
	return c.Render(http.StatusOK, "hopping_boiling.html", map[string]interface{}{
		"Title":       "Hopping " + r.recipe.Name,
		"Subtitle":    "2. Boil",
		"RecipeID":    id,
		"Ingredients": r.ingredients,
	})
}

// getEndHoppingHandler returns the handler for the end hopping route
func (r *HoppingRouter) getEndHoppingHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return errors.New("no recipe id provided")
	}
	r.addTimelineEvent("Finished Hopping")
	return c.Render(http.StatusOK, "hopping_end.html", map[string]interface{}{
		"Title":    "Hopping " + r.recipe.Name,
		"Subtitle": "4. Measure volume after boiling",
		"RecipeID": id,
	})
}

// postEndHoppingHandler returns the handler for the end hopping route
func (r *HoppingRouter) postEndHoppingHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return errors.New("no recipe id provided")
	}
	if r.recipe == nil {
		return errors.New("no recipe loaded")
	}
	r.addTimelineEvent("Start heating up")
	var req ReqPostEndHopping
	err := c.Bind(&req)
	if err != nil {
		return err
	}
	r.addSummaryMeasuredVolume("Measured volume after boiling", req.FinalVolume, req.Notes)
	evap := tools.CalculateEvaporation(r.initialVol, req.FinalVolume, r.recipe.Hopping.TotalCookingTime)
	r.addSummaryEvaporation(evap)
	return c.Redirect(http.StatusFound, c.Echo().Reverse("getCooling", id))
}

// getHoppingHandler returns the handler for the hopping route
func (r *HoppingRouter) getHoppingHandler(c echo.Context) error {
	id := c.Param("recipe_id")
	if id == "" {
		return errors.New("no recipe id provided")
	}
	if r.recipe == nil {
		return errors.New("no recipe loaded")
	}
	ingrNumStr := c.Param("ingr_num")
	if ingrNumStr == "" {
		return errors.New("no hop number provided")
	}
	ingrNum, err := strconv.Atoi(ingrNumStr)
	if err != nil {
		return err
	}
	var cookingTime float32
	if ingrNum >= len(r.ingredients) {
		r.addTimelineEvent("Boil finished")
		return c.Redirect(http.StatusFound, c.Echo().Reverse("getEndHopping", id))
	}
	if ingrNum == 0 {
		r.addTimelineEvent("Boil started")
		cookingTime = r.recipe.Hopping.TotalCookingTime
	} else {
		cookingTime = r.ingredients[ingrNum-1].Duration
	}
	ingredient := r.ingredients[ingrNum]
	return c.Render(http.StatusOK, "hopping_hop.html", map[string]interface{}{
		"Title":            "Hopping " + r.recipe.Name,
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
		return errors.New("no recipe id provided")
	}
	if r.recipe == nil {
		return errors.New("no recipe loaded")
	}
	ingrNumStr := c.Param("ingr_num")
	if ingrNumStr == "" {
		return errors.New("no hop number provided")
	}
	ingrNum, err := strconv.Atoi(ingrNumStr)
	if err != nil {
		return err
	}
	if ingrNum < 0 || ingrNum > len(r.ingredients) {
		return errors.New("invalid hop number")
	} else if ingrNum < len(r.recipe.Hopping.Hops) {
		ingredient := r.ingredients[ingrNum]
		var req ReqPostHopping
		err = c.Bind(&req)
		if err != nil {
			return err
		}
		r.addSummaryHopping(ingredient.Name, req.RealAmount, req.RealAlpha, "")
	}
	return c.Redirect(http.StatusFound, c.Echo().Reverse("getHopping", id, ingrNum+1))
}
