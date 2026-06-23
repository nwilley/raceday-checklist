## ADDED Requirements

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
