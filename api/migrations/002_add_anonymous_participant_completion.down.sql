ALTER TABLE event_checklist_item_completions
  DROP FOREIGN KEY fk_event_checklist_item_completions_participant,
  DROP INDEX uq_event_checklist_item_completions_participant_item,
  DROP INDEX idx_event_checklist_item_completions_participant,
  DROP COLUMN participant_id,
  ADD UNIQUE KEY uq_event_checklist_item_completions_event_item (event_id, item_id);

DROP TABLE IF EXISTS raceday_participants;
