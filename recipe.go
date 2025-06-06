package main

import (
	"github.com/aquilax/cooklang-go"
)

func parseCooklangRecipe(content string) (*cooklang.RecipeV2, error) {
	parser := cooklang.NewParserV2(&cooklang.ParseV2Config{})
	return parser.ParseString(content)
}
