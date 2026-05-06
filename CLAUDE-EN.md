# CLAUDE-EN.md

This file is the English companion to [CLAUDE.md](./CLAUDE.md). The Chinese file is the canonical project rule set. If the two versions disagree, follow `CLAUDE.md`.

## 1. Project Purpose

GoBaseProject is a Go + Vue 3 admin panel starter. Its goal is to provide a reusable management-system base that can be started with Docker Compose and extended safely.

The project currently includes:

1. A Go backend service
2. A Vue 3 / TypeScript frontend
3. MySQL initialization scripts
4. Redis cache configuration
5. Docker Compose modes for remote and local dependencies

`development-plan/` is local planning and design reference only. It is ignored by Git. Production source code, scripts, Docker configuration, and committed documentation must not depend on it.

## 2. Common Commands

Copy the environment template from the repository root:

```bash
cp .env.example .env
```

Windows PowerShell:

```powershell
Copy-Item .env.example .env
```

Start with remote MySQL / Redis:

```bash
docker compose up -d --build
```

Start with local MySQL / Redis:

```bash
docker compose -f docker-compose.yaml -f docker-compose.local.yaml up -d --build
```

Backend tests:

```bash
cd server
go test ./...
```

Frontend build check:

```bash
cd web
npm install
npm run build
```

Compose validation:

```bash
docker compose config
```

If Docker, Go, or Node.js is unavailable in the current environment, state which checks were not run in the final response.

## 3. Repository Layout

```text
.
|-- AGENTS.md                 # Agent entry point
|-- CLAUDE.md                 # Canonical Chinese rules
|-- CLAUDE-EN.md              # English rules
|-- README.md                 # Quick user-facing guide
|-- docker-compose.yaml       # Default Compose configuration
|-- docker-compose.local.yaml # Local MySQL / Redis override
|-- deploy                    # Deployment, initialization, Nginx resources
|-- docs                      # Committed documentation
|-- server                    # Go backend
|-- web                       # Vue frontend
`-- development-plan          # ignored, local planning and design reference only
```

Committed runtime documentation belongs in `README.md` and `docs/`. Deployment scripts, initialization scripts, and Nginx configuration belong in `deploy/`.

## 4. Secret Handling

1. `.env` is local only and must not be committed.
2. `.env.example` may contain placeholders only, such as `change-me`.
3. Do not write real passwords, tokens, or keys into code, tests, logs, Markdown, SQL, Dockerfiles, or committed config templates.
4. Demo credentials may be documented, but production accounts and real secrets must not be documented.
5. Use `rg` before and after changes when checking for accidental leaks:

```bash
rg "<REAL_SECRET_TO_CHECK>" -g "!.env" -g "!development-plan/**" -g "!web/node_modules/**" -g "!web/dist/**" .
```

No output means that keyword was not found.

## 5. Backend Rules

The backend lives in `server/`.

The current backend uses Go 1.22, standard-library `net/http`, MySQL, JWT, bcrypt, and YAML config loading. Add dependencies only when they clearly reduce complexity and match the existing style.

Recommended structure:

```text
server
|-- cmd
|   `-- gobaseproject          # main entry point
|-- configs                    # config templates
|-- internal
|   |-- handler                # HTTP handlers
|   |-- service                # business logic
|   |-- repository             # data access
|   |-- model                  # domain types, DTOs, errors
|   |-- middleware             # permission, audit, CORS
|   `-- infra                  # config, database, cache, logging
|-- pkg                        # reusable packages
`-- tests                      # backend tests
```

Layer responsibilities:

```text
handler: parse requests, validate input, call services, write responses
service: business rules, permissions, transactions, repository coordination
repository: database access, query building, persistence details
model: domain models, DTOs, enums, errors
middleware: request-chain auth, audit, context handling
infra: config, database, Redis, logging
pkg: reusable code weakly coupled to business modules
```

Backend rules:

1. Place new business code through `handler -> service -> repository -> model` by default.
2. Cross-module calls should go through services, not directly into another module's repository.
3. Keep API responses consistent.
4. Paginate list endpoints by default and cap `page_size`.
5. Store password hashes only. Prefer bcrypt or argon2id.
6. Store only token digests in JWT blocklists, never raw tokens.
7. Record audit logs for failed login, logout, password reset, permission changes, and write operations.
8. Mask passwords, tokens, keys, and similar fields in audit logs.
9. Run `gofmt` before committing. Run `go test ./...` when backend logic changes.

Standard response example:

```json
{
  "code": 0,
  "message": "ok",
  "data": {},
  "trace_id": "request-trace-id"
}
```

## 6. Frontend Rules

The frontend lives in `web/`.

The current frontend uses Vue 3, TypeScript, Vite, Element Plus, Pinia, Vue Router, and Axios. New pages and components should match the existing admin style: moderate information density, clear operations, and no marketing-page treatment.

Recommended structure:

```text
web/src
|-- api
|   |-- request                # Axios instance, interceptors, errors
|   |-- auth
|   |-- user
|   `-- ...
|-- assets
|   `-- styles                 # global styles, theme variables
|-- components
|   |-- common                 # shared components
|   |-- business               # business components
|   `-- layout                 # layout components
|-- composables                # useXxx composables
|-- directives                 # v-auth and similar directives
|-- layouts                    # page shells
|-- permissions                # permission checks
|-- router                     # static and dynamic routes
|-- store
|   `-- modules                # Pinia modules
|-- types                      # DTOs, enums, type declarations
|-- utils                      # pure utility functions
`-- views                      # page-level components
```

Frontend rules:

1. Page component entry files should be named `index.vue`.
2. Business components use PascalCase, for example `RoleSelect.vue`.
3. Composables use `useXxx.ts`.
4. API files are named by resource, for example `api/user/index.ts`.
5. All API requests must use the shared Axios instance.
6. Frontend permissions control display only and never replace backend authorization.
7. Prefer Element Plus and existing project components before adding new UI primitives.
8. When types change, update `types/`, `api/`, `store/`, and all page usages together.
9. For frontend logic changes, run at least `npm run build`; run `npm run test` when tests are relevant.

## 7. Database And Initialization

Database initialization scripts live in:

```text
deploy/mysql/init
```

The current committed initialization file is:

```text
deploy/mysql/init/init.sql
```

Database name:

```text
go_base_project
```

Rules:

1. Do not use `development-plan/sql` as a runtime script source.
2. Initialization scripts should be reviewable. Explain the impact of schema or seed-data changes.
3. Future production database changes should use a migration mechanism and document the process in README or `docs/`.
4. After SQL changes, validate startup and initialization with a local container or test database when possible.

## 8. Docker Rules

`docker-compose.yaml` is the default entry point and connects to the external MySQL / Redis configured in `.env`.

`docker-compose.local.yaml` is the local development override for running MySQL / Redis locally.

Rules:

1. Run `docker compose config` after Compose changes.
2. When ports, service names, or environment variables change, update `.env.example`, README, and `docs/RUNNING.md`.
3. Do not write real passwords into Compose files.
4. Do not commit local volumes, build artifacts, or temporary logs.

## 9. Documentation Rules

1. `README.md` is user-facing and should contain quick start, configuration, endpoints, and project overview.
2. `docs/RUNNING.md` contains detailed runtime instructions.
3. `CLAUDE.md` contains the long-lived Chinese development rules.
4. `CLAUDE-EN.md` is the English translation for English-speaking agents or collaborators.
5. `AGENTS.md` stays lightweight and only contains the agent entry point and quick context.
6. Use ASCII tree characters in Markdown when possible to reduce Windows terminal encoding issues.

## 10. Files And Workflow

1. New files should be UTF-8.
2. Prefer ASCII punctuation and path symbols in scripts, config, and Markdown.
3. Do not commit build artifacts: `web/dist/`, `web/node_modules/`, `server/bin/`.
4. Do not delete `.gitkeep` unless the directory already contains real files.
5. When changing directory structure, update README, docs, and Docker configuration together.
6. The working tree may contain other people's changes. Check `git status` before editing and do not revert changes you did not make.
7. In the final response, state which files changed, which checks ran, and which checks were not run.

## 11. Suggested Development Priority

Recommended order:

1. Config loading, logging, standard responses
2. Database connection and initialization / migration strategy
3. Redis connection
4. Authentication and authorization
5. Users, roles, menus, permissions
6. Audit logs
7. Files, dictionaries, system parameters
8. Vue login page, base layout, dynamic routes, button permissions
9. CRUD pages and UX polish

Every completed feature area should include a minimal verification path.
