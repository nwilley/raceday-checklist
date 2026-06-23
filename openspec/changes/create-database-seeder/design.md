## Context

The API now has a MySQL schema for checklist templates, sections, items, raceday events, and event completion state, and `/api/checklist` reads the default event checklist through a repository. A fresh database can satisfy the schema but still return no checklist because no template, event, or item rows exist.

This change adds an explicit seeding path for starter data. The seeder should be safe to run repeatedly, should reuse the existing database configuration and connection code, and should remain separate from normal API server startup so production environments decide when seed data is appropriate.

## Goals / Non-Goals

**Goals:**

- Provide a Go package or command that seeds starter raceday checklist data into MySQL.
- Seed one default checklist template with the requested ordered sections and starter maintenance checklist items.
- Seed one initial raceday event that uses the default template so `/api/checklist` has a deterministic default event to read.
- Make seeding idempotent by using stable slugs and database uniqueness constraints.
- Add tests or validation for seed statements, ordering, idempotency behavior, and created data shape without requiring a live MySQL server.

**Non-Goals:**

- Run seed logic automatically from the API server on every startup.
- Add user-configurable seed files or a general fixture framework.
- Add checklist mutation APIs.
- Add environment-specific production data management.
- Change the `/api/checklist` response shape.

## Decisions

1. Add a dedicated seed command under the API module.

   A command such as `api/cmd/seed` can load the same environment configuration as the API server, connect through `database.Connect`, run seed logic, and exit. This keeps seed execution explicit and scriptable for local setup and deployments.

   Alternative considered: seed automatically during `api/cmd/server` startup. That is convenient locally, but it hides data mutation in server boot and can surprise production deployments.

2. Keep seed behavior in an internal package.

   The command should be thin and delegate to an internal seeding package that accepts `context.Context` and `*sql.DB`. This makes the seeding behavior testable and reusable without coupling tests to command-line process behavior.

   Alternative considered: put all SQL directly in `cmd/seed/main.go`. That is quick, but it makes unit testing and future seed changes harder.

3. Use stable slugs and upsert-style writes for idempotency.

   The seeder should insert the default template, sections, items, and event using stable slugs. For rows protected by unique keys, the SQL should either no-op or update mutable display fields when the row already exists. Re-running the seeder should not create duplicates.

   Alternative considered: delete and recreate seed data. That is simpler, but it can destroy event completion state and makes repeated local use risky.

4. Seed an initial event separately from reusable template data.

   The schema separates reusable checklist definitions from event completion state. The seeder should create a reusable template and an initial event that references it, while leaving completion rows absent so items default to incomplete.

   Alternative considered: seed event completion rows for every item. That is unnecessary because the repository already treats missing completion rows as not done.

5. Seed the requested race-day maintenance workflow exactly.

   Seed four ordered sections with stable section slugs:
   - `pre-practice`: Set Droop, Set Camber, Set Ride Height
   - `pre-qualifying`: Swap Diffs, Bleed Shocks, Swap Hinge Pins, Bolt Check, Clean Car, Set Droop, Set Camber, Set Ride Height
   - `mid-qualifying-maintenance`: Rebuild first diff set, Polish first hinge pin set, Check droop before Q2, Check camber before Q2, Check ride height before Q2
   - `pre-main`: Swap in fresh diffs, Swap in fresh hinge pins, Bleed Shocks, Bolt Check, Clean Car, Set Droop, Set Camber, Set Ride Height

   Item slugs should be stable within their section scope, so repeated natural actions such as `set-droop`, `set-camber`, `set-ride-height`, `bleed-shocks`, `bolt-check`, and `clean-car` can appear in multiple sections without violating schema constraints.

   Alternative considered: introduce large sample race-day data. That would be better handled by separate demo fixtures after the core seeder exists.

## Risks / Trade-offs

- Seed data changes after users edit data -> Use stable slugs and conservative upserts; avoid overwriting completion state.
- No automatic seed means local setup has an extra step -> Document or expose a clear command, and keep the command fast and repeatable.
- Upsert SQL can become verbose -> Keep the seed set small and isolate SQL in one package.
- Tests without live MySQL may miss syntax issues -> Validate generated SQL shape in tests and rely on migration/schema tests until integration tests are introduced.

## Migration Plan

1. Apply the existing checklist schema migrations.
2. Run the new seed command against the target MySQL database.
3. Start the API server and verify `/api/checklist` returns seeded starter items.
4. Roll back by deleting seeded rows manually in dependency-safe order if seed data is no longer desired; this change does not require a schema rollback.

## Open Questions

- What slug and display name should the default event use long term?
- Should the initial event date be fixed for idempotency or configurable at seed time?
- Should seed command documentation live in the README, Makefile, or both?
