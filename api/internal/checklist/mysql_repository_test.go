package checklist

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"testing"
)

func TestMySQLRepositoryListDefaultEventItems(t *testing.T) {
	database := &fakeDatabase{
		execResults: []sql.Result{fakeResult{id: 42}},
		queryRows: []rows{
			&fakeRows{
				values: [][]any{
					{"pre-practice", "set-droop", "Set Droop", "Pre-practice", true},
					{"pre-qualifying", "set-droop", "Set Droop", "Pre-qualifying", false},
				},
			},
		},
	}
	repository := newMySQLRepositoryWithDatabase(database)

	items, err := repository.ListDefaultEventItems(context.Background(), "client-1")
	if err != nil {
		t.Fatalf("list default event items: %v", err)
	}

	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}
	if items[0] != (Item{ID: "set-droop", SectionID: "pre-practice", ItemID: "set-droop", Title: "Set Droop", Category: "Pre-practice", Done: true}) {
		t.Fatalf("unexpected first item: %#v", items[0])
	}
	if items[1] != (Item{ID: "set-droop", SectionID: "pre-qualifying", ItemID: "set-droop", Title: "Set Droop", Category: "Pre-qualifying", Done: false}) {
		t.Fatalf("unexpected second item: %#v", items[1])
	}

	assertQueryContains(t, database.execCalls[0].query, "INSERT INTO raceday_participants")
	assertQueryContains(t, database.execCalls[0].query, "SELECT raceday_events.id, ?")
	assertQueryContains(t, database.execCalls[0].query, "ORDER BY event_date DESC, raceday_events.id DESC")
	assertQueryContains(t, database.execCalls[0].query, "ON DUPLICATE KEY UPDATE raceday_participants.id = LAST_INSERT_ID(raceday_participants.id)")
	assertQueryContains(t, database.queryCalls[0].query, "event_checklist_item_completions.participant_id = ?")
	assertQueryContains(t, database.queryCalls[0].query, "checklist_sections.slug")
	assertQueryContains(t, database.queryCalls[0].query, "ORDER BY checklist_sections.display_order, checklist_items.display_order, checklist_items.id")
	if database.queryCalls[0].args[0] != int64(42) {
		t.Fatalf("expected participant id query arg 42, got %#v", database.queryCalls[0].args[0])
	}
}

func TestMySQLRepositoryListDefaultEventItemsUnavailableWhenNoParticipantCreated(t *testing.T) {
	repository := newMySQLRepositoryWithDatabase(&fakeDatabase{execResults: []sql.Result{fakeResult{id: 0}}})

	_, err := repository.ListDefaultEventItems(context.Background(), "client-1")
	if !errors.Is(err, ErrUnavailable) {
		t.Fatalf("expected ErrUnavailable, got %v", err)
	}
}

func TestMySQLRepositoryUpdateDefaultEventItemCompletion(t *testing.T) {
	database := &fakeDatabase{
		execResults: []sql.Result{
			fakeResult{id: 42},
			fakeResult{id: 0},
		},
		queryRows: []rows{
			&fakeRows{values: [][]any{{int64(7), int64(99)}}},
		},
	}
	repository := newMySQLRepositoryWithDatabase(database)

	err := repository.UpdateDefaultEventItemCompletion(context.Background(), "client-1", CompletionUpdate{
		SectionID: "pre-practice",
		ItemID:    "set-droop",
		Done:      true,
	})
	if err != nil {
		t.Fatalf("update completion: %v", err)
	}

	if len(database.execCalls) != 2 {
		t.Fatalf("expected 2 exec calls, got %d", len(database.execCalls))
	}
	assertQueryContains(t, database.queryCalls[0].query, "checklist_sections.slug = ?")
	assertQueryContains(t, database.queryCalls[0].query, "checklist_items.slug = ?")
	assertQueryContains(t, database.execCalls[1].query, "INSERT INTO event_checklist_item_completions")
	assertQueryContains(t, database.execCalls[1].query, "ON DUPLICATE KEY UPDATE")

	wantArgs := []any{int64(7), int64(42), int64(99), true, true}
	for i, want := range wantArgs {
		if database.execCalls[1].args[i] != want {
			t.Fatalf("expected completion arg %d to be %#v, got %#v", i, want, database.execCalls[1].args[i])
		}
	}
}

func TestMySQLRepositoryUpdateDefaultEventItemCompletionNotFound(t *testing.T) {
	database := &fakeDatabase{
		execResults: []sql.Result{fakeResult{id: 42}},
		queryRows:   []rows{&fakeRows{}},
	}
	repository := newMySQLRepositoryWithDatabase(database)

	err := repository.UpdateDefaultEventItemCompletion(context.Background(), "client-1", CompletionUpdate{SectionID: "bad", ItemID: "missing"})
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestMySQLRepositoryPropagatesQueryError(t *testing.T) {
	expected := errors.New("query failed")
	repository := newMySQLRepositoryWithDatabase(&fakeDatabase{
		execResults: []sql.Result{fakeResult{id: 42}},
		queryErr:    expected,
	})

	_, err := repository.ListDefaultEventItems(context.Background(), "client-1")
	if !errors.Is(err, expected) {
		t.Fatalf("expected query error, got %v", err)
	}
}

func assertQueryContains(t *testing.T, query string, want string) {
	t.Helper()
	if !strings.Contains(query, want) {
		t.Fatalf("expected query to contain %q, got:\n%s", want, query)
	}
}

type databaseCall struct {
	query string
	args  []any
}

type fakeDatabase struct {
	execCalls   []databaseCall
	queryCalls  []databaseCall
	execResults []sql.Result
	queryRows   []rows
	execErr     error
	queryErr    error
}

func (database *fakeDatabase) ExecContext(_ context.Context, query string, args ...any) (sql.Result, error) {
	database.execCalls = append(database.execCalls, databaseCall{query: query, args: append([]any(nil), args...)})
	if database.execErr != nil {
		return nil, database.execErr
	}
	if len(database.execResults) == 0 {
		return fakeResult{}, nil
	}
	result := database.execResults[0]
	database.execResults = database.execResults[1:]
	return result, nil
}

func (database *fakeDatabase) QueryContext(_ context.Context, query string, args ...any) (rows, error) {
	database.queryCalls = append(database.queryCalls, databaseCall{query: query, args: append([]any(nil), args...)})
	if database.queryErr != nil {
		return nil, database.queryErr
	}
	if len(database.queryRows) == 0 {
		return &fakeRows{}, nil
	}
	result := database.queryRows[0]
	database.queryRows = database.queryRows[1:]
	return result, nil
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
		case *int64:
			*target = current[i].(int64)
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

type fakeResult struct {
	id int64
}

func (result fakeResult) LastInsertId() (int64, error) {
	return result.id, nil
}

func (result fakeResult) RowsAffected() (int64, error) {
	return 1, nil
}
