# AGENTS.md

Project instructions live in [CLAUDE.md](./CLAUDE.md).

All agents and assistants working in this repository must read and follow `CLAUDE.md` first. Do not duplicate long-lived rules here; keep this file as a stable pointer so the project has a single source of truth.

## Quick context

- **Stack**: Go (standard-library HTTP) backend + Vue 3 / TypeScript frontend.
- **Local start**: `docker compose -f docker-compose.yaml -f docker-compose.local.yaml up -d --build`
- **Remote-DB start**: `docker compose up -d --build` (configure `.env` first).
- **Backend tests**: `cd server && go test ./...`
- **Frontend build check**: `cd web && npm run build`
- **Default login**: `admin` / `Admin@123456`
