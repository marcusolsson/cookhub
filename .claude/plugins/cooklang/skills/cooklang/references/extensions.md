# Cooklang Extensions Reference

This document details the optional Cooklang syntax extensions that enhance the base specification.

## Table of Contents

1. [Modifiers](#modifiers)
2. [Component Aliases](#component-aliases)
3. [Intermediate Preparations](#intermediate-preparations)
4. [Advanced Units](#advanced-units)
5. [Modes](#modes)
6. [Range Values](#range-values)
7. [Inline Quantities](#inline-quantities)

---

## Modifiers

Modifiers alter ingredient and cookware behavior. Place them between `@`/`#` and the name.

### Reference (`&`)

References a previously defined ingredient. Quantities can be combined.

```cooklang
Add @flour{200%g} to the bowl.
Later, add more @&flour{100%g}.
-- Total flour: 300g
```

### Hidden (`-`)

Ingredient appears inline but not in the shopping list.

```cooklang
Season with @-salt and @-pepper to taste.
```

### Optional (`?`)

Marks ingredient as optional.

```cooklang
Add @?chili flakes{1%tsp} if you like heat.
```

### Recipe Reference (`@@`)

References another recipe file by name.

```cooklang
Serve with @@tomato sauce{200%ml}.
```

### New (`+`)

Forces creation of a new ingredient (used with modes).

```cooklang
>> [duplicate]: ref
@flour{200%g}
@+flour{100%g}  -- creates new entry instead of reference
```

### Combining Modifiers

Modifiers can be combined (except `@@` which replaces `@`):

```cooklang
@-?salt  -- hidden and optional
@&-flour -- reference that's hidden from list
```

---

## Component Aliases

Display a different name than the actual ingredient name using `|`.

```cooklang
@white wine|wine{100%ml}     -- displays as "wine"
@tipo zero flour|flour{200%g} -- displays as "flour"
#stand mixer|mixer{}          -- works for cookware too
```

Useful for references to display shorter names:

```cooklang
@tipo zero flour{500%g}
@&tipo zero flour|flour{200%g}  -- references show "flour"
```

---

## Intermediate Preparations

Reference the result of previous steps as ingredients.

### Relative Step Reference

```cooklang
Mix @flour{200%g} and @water{100%ml} until smooth.

Let the @&(~1)dough{} rest for ~{1%hour}.
-- References what was made 1 step back
```

### Absolute Step Reference

```cooklang
@&(2)mixture{}   -- references step number 2
```

### Section Reference

```cooklang
@&(=2)filling{}   -- references section number 2
@&(=~1)base{}     -- references 1 section back
```

### Rules

- Only past steps from current section can be relatively referenced
- Text steps cannot be referenced
- In relative references, text steps are ignored
- Can only combine with optional (`?`) modifier
- These ingredients don't appear in shopping list

---

## Advanced Units

When enabled, modifies quantity parsing and adds validation.

### Space as Unit Separator

When value is numeric, space can replace `%`:

```cooklang
@water{1 L}      -- same as @water{1%L}
@flour{200 g}    -- same as @flour{200%g}
@eggs{3}         -- still works (no unit)
```

**Note:** Text values still need explicit format:
```cooklang
@salt{to taste}  -- text value, no separator needed
```

### Validation

- Checks unit compatibility between references
- Verifies timers have time units

---

## Modes

Special metadata keys control parsing behavior. Use square brackets.

### Define Mode

Controls how components are parsed.

```cooklang
>> [mode]: ingredients
-- or
>> [define]: components
```

**Values:**
- `all` / `default` - Normal parsing
- `ingredients` / `components` - Only parse components, ignore text
- `steps` - All ingredients are references unless marked with `+`
- `text` - All steps become text blocks

**Use case:** Define ingredient list at recipe start:

```cooklang
>> [mode]: ingredients
@flour{500%g}
@sugar{200%g}
@eggs{3}

>> [mode]: steps
Mix the @&flour with @&sugar...
```

### Duplicate Mode

Controls behavior when same ingredient name appears.

```cooklang
>> [duplicate]: reference
```

**Values:**
- `new` / `default` - Each occurrence creates new ingredient
- `reference` / `ref` - Same names become references automatically

```cooklang
>> [duplicate]: ref
@water{1%cup} @water{2%cups}
-- Equivalent to:
@water{1%cup} @&water{2%cups}
-- Total: 3 cups
```

---

## Range Values

Express approximate or variable quantities.

```cooklang
@eggs{2-4}                    -- 2 to 4 eggs
@tomato sauce{200-300%ml}     -- 200-300 ml
@water{1.5-2%L}               -- decimals work
```

### With References

```cooklang
@flour{100%g}
@&flour{200-400%g}
-- Total: 300-500 g
```

---

## Inline Quantities

Temperatures and other quantities found in text are parsed automatically.

```cooklang
Preheat the #oven to 180°C.
-- 180°C is detected as an inline quantity
```

Supported formats:
- `180°C`, `350°F`
- `180 °C`, `350 °F`
- `180℃`, `350℉`

---

## Extension Compatibility

| Extension | Requires |
|-----------|----------|
| Intermediate preparations | Modifiers |
| Modes | Modifiers (for `+` modifier) |
| Advanced units | None |
| Range values | None |
| Inline quantities | None |
| Aliases | None |
