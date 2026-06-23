package checklist

import (
	"context"
	"database/sql"
	"errors"
)

const listDefaultEventItemsQuery = `
SELECT
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
  AND event_checklist_item_completions.item_id = checklist_items.id
WHERE raceday_events.id = (
  SELECT id
  FROM raceday_events
  ORDER BY event_date DESC, id DESC
  LIMIT 1
)
ORDER BY checklist_sections.display_order, checklist_items.display_order, checklist_items.id`

type queryer interface {
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

type MySQLRepository struct {
	queryer queryer
}

func NewMySQLRepository(db *sql.DB) MySQLRepository {
	return MySQLRepository{queryer: NewSQLDB(db)}
}

func newMySQLRepositoryWithQueryer(queryer queryer) MySQLRepository {
	return MySQLRepository{queryer: queryer}
}

func (repository MySQLRepository) ListDefaultEventItems(ctx context.Context) ([]Item, error) {
	resultRows, err := repository.queryer.QueryContext(ctx, listDefaultEventItemsQuery)
	if err != nil {
		return nil, err
	}
	defer resultRows.Close()

	var items []Item
	for resultRows.Next() {
		var item Item
		if err := resultRows.Scan(&item.ID, &item.Title, &item.Category, &item.Done); err != nil {
			return nil, err
		}
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

func IsUnavailable(err error) bool {
	return errors.Is(err, ErrUnavailable)
}
