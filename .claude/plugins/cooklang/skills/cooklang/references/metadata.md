# Cooklang Metadata Reference

This document details metadata keys, formats, and parsing behavior.

## Table of Contents

1. [Metadata Formats](#metadata-formats)
2. [Standard Keys](#standard-keys)
3. [Name and URL Parsing](#name-and-url-parsing)
4. [Time Parsing](#time-parsing)
5. [Servings Format](#servings-format)
6. [Tags Format](#tags-format)
7. [Custom Metadata](#custom-metadata)

---

## Metadata Formats

### YAML Frontmatter (Preferred)

Place at the very beginning of the file, between `---` delimiters:

```cooklang
---
title: Recipe Title
servings: 4
tags: [tag1, tag2]
author: Name <https://url.com>
---

Recipe content starts here...
```

**Benefits:**
- YAML standard formatting
- Multi-line values
- Proper arrays and nested structures
- Better editor support

### Legacy Inline Format

```cooklang
>> key: value
>> another key: another value
```

**Note:** Legacy format is deprecated but still supported. Prefer YAML frontmatter for new recipes.

---

## Standard Keys

### title

Recipe name.

```yaml
title: Chocolate Chip Cookies
```

### description / introduction

Recipe description or introduction text.

```yaml
description: A classic cookie recipe that's crispy on the outside and chewy inside.
```

Aliases: `description`, `introduction`

### servings / serves / portions

How many servings the recipe makes.

```yaml
# As number (enables scaling)
servings: 4

# As text (no automatic scaling)
servings: "2-4 portions"
```

Aliases: `servings`, `serves`, `portions`

### tags / tag

Recipe categories or labels.

```yaml
# As array
tags: [dessert, chocolate, cookies, baking]

# As single value
tags: dessert
```

Aliases: `tags`, `tag`

### author

Recipe author with optional URL.

```yaml
# Name only
author: Jane Doe

# Name with URL
author: Jane Doe <https://janescooking.com>

# URL only
author: <https://janescooking.com>
author: https://janescooking.com
```

### source

Recipe source with optional URL.

```yaml
# Name only
source: Mom's Recipe Box

# Name with URL  
source: NYT Cooking <https://cooking.nytimes.com/recipe>

# URL only
source: https://cooking.nytimes.com/recipe
```

### time

Total recipe time.

```yaml
time: 45 min
time: 1 hour 30 minutes
time: 1h 15m
```

### prep_time / prep time

Preparation time before cooking.

```yaml
prep_time: 20 min
```

Aliases: `prep_time`, `prep time`, `prep-time`

### cook_time / cook time

Active cooking time.

```yaml
cook_time: 30 min
```

Aliases: `cook_time`, `cook time`, `cook-time`

### difficulty

Recipe difficulty level.

```yaml
difficulty: easy
difficulty: intermediate  
difficulty: advanced
```

### cuisine

Cuisine type or origin.

```yaml
cuisine: Italian
cuisine: Mexican
cuisine: Japanese
```

### diet

Dietary information.

```yaml
diet: vegetarian
diet: vegan
diet: gluten-free
```

### images

Recipe images.

```yaml
# Single image
images: recipe.jpg

# Multiple images
images: [main.jpg, step1.jpg, step2.jpg]
```

### locale

Language/region for the recipe.

```yaml
locale: en-US
locale: [en, US]
```

---

## Name and URL Parsing

For `author` and `source` fields, URLs are detected in angle brackets:

| Input | Parsed Name | Parsed URL |
|-------|-------------|------------|
| `Name <https://url.com>` | Name | https://url.com |
| `Name` | Name | (none) |
| `https://url.com` | (none) | https://url.com |
| `<https://url.com>` | (none) | https://url.com |
| `Name <invalid>` | Name <invalid> | (none) |

**Examples:**

```yaml
# Full attribution
author: Julia Child <https://juliachildfoundation.org>

# Just name
author: Grandma

# Just URL
source: https://seriouseats.com/recipe

# Complex names work
author: "Rachel R. Peterson <https://rachel.url>"
author: "Mom's Cookbook <https://moms-cookbook.url>"
```

---

## Time Parsing

Time values support various formats:

```yaml
# Minutes
time: 30 min
time: 30 minutes
time: 30m

# Hours
time: 2 hours
time: 2 hour
time: 2h

# Combined
time: 1 hour 30 minutes
time: 1h 30m
time: 1h30m

# Seconds (less common)
time: 90 seconds
time: 90s
```

---

## Servings Format

### Numeric (Enables Scaling)

```yaml
servings: 4
servings: 8
```

When servings is numeric, ingredient quantities can be scaled proportionally.

### Text (No Automatic Scaling)

```yaml
servings: "2-4 portions"
servings: "serves 6-8"
servings: "makes 24 cookies"
```

Text servings describe the yield but don't enable automatic scaling.

---

## Tags Format

### YAML Array

```yaml
tags: [breakfast, eggs, quick, easy]
```

### YAML List

```yaml
tags:
  - breakfast
  - eggs
  - quick
  - easy
```

### Single Tag

```yaml
tags: breakfast
```

### Legacy Format

```cooklang
>> tags: breakfast, eggs, quick
```

---

## Custom Metadata

Any key not in the standard list is preserved as custom metadata:

```yaml
---
title: My Recipe
nutrition: 350 calories per serving
rating: 5 stars
wine_pairing: Pinot Noir
equipment_needed: stand mixer, 9x13 pan
---
```

Custom keys are accessible via the metadata map but don't have special parsing rules.

---

## Metadata Validation

Standard keys may produce warnings if values don't match expected formats:

| Key | Expected | Warning Example |
|-----|----------|-----------------|
| `servings` | number or text | (accepts both) |
| `time` | duration string | "invalid" → warning |
| `tags` | array or string | (accepts both) |
| `author` | name/url string | (accepts any string) |

Warnings don't prevent parsing—they alert to potential issues.
