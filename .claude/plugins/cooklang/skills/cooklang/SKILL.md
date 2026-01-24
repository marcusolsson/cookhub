---
name: cooklang
description: Parse, create, and work with Cooklang recipe files (.cook). Use when the user asks to create recipes, parse recipe files, work with cooking instructions, generate shopping lists from recipes, or convert recipes to/from Cooklang format. Cooklang is a markup language for recipes that makes ingredients, cookware, and timers machine-readable while keeping recipes human-readable.
---

# Cooklang Skill

Cooklang is a markup language for writing recipes. It uses simple markers (`@` for ingredients, `#` for cookware, `~` for timers) to make recipe components machine-readable while keeping the text human-readable.

## Core Syntax

### Ingredients (`@`)

```cooklang
@ingredient              -- single word, implicit "some" quantity
@multi word ingredient{} -- multi-word with empty braces
@flour{200%g}            -- quantity with unit (% separates value from unit)
@eggs{3}                 -- quantity without unit
@milk{1/2%cup}           -- fractions supported
@butter{1 1/2%tbsp}      -- mixed numbers supported
@water{1-2%cups}         -- ranges supported
@salt{to taste}          -- text quantities
```

### Cookware (`#`)

```cooklang
#pot                     -- single word
#frying pan{}            -- multi-word with braces
#bowl{2}                 -- with quantity
#baking sheet{large}     -- with text descriptor
```

### Timers (`~`)

```cooklang
~{10%minutes}            -- unnamed timer
~baking{30%minutes}      -- named timer
~rest                    -- single word (no duration)
```

### Comments

```cooklang
-- This is a comment (ignored by parser)
@flour{200%g} -- inline comments work too
```

### Sections

```cooklang
= Section Name
== Section With Trailing Equals ==
```

### Notes/Text Blocks

```cooklang
> This is a note or tip that won't be treated as a cooking step.
```

### Metadata

Two formats supported:

**YAML frontmatter (preferred):**
```cooklang
---
title: Chocolate Cake
servings: 8
tags: [dessert, chocolate, baking]
author: Jane Doe <https://example.com>
prep_time: 20 min
cook_time: 45 min
---
```

**Legacy inline metadata:**
```cooklang
>> title: Chocolate Cake
>> servings: 8
```

#### Standard Metadata Keys

| Key | Format | Example |
|-----|--------|---------|
| `title` | string | `"Chocolate Cake"` |
| `description` | string | `"A rich dessert"` |
| `servings` | number or text | `4` or `"2-4 portions"` |
| `tags` | array | `[dessert, easy]` |
| `author` | name + optional URL | `"Jane <https://...>"` |
| `source` | name + optional URL | `"Cookbook <https://...>"` |
| `time` | duration | `"45 min"` |
| `prep_time` | duration | `"15 min"` |
| `cook_time` | duration | `"30 min"` |
| `difficulty` | string | `"easy"` |
| `cuisine` | string | `"Italian"` |
| `diet` | string | `"vegetarian"` |

## Complete Example

```cooklang
---
title: Simple Pasta
servings: 2
tags: [quick, italian, easy]
prep_time: 5 min
cook_time: 15 min
---

= Cooking the Pasta

Bring @water{2%L} to boil in a #large pot.

Add @salt{1%tbsp} and @pasta{200%g}. Cook for ~{10%minutes} until al dente.

= Making the Sauce

While pasta cooks, heat @olive oil{2%tbsp} in a #pan.

Add @garlic{2%cloves}(minced) and sauté for ~{1%minute}.

> Tip: Don't let the garlic burn or it will taste bitter.

Drain pasta, toss with sauce, and serve immediately.
```

## Quantity Syntax Details

```
@ingredient{VALUE%UNIT}
         │     │   │
         │     │   └── Unit (optional): g, ml, cups, tbsp, etc.
         │     └────── Separator: % separates value from unit
         └──────────── Value: number, fraction, range, or text
```

**Value formats:**
- Integer: `3`
- Decimal: `1.5`
- Fraction: `1/2`
- Mixed number: `1 1/2`
- Range: `2-3`
- Text: `some`, `to taste`, `a pinch`

**With advanced units extension**, space can replace `%` when value is numeric:
```cooklang
@water{1 L}     -- equivalent to @water{1%L}
```

## Working with Cooklang Files

### Creating a Recipe

1. Start with metadata (YAML frontmatter recommended)
2. Write cooking steps as natural prose
3. Mark ingredients with `@`, cookware with `#`, timers with `~`
4. Use sections (`=`) to organize longer recipes
5. Add notes (`>`) for tips and variations

### Parsing Considerations

When parsing Cooklang:
- Empty braces `{}` mean "some" quantity
- Missing braces on single word = single word ingredient
- Punctuation terminates single-word components
- Unicode whitespace and punctuation are handled
- Fractions convert to decimals (1/2 → 0.5)

## Extensions

Cooklang supports optional extensions for advanced features. See `references/extensions.md` for:
- **Modifiers**: `@&` (reference), `@-` (hidden), `@?` (optional), `@@` (recipe ref)
- **Aliases**: `@white wine|wine{}` displays as "wine"
- **Intermediate preparations**: Reference previous steps
- **Modes**: Control parsing behavior
- **Range values**: `@eggs{2-4}`

## Additional Resources

- `references/extensions.md` - Complete extension documentation
- `references/metadata.md` - Detailed metadata key specifications
- `references/examples.md` - Full recipe examples for various cuisines
