## Why

The backend can now read `/api/checklist` from MySQL, but an empty database has no default raceday event or checklist data to serve. A repeatable database seeder is needed so local development and fresh environments can create starter checklist data explicitly instead of relying on hard-coded router fallback behavior.

## What Changes

- Add a database seeding capability for creating the default checklist template, ordered sections, ordered items, and an initial raceday event.
- Make the seed operation idempotent so it can be run repeatedly without duplicating templates, sections, items, events, or completion rows.
- Provide an executable backend entry point for running the seeder against the configured MySQL database.
- Seed data that matches the requested race-day workflow sections: Pre-practice, Pre-qualifying, Mid-qualifying maintenance, and Pre-main.
- Keep seeding separate from API server startup and request handling.

## Capabilities

### New Capabilities

- `database-seeding`: Defines repeatable MySQL seed behavior for starter raceday checklist data and an executable backend seeder entry point.

### Modified Capabilities

- `mysql-checklist-storage`: Adds requirements that seeded data conform to the persistent checklist schema and support the default checklist read path.

## Impact

- Affected code: new API seed command or package, checklist seed data definitions, database connection reuse, and tests for generated SQL/seed behavior.
- Affected APIs: no external HTTP API contract changes.
- Dependencies: no new external runtime dependency is expected.
- Systems: local and deployed environments can run the seeder after migrations and before relying on database-backed checklist reads.
