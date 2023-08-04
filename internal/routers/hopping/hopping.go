package hopping

import (
	"brewday/internal/recipe"
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
	if ingrNum >= len(r.recipe.Hopping.Hops) {
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
	if ingrNum < 0 || ingrNum > len(r.recipe.Hopping.Hops) {
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
