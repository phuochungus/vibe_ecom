# Repository Guidelines

## Project Structure & Module Organization
This repository is a Vite + React 19 + TypeScript customer storefront. Main source code lives in `src/`. Use `src/pages/` for route-level screens, `src/components/` for reusable UI, `src/services/` for API-facing business calls, `src/lib/` for shared providers and utilities, and `src/types/` for shared types. Static assets belong in `public/` or `src/assets/`. Build output is generated in `dist/` and should not be edited manually.

Routing is defined in `src/router.tsx`. Application providers are wired in `src/App.tsx`, and the Ant Design theme is centralized in `src/theme.ts`. The Vite alias `@` resolves to `src`.

## Build, Test, and Development Commands
- `npm install`: install dependencies.
- `npm run dev`: start the local dev server on `http://localhost:3001`.
- `npm run build`: run TypeScript project checks and create the production bundle in `dist/`.
- `npm run lint`: run ESLint across the repository.
- `npm run preview`: serve the built app locally for a final smoke check.

The frontend expects the backend API at `http://localhost:8080`; `/api` is proxied by Vite during local development.

## Coding Style & Naming Conventions
Follow the existing TypeScript React style: 2-space indentation, single quotes, and no semicolons. Prefer functional components and keep files focused on one responsibility. Use `PascalCase` for component and page symbols, `camelCase` for hooks, helpers, and service functions, and kebab-case filenames such as `product-card.tsx` or `payment-return.tsx`.

Run `npm run lint` before opening a PR. ESLint is configured in `eslint.config.js`; do not bypass warnings without a clear reason.

## Testing Guidelines
There is no automated test runner configured yet. For now, every change should pass `npm run lint` and `npm run build`, plus a manual check of the affected route or flow. If you add tests, colocate them with the feature as `*.test.ts` or `*.test.tsx`.

## Commit & Pull Request Guidelines
Recent history favors short, imperative commit messages with conventional prefixes, especially lowercase `refactor:`. Follow the same pattern when possible, for example `fix: handle expired refresh token`.

PRs should include a concise summary, impacted pages or services, any environment changes, and screenshots for visible UI updates. Link the related issue or task when one exists.

## Security & Configuration Tips
Copy `.env.example` when setting up local configuration. Keep secrets out of source control, and document any new `VITE_` variables in `.env.example`.
