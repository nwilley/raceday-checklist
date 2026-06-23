## Context

The backend currently starts by opening a MySQL connection, but `api/internal/server.NewRouter()` does not receive that connection and still serves `/api/checklist` from hard-coded package-local data. The recently added MySQL schema defines reusable checklist templates, ordered sections, ordered items, raceday events, and event-specific completion state, but no application layer reads from it yet.

This change should introduce the first persistence-backed checklist read path while keeping the Gin router focused on HTTP concerns. The route can continue returning the current flat `items` payload so frontend clients do not need to change in the same step.

## Goals / Non-Goals

**Goals:**

- Move checklist retrieval out of `api/internal/server/router.go` and behind an injected service boundary.
- Add a repository implementation that reads ordered checklist data and event completion state from MySQL.
- Wire `api/cmd/server` so the router receives a database-backed checklist service.
- Preserve the current `/api/checklist` JSON response shape.
- Keep router tests independent from MySQL by using a fake or stub checklist service.
- Add repository/service tests that validate SQL behavior without requiring a live MySQL server.

**Non-Goals:**

- Add new public checklist endpoints.
- Change the frontend API contract to expose nested sections.
- Add checklist write, event creation, template management, or completion mutation APIs.
- Add a migration runner or seed workflow.
- Add authentication, ownership, or multi-tenant filtering.

## Decisions

1. Introduce a checklist service interface at the router boundary.

   `server.NewRouter` should accept a dependency that can list checklist items for the existing endpoint. The router should translate service results into HTTP responses and status codes, but it should not know about SQL tables or default item construction.

   Alternative considered: pass `*sql.DB` directly into the router. That is faster to wire, but it pushes query construction and storage errors into the HTTP layer and makes router tests depend on database concerns.

2. Keep the service model close to the current response shape.

   The service can return a flat item model with `ID`, `Title`, `Category`, and `Done` fields. The repository can map `checklist_sections.title` to `Category`, `checklist_items.slug` to `ID`, `checklist_items.title` to `Title`, and event completion state to `Done`.

   Alternative considered: introduce nested section DTOs now. The schema supports that, but changing the API shape would broaden the change beyond the proposal and require frontend coordination.

3. Use a MySQL repository under the service, not as the service itself.

   The repository should own SQL queries and row scanning. The service should own use-case behavior, including choosing the active/default checklist event if the route does not yet accept an event identifier. This leaves room for future endpoints to add event selection without rewriting the router.

   Alternative considered: a single `ChecklistStore` interface used directly by HTTP handlers. That is acceptable for a small codebase, but a service layer gives a clearer place for fallback policy, not-found behavior, and future orchestration.

4. Select the default checklist event deterministically until the API exposes event selection.

   Because `/api/checklist` currently has no event parameter, the repository should read one deterministic event, such as the most recent `raceday_events` row by `event_date` and `id`, then return that event's template items. If no event exists, the service should return a not-found style error that maps to an HTTP error instead of silently returning hard-coded data.

   Alternative considered: always read the first template and default all items to incomplete. That can be useful later, but it blurs the difference between reusable template data and event-specific completion state.

5. Preserve hard-coded checklist data only as test fixtures or optional seed data outside the router.

   The existing defaults are useful examples, but keeping them inside `router.go` would preserve the current coupling. If seed data is needed, it should be added through a separate seed mechanism or future migration decision, not hidden in HTTP code.

   Alternative considered: keep hard-coded fallback when MySQL has no rows. That makes local demos easy, but it can hide failed database setup and contradicts the goal of database-backed retrieval.

6. Use SQL-focused tests without a live MySQL dependency.

   Repository tests should use a SQL mocking package or a small query abstraction so they can verify query arguments, row ordering, completion merging, and error paths without requiring MySQL. Router tests should inject a fake checklist service and assert HTTP behavior.

   Alternative considered: integration tests against a real MySQL container. Those would be valuable later, but they add environment requirements that are disproportionate for this first service/repository boundary.

## Risks / Trade-offs

- Default event selection may not match real user workflows -> Keep the selection deterministic and document it as temporary until event-aware endpoints exist.
- Preserving the flat API may lose section metadata -> Use section title as `Category` for compatibility and leave nested sections for a future API change.
- Returning an error when no event exists may make local startup feel less complete -> This is preferable to silently serving stale hard-coded data; seed data can be handled explicitly later.
- SQL mock tests may miss MySQL-specific behavior -> Keep migration validation in place and consider container-backed integration tests when write paths arrive.
- Service and repository layers add structure before many endpoints exist -> Keep interfaces narrow and package-local where possible so the abstraction remains proportional.

## Migration Plan

1. Add checklist domain/service/repository code without changing the public route path or response shape.
2. Update `main.go` to build a MySQL repository and service from the existing `*sql.DB`, then pass the service into `server.NewRouter`.
3. Update router tests to inject fake service behavior.
4. Add repository/service tests for ordered reads, completion state, no-event behavior, and database error propagation.
5. Rollback by reverting the service/repository wiring; no database schema change is required by this change.

## Open Questions

- What should identify the default raceday event once multiple events exist: most recent event date, latest created row, configured active event, or an explicit request parameter?
- Should `/api/checklist` return `404`, `409`, or `500` when the database is reachable but no checklist event has been configured?
- Should seed data be introduced as a separate change before database-backed reads become the default in local development?
