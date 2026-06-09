# Web

React + TypeScript web client for the design-youtube platform. Built with Vite, served via nginx in production.

## Development

```bash
npm install
npm run dev      # dev server at http://localhost:5173
npm run build    # production build → dist/
```

Run via Docker Compose alongside the full backend stack:

```bash
# From repo root
docker compose up --build
# Available at http://localhost:3000
```

## Tech Stack

- [React](https://react.dev/) + [TypeScript](https://www.typescriptlang.org/)
- [Vite](https://vite.dev/) — dev server and bundler
- [nginx](https://nginx.org/) — production static file server
