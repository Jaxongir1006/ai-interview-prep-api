# Migration

Uses **goose** package.

## Commands

| Command               | Description                  |
| --------------------- | ---------------------------- |
| `make migrate-create` | Create new migration         |
| `make migrate-up`     | Apply all pending migrations |
| `make migrate-down`   | Rollback last migration      |

## File Naming

Prefix with module name, use snake_case: `auth_init_schema`, `auth_add_user_roles`, `platform_init_taskmill`

## Migration Rules

- **Single folder** — all migrations in `./migrations`, no subfolders
- **Auto-execution** — runs on application startup (including production)
- **Queries order** — 1. CREATE TABLE → 2. CREATE INDEX → 3. ALTER TABLE (foreign keys)
- **Rollback** — one at a time, reverse order, use `IF EXISTS`

## Column Types

- Timestamps: always use `timestamptz`

## Applying Migrations

Migrations are applied on application startup automatically.
