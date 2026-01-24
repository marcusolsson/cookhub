---
description: Scale a Cooklang recipe to different servings. Specify file and target servings.
allowed-tools: Read, Write, Edit
argument-hint: <recipe.cook> <target-servings>
---

Scale a Cooklang recipe to the specified number of servings: $ARGUMENTS

## Process

1. **Read the recipe** and extract current servings from metadata

2. **Calculate scaling factor**:
   - factor = target_servings / current_servings

3. **Scale quantities**:
   - Multiply numeric values by the scaling factor
   - Round to sensible precision:
     - Whole numbers for items (eggs, cloves)
     - One decimal for most measurements
     - Fractions where appropriate (1/2, 1/4, 3/4)
   - Preserve text quantities unchanged ("to taste", "some")
   - Handle ranges by scaling both ends

4. **Update metadata**:
   - Update `servings:` field
   - Keep all other metadata unchanged

5. **Output**:
   - Save scaled recipe to `<original-name>-scaled.cook`
   - Or overwrite original if user confirms

## Scaling Rules

| Original | Factor | Result |
|----------|--------|--------|
| 2 | 2x | 4 |
| 1/2 | 2x | 1 |
| 200g | 0.5x | 100g |
| 1-2 | 2x | 2-4 |
| "some" | any | "some" |
| "to taste" | any | "to taste" |

## Notes

- Warn if scaling creates unusual quantities (0.3 eggs)
- Suggest rounding to practical amounts
- Timers typically don't scale (cooking time ≠ servings)
- Some ingredients scale non-linearly (salt, spices) - note this
