CREATE TABLE raceday_participants (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  event_id BIGINT UNSIGNED NOT NULL,
  client_id VARCHAR(64) NOT NULL,
  display_name VARCHAR(255) NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uq_raceday_participants_event_client (event_id, client_id),
  CONSTRAINT fk_raceday_participants_event
    FOREIGN KEY (event_id) REFERENCES raceday_events (id)
    ON UPDATE CASCADE
    ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

ALTER TABLE event_checklist_item_completions
  DROP INDEX uq_event_checklist_item_completions_event_item,
  ADD COLUMN participant_id BIGINT UNSIGNED NULL AFTER event_id,
  ADD UNIQUE KEY uq_event_checklist_item_completions_participant_item (event_id, participant_id, item_id),
  ADD KEY idx_event_checklist_item_completions_participant (participant_id),
  ADD CONSTRAINT fk_event_checklist_item_completions_participant
    FOREIGN KEY (participant_id) REFERENCES raceday_participants (id)
    ON UPDATE CASCADE
    ON DELETE CASCADE;
