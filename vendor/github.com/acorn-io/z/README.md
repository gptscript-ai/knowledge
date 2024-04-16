# z

`z` exposes a curated set of utility functions.

## Goals 

Its primary goals are two-fold:

1. Make the best commonly used utilities discoverable; "Ugh! What's the import for that helper again?"
2. Reduce hand strain; "I can finally stop writing `func must(...)` everywhere!"

## Why call it 'z'?

First off, it's a much better name than 'y'.

Besides that, the choice in name is supported by a few rules-of-thumb:

- A smaller package name is more ergonomic (`z.Pointer("foo")` vs. `utils.Pointer("foo")`)
- The last letter of the alphabet is _probably definitely_ less likely to be shadowed by dependents 
- It's memorable (and fun to say out loud while typing!)

## Adding to this Module

This module has two distinct classes of utilities, each with their own purpose and threshold for acceptance.

### A) External

_External_ utilities are curated components of external packages that are re-exported by `z` to make them more easily discoverable.

The bar for accepting PRs adding external utilities should be **low**, since the source _is not_ maintained in `z`.

### B) Local

_Local_ utilities are home baked, and should consist of the most often reimplemented helpers. 

The bar for accepting PRs adding local utilities should be **high**, since the source _is_ maintained in `z`.

