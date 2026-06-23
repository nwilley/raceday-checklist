## ADDED Requirements

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
