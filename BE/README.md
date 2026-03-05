# Stock Management API

Simple Go backend for stock management with PostgreSQL.

## Requirements

- Go `1.23+`
- PostgreSQL `16+` (or use Docker)
- `make`
- [`migrate` CLI](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

## Installation

1. Clone this repository.
2. Install dependencies:

```bash
go mod tidy
```

3. Create environment file:

```bash
cp .env.example .env
```

4. Start PostgreSQL (example using Docker):

```bash
docker compose up -d postgres
```

## Run Migration

Load env first, then run migration:

```bash
set -a; source .env; set +a
make migrate-up
```

Optional:

```bash
make migrate-down
make migrate-status
```

## Run Seeder

```bash
set -a; source .env; set +a
make seed
```

## Run Project

```bash
make run
```

API will run on `http://localhost:8080` (default from `.env`).

## Quick Run with Docker Compose

Run API + PostgreSQL + migration:

```bash
docker compose up --build
```

Run seeder profile:

```bash
docker compose --profile seed up --build seed
```
