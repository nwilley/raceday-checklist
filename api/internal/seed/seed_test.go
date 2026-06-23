package seed

import (
	"context"
	"database/sql"
	"strings"
	"testing"
)

func TestDefaultPlanMatchesRequestedChecklist(t *testing.T) {
	sections := defaultPlan.Template.Sections
	if len(sections) != 4 {
		t.Fatalf("expected 4 sections, got %d", len(sections))
	}

	assertSection(t, sections[0], "pre-practice", "Pre-practice", []string{
		"Set Droop",
		"Set Camber",
		"Set Ride Height",
	})
	assertSection(t, sections[1], "pre-qualifying", "Pre-qualifying", []string{
		"Swap Diffs",
		"Bleed Shocks",
		"Swap Hinge Pins",
		"Bolt Check",
		"Clean Car",
		"Set Droop",
		"Set Camber",
		"Set Ride Height",
	})
	assertSection(t, sections[2], "mid-qualifying-maintenance", "Mid-qualifying maintenance", []string{
		"Rebuild first diff set",
		"Polish first hinge pin set",
		"Check droop before Q2",
		"Check camber before Q2",
		"Check ride height before Q2",
	})
	assertSection(t, sections[3], "pre-main", "Pre-main", []string{
		"Swap in fresh diffs",
		"Swap in fresh hinge pins",
		"Bleed Shocks",
		"Bolt Check",
		"Clean Car",
		"Set Droop",
		"Set Camber",
		"Set Ride Height",
	})
}

func TestRunUsesIdempotentUpsertsInDependencyOrder(t *testing.T) {
	executor := &fakeExecutor{}

	if err := run(context.Background(), executor, defaultPlan); err != nil {
		t.Fatalf("run seed: %v", err)
	}

	if len(executor.calls) != 30 {
		t.Fatalf("expected 30 seed statements, got %d", len(executor.calls))
	}

	assertCall(t, executor.calls[0], "INSERT INTO checklist_templates", defaultTemplateSlug)
	assertCall(t, executor.calls[1], "INSERT INTO checklist_sections", int64(1), "pre-practice", "Pre-practice", 1)
	assertCall(t, executor.calls[2], "INSERT INTO checklist_items", int64(2), "set-droop", "Set Droop", 1)
	assertCall(t, executor.calls[29], "INSERT INTO raceday_events", int64(1), defaultEventSlug, defaultEventName, defaultEventDate)

	for _, call := range executor.calls {
		if !strings.Contains(call.query, "ON DUPLICATE KEY UPDATE") {
			t.Fatalf("expected idempotent upsert query, got:\n%s", call.query)
		}
		if strings.Contains(call.query, "event_checklist_item_completions") {
			t.Fatalf("seed must not touch completion rows, got:\n%s", call.query)
		}
	}
}

func TestRunReturnsExecutionError(t *testing.T) {
	executor := &fakeExecutor{err: sql.ErrConnDone}

	err := run(context.Background(), executor, defaultPlan)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "seed template") {
		t.Fatalf("expected template context in error, got %v", err)
	}
}

func assertSection(t *testing.T, section Section, slug string, title string, itemTitles []string) {
	t.Helper()

	if section.Slug != slug {
		t.Fatalf("expected section slug %q, got %q", slug, section.Slug)
	}
	if section.Title != title {
		t.Fatalf("expected section title %q, got %q", title, section.Title)
	}
	if len(section.Items) != len(itemTitles) {
		t.Fatalf("expected %d items in %s, got %d", len(itemTitles), section.Title, len(section.Items))
	}
	for i, want := range itemTitles {
		item := section.Items[i]
		if item.Title != want {
			t.Fatalf("expected item %d title %q, got %q", i+1, want, item.Title)
		}
		if item.DisplayOrder != i+1 {
			t.Fatalf("expected item %q display order %d, got %d", item.Title, i+1, item.DisplayOrder)
		}
		if item.Slug == "" {
			t.Fatalf("expected item %q to have stable slug", item.Title)
		}
	}
}

func assertCall(t *testing.T, call execCall, queryPart string, args ...any) {
	t.Helper()

	if !strings.Contains(call.query, queryPart) {
		t.Fatalf("expected query to contain %q, got:\n%s", queryPart, call.query)
	}
	if len(call.args) < len(args) {
		t.Fatalf("expected at least %d args, got %d", len(args), len(call.args))
	}
	for i, want := range args {
		if call.args[i] != want {
			t.Fatalf("expected arg %d to be %#v, got %#v", i, want, call.args[i])
		}
	}
}

type execCall struct {
	query string
	args  []any
}

type fakeExecutor struct {
	calls  []execCall
	nextID int64
	err    error
}

func (executor *fakeExecutor) ExecContext(_ context.Context, query string, args ...any) (sql.Result, error) {
	if executor.err != nil {
		return nil, executor.err
	}

	executor.nextID++
	callArgs := append([]any(nil), args...)
	executor.calls = append(executor.calls, execCall{query: query, args: callArgs})

	return fakeResult{id: executor.nextID}, nil
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
