## Purpose

Defines persistent MySQL storage for raceday checklist templates, ordered checklist content, raceday events, event-specific completion state, and backend database connectivity.

## Requirements

### Requirement: MySQL checklist schema

The system SHALL provide a MySQL schema that stores reusable raceday checklist templates, ordered checklist sections, ordered checklist items, raceday events, and per-event checklist completion state.

#### Scenario: Initialize an empty database

- **WHEN** the schema migration is applied to an empty MySQL database
- **THEN** the database contains tables for checklist templates, checklist sections, checklist items, raceday events, and event item completion state

#### Scenario: Reuse checklist items across events

- **WHEN** two raceday events are created from the same checklist template
- **THEN** each event can store independent completion state without changing the reusable template item rows

### Requirement: Schema integrity constraints

The schema SHALL enforce referential integrity and uniqueness constraints that prevent orphaned checklist records and duplicate slugs within the same parent scope.

#### Scenario: Reject orphaned checklist item

- **WHEN** a checklist item references a section that does not exist
- **THEN** MySQL rejects the row

#### Scenario: Reject duplicate item slug in section

- **WHEN** two checklist items in the same section use the same slug
- **THEN** MySQL rejects the duplicate row

#### Scenario: Allow same item slug in different sections

- **WHEN** two checklist items in different sections use the same slug
- **THEN** MySQL accepts both rows

### Requirement: Checklist display ordering

The schema SHALL store section display order within a checklist template and item display order within each checklist section.

#### Scenario: Display ordered sections

- **WHEN** a checklist template contains three sections with distinct display order values
- **THEN** the backend can retrieve the sections in that configured order

#### Scenario: Display ordered items within a section

- **WHEN** a checklist section contains five items with distinct display order values
- **THEN** the backend can retrieve the items in that configured order within the section

### Requirement: Ordered migrations

The system SHALL include ordered SQL migration files that create and drop the MySQL checklist schema.

#### Scenario: Apply schema migration

- **WHEN** the `up` migration is executed
- **THEN** all checklist storage tables, indexes, and foreign keys are created in dependency-safe order

#### Scenario: Roll back schema migration

- **WHEN** the `down` migration is executed after the `up` migration
- **THEN** all checklist storage tables are dropped in dependency-safe order

### Requirement: Backend MySQL connection

The backend SHALL expose a database connection function that uses the existing database configuration to open and verify a MySQL connection.

#### Scenario: Connect with valid configuration

- **WHEN** `database.Connect` is called with valid MySQL host, port, user, password, and database name
- **THEN** it opens a MySQL connection and verifies connectivity with the provided context

#### Scenario: Return connection error

- **WHEN** `database.Connect` cannot open or ping MySQL
- **THEN** it returns an error to the caller

### Requirement: Read default event checklist

The system SHALL support reading checklist items for a deterministic default raceday event from the persistent MySQL schema.

#### Scenario: Read most recent event checklist

- **WHEN** multiple raceday events exist
- **THEN** the repository selects a deterministic default event using event date and row identity ordering

#### Scenario: Return no default event

- **WHEN** no raceday event exists
- **THEN** the repository reports that no default event checklist is available

### Requirement: Read ordered checklist items

The system SHALL read checklist sections and items in their configured display order.

#### Scenario: Ordered sections and items

- **WHEN** a default event's template contains multiple sections and items with display order values
- **THEN** repository results are ordered by section display order and then item display order

#### Scenario: Section title as response category

- **WHEN** checklist items are mapped to the existing flat API item model
- **THEN** each item's category value comes from its checklist section title

### Requirement: Read event completion state

The system SHALL combine reusable checklist item rows with per-event completion state when reading a checklist.

#### Scenario: Completed item

- **WHEN** an event completion row marks an item done
- **THEN** the repository returns that item with completion state set to done

#### Scenario: Item without completion row

- **WHEN** a template item has no completion row for the default event
- **THEN** the repository returns that item with completion state set to not done

#### Scenario: Independent event completion state

- **WHEN** two events use the same checklist template item
- **THEN** reading the default event uses only completion rows for the selected event

### Requirement: Seeded data supports default checklist reads

The persistent checklist schema SHALL support seeded starter data that can be read by the default checklist repository path.

#### Scenario: Read seeded checklist

- **WHEN** migrations have been applied and the seed command has completed
- **THEN** the default checklist read returns the seeded Pre-practice, Pre-qualifying, Mid-qualifying maintenance, and Pre-main items

#### Scenario: Seeded items default incomplete

- **WHEN** seeded starter items have no event completion rows
- **THEN** the default checklist read returns those items with completion state set to not done

#### Scenario: Seeded ordering

- **WHEN** the default checklist read returns seeded starter items
- **THEN** the items are ordered by their seeded section display order and item display order

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
