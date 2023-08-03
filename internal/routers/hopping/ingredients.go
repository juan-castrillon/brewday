package hopping

import (
	"brewday/internal/recipe"
	"sort"
)

// ingredient represents an ingredient to add
type ingredient struct {
	Name     string
	Amount   float32
	Duration float32
	Alpha    float32 // only for hops
	IsHop    bool
}

// ingredientList is a list of ingredients that implements the sort interface
// it is sorted by duration
type ingredientList []ingredient

func (l ingredientList) Len() int {
	return len(l)
}

func (l ingredientList) Less(i, j int) bool {
	return l[i].Duration >= l[j].Duration
}

func (l ingredientList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

// organizeIngredients organizes the ingredients by time of addition
func organizeIngredients(re *recipe.Recipe) ingredientList {
	var ings ingredientList
	for _, h := range re.Hopping.Hops {
		if !h.DryHop && !h.Vorderwuerze {
			ings = append(ings, ingredient{
				Name:     h.Name,
				Amount:   h.Amount,
				Duration: h.Duration,
				Alpha:    h.Alpha,
				IsHop:    true,
			})
		}
	}
	for _, a := range re.Hopping.AdditionalIngredients {
		ings = append(ings, ingredient{
			Name:     a.Name,
			Amount:   a.Amount,
			Duration: a.Duration,
			IsHop:    false,
		})
	}
	sort.Sort(ings)
	return ings
}
