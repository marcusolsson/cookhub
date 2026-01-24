package views

import (
	"cmp"
	"fmt"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/aquilax/cooklang-go"
)

type RecipeMetadata struct {
	Recipe *cooklang.RecipeV2

	normalized map[string]any
}

func ParseCanonicalMetadata(recipe *cooklang.RecipeV2) *RecipeMetadata {
	normalized := make(map[string]any)

	for key, value := range recipe.Metadata {
		normalized[strings.ToLower(key)] = value
	}

	return &RecipeMetadata{
		Recipe:     recipe,
		normalized: normalized,
	}
}

func (rm RecipeMetadata) Title(path string) string {
	return cmp.Or(
		rm.getStringProperty("title"),
		rm.getStringProperty("name"),
		strings.TrimSuffix(
			filepath.Base(path),
			filepath.Ext(path),
		),
	)
}

func (rm *RecipeMetadata) Description() string {
	return cmp.Or(
		rm.getStringProperty("introduction"),
		rm.getStringProperty("description"),
	)
}

func (rm *RecipeMetadata) ImageURL() string {
	return cmp.Or(
		rm.getStringProperty("image"),
		rm.getStringProperty("picture"),
	)
}

func (rm *RecipeMetadata) Servings() float64 {
	return cmp.Or(
		rm.getNumericProperty("servings"),
		rm.getNumericProperty("serves"),
		rm.getNumericProperty("yield"),
	)
}

func (rm *RecipeMetadata) Cuisine() string {
	return rm.getStringProperty("cuisine")
}

func (rm *RecipeMetadata) Course() string {
	return rm.getStringProperty("course")
}

func (rm *RecipeMetadata) Category() string {
	return rm.getStringProperty("category")
}

func (rm *RecipeMetadata) Source() *url.URL {
	source := rm.getStringProperty("source")

	if source == "" {
		return nil
	}

	u, err := url.Parse(source)
	if err != nil {
		return nil
	}

	return u
}

func (rm *RecipeMetadata) Tags() []string {
	if stringTags, ok := rm.Recipe.Metadata["tags"].(string); ok {
		return strings.Split(stringTags, ",")
	}

	if anyTags, ok := rm.Recipe.Metadata["tags"].([]any); ok {
		var tags []string
		for _, tag := range anyTags {
			if str, ok := tag.(string); ok {
				tags = append(tags, str)
			}
		}
		return tags
	}

	return []string{}
}

func (rm *RecipeMetadata) TotalTime() string {
	return cmp.Or(
		rm.getStringProperty("time required"),
		rm.getStringProperty("time"),
		rm.getStringProperty("duration"),
		rm.getStringProperty("total time"),
	)
}

func (rm *RecipeMetadata) PrepTime() string {
	return cmp.Or(
		rm.getStringProperty("prep time"),
	)
}

func (rm *RecipeMetadata) CookTime() string {
	return cmp.Or(
		rm.getStringProperty("cook time"),
	)
}

func (rm *RecipeMetadata) getStringProperty(name string) string {
	if str, ok := rm.normalized[name].(string); ok {
		return str
	}
	return ""
}

func (rm *RecipeMetadata) getNumericProperty(name string) float64 {
	if val, ok := rm.normalized[name].(int); ok {
		return float64(val)
	}
	if val, ok := rm.normalized[name].(float64); ok {
		return val
	}
	if val, ok := rm.normalized[name].(string); ok {
		v, err := strconv.ParseFloat(val, 64)
		if err == nil {
			return v
		}
	}
	return 0.0
}

func IngredientAsString(ingredient cooklang.IngredientV2) string {
	switch {
	case ingredient.Units == "" && ingredient.Quantity == 0:
		return ingredient.Name
	case ingredient.Units == "":
		return fmt.Sprintf("%v %s", ingredient.Quantity, ingredient.Name)
	default:
		return fmt.Sprintf("%v %s %s", ingredient.Quantity, ingredient.Units, ingredient.Name)
	}
}

func TimerAsDuration(timer cooklang.TimerV2) time.Duration {
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

// Slugify converts a string to a URL-friendly slug.
// "Fluffy Pancake" becomes "fluffy-pancake"
func Slugify(s string) string {
	var result strings.Builder

	prevDash := false
	for _, r := range strings.ToLower(s) {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			result.WriteRune(r)
			prevDash = false
		case unicode.IsSpace(r) || r == '-' || r == '_':
			if !prevDash && result.Len() > 0 {
				result.WriteRune('-')
				prevDash = true
			}
		}
	}

	// Trim trailing dash
	str := result.String()
	return strings.TrimSuffix(str, "-")
}
