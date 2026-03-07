# Welcome to [Slidev](https://github.com/slidevjs/slidev)!

To start the slide show:

- `pnpm install`
- `pnpm dev`
- visit <http://localhost:3030>

Edit the [slides.md](./slides.md) to see the changes.

## Multi-Deck Setup

- Index deck (dev): `pnpm dev`
- Scoped deck 1 (dev): `npm run dev:deck:1`
- Scoped deck 2 (dev): `npm run dev:deck:2`

Build outputs:

- Index deck only: `npm run build:index` -> `dist/`
- Scoped decks only: `npm run build:decks` -> `dist/example-deck-1` and `dist/example-deck-2`
- Index + scoped decks: `npm run build` (same as `npm run build:all`)

## Deploy To Cloudflare Pages

- Root directory: `slides` (if your Pages project points to the monorepo root)
- Build command: `npm run build`
- Build output directory: `dist`
- SPA routing is handled by `public/_redirects`, including sub-routes for `/example-deck-1/*` and `/example-deck-2/*`

If you deploy with Wrangler:

- `npx wrangler pages deploy dist --project-name <your-project-name>`

Learn more about Slidev at the [documentation](https://sli.dev/).
