# Cooklang Plugin for Claude Code

A comprehensive plugin for working with the [Cooklang](https://cooklang.org) recipe markup language in Claude Code.

## Features

### Subagent: `cooklang-expert`
A specialized AI agent that automatically activates when working with recipes. Handles:
- Creating well-structured Cooklang recipes
- Parsing and analyzing existing `.cook` files
- Extracting shopping lists
- Scaling recipes
- Converting between formats

### Slash Commands

| Command | Description |
|---------|-------------|
| `/recipe <dish>` | Create a new Cooklang recipe file |
| `/shopping-list [files...]` | Extract combined shopping list from recipes |
| `/scale-recipe <file> <servings>` | Scale a recipe to different servings |
| `/convert-to-cooklang <file or text>` | Convert recipes to Cooklang format |

### Agent Skill: `cooklang`
Comprehensive Cooklang syntax reference that loads automatically when needed.

## Installation

### From GitHub
```bash
claude plugin install https://github.com/your-username/cooklang-plugin
```

### From Local Directory
```bash
claude plugin install ./path/to/cooklang-plugin
```

### Manual Installation
Copy the plugin directory to `~/.claude/plugins/cooklang/`

## Cooklang Quick Reference

```cooklang
---
title: Recipe Name
servings: 4
prep_time: 10 min
cook_time: 20 min
tags: [quick, easy]
---

= Section Name

Add @ingredient{200%g} to #cookware and cook for ~{10%minutes}.

> This is a helpful tip or note.
```

### Markers
- `@` - Ingredients: `@flour{200%g}`, `@eggs{3}`, `@salt{to taste}`
- `#` - Cookware: `#pan`, `#bowl{2}`, `#oven`
- `~` - Timers: `~{10%minutes}`, `~baking{30%min}`
- `=` - Sections: `= Dough`, `== Filling ==`
- `>` - Notes: `> Tip: Let rest before serving`
- `--` - Comments: `-- This is ignored`

### Quantity Formats
- Numbers: `{3}`, `{200%g}`
- Fractions: `{1/2%cup}`, `{1 1/2%tbsp}`
- Ranges: `{2-4}`, `{100-200%ml}`
- Text: `{to taste}`, `{some}`

## Examples

### Create a Recipe
```
/recipe chicken stir fry with vegetables
```

### Generate Shopping List
```
/shopping-list dinner.cook dessert.cook
```

### Scale a Recipe
```
/scale-recipe pasta-carbonara.cook 8
```

### Convert Existing Recipe
```
/convert-to-cooklang grandmas-cookies.txt
```

## License

MIT
