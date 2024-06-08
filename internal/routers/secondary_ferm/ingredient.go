package secondaryferm

import (
	"brewday/internal/recipe"
	"fmt"
	"html/template"
	"regexp"
	"strings"
	"sync"
)

type ingredient struct {
	// THis serves to distinguish between two dry hops of the same ingredient e.g put some hops and the next day more of the same. SO names would be Hop1(1) and Hop1(2)
	Name             string
	SanitizedName    string
	Amount           float32
	TimeElapsed      float32
	IsHop            bool
	StartClickedOnce bool
	NameJS           template.JS
	// Eventually duration if i decide to support it
}

type ingredientCache struct {
	cache map[string][]ingredient
	lock  sync.Mutex
}

var sanitationRegex = regexp.MustCompile(`\s|[()]`)

func sanitizeName(name string) string {
	return strings.TrimSuffix(sanitationRegex.ReplaceAllLiteralString(name, "_"), "_")
}

func getIngredientList(re *recipe.Recipe) []ingredient {
	result := []ingredient{}
	done := make(map[string]int)
	for _, hop := range re.Hopping.Hops {
		if hop.DryHop {
			done[hop.Name]++
			index := done[hop.Name]
			name := fmt.Sprintf("%s(%d)", hop.Name, index)
			sanitized := sanitizeName(name)
			result = append(result, ingredient{
				Name:          name,
				SanitizedName: sanitized,
				NameJS:        template.JS(sanitized),
				Amount:        hop.Amount,
				IsHop:         true,
			})
		}
	}
	for _, ad := range re.Fermentation.AdditionalIngredients {
		done[ad.Name]++
		index := done[ad.Name]
		name := fmt.Sprintf("%s(%d)", ad.Name, index)
		sanitized := sanitizeName(name)
		result = append(result, ingredient{
			Name:          name,
			SanitizedName: sanitized,
			NameJS:        template.JS(sanitized),
			Amount:        ad.Amount,
			IsHop:         false,
		})
	}
	return result
}

// getIngredients returns the ingredients for the given recipe from the cache
// If the ingredients are not in the cache, it calculates them and stores them in the cache
func (c *ingredientCache) getIngredients(id string, re *recipe.Recipe) []ingredient {
	if c == nil {
		c = &ingredientCache{
			cache: make(map[string][]ingredient),
		}
		c.cache[id] = getIngredientList(re)
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.cache == nil {
		c.cache = make(map[string][]ingredient)
	}
	if ingredients, ok := c.cache[id]; ok {
		return ingredients
	}
	ingredients := getIngredientList(re)
	c.cache[id] = ingredients
	return ingredients
}
