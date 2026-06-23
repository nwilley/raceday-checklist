## Why

Checklist completion currently changes only in frontend React state, so a refresh loses progress during a race day. Users should be able to check items off and keep that state throughout the day without creating accounts or logging in.

## What Changes

- Add anonymous participant identity based on a browser-generated client ID instead of full authentication.
- Persist checklist item completion state per default raceday event and anonymous participant.
- Add stable item identity to checklist responses using section and item slugs so repeated actions can be updated unambiguously.
- Add a backend endpoint for setting a checklist item's completion state using `sectionId`, `itemId`, and `done`.
- Update the frontend to optimistically update checkbox state, then persist the change to the backend.
- Handle failed saves with visible unsaved/error state instead of silently losing user intent.
- Keep login, accounts, passwords, and cross-device identity out of scope for this change.

## Capabilities

### New Capabilities

- `anonymous-participant-state`: Defines account-free participant identity and per-participant checklist completion persistence.

### Modified Capabilities

- `checklist-service-repository`: Adds requirements for stable checklist item keys, completion update endpoints, and optimistic frontend persistence behavior.
- `mysql-checklist-storage`: Adds requirements for participant-scoped completion rows and durable completion updates.

## Impact

- Affected code: `web/src/App.tsx`, frontend API helpers or state handling, `api/internal/server`, `api/internal/checklist`, MySQL migrations, and repository tests.
- Affected APIs: `GET /api/checklist` response gains stable section/item identifiers; a new completion update endpoint is added.
- Dependencies: no external auth dependency is expected.
- Systems: completion state persists for the same browser/device through a client ID stored locally; clearing browser storage or switching devices creates a separate anonymous participant.
