## 1. Downgrade TypeScript

- [x] 1.1 In `frontend/web`, run `npm install --save-dev typescript@~5.8.0 --legacy-peer-deps` to update `package.json` and `package-lock.json`
- [x] 1.2 Verify `npm install` (no flags) succeeds cleanly

## 2. Update Dockerfile

- [x] 2.1 In `frontend/web/Dockerfile`, change `RUN npm ci --legacy-peer-deps` to `RUN npm ci`

## 3. Verify

- [x] 3.1 Run `npm run build` in `frontend/web` — confirm it compiles with TypeScript 5.8 without errors
- [x] 3.2 Run `npm run dev` in `frontend/web` — confirm the dev server starts
