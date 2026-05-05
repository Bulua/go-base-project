# GoBaseProject

Go + Vue3 admin base project. Runs entirely in Docker; no local Go or Node toolchain required.

Default login: `admin` / `Admin@123456`

---

## Option A — Full local (recommended for development)

All services (server, web, MySQL, Redis) run as containers on your machine.

**Step 1 — Copy env file**

```bash
cp .env.example .env
```

Windows PowerShell:

```powershell
Copy-Item .env.example .env
```

**Step 2 — Set passwords** (edit `.env`)

Change at minimum:

```text
JWT_SECRET=any-random-string
MYSQL_ROOT_PASSWORD=any-local-password
DB_PASSWORD=any-local-password   # must match MYSQL_ROOT_PASSWORD
REDIS_PASSWORD=any-local-password
```

> `change-me` values work for a quick smoke test but are not safe even locally.

**Step 3 — Start**

```bash
docker compose -f docker-compose.yaml -f docker-compose.local.yaml up -d --build
```

Open:

```text
Frontend : http://localhost:5173
API health: http://localhost:8080/api/v1/health
```

---

## Option B — Remote database (default)

Use when you already have MySQL and Redis running on a remote server.

**Step 1 — Copy env file and edit**

```bash
cp .env.example .env
```

Set the real values in `.env`:

```text
DB_HOST=<your-db-host>
DB_PORT=<your-db-port>
DB_USER=<db-user>
DB_PASSWORD=<db-password>

REDIS_HOST=<your-redis-host>
REDIS_PORT_IN_CONTAINER=<redis-port>
REDIS_PASSWORD=<redis-password>

JWT_SECRET=<random-secret>
```

**Step 2 — Start**

```bash
docker compose up -d --build
```

Only `server` and `web` containers start; no local MySQL or Redis is created.

---

## Stopping / cleanup

```bash
# Stop containers
docker compose down

# Stop and remove volumes (deletes local DB data)
docker compose down -v
```

---

## Repository layout

```text
.
|-- docker-compose.yaml         main compose file
|-- docker-compose.local.yaml   full-local override (MySQL + Redis as containers)
|-- .env.example                environment template
|-- deploy/
|   |-- mysql/init/             SQL init scripts (auto-run on first MySQL start)
|   `-- nginx/                  Nginx config for the web container
|-- docs/                       additional documentation
|-- server/                     Go backend
`-- web/                        Vue3 frontend
```

## Database initialization

`deploy/mysql/init/init.sql` is mounted into the MySQL container and runs automatically on the first start. If you need to re-initialize, remove the volume:

```bash
docker compose down -v
docker compose -f docker-compose.yaml -f docker-compose.local.yaml up -d --build
```
