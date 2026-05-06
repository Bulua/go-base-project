# AGENTS-EN.md

This file is the English companion to [AGENTS.md](./AGENTS.md). The Chinese file is canonical. If the two versions disagree, follow `AGENTS.md`.

Long-lived project rules live in [CLAUDE.md](./CLAUDE.md). Read and follow it before changing code, scripts, configuration, or documentation.

## Quick Context

- **Project type**: Go + Vue 3 admin panel starter.
- **Backend**: Go 1.22, standard-library `net/http`, MySQL, JWT, bcrypt.
- **Frontend**: Vue 3, TypeScript, Vite, Element Plus, Pinia, Vue Router, Axios.
- **Default remote-dependency start**: `docker compose up -d --build`
- **Full local start**: `docker compose -f docker-compose.yaml -f docker-compose.local.yaml up -d --build`
- **Backend tests**: `cd server && go test ./...`
- **Frontend build check**: `cd web && npm run build`
- **Frontend tests**: `cd web && npm run test`
- **Default demo login**: `admin` / `Admin@123456`

## Agent Rules

1. Read [CLAUDE.md](./CLAUDE.md) before changing any file.
2. Do not copy long-lived rules into this file; avoid drift from `CLAUDE.md`.
3. Do not commit or write real passwords, tokens, or keys.
4. Do not make production code, scripts, Docker configuration, or committed documentation depend on `development-plan/`.
5. The working tree may contain other people's changes. Check status before editing and do not revert changes you did not create.
