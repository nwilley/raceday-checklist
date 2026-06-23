## 1. Checklist Domain and Service

- [x] 1.1 Create a checklist package with the flat item model used by `/api/checklist`.
- [x] 1.2 Define a narrow repository interface for reading the default event checklist.
- [x] 1.3 Implement a checklist service that delegates to the repository and returns a distinct missing checklist error.
- [x] 1.4 Add service unit tests for successful reads, missing checklist behavior, and repository error propagation.

## 2. MySQL Repository

- [x] 2.1 Implement a MySQL checklist repository that selects a deterministic default raceday event.
- [x] 2.2 Query ordered sections and items using section display order and item display order.
- [x] 2.3 Join event completion rows so completed items are returned as done and missing completion rows default to not done.
- [x] 2.4 Map section titles to the flat response category field and item slugs/titles to response item fields.
- [x] 2.5 Add repository tests that validate default event selection, ordering, completion merging, no-event behavior, and database errors without requiring a live MySQL server.

## 3. Router Integration

- [x] 3.1 Update `server.NewRouter` to accept an injected checklist service dependency.
- [x] 3.2 Refactor `/api/checklist` to call the service and return the existing top-level `items` response on success.
- [x] 3.3 Map missing checklist errors to a client-visible unavailable checklist response.
- [x] 3.4 Map unexpected service errors to a generic server error response without leaking database details.
- [x] 3.5 Update router tests to use a fake checklist service for success, missing checklist, and unexpected error cases.

## 4. Application Wiring and Verification

- [x] 4.1 Wire `api/cmd/server` to construct the MySQL repository and checklist service from the existing `*sql.DB`.
- [x] 4.2 Remove hard-coded checklist fallback data from `api/internal/server/router.go`.
- [x] 4.3 Confirm `GET /api/checklist` preserves `items[].id`, `items[].title`, `items[].category`, and `items[].done`.
- [x] 4.4 Run `go test ./...` from the API module and fix any compile or test failures.
- [x] 4.5 Confirm the OpenSpec requirements for `checklist-service-repository` and the added `mysql-checklist-storage` read behavior are satisfied.
