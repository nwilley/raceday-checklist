package checklist

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestMySQLRepositoryListDefaultEventItems(t *testing.T) {
	queryer := &fakeQueryer{
		rows: &fakeRows{
			values: [][]any{
				{"fuel", "Fuel and fluids checked", "Car", true},
				{"helmet", "Helmet packed", "Driver", false},
			},
		},
	}
	repository := newMySQLRepositoryWithQueryer(queryer)

	items, err := repository.ListDefaultEventItems(context.Background())
	if err != nil {
		t.Fatalf("list default event items: %v", err)
	}

	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}
	if items[0] != (Item{ID: "fuel", Title: "Fuel and fluids checked", Category: "Car", Done: true}) {
		t.Fatalf("unexpected first item: %#v", items[0])
	}
	if items[1] != (Item{ID: "helmet", Title: "Helmet packed", Category: "Driver", Done: false}) {
		t.Fatalf("unexpected second item: %#v", items[1])
	}

	assertQueryContains(t, queryer.query, "ORDER BY event_date DESC, id DESC")
	assertQueryContains(t, queryer.query, "LEFT JOIN event_checklist_item_completions")
	assertQueryContains(t, queryer.query, "COALESCE(event_checklist_item_completions.done, checklist_items.default_done)")
	assertQueryContains(t, queryer.query, "ORDER BY checklist_sections.display_order, checklist_items.display_order, checklist_items.id")
}

func TestMySQLRepositoryListDefaultEventItemsUnavailable(t *testing.T) {
	repository := newMySQLRepositoryWithQueryer(&fakeQueryer{rows: &fakeRows{}})

	_, err := repository.ListDefaultEventItems(context.Background())
	if !errors.Is(err, ErrUnavailable) {
		t.Fatalf("expected ErrUnavailable, got %v", err)
	}
}

func TestMySQLRepositoryListDefaultEventItemsQueryError(t *testing.T) {
	expected := errors.New("query failed")
	repository := newMySQLRepositoryWithQueryer(&fakeQueryer{err: expected})

	_, err := repository.ListDefaultEventItems(context.Background())
	if !errors.Is(err, expected) {
		t.Fatalf("expected query error, got %v", err)
	}
}

func TestMySQLRepositoryListDefaultEventItemsRowsError(t *testing.T) {
	expected := errors.New("rows failed")
	repository := newMySQLRepositoryWithQueryer(&fakeQueryer{rows: &fakeRows{err: expected}})

	_, err := repository.ListDefaultEventItems(context.Background())
	if !errors.Is(err, expected) {
		t.Fatalf("expected rows error, got %v", err)
	}
}

func assertQueryContains(t *testing.T, query string, want string) {
	t.Helper()
	if !strings.Contains(query, want) {
		t.Fatalf("expected query to contain %q, got:\n%s", want, query)
	}
}

type fakeQueryer struct {
	query string
	rows  rows
	err   error
}

func (queryer *fakeQueryer) QueryContext(_ context.Context, query string, _ ...any) (rows, error) {
	queryer.query = query
	return queryer.rows, queryer.err
}

type fakeRows struct {
	values [][]any
	index  int
	err    error
}

func (rows *fakeRows) Next() bool {
	return rows.index < len(rows.values)
}

func (rows *fakeRows) Scan(dest ...any) error {
	current := rows.values[rows.index]
	rows.index++

	for i := range dest {
		switch target := dest[i].(type) {
		case *string:
			*target = current[i].(string)
		case *bool:
			*target = current[i].(bool)
		default:
			return errors.New("unsupported scan target")
		}
	}

	return nil
}

func (rows *fakeRows) Close() error {
	return nil
}

func (rows *fakeRows) Err() error {
	return rows.err
}
