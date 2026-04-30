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

## 2. Remote Database Mode

The default `.env.example` uses:

```text
DB_HOST=120.53.251.75
DB_PORT=13307
DB_NAME=go_base_project
REDIS_HOST=120.53.251.75
REDIS_PORT_IN_CONTAINER=16379
```

Start app services:

```bash
docker compose up -d --build
```

## 3. Local Database Mode

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

## 4. Local Redis Mode

To start Redis locally with Docker:

```bash
docker compose --profile local-cache up -d --build
```

Set these values in `.env`:

```text
REDIS_HOST=redis
REDIS_PORT_IN_CONTAINER=6379
```

## 5. Endpoints

```text
Frontend: http://localhost:5173
Backend health: http://localhost:8080/api/v1/health
```

## 6. Useful Commands

```bash
docker compose ps
docker compose logs -f server
docker compose logs -f web
docker compose down
```
