## ADDED Requirements

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
