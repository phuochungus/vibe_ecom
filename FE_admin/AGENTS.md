# Repository Guidelines

## Project Structure & Module Organization
This repository is a Vite + React + TypeScript admin frontend. Application bootstrap lives in `src/main.tsx`, with routing in `src/router.tsx`. Keep feature screens under `src/pages/` (for example `src/pages/products/product-list.tsx`), shared layout in `src/components/layout/`, and reusable shadcn/ui primitives in `src/components/ui/`. Put API clients and auth helpers in `src/lib/` and `src/services/`, shared types in `src/types/`, hooks in `src/hooks/`, static assets in `src/assets/`, and public files in `public/`.

## Build, Test, and Development Commands
Use `npm install` to install dependencies. Run `npm run dev` to start the Vite dev server; `/api` requests are proxied to `http://localhost:8080`. Use `npm run build` to run `tsc -b` and create a production build, `npm run preview` to serve the built app locally, and `npm run lint` to run ESLint across the project. For routine changes, run lint and build before opening a PR.

## Coding Style & Naming Conventions
Follow the existing code style: TypeScript, 2-space indentation, single quotes, and no semicolons. Prefer the `@/` import alias over long relative paths. Use PascalCase for React components and page exports (`DashboardPage`), kebab-case for most feature filenames (`admin-layout.tsx`, `order-detail.tsx`), and camelCase for variables, helpers, and service methods. Respect strict TypeScript settings; do not leave unused locals or parameters behind.

## Testing Guidelines
There is no automated test runner configured in `package.json` yet, and no coverage gate is enforced. Until that changes, every contribution should pass `npm run lint` and `npm run build` and include manual verification notes for affected flows such as login, products, orders, or revenue dashboards. If you add tests, colocate them with the feature and use `*.test.ts` or `*.test.tsx`.

## Commit & Pull Request Guidelines
Recent commits follow short conventional-style subjects such as `refactor: update product pricing...`. Continue using `<type>: <imperative summary>` where `type` is typically `feat`, `fix`, `refactor`, or `chore`. PRs should include a concise description, linked issue or task, screenshots for UI changes, manual test steps, and notes about any backend API contract changes.

## Security & Configuration Tips
Do not commit secrets, real tokens, or environment-specific URLs. The app stores auth tokens in `localStorage` and expects backend endpoints under `/api/v1`; keep auth and request changes centralized in `src/lib/api.ts` and related service modules.
