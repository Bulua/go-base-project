# GoBaseProject

GoBaseProject is a Go + Vue3 admin base project.

## Quick Start

1. Copy environment variables:

```bash
cp .env.example .env
```

On Windows PowerShell:

```powershell
Copy-Item .env.example .env
```

2. Edit `.env`, especially database values.

3. Start the project:

```bash
docker compose up -d --build
```

Open:

```text
Frontend: http://localhost:5173
Backend health: http://localhost:8080/api/v1/health
```

## Database

By default `.env.example` points `DB_HOST` to `120.53.251.75`, because the shared database is expected to live on that server.
The default MySQL port is `13307`; the default Redis port is `16379`.

For a fully local Docker database:

```bash
docker compose --profile local-db up -d --build
```

Then set:

```text
DB_HOST=mysql
DB_PORT=3306
DB_NAME=go_base_project
DB_USER=root
DB_PASSWORD=<MYSQL_ROOT_PASSWORD>
```

Local MySQL initialization scripts are kept in `deploy/mysql/init/`.

For a local Redis container as well:

```bash
docker compose --profile local-cache up -d --build
```

Set:

```text
REDIS_HOST=redis
REDIS_PORT_IN_CONTAINER=6379
```

## Repository Layout

```text
.
|-- docker-compose.yaml
|-- deploy
|   |-- mysql/init
|   `-- nginx
|-- docs
|-- server
`-- web
```

`development-plan/` is local planning material and is intentionally ignored by git.
