---
description: Extract and combine shopping lists from Cooklang recipe files. Optionally specify files or use all .cook files in current directory.
allowed-tools: Read, Glob, Bash
argument-hint: [file1.cook file2.cook ...] or leave empty for all
---

Generate a shopping list from Cooklang recipes: $ARGUMENTS

## Process

1. **Find recipes**:
   - If files specified, use those
   - Otherwise, find all `.cook` files in current directory

2. **Extract ingredients**:
   - Parse each recipe for `@ingredient{quantity%unit}` patterns
   - Handle references (`@&ingredient`) by combining quantities
   - Skip hidden ingredients (`@-ingredient`)
   - Mark optional ingredients (`@?ingredient`)

3. **Combine and group**:
   - Combine same ingredients across recipes
   - Add quantities with matching units
   - Group by category when possible:
     - Produce
     - Dairy & Eggs
     - Meat & Seafood
     - Pantry & Dry Goods
     - Spices & Seasonings
     - Frozen
     - Other

4. **Output format**:
```
SHOPPING LIST
=============
From: [list of recipe titles]
Servings: [if scaling applied]

PRODUCE
- ingredient: quantity unit
- ingredient: quantity unit

DAIRY & EGGS
- ingredient: quantity unit

[etc.]

OPTIONAL
- ingredient: quantity unit (from: recipe name)
```

## Notes

- Combine compatible units (g with g, cups with cups)
- Flag incompatible units that can't be combined
- Include recipe source for each ingredient if from multiple recipes
- Sort alphabetically within categories
