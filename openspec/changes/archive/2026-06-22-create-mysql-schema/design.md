## Context

The Go API is structured around Gin and already loads MySQL connection settings through `config.DatabaseConfig`. `api/cmd/server/main.go` imports `api/internal/database` and calls `database.Connect(ctx, cfg.Database)`, but the database package and schema are not present in the current tree. The existing checklist endpoint returns a small hard-coded set of default items, so persistent storage can be added without changing the public HTTP response shape in this change.

This change establishes the durable MySQL foundation for raceday checklist data. It should create the database connection package and migration files needed to initialize an empty database, while leaving higher-level API repository integration for a later change.

## Goals / Non-Goals

**Goals:**

- Define a normalized MySQL schema for raceday checklist templates, ordered sections, ordered items, events, and per-event item completion.
- Provide ordered migration files that can create and drop the schema.
- Implement `database.Connect(ctx, config.DatabaseConfig)` so backend startup can validate and open a MySQL connection using existing configuration.
- Add focused tests or validation for connection string construction and migration SQL shape.

**Non-Goals:**

- Replace the `/api/checklist` hard-coded response with live database reads.
- Add user authentication, user ownership, or multi-tenant authorization.
- Add a migration runner service or production deployment automation.
- Add seed data beyond what is needed to validate schema behavior.

## Decisions

1. Store reusable checklist structure separately from raceday event completion state.

   The schema will use checklist templates, ordered sections, and ordered section items for reusable checklist definitions, plus events and event item state for per-raceday completion. This avoids overwriting a reusable checklist item when a user marks it done for a specific event.

   Alternative considered: a single `checklist_items` table with a `done` column. That is simpler, but it mixes reusable checklist definitions with event-specific state and would make repeated raceday usage lossy.

2. Use integer primary keys and stable string slugs where API-facing identifiers are useful.

   Internal joins will use unsigned integer IDs for compact indexes and clear foreign keys. Template, section, and item rows will include stable slugs so the backend can preserve simple identifiers like `fuel`, `tires`, and `helmet` without relying on display titles.

   Alternative considered: UUID-only primary keys. UUIDs are useful for distributed creation, but this backend is starting with a single MySQL database and does not need that complexity yet.

3. Model sections as ordered display groups, not as item categories.

   Each checklist template can contain multiple sections, and each section has its own display order. Checklist items belong to exactly one section and have a display order within that section. This matches the UI model where a checklist is shown as several sections, each with its own ordered list of items.

   Alternative considered: storing a free-form `section` string directly on each item. That would reduce table count, but it would make section ordering, section renaming, and section-level uniqueness harder to enforce.

4. Use explicit foreign keys with restrictive deletes for template structure and cascading deletes for event state.

   Sections and items should not disappear accidentally while events reference them. Event-specific rows can cascade when an event is deleted because those rows have no purpose without the event.

   Alternative considered: no foreign keys with application-enforced integrity. That reduces database constraints, but it makes data corruption easier and undermines the value of defining a schema now.

5. Keep migrations as plain SQL files in the API tree.

   Plain SQL keeps the schema reviewable and can be used by common migration tools later. The implementation should place ordered `up` and `down` files under a clear migrations directory such as `api/migrations`.

   Alternative considered: embedding schema creation directly in Go startup. That is convenient locally, but it couples application boot to schema ownership and makes rollback harder.

6. Use the standard `database/sql` package with the Go MySQL driver.

   `database.Connect` should return a `*sql.DB`, build a DSN from `config.DatabaseConfig`, apply sane pool settings, and ping using the provided context. This satisfies the current startup call without introducing a repository abstraction prematurely.

   Alternative considered: adopting an ORM. The immediate need is schema and connectivity, and an ORM would add unnecessary conventions before persistence workflows are defined.

## Risks / Trade-offs

- Schema guesses future API needs incorrectly -> Keep the schema focused on checklist templates and event completion, and avoid user/account fields until authentication requirements exist.
- Migration files are added without a runner -> Document or test the SQL shape now, and leave runner selection to a later deployment-focused change.
- Startup requires a reachable database -> Preserve current config validation, ping with context, and keep tests isolated from a live MySQL server unless explicitly configured.
- Slugs and ordering constraints are too strict -> Use uniqueness scoped to the parent entity, so different templates and sections can reuse natural item slugs and display order values.

## Migration Plan

1. Add the database package and SQL migration files.
2. Apply the `up` migration to local and deployed MySQL databases before starting the API with database-backed features.
3. Roll back by applying the matching `down` migration in environments where no production data needs to be preserved, or by taking a backup and manually migrating data before destructive rollback.

## Open Questions

- Which migration runner should be standardized for local development and deployment?
- Should starter checklist templates be seeded by migration, by a separate seed command, or by application logic in a future change?
