---
description: Create a new Cooklang recipe file. Provide a dish name or description.
allowed-tools: Write, Read
argument-hint: <dish name or description>
---

Create a Cooklang recipe based on the user's request: $ARGUMENTS

## Requirements

1. **File format**: Create a `.cook` file with valid Cooklang syntax
2. **Metadata**: Include complete YAML frontmatter:
   - title, servings, prep_time, cook_time
   - tags (at least 2-3 relevant tags)
   - difficulty (easy/intermediate/advanced)
   - cuisine (if applicable)

3. **Structure**:
   - Use sections (`=`) for recipes with distinct phases
   - Write natural prose with embedded `@ingredients`, `#cookware`, and `~timers`
   - Include helpful notes (`>`) for tips

4. **Components**:
   - All ingredients marked with `@name{quantity%unit}`
   - All cookware marked with `#name{}`
   - All timing marked with `~name{duration%unit}`

5. **Quality**:
   - Realistic quantities and times
   - Clear, actionable instructions
   - Proper unit consistency (metric preferred)

## Output

Save the recipe to `<dish-name>.cook` in the current directory.
Confirm the file was created and summarize the recipe.
