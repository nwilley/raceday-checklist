## Purpose

Defines checklist retrieval through backend service and repository boundaries for the existing `/api/checklist` endpoint.

## Requirements

### Requirement: Checklist route delegates to service

The backend SHALL route `/api/checklist` requests through an injected checklist service rather than constructing checklist items inside the HTTP router.

#### Scenario: Successful checklist response

- **WHEN** the checklist service returns checklist items
- **THEN** `GET /api/checklist` responds with HTTP 200 and a JSON object containing those items

#### Scenario: Router test with injected service

- **WHEN** router tests create a router with a fake checklist service
- **THEN** the route behavior can be verified without opening a MySQL connection

#### Scenario: Service failure response

- **WHEN** the checklist service returns an unexpected error
- **THEN** `GET /api/checklist` responds with a server error and does not expose internal database details

### Requirement: Checklist response compatibility

The backend SHALL preserve the existing flat `/api/checklist` response shape while moving data retrieval behind the service and repository boundary.

#### Scenario: Item fields preserved

- **WHEN** `GET /api/checklist` returns checklist items
- **THEN** each item contains `id`, `title`, `category`, and `done` fields

#### Scenario: Items wrapper preserved

- **WHEN** `GET /api/checklist` succeeds
- **THEN** the response body contains an `items` array at the top level

### Requirement: Checklist service handles missing configured checklist

The checklist service SHALL return a distinct not-found style error when no default raceday event or checklist data is configured.

#### Scenario: No configured checklist

- **WHEN** the repository reports that no default checklist event exists
- **THEN** the service returns a missing checklist error instead of hard-coded fallback data

#### Scenario: Missing checklist HTTP response

- **WHEN** `GET /api/checklist` is requested and no default checklist event exists
- **THEN** the API responds with a client-visible error status that identifies the checklist as unavailable

### Requirement: Checklist service owns use-case behavior

The backend SHALL keep checklist retrieval use-case behavior in a service layer separate from HTTP routing and SQL persistence.

#### Scenario: Router has no SQL dependency

- **WHEN** the router handles `GET /api/checklist`
- **THEN** it calls the checklist service without referencing SQL tables, SQL queries, or `*sql.DB`

#### Scenario: Repository error propagation

- **WHEN** the repository cannot retrieve checklist data
- **THEN** the service returns an error that the router can map to an HTTP response

### Requirement: Stable checklist item keys

The checklist API SHALL expose stable section and item identifiers that allow clients to update repeated item slugs unambiguously.

#### Scenario: Checklist response includes stable keys

- **WHEN** `GET /api/checklist` returns checklist items
- **THEN** each item includes `sectionId` and `itemId` values derived from checklist section and item slugs

#### Scenario: Repeated item slug in different sections

- **WHEN** the checklist contains the same item slug in multiple sections
- **THEN** each item remains uniquely addressable by the combination of `sectionId` and `itemId`

### Requirement: Checklist completion update endpoint

The backend SHALL expose an endpoint that updates one checklist item's completion state for the anonymous participant.

#### Scenario: Mark item done

- **WHEN** a client sends a completion update for `sectionId`, `itemId`, and `done: true`
- **THEN** the backend persists that item as done for the request's anonymous participant

#### Scenario: Mark item not done

- **WHEN** a client sends a completion update for `sectionId`, `itemId`, and `done: false`
- **THEN** the backend persists that item as not done for the request's anonymous participant

#### Scenario: Unknown item key

- **WHEN** a client sends a completion update for a section or item key that does not exist in the default event checklist
- **THEN** the backend returns a client-visible not-found response

### Requirement: Optimistic frontend completion persistence

The frontend SHALL update checkbox state immediately on user interaction and then persist that change to the backend.

#### Scenario: Successful optimistic save

- **WHEN** a user checks or unchecks an item and the backend save succeeds
- **THEN** the UI keeps the updated state and clears any pending save indicator for that item

#### Scenario: Failed optimistic save

- **WHEN** a user checks or unchecks an item and the backend save fails
- **THEN** the UI clearly indicates the save failure and does not silently pretend the item was durably saved

#### Scenario: Completion state after reload

- **WHEN** the page reloads after successful completion saves
- **THEN** the frontend displays the completion state returned by the backend for the same anonymous client ID
