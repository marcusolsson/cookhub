package views

import (
	"fmt"
	"strings"
	"time"

	"github.com/aquilax/cooklang-go"
)

func getCategory(recipe *cooklang.RecipeV2) string {
	if category, ok := recipe.Metadata["course"].(string); ok {
		return category
	}
	return ""
}

func getRecipeYield(recipe *cooklang.RecipeV2) string {
	if servings, ok := recipe.Metadata["servings"].(string); ok {
		return servings
	}
	return ""
}

func getIngredients(recipe *cooklang.RecipeV2) []string {
	var res []string

	for _, section := range recipe.Steps {
		for _, step := range section {
			switch val := step.(type) {
			case cooklang.IngredientV2:
				switch {
				case val.Units == "" && val.Quantity == 0:
					res = append(res, val.Name)
				case val.Units == "":
					res = append(res, fmt.Sprintf("%v %s", val.Quantity, val.Name))
				default:
					res = append(res, fmt.Sprintf("%v %s %s", val.Quantity, val.Units, val.Name))
				}
			}
		}
	}
	return res
}

func ingredientAsString(ingredient cooklang.IngredientV2) string {
	switch {
	case ingredient.Units == "" && ingredient.Quantity == 0:
		return ingredient.Name
	case ingredient.Units == "":
		return fmt.Sprintf("%v %s", ingredient.Quantity, ingredient.Name)
	default:
		return fmt.Sprintf("%v %s %s", ingredient.Quantity, ingredient.Units, ingredient.Name)
	}
}

func timerAsDuration(timer cooklang.TimerV2) time.Duration {
	switch timer.Unit {
	case "s", "sec", "second", "seconds":
		return time.Duration(timer.Quantity) * time.Second
	case "m", "min", "minute", "minutes":
		return time.Duration(timer.Quantity) * time.Minute
	case "h", "hour", "hours":
		return time.Duration(timer.Quantity) * time.Hour
	default:
		return 0
	}
}

func DurationToISO8601(d time.Duration) string {
	if d == 0 {
		return "PT0S"
	}

	var result strings.Builder
	result.WriteString("PT")

	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if hours > 0 {
		result.WriteString(fmt.Sprintf("%dH", hours))
	}
	if minutes > 0 {
		result.WriteString(fmt.Sprintf("%dM", minutes))
	}
	if seconds > 0 {
		result.WriteString(fmt.Sprintf("%dS", seconds))
	}

	return result.String()
}
