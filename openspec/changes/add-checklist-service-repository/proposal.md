## Why

The `/api/checklist` route still owns hard-coded checklist data even though backend startup now establishes a MySQL connection. Moving checklist reads behind a service and repository boundary will make the router thin, make persistence integration testable, and prepare the endpoint to serve database-backed checklist state without changing the HTTP contract prematurely.

## What Changes

- Add a checklist service/repository capability for retrieving checklist data through an explicit backend boundary.
- Refactor the Gin router so `/api/checklist` delegates to an injected checklist service instead of calling package-local default data directly.
- Add a MySQL-backed repository that reads checklist templates, sections, items, raceday events, and completion state from the existing schema.
- Preserve the current flat `/api/checklist` response shape for this change unless a missing database-backed checklist requires an error response.
- Keep hard-coded fallback/default checklist data out of the router so tests can exercise handler behavior with injected dependencies.

## Capabilities

### New Capabilities

- `checklist-service-repository`: Defines checklist retrieval through backend service and repository boundaries, including MySQL-backed reads for the existing `/api/checklist` endpoint.

### Modified Capabilities

- `mysql-checklist-storage`: Adds requirements for reading ordered checklist data and event-specific completion state from the persistent MySQL schema.

## Impact

- Affected code: `api/internal/server`, new or updated checklist service/repository packages, and `api/cmd/server` dependency wiring.
- Affected APIs: `/api/checklist` should keep its existing response shape for clients.
- Dependencies: no new external dependency is expected beyond the existing Go MySQL driver.
- Tests: router tests should use injected checklist behavior, and repository/service tests should validate ordering, completion merging, and error handling without requiring a live MySQL server where feasible.
