<div align="center">

# GoBaseProject

A production-ready admin panel starter built with Go and Vue 3.

[![Go](https://img.shields.io/badge/Go-1.22-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![Vue](https://img.shields.io/badge/Vue-3-4FC08D?logo=vue.js&logoColor=white)](https://vuejs.org/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5-3178C6?logo=typescript&logoColor=white)](https://www.typescriptlang.org/)
[![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?logo=docker&logoColor=white)](https://docs.docker.com/compose/)
[![License](https://img.shields.io/badge/License-MIT-green)](./LICENSE)

Clone → configure `.env` → `docker compose up`. That's it.

</div>

---

## Features

- **Authentication** — JWT access + refresh tokens, blocklist on logout, forced password change
- **RBAC** — roles, menus, API-level permission policies; super-admin auto-granted on startup
- **Dynamic routing** — menus drive the Vue Router; keep-alive and affix tabs per menu config
- **Audit logs** — every login attempt and mutating API call is recorded with IP and user-agent
- **File management** — upload, preview, soft-delete; files served through an authenticated endpoint
- **Dictionary & system params** — key-value stores for drop-downs and runtime configuration
- **API management** — auto-sync of all registered routes to the database on each startup
- **Code generation** — creating a menu with a component path auto-scaffolds the Vue file
- **Dark mode** — system-aware toggle with persistent preference
- **Multi-tab navigation** — browser-like tab bar with keep-alive, affix tabs, and right-click menu

---

## Tech Stack

| Layer | Technology |
|---|---|
| Backend | Go 1.22 · standard `net/http` · JWT · bcrypt |
| Frontend | Vue 3 · TypeScript · Vite · Element Plus · Pinia · Vue Router |
| Database | MySQL 8.0 |
| Cache | Redis 7 |
| Proxy | Nginx (inside the web container) |
| Runtime | Docker Compose |

---

## Getting Started

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/) with Compose v2 (`docker compose version`)
- No local Go or Node.js installation required

### Option A — Full local (all services in containers)

```bash
# 1. Copy the environment template
cp .env.example .env          # macOS / Linux
Copy-Item .env.example .env   # Windows PowerShell

# 2. Set your secrets in .env (see Configuration below)

# 3. Start everything
docker compose -f docker-compose.yaml -f docker-compose.local.yaml up -d --build
```

### Option B — Remote database

Use when MySQL and Redis are already running on a remote server.

```bash
cp .env.example .env
# Edit .env: set DB_HOST, DB_PORT, DB_USER, DB_PASSWORD,
#            REDIS_HOST, REDIS_PORT_IN_CONTAINER, REDIS_PASSWORD
docker compose up -d --build
```

When the server runs inside Docker, `localhost` means the server container itself.
Use service names such as `mysql` and `redis` for Compose-managed containers.
On Docker Desktop, use `host.docker.internal` when a container must connect to a database or Redis running directly on the host machine.

### Access

| Service | URL |
|---|---|
| Frontend | http://localhost:5173 |
| Backend health | http://localhost:8080/api/v1/health |

**Default credentials:** `admin` / `Admin@123456`

---

## Configuration

Copy `.env.example` to `.env` and edit the values below.

```dotenv
# ── Application ────────────────────────────────────────────────────────────────
APP_NAME=GoBaseProject
APP_ENV=development
JWT_SECRET=replace-with-a-long-random-string   # required

# ── Database ───────────────────────────────────────────────────────────────────
DB_HOST=localhost       # direct host run; docker-compose.local.yaml overrides this to "mysql"
DB_PORT=3306
DB_NAME=go_base_project
DB_USER=root
DB_PASSWORD=change-me

# ── Redis ──────────────────────────────────────────────────────────────────────
REDIS_HOST=localhost    # direct host run; docker-compose.local.yaml overrides this to "redis"
REDIS_PORT_IN_CONTAINER=6379
REDIS_PASSWORD=change-me

# ── Local containers (only used with docker-compose.local.yaml) ────────────────
MYSQL_ROOT_PASSWORD=change-me   # must match DB_PASSWORD
```

> `docker-compose.local.yaml` overrides `DB_HOST` and `REDIS_HOST` automatically — you only need to keep passwords consistent.
> Do not use `localhost` inside a container to reach host services. Use `host.docker.internal` on Docker Desktop, or a real network address.

---

## Project Structure

```text
.
├── docker-compose.yaml          # base Compose file (server + web)
├── docker-compose.local.yaml    # override: adds local MySQL + Redis
├── .env.example                 # environment template
│
├── server/                      # Go backend
│   ├── cmd/gobaseproject/       # main entry point
│   ├── internal/
│   │   ├── handler/             # HTTP handlers (auth, user, role, menu, …)
│   │   ├── service/             # business logic
│   │   ├── repository/          # database access
│   │   ├── model/               # domain types, DTOs, errors
│   │   ├── middleware/          # permission, operation-log, CORS
│   │   └── infra/config/        # config loader (YAML + env)
│   └── pkg/                     # reusable packages (response, routereg, …)
│
├── web/                         # Vue 3 frontend
│   └── src/
│       ├── api/                 # Axios request modules
│       ├── views/               # page components
│       ├── layouts/             # app shell (sidebar, navbar, tabs)
│       ├── store/modules/       # Pinia stores (auth, tabs)
│       ├── router/              # static + dynamic route registration
│       ├── components/          # shared UI components
│       └── types/               # TypeScript interfaces
│
└── deploy/
    ├── mysql/init/init.sql      # database schema + seed data
    └── nginx/                   # Nginx config for the web container
```

---

## Database Initialization

`deploy/mysql/init/init.sql` is mounted into the MySQL container and executed automatically on the first start.

To reset the database and start fresh:

```bash
docker compose down -v   # removes volumes
docker compose -f docker-compose.yaml -f docker-compose.local.yaml up -d --build
```

---

## Development

### Backend

```bash
cd server
go test ./...   # run tests
gofmt -w .      # format before committing
```

### Frontend

```bash
cd web
npm install
npm run dev     # dev server with HMR
npm run build   # production build (type-checks included)
```

---

## License

[MIT](./LICENSE)
