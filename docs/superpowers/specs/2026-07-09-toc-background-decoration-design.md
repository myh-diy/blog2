# TOC, Global Background Image & Decoration Character Redesign

## Goals
1. Move the post table of contents (TOC) closer to the article body on desktop.
2. Add an admin-uploadable global background image with readable overlay.
3. Redraw the right-side decoration character in a cuter anime-girl style.

## Design

### TOC Placement
- In `PostDetailView.vue`, wrap `<article>` and the TOC in a shared flex container (`max-w-6xl mx-auto gap-8 items-start`).
- Article takes `flex-1 min-w-0`; TOC is a `w-56 sticky top-28 shrink-0` aside.
- TOC is hidden below `xl` breakpoint to keep mobile layout unchanged.
- Remove the previous viewport-fixed `right-6` positioning and `xl:mr-72` offset.

### Global Background Image
- New composable `useBackgroundImage.ts`:
  - Persists image URL in `localStorage` key `blog-bg-image`.
  - Exposes `backgroundImage` ref, `setBackground(url)`, `clearBackground()`.
  - Applies CSS variable `--bg-image` to `:root`.
- `DefaultLayout.vue` renders a fixed full-screen background layer:
  - `backgroundImage: url(...)` when set.
  - Semi-transparent overlay (`bg-white/80 dark:bg-slate-950/80`) so text stays readable.
  - Falls back to existing solid `bg-gray-50 dark:bg-slate-950` when no image is set.
- `AdminView.vue` adds a "Global Background" card:
  - File input using existing `/api/admin/upload-image` endpoint.
  - Live preview of the uploaded/current background.
  - Clear / reset to default button.

### Decoration Character
- Rewrite `PageDecorations.vue` SVG character to an anime-girl style:
  - Large eyes with highlights, soft skin tone, long hair, small smile, blushing cheeks.
  - Theme-aware accents (hair clip / eyes use brand color subtly) but keep the face pleasant on all themes.
  - Keep float/bounce animation and surrounding stars/hearts.

## Files Changed
- `frontend/src/views/PostDetailView.vue`
- `frontend/src/composables/useBackgroundImage.ts` (new)
- `frontend/src/layouts/DefaultLayout.vue`
- `frontend/src/views/AdminView.vue`
- `frontend/src/components/PageDecorations.vue`
