# Cooklang Recipe Examples

Complete examples demonstrating various Cooklang features.

## Table of Contents

1. [Basic Recipe](#basic-recipe)
2. [Recipe with Sections](#recipe-with-sections)
3. [Recipe with All Metadata](#recipe-with-all-metadata)
4. [Recipe with Extensions](#recipe-with-extensions)
5. [Baking Recipe](#baking-recipe)
6. [Multi-Component Recipe](#multi-component-recipe)

---

## Basic Recipe

Simple recipe with core syntax only.

```cooklang
---
title: Scrambled Eggs
servings: 2
time: 10 min
tags: [breakfast, eggs, quick]
---

Crack @eggs{4} into a #bowl and beat with a #fork.

Add @salt{1/4%tsp} and @black pepper{} to taste.

Heat @butter{1%tbsp} in a #non-stick pan over medium heat.

Pour in eggs and cook, stirring gently, for ~{3-4%minutes} until just set.

Serve immediately with toast.
```

---

## Recipe with Sections

Organized recipe using sections for different components.

```cooklang
---
title: Classic Burger
servings: 4
prep_time: 15 min
cook_time: 10 min
tags: [dinner, american, grilling]
---

= Burger Patties

Mix @ground beef{500%g} with @salt{1%tsp} and @pepper{1/2%tsp}.

Form into 4 equal patties, slightly larger than your buns.

Make a small indent in the center of each patty.

= Cooking

Heat #grill or #cast iron pan to high heat.

Cook patties for ~{4%minutes} per side for medium.

> For food safety, internal temperature should reach 71°C.

Add @cheese{4%slices} in the last minute of cooking.

= Assembly

Toast @burger buns{4} cut-side down on the grill.

Layer: bottom bun, @lettuce{4%leaves}, patty with cheese, @tomato{4%slices}, @onion{4%rings}, top bun.

Serve with your favorite condiments.
```

---

## Recipe with All Metadata

Demonstrating comprehensive metadata usage.

```cooklang
---
title: Authentic Margherita Pizza
description: A classic Neapolitan pizza with San Marzano tomatoes and fresh mozzarella. This recipe uses a long cold fermentation for the best flavor.
servings: 4
author: Antonio Rossi <https://authenticitalianrecipes.com>
source: Traditional Neapolitan Recipe <https://pizzanapoletana.org>
prep_time: 30 min
cook_time: 2 min
time: 24 hours
difficulty: intermediate
cuisine: Italian
diet: vegetarian
tags: [pizza, italian, vegetarian, baking]
images: [margherita.jpg, dough-process.jpg]
locale: en-US
---

= Dough (Day Before)

Dissolve @instant yeast{1%tsp} in @warm water{325%ml}.

Mix @bread flour{500%g} and @salt{10%g} in a #large bowl.

Add water mixture and combine until shaggy dough forms.

Knead for ~{10%minutes} until smooth and elastic.

Cover and refrigerate for ~{24%hours}.

= Pizza Assembly

Remove dough from fridge ~{2%hours} before baking.

Divide into 4 balls and let rest at room temperature.

Preheat #pizza stone in oven at maximum temperature for ~{1%hour}.

Stretch each ball into 12-inch round, leaving edges thicker.

Top with @San Marzano tomatoes{200%g}(crushed), @fresh mozzarella{200%g}(torn), @olive oil{2%tbsp}, @fresh basil{8%leaves}.

Bake for ~{90%seconds} to ~{2%minutes} until charred in spots.
```

---

## Recipe with Extensions

Demonstrating extension features.

```cooklang
---
title: Layered Tiramisu
servings: 8
time: 30 min plus chilling
tags: [dessert, italian, no-bake]
---

>> [duplicate]: ref

= Coffee Mixture

Brew @espresso{300%ml} and let cool.

Stir in @coffee liqueur{2%tbsp}(optional) and @sugar{2%tbsp}.

= Mascarpone Cream

Separate @eggs{4} into yolks and whites.

Beat @&eggs|yolks{} with @sugar{100%g} until pale and thick.

> Use pasteurized eggs if concerned about raw eggs.

Add @mascarpone{500%g} and fold gently until smooth.

In a #clean bowl with #electric mixer, beat @&eggs|egg whites{} until stiff peaks form.

Fold the @&(~1)meringue{} into mascarpone mixture.

= Assembly

Quickly dip @ladyfingers{200%g} into coffee mixture.

Layer in a #9x13 inch dish{}.

Spread half the @&(=2)cream{} over ladyfingers.

Repeat with remaining @-ladyfingers and @-&cream.

= Finishing

Cover and refrigerate for ~{4%hours} minimum, preferably overnight.

Dust with @cocoa powder{2%tbsp} before serving.
```

---

## Baking Recipe

Precise measurements typical in baking.

```cooklang
---
title: Classic Chocolate Chip Cookies
servings: 24 cookies
prep_time: 20 min
cook_time: 12 min
difficulty: easy
tags: [dessert, cookies, baking, chocolate]
---

Preheat #oven to 180°C. Line #baking sheets{2} with parchment paper.

= Dry Ingredients

Whisk together @all-purpose flour{280%g}, @baking soda{1%tsp}, and @salt{1/2%tsp} in a #medium bowl.

= Wet Ingredients

Beat @unsalted butter{225%g}(softened) with @granulated sugar{100%g} and @brown sugar{200%g} in #stand mixer until fluffy, about ~{3%minutes}.

Add @eggs{2} one at a time, beating well after each.

Mix in @vanilla extract{2%tsp}.

= Combining

Add dry ingredients to wet in three additions, mixing on low until just combined.

Fold in @chocolate chips{340%g} with a #rubber spatula.

> For best results, chill dough for 30 minutes before baking.

= Baking

Scoop @?dough{2%tbsp} portions onto prepared sheets, spacing 2 inches apart.

Bake for ~{10-12%minutes} until edges are golden but centers look slightly underdone.

Cool on pan for ~{5%minutes}, then transfer to #wire rack.

> Cookies will firm up as they cool.
```

---

## Multi-Component Recipe

Complex recipe with multiple sub-recipes.

```cooklang
---
title: Eggs Benedict
servings: 4
prep_time: 20 min
cook_time: 20 min
difficulty: advanced
cuisine: American
tags: [breakfast, eggs, brunch]
---

= Hollandaise Sauce

> Make this first and keep warm. It cannot be reheated.

Melt @unsalted butter{115%g} in #small saucepan.

In #blender, combine @egg yolks{3}, @lemon juice{1%tbsp}, @salt{1/4%tsp}, and @cayenne pepper{}.

With blender running, slowly drizzle in hot @&butter until thick and emulsified.

Transfer to #small bowl set over #pot of warm water to keep warm.

= Poached Eggs

Fill #large pot with water and bring to gentle simmer.

Add @white vinegar{1%tbsp} to the water.

Crack @eggs{4} individually into #small cups.

Create gentle whirlpool in water, slide in egg.

Poach for ~{3%minutes} for runny yolk.

Remove with #slotted spoon to #plate lined with paper towels.

= Assembly

Toast @English muffins{4}(split) until golden.

Place @Canadian bacon{8%slices} on #skillet over medium heat for ~{1%minute} per side.

Stack: muffin half, 2 slices @&Canadian bacon|bacon{}, @&(=2)poached egg{}, generous spoon of @&(=1)hollandaise{}.

Top with @chives{}(chopped) and @&cayenne pepper{}.

Serve immediately.
```

---

## Output Formats

### Shopping List Format

When generating shopping lists from recipes, combine and group ingredients:

```
PRODUCE
- eggs: 4
- lemon juice: 1 tbsp
- chives: some

DAIRY
- unsalted butter: 115g
- Canadian bacon: 8 slices

BAKERY
- English muffins: 4

PANTRY
- white vinegar: 1 tbsp
- salt: 1/4 tsp
- cayenne pepper: some
```

### Plain Text Format

For human-readable output without markup:

```
Eggs Benedict
Serves 4

HOLLANDAISE SAUCE
Make this first and keep warm. It cannot be reheated.
Melt 115g unsalted butter in small saucepan...
```
