## Why

The backend currently exposes checklist data from hard-coded in-memory defaults while startup already expects a MySQL connection. A durable schema is needed so raceday checklists, ordered display sections, ordered items, and completion state can be stored consistently and used by future API endpoints.

## What Changes

- Add a MySQL database schema for checklist templates, checklist sections, section items, and per-event item completion state.
- Add migration artifacts that can initialize the backend database from an empty MySQL schema.
- Add a backend database package contract for connecting to MySQL with the existing `config.DatabaseConfig`.
- Establish tests or validation that prove the schema can be applied and the connection configuration is built correctly.

## Capabilities

### New Capabilities

- `mysql-checklist-storage`: Defines persistent MySQL storage for raceday checklist data and backend database connectivity.

### Modified Capabilities

- None.

## Impact

- Affected code: `api/internal/database`, backend startup in `api/cmd/server`, and any migration or schema files added for MySQL.
- Affected APIs: no external HTTP contract changes are required for this change; existing checklist responses may continue to use defaults until repository/query integration is added separately.
- Dependencies: the Go API will need a MySQL driver if one is not already present.
- Systems: local and deployed environments must provide MySQL credentials through the existing `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, and `DB_NAME` configuration.
