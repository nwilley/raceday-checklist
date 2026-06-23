## Purpose

Defines anonymous participant identity and participant-scoped checklist state without requiring user accounts.

## Requirements

### Requirement: Anonymous participant identity

The system SHALL identify checklist participants without requiring login, accounts, passwords, or email.

#### Scenario: First browser visit

- **WHEN** a user opens the checklist in a browser with no existing race-day client ID
- **THEN** the frontend creates a stable anonymous client ID and stores it locally

#### Scenario: Returning browser visit

- **WHEN** a user returns in the same browser with an existing race-day client ID
- **THEN** the frontend reuses that client ID for checklist read and write requests

#### Scenario: No account prompt

- **WHEN** a user checks off checklist items
- **THEN** the system does not require account creation or login

### Requirement: Participant-scoped checklist state

The system SHALL keep checklist completion state separate for each anonymous participant within a raceday event.

#### Scenario: Independent anonymous participants

- **WHEN** two browsers use different anonymous client IDs for the same raceday event
- **THEN** each browser sees and updates independent checklist completion state

#### Scenario: Same anonymous participant returns

- **WHEN** the same browser returns with the same anonymous client ID
- **THEN** the checklist reflects completion state previously saved for that participant

### Requirement: Client identity is convenience identity only

The system SHALL treat anonymous client IDs as convenience identity, not secure authentication.

#### Scenario: Client ID supplied by request

- **WHEN** the backend receives a race-day client ID header
- **THEN** it uses the client ID only to scope checklist completion state and does not treat it as proof of secure user ownership
