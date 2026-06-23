## Context

The frontend currently toggles checklist items only in React state, so progress is lost on refresh. The backend already reads completion state from `event_checklist_item_completions`, but there is no write endpoint and no way to distinguish one user's completion state from another user's without accounts.

The product direction is intentionally account-free. A racer should be able to open the app at the track, check items off, and keep progress on that same browser/device without creating a login. The current API response also needs stronger item identity because item slugs repeat across sections, such as `set-droop` appearing in multiple phases.

## Goals / Non-Goals

**Goals:**

- Persist checklist completion state through the backend after optimistic frontend checkbox updates.
- Identify anonymous participants by a browser-generated client ID stored locally.
- Scope completion state to the selected/default raceday event and anonymous participant.
- Add stable checklist item identity to API responses using section slug plus item slug.
- Add a completion update endpoint that accepts `sectionId`, `itemId`, and `done`.
- Preserve the current no-login user experience.
- Show failed save state clearly enough that users know a change did not persist.

**Non-Goals:**

- Add full authentication, passwords, email login, passkeys, or account management.
- Guarantee state follows a user across devices or cleared browser storage.
- Add collaborative realtime updates.
- Add event creation or event selection UI.
- Change checklist display from the current sectioned frontend layout.

## Decisions

1. Use anonymous browser client IDs instead of accounts.

   The frontend will create a UUID once and store it in local storage. Requests that read or write checklist completion state will include that value in a header such as `X-Raceday-Client`. The backend will find or create a participant row for the default raceday event and that client ID.

   Alternative considered: full authentication. That would solve cross-device continuity, but it adds friction exactly where the product should stay fast and track-friendly.

2. Add `raceday_participants` and scope completions by participant.

   Add a participant table with `event_id`, `client_id`, optional `display_name`, and timestamps. Completion rows should include `participant_id` so two anonymous users checking the same item do not overwrite each other.

   Alternative considered: store `client_id` directly on `event_checklist_item_completions`. That is simpler, but a participant row gives a natural place for optional display name and future participant metadata without changing every completion row again.

3. Keep section slug plus item slug as the public item key.

   `GET /api/checklist` should continue returning `id`, `title`, `category`, and `done`, while also adding stable identifiers such as `sectionId` and `itemId`. The frontend should use `sectionId + itemId` for writes because item slugs are only unique within a section.

   Alternative considered: expose internal database IDs. That is easy, but it couples clients to storage internals and makes seeded slug-based behavior less transparent.

4. Add a targeted completion update endpoint.

   Use an endpoint shaped around item completion, for example `PATCH /api/checklist/items/:sectionId/:itemId`, with body `{ "done": true }`. The handler resolves the default event, anonymous participant, section, and item, then upserts the completion row.

   Alternative considered: send the entire checklist after every change. That increases payload size, makes conflict handling less precise, and complicates error recovery.

5. Use optimistic frontend updates with per-item save state.

   On checkbox change, the frontend should update local item state immediately and send the PATCH request. While pending, the item can be treated as saving. On success, clear the pending state. On failure, either revert the item or show it as unsaved; for this app, reverting with an inline error is the clearer first behavior because backend state remains authoritative.

   Alternative considered: block the checkbox until the backend responds. That is safer but feels sluggish in a race-day workflow.

6. Keep backend state authoritative.

   The frontend client ID is stored locally, but completion state is not primarily stored in local storage. On load, the app fetches the backend checklist for the current client ID and uses that response as the source of truth.

   Alternative considered: hybrid local storage overlay with background sync. That helps offline usage, but introduces reconciliation complexity before offline behavior is explicitly needed.

## Risks / Trade-offs

- Browser storage is cleared -> The participant identity is lost; users get a fresh anonymous checklist. This is acceptable for no-account mode and should be documented if needed.
- Client ID header can be spoofed -> This is not authentication. Treat it as convenience identity, not security or ownership. Avoid using it for sensitive data.
- Existing completion rows lack participant scope -> Migration must add participant-aware completion storage carefully. Because existing seeded data defaults incomplete, this can be handled with a new table shape or migration that creates a default anonymous participant if needed.
- Duplicate item slugs across sections -> Require `sectionId + itemId` for writes and add tests for repeated slugs.
- Failed save after optimistic update -> Revert or surface unsaved state immediately; avoid silently diverging from backend state.

## Migration Plan

1. Add a migration for anonymous participants and participant-scoped completion uniqueness.
2. Update repository reads to resolve/create participant by client ID and return completion state for that participant.
3. Add stable `sectionId` and `itemId` fields to checklist response items.
4. Add the completion update endpoint and repository upsert.
5. Update the frontend to create/store a client ID, include it in GET/PATCH requests, and optimistically persist checkbox changes.
6. Roll back by reverting endpoint/frontend changes and restoring the prior completion table shape if no participant-scoped data must be preserved.

## Open Questions

- Should a missing or invalid client ID create a new participant automatically, or should the backend reject it and let the frontend generate one?
- Should failed optimistic saves revert immediately or remain visually checked with an unsaved marker and retry option?
- Should participant display names be included now as optional metadata, or deferred until a UI actually needs them?
