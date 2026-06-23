## 1. Seed Data Model

- [x] 1.1 Create an internal seed package for checklist seed data and seeding behavior.
- [x] 1.2 Define stable seed constants for the default template, initial raceday event, and ordered sections `pre-practice`, `pre-qualifying`, `mid-qualifying-maintenance`, and `pre-main`.
- [x] 1.3 Define seeded items for Pre-practice: Set Droop, Set Camber, and Set Ride Height.
- [x] 1.4 Define seeded items for Pre-qualifying: Swap Diffs, Bleed Shocks, Swap Hinge Pins, Bolt Check, Clean Car, Set Droop, Set Camber, and Set Ride Height.
- [x] 1.5 Define seeded items for Mid-qualifying maintenance: Rebuild first diff set, Polish first hinge pin set, Check droop before Q2, Check camber before Q2, and Check ride height before Q2.
- [x] 1.6 Define seeded items for Pre-main: Swap in fresh diffs, Swap in fresh hinge pins, Bleed Shocks, Bolt Check, Clean Car, Set Droop, Set Camber, and Set Ride Height.

## 2. Seed Execution

- [x] 2.1 Implement idempotent template insertion using stable slugs and existing uniqueness constraints.
- [x] 2.2 Implement idempotent section insertion with template-scoped slugs and display ordering.
- [x] 2.3 Implement idempotent item insertion with section-scoped slugs and display ordering.
- [x] 2.4 Implement idempotent initial raceday event insertion linked to the seeded template.
- [x] 2.5 Avoid inserting or resetting event completion rows during seeding.

## 3. Seed Command

- [x] 3.1 Add an API seed command entry point that loads existing config and opens a MySQL connection.
- [x] 3.2 Wire the seed command to run the internal seed package and exit with an error on connection or seed failure.
- [x] 3.3 Keep API server startup unchanged so seeding only runs through the explicit seed command.

## 4. Validation

- [x] 4.1 Add unit tests or SQL-shape tests for seed ordering and idempotent upsert behavior without requiring a live MySQL server.
- [x] 4.2 Add tests confirming the exact four seeded sections and their item titles/display order match the requested checklist.
- [x] 4.3 Add tests confirming seeded data supports the default checklist read path with incomplete item state when no completion rows exist.
- [x] 4.4 Add validation that the seed command does not touch event completion state.
- [x] 4.5 Run `go test ./...` from the API module and fix any compile or test failures.

## 5. Documentation and OpenSpec Confirmation

- [x] 5.1 Document how to run the seed command after migrations.
- [x] 5.2 Confirm the OpenSpec requirements for `database-seeding` are satisfied.
- [x] 5.3 Confirm the added `mysql-checklist-storage` seeded-read requirements are satisfied.
