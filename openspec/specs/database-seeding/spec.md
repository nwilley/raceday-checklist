## Purpose

Defines repeatable MySQL seed behavior for starter race-day checklist data and an explicit backend seeder command.

## Requirements

### Requirement: Explicit database seed command

The system SHALL provide an explicit backend command for seeding starter checklist data into the configured MySQL database.

#### Scenario: Run seed command

- **WHEN** the seed command is run with valid database configuration
- **THEN** it connects to MySQL and applies the starter checklist seed data

#### Scenario: Seed command failure

- **WHEN** the seed command cannot connect or apply seed data
- **THEN** it exits with an error instead of starting the API server

### Requirement: Idempotent seed operation

The database seeder SHALL be safe to run multiple times without creating duplicate checklist templates, sections, items, events, or completion rows.

#### Scenario: Re-run seed command

- **WHEN** the seed command is run after seed data already exists
- **THEN** the database still contains one seeded template, one seeded default event, and one row per seeded section and item slug in its parent scope

#### Scenario: Preserve completion state

- **WHEN** event completion rows already exist for seeded items
- **THEN** the seed command does not reset those completion states

### Requirement: Starter checklist seed data

The database seeder SHALL create starter raceday checklist data that matches the requested race-day maintenance workflow.

#### Scenario: Seed ordered sections

- **WHEN** the seed operation completes
- **THEN** the database contains ordered sections titled `Pre-practice`, `Pre-qualifying`, `Mid-qualifying maintenance`, and `Pre-main`

#### Scenario: Seed pre-practice items

- **WHEN** the seed operation completes
- **THEN** the `Pre-practice` section contains `Set Droop`, `Set Camber`, and `Set Ride Height` in that order

#### Scenario: Seed pre-qualifying items

- **WHEN** the seed operation completes
- **THEN** the `Pre-qualifying` section contains `Swap Diffs`, `Bleed Shocks`, `Swap Hinge Pins`, `Bolt Check`, `Clean Car`, `Set Droop`, `Set Camber`, and `Set Ride Height` in that order

#### Scenario: Seed mid-qualifying maintenance items

- **WHEN** the seed operation completes
- **THEN** the `Mid-qualifying maintenance` section contains `Rebuild first diff set`, `Polish first hinge pin set`, `Check droop before Q2`, `Check camber before Q2`, and `Check ride height before Q2` in that order

#### Scenario: Seed pre-main items

- **WHEN** the seed operation completes
- **THEN** the `Pre-main` section contains `Swap in fresh diffs`, `Swap in fresh hinge pins`, `Bleed Shocks`, `Bolt Check`, `Clean Car`, `Set Droop`, `Set Camber`, and `Set Ride Height` in that order

#### Scenario: Seed ordered categories

- **WHEN** the seed operation creates starter sections and items
- **THEN** each seeded item belongs to an ordered section whose title can be returned as the flat API `category`

#### Scenario: Reuse action slugs across sections

- **WHEN** repeated actions such as `Set Droop`, `Set Camber`, `Set Ride Height`, `Bleed Shocks`, `Bolt Check`, or `Clean Car` appear in multiple sections
- **THEN** the seeder stores each repeated action with a stable item slug scoped to its own section

#### Scenario: Seed initial default event

- **WHEN** the seed operation completes
- **THEN** the database contains an initial raceday event that references the seeded checklist template

### Requirement: Seeder is separate from API startup

The database seeder SHALL run only through its explicit seed entry point and not as an implicit side effect of API server startup.

#### Scenario: Start API server

- **WHEN** the API server starts
- **THEN** it does not automatically insert or update seed data
