# Task 2 Fix Report

## Status: DONE

## Changes Made

### 1. Added `.gitkeep` to empty directories (Important)
Empty directories are not tracked by git. Added `.gitkeep` files to:
- `frontend/src/views/.gitkeep`
- `frontend/src/stores/.gitkeep`
- `frontend/src/router/.gitkeep`
- `frontend/src/utils/.gitkeep`

### 2. Fixed index.html title (Minor)
Changed `<title>frontend</title>` to `<title>Blog</title>` in `frontend/index.html`.

### 3. Removed Vite scaffold leftovers (Minor)
Removed the following scaffold files:
- `frontend/src/components/HelloWorld.vue`
- `frontend/src/assets/hero.png`
- `frontend/src/assets/vite.svg`
- `frontend/src/assets/vue.svg`

## Verification
- All four `.gitkeep` files exist in their respective directories
- `index.html` title correctly displays "Blog"
- `HelloWorld.vue` and asset demo files no longer exist
