---
description: Convert a recipe from plain text, markdown, or other formats to Cooklang. Provide file path or paste text.
allowed-tools: Read, Write
argument-hint: <file-path> or paste recipe text
---

Convert a recipe to Cooklang format: $ARGUMENTS

## Input Handling

1. **If file path provided**: Read the file content
2. **If text provided**: Use the pasted/typed content
3. **Supported input formats**:
   - Plain text recipes
   - Markdown recipes
   - JSON recipe formats
   - HTML (basic parsing)

## Conversion Process

1. **Extract metadata**:
   - Title (from heading or filename)
   - Servings/yield
   - Prep time, cook time
   - Source/author if present
   - Tags from categories/keywords

2. **Identify ingredients**:
   - Look for ingredient lists (usually bulleted or at start)
   - Parse quantities, units, and names
   - Convert to `@name{quantity%unit}` format

3. **Identify cookware**:
   - Find equipment mentions in instructions
   - Convert to `#name{}` format

4. **Identify timers**:
   - Find time durations in instructions
   - Convert to `~name{duration%unit}` format

5. **Structure instructions**:
   - Break into logical sections if recipe has phases
   - Preserve natural language flow
   - Embed components inline where they appear

## Output

Create `<recipe-name>.cook` with:
- Complete YAML frontmatter
- Well-structured Cooklang content
- Notes for any tips/variations found

## Example Conversion

**Input:**
```
Pasta with Garlic
Serves 4

Ingredients:
- 400g spaghetti
- 4 cloves garlic, minced
- 3 tbsp olive oil
- Salt to taste

Boil pasta in salted water for 10 minutes.
Sauté garlic in olive oil until golden.
Toss pasta with garlic oil and serve.
```

**Output:**
```cooklang
---
title: Pasta with Garlic
servings: 4
tags: [pasta, quick, italian]
---

Boil @spaghetti{400%g} in salted water for ~{10%minutes}.

Sauté @garlic{4%cloves}(minced) in @olive oil{3%tbsp} in a #pan until golden.

Toss pasta with garlic oil, add @salt{to taste}, and serve.
```
