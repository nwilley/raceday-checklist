## ADDED Requirements

### Requirement: Anonymous participants storage

The MySQL schema SHALL store anonymous participants for raceday events using browser-provided client IDs.

#### Scenario: Store participant for event

- **WHEN** a checklist read or write request includes a new anonymous client ID
- **THEN** the backend can store a participant row scoped to the selected raceday event and client ID

#### Scenario: Reuse participant for event

- **WHEN** a checklist read or write request includes an existing anonymous client ID for the selected raceday event
- **THEN** the backend reuses the existing participant row

### Requirement: Participant-scoped completion rows

The MySQL schema SHALL scope event checklist item completion rows by anonymous participant.

#### Scenario: Persist participant completion

- **WHEN** an anonymous participant marks an item done
- **THEN** MySQL stores the completion state for that participant and item within the raceday event

#### Scenario: Prevent duplicate participant item state

- **WHEN** the same anonymous participant updates the same item more than once
- **THEN** MySQL maintains one completion row for that participant and item within the raceday event

#### Scenario: Preserve independent participant state

- **WHEN** two anonymous participants update the same checklist item in the same raceday event
- **THEN** MySQL stores independent completion state for each participant

### Requirement: Durable completion updates

The backend SHALL upsert completion state changes into MySQL and return the saved state on subsequent checklist reads.

#### Scenario: Read after save

- **WHEN** a participant marks an item done and then reloads the checklist with the same client ID
- **THEN** the checklist read returns that item as done

#### Scenario: Clear completion

- **WHEN** a participant marks a previously done item as not done
- **THEN** subsequent checklist reads for that participant return that item as not done
