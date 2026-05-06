# Running GoBaseProject

This project is designed to start from the repository root with Docker Compose.

## 1. Prepare Environment

```bash
cp .env.example .env
```

Windows PowerShell:

```powershell
Copy-Item .env.example .env
```

Edit `.env` before starting services.

## 2. Full Local Docker Mode

Use this when MySQL, Redis, server, and web all run in Docker:

```bash
docker compose -f docker-compose.yaml -f docker-compose.local.yaml up -d --build
```

`docker-compose.local.yaml` overrides the server container to use Docker service names:

```text
DB_HOST=mysql
DB_PORT=3306
REDIS_HOST=redis
REDIS_PORT_IN_CONTAINER=6379
```

## 3. External Database Mode

Use this when MySQL and Redis already run outside this Compose project.

When the server runs inside Docker, `localhost` means the server container itself. Do not use `localhost` in `.env` to reach host services from the container.

Use one of these values instead:

```text
Docker Compose service: DB_HOST=mysql, REDIS_HOST=redis
Docker Desktop host service: DB_HOST=host.docker.internal, REDIS_HOST=host.docker.internal
Remote service: DB_HOST=<server-ip-or-domain>, REDIS_HOST=<server-ip-or-domain>
```

Then start app services:

```bash
docker compose up -d --build
```

## 4. Local Database Profile

To start MySQL locally with Docker:

```bash
docker compose --profile local-db up -d --build
```

Set these values in `.env`:

```text
DB_HOST=mysql
DB_PORT=3306
DB_NAME=go_base_project
DB_USER=root
DB_PASSWORD=<same as MYSQL_ROOT_PASSWORD>
```

The local database is initialized from:

```text
deploy/mysql/init
```

## 5. Local Redis Profile

To start Redis locally with Docker:

```bash
docker compose --profile local-cache up -d --build
```

Set these values in `.env`:

```text
REDIS_HOST=redis
REDIS_PORT_IN_CONTAINER=6379
```

## 6. Endpoints

```text
Frontend: http://localhost:5173
Backend health: http://localhost:8080/api/v1/health
```

## 7. Useful Commands

```bash
docker compose ps
docker compose logs -f server
docker compose logs -f web
docker compose down
```
