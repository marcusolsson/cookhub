package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/aquilax/cooklang-go"
)

type SchemaOrgRecipe struct {
	Context            string               `json:"@context"`
	Type               string               `json:"@type"`
	Name               string               `json:"name,omitempty"`
	RecipeCategory     string               `json:"recipeCategory,omitempty"`
	RecipeIngredient   []string             `json:"recipeIngredient,omitempty"`
	RecipeInstructions []SchemaOrgHowToStep `json:"recipeInstructions,omitempty"`
	RecipeYield        string               `json:"recipeYield,omitempty"`
	PrepTime           string               `json:"prepTime,omitempty"`
	CookTime           string               `json:"cookTime,omitempty"`
	TotalTime          string               `json:"totalTime,omitempty"`
}

type SchemaOrgHowToStep struct {
	Type  string `json:"@type"`
	Text  string `json:"text,omitempty"`
	Image string `json:"image,omitempty"`
}

func ParseCooklangRecipe(name, content string) (*cooklang.RecipeV2, error) {
	parser := cooklang.NewParserV2(&cooklang.ParseV2Config{})

	return parser.ParseString(content)
}

func ConvertCooklangToSchemaOrg(name string, recipe *cooklang.RecipeV2) SchemaOrgRecipe {
	return SchemaOrgRecipe{
		Context:            "https://schema.org",
		Type:               "Recipe",
		Name:               name,
		RecipeCategory:     getCategory(recipe),
		RecipeIngredient:   getIngredients(recipe),
		RecipeInstructions: getInstructions(recipe),
		RecipeYield:        getRecipeYield(recipe),
	}
}

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

func getInstructions(recipe *cooklang.RecipeV2) []SchemaOrgHowToStep {
	var res []SchemaOrgHowToStep
	var stepBuilder strings.Builder

	for _, step := range recipe.Steps {
		for _, component := range step {
			switch val := component.(type) {
			case cooklang.TextV2:
				stepBuilder.WriteString(val.Value)
				if strings.HasSuffix(val.Value, ".") {
					res = append(res, SchemaOrgHowToStep{
						Type: "HowToStep",
						Text: stepBuilder.String(),
					})
					stepBuilder.Reset()
				}
			case cooklang.IngredientV2:
				stepBuilder.WriteString(val.Name)
			case cooklang.CookwareV2:
				stepBuilder.WriteString(val.Name)
			case cooklang.TimerV2:
				stepBuilder.WriteString(fmt.Sprintf("%v %s", val.Quantity, val.Unit))
			}
		}
	}

	return res
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
