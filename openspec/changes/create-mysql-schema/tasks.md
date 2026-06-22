## 1. Database Schema

- [x] 1.1 Create an API migrations directory for ordered MySQL SQL files.
- [x] 1.2 Add the `up` migration that creates checklist templates, ordered sections, ordered section items, raceday events, and event item completion tables.
- [x] 1.3 Add indexes, unique constraints, and foreign keys for template-scoped sections, section-scoped items, display ordering, and referential integrity.
- [x] 1.4 Add the matching `down` migration that drops tables in dependency-safe reverse order.

## 2. Backend Connectivity

- [x] 2.1 Add the Go MySQL driver dependency to `api/go.mod`.
- [x] 2.2 Create `api/internal/database` with `Connect(ctx, config.DatabaseConfig) (*sql.DB, error)`.
- [x] 2.3 Build the MySQL DSN from the existing database config and include sensible connection options for time parsing and UTF-8 support.
- [x] 2.4 Configure basic connection pool settings and verify connectivity with `PingContext`.

## 3. Validation

- [x] 3.1 Add unit tests for DSN construction or connection configuration behavior without requiring a live MySQL server.
- [x] 3.2 Add migration validation that checks the expected tables, constraints, and rollback order are present in the SQL files.
- [x] 3.3 Run `go test ./...` from the API module and fix any compile or test failures.
- [x] 3.4 Confirm the OpenSpec requirements for `mysql-checklist-storage` are satisfied by the implementation.
