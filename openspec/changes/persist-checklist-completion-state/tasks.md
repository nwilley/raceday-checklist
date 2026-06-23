## 1. Database Migration

- [x] 1.1 Add an ordered `up` migration for `raceday_participants` with event-scoped anonymous `client_id` uniqueness.
- [x] 1.2 Update completion storage to support participant-scoped item completion uniqueness.
- [x] 1.3 Add a matching `down` migration that rolls back participant-scoped completion storage in dependency-safe order.
- [x] 1.4 Add migration validation tests for participant table, participant foreign keys, uniqueness constraints, and rollback order.

## 2. Checklist Repository and Service

- [x] 2.1 Extend checklist item models to include `sectionId` and `itemId` while preserving `id`, `title`, `category`, and `done`.
- [x] 2.2 Add repository logic to resolve or create an anonymous participant for the default raceday event and client ID.
- [x] 2.3 Update default checklist reads to return completion state scoped to the resolved anonymous participant.
- [x] 2.4 Add repository logic to upsert completion state by default event, anonymous participant, section slug, and item slug.
- [x] 2.5 Add service methods for reading participant-scoped checklist state and updating item completion.
- [x] 2.6 Add repository and service tests for participant creation/reuse, repeated item slugs, independent participant state, save, clear, and not-found item behavior.

## 3. Backend HTTP API

- [x] 3.1 Add request handling for the `X-Raceday-Client` anonymous client ID header.
- [x] 3.2 Update `GET /api/checklist` to use the anonymous client ID and return stable `sectionId` and `itemId` fields.
- [x] 3.3 Add `PATCH /api/checklist/items/:sectionId/:itemId` with body `{ "done": boolean }`.
- [x] 3.4 Map unknown section or item keys to a client-visible not-found response.
- [x] 3.5 Map validation and unexpected persistence errors to appropriate non-leaking HTTP responses.
- [x] 3.6 Add router tests for missing client ID, successful reads, successful updates, unknown item keys, and unexpected errors.

## 4. Frontend Optimistic Persistence

- [x] 4.1 Add frontend client ID creation and reuse through local storage.
- [x] 4.2 Include the anonymous client ID header on checklist read and write requests.
- [x] 4.3 Update frontend item types to use `sectionId` and `itemId` as stable completion keys.
- [x] 4.4 Persist checkbox changes through the new PATCH endpoint after optimistic UI updates.
- [x] 4.5 Track per-item pending save state while completion updates are in flight.
- [x] 4.6 Handle failed saves by clearly surfacing the failure and avoiding silent divergence from backend state.

## 5. Verification

- [x] 5.1 Run `go test ./...` from the API module and fix any compile or test failures.
- [x] 5.2 Run the frontend build and fix any TypeScript or Vite failures.
- [x] 5.3 Confirm no login, account, password, or email flow is introduced.
- [x] 5.4 Confirm OpenSpec requirements for `anonymous-participant-state` are satisfied.
- [x] 5.5 Confirm added `checklist-service-repository` and `mysql-checklist-storage` requirements are satisfied.
