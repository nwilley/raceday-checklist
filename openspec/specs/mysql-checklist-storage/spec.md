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
