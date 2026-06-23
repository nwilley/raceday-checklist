package checklist

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

const listDefaultEventItemsQuery = `
SELECT
  checklist_sections.slug,
  checklist_items.slug,
  checklist_items.title,
  checklist_sections.title AS category,
  COALESCE(event_checklist_item_completions.done, checklist_items.default_done) AS done
FROM raceday_events
JOIN checklist_templates
  ON checklist_templates.id = raceday_events.template_id
JOIN checklist_sections
  ON checklist_sections.template_id = checklist_templates.id
JOIN checklist_items
  ON checklist_items.section_id = checklist_sections.id
LEFT JOIN event_checklist_item_completions
  ON event_checklist_item_completions.event_id = raceday_events.id
  AND event_checklist_item_completions.participant_id = ?
  AND event_checklist_item_completions.item_id = checklist_items.id
WHERE raceday_events.id = (
  SELECT id
  FROM raceday_events
  ORDER BY event_date DESC, id DESC
  LIMIT 1
)
ORDER BY checklist_sections.display_order, checklist_items.display_order, checklist_items.id`

const upsertDefaultEventParticipantQuery = `
INSERT INTO raceday_participants (event_id, client_id)
SELECT id, ?
FROM raceday_events
ORDER BY event_date DESC, id DESC
LIMIT 1
ON DUPLICATE KEY UPDATE id = LAST_INSERT_ID(id)`

const findDefaultEventItemQuery = `
SELECT raceday_events.id, checklist_items.id
FROM raceday_events
JOIN checklist_templates
  ON checklist_templates.id = raceday_events.template_id
JOIN checklist_sections
  ON checklist_sections.template_id = checklist_templates.id
JOIN checklist_items
  ON checklist_items.section_id = checklist_sections.id
WHERE raceday_events.id = (
  SELECT id
  FROM raceday_events
  ORDER BY event_date DESC, id DESC
  LIMIT 1
)
  AND checklist_sections.slug = ?
  AND checklist_items.slug = ?
LIMIT 1`

const upsertCompletionQuery = `
INSERT INTO event_checklist_item_completions (event_id, participant_id, item_id, done, completed_at)
VALUES (?, ?, ?, ?, CASE WHEN ? THEN CURRENT_TIMESTAMP ELSE NULL END)
ON DUPLICATE KEY UPDATE
  done = VALUES(done),
  completed_at = VALUES(completed_at)`

type databaseExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (rows, error)
}

type rows interface {
	Next() bool
	Scan(dest ...any) error
	Close() error
	Err() error
}

type SQLDB struct {
	db *sql.DB
}

func NewSQLDB(db *sql.DB) SQLDB {
	return SQLDB{db: db}
}

func (database SQLDB) QueryContext(ctx context.Context, query string, args ...any) (rows, error) {
	return database.db.QueryContext(ctx, query, args...)
}

func (database SQLDB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return database.db.ExecContext(ctx, query, args...)
}

type MySQLRepository struct {
	database databaseExecutor
}

func NewMySQLRepository(db *sql.DB) MySQLRepository {
	return MySQLRepository{database: NewSQLDB(db)}
}

func newMySQLRepositoryWithDatabase(database databaseExecutor) MySQLRepository {
	return MySQLRepository{database: database}
}

func (repository MySQLRepository) ListDefaultEventItems(ctx context.Context, clientID string) ([]Item, error) {
	participantID, err := repository.resolveDefaultEventParticipant(ctx, clientID)
	if err != nil {
		return nil, err
	}

	resultRows, err := repository.database.QueryContext(ctx, listDefaultEventItemsQuery, participantID)
	if err != nil {
		return nil, err
	}
	defer resultRows.Close()

	var items []Item
	for resultRows.Next() {
		var item Item
		if err := resultRows.Scan(&item.SectionID, &item.ItemID, &item.Title, &item.Category, &item.Done); err != nil {
			return nil, err
		}
		item.ID = item.ItemID
		items = append(items, item)
	}

	if err := resultRows.Err(); err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, ErrUnavailable
	}

	return items, nil
}

func (repository MySQLRepository) UpdateDefaultEventItemCompletion(ctx context.Context, clientID string, update CompletionUpdate) error {
	participantID, err := repository.resolveDefaultEventParticipant(ctx, clientID)
	if err != nil {
		return err
	}

	eventID, itemID, err := repository.findDefaultEventItem(ctx, update.SectionID, update.ItemID)
	if err != nil {
		return err
	}

	if _, err := repository.database.ExecContext(
		ctx,
		upsertCompletionQuery,
		eventID,
		participantID,
		itemID,
		update.Done,
		update.Done,
	); err != nil {
		return fmt.Errorf("upsert completion: %w", err)
	}

	return nil
}

func (repository MySQLRepository) resolveDefaultEventParticipant(ctx context.Context, clientID string) (int64, error) {
	result, err := repository.database.ExecContext(ctx, upsertDefaultEventParticipantQuery, clientID)
	if err != nil {
		return 0, fmt.Errorf("upsert participant: %w", err)
	}

	participantID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("read participant id: %w", err)
	}
	if participantID == 0 {
		return 0, ErrUnavailable
	}

	return participantID, nil
}

func (repository MySQLRepository) findDefaultEventItem(ctx context.Context, sectionID string, itemID string) (int64, int64, error) {
	resultRows, err := repository.database.QueryContext(ctx, findDefaultEventItemQuery, sectionID, itemID)
	if err != nil {
		return 0, 0, err
	}
	defer resultRows.Close()

	if !resultRows.Next() {
		if err := resultRows.Err(); err != nil {
			return 0, 0, err
		}
		return 0, 0, ErrNotFound
	}

	var eventID int64
	var databaseItemID int64
	if err := resultRows.Scan(&eventID, &databaseItemID); err != nil {
		return 0, 0, err
	}
	if resultRows.Next() {
		return 0, 0, errors.New("expected one default event item")
	}
	if err := resultRows.Err(); err != nil {
		return 0, 0, err
	}

	return eventID, databaseItemID, nil
}

func IsUnavailable(err error) bool {
	return errors.Is(err, ErrUnavailable)
}
