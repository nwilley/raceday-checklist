package database

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestChecklistStorageUpMigration(t *testing.T) {
	sql := readMigration(t, "001_create_checklist_storage.up.sql")

	for _, want := range []string{
		"CREATE TABLE checklist_templates",
		"CREATE TABLE checklist_sections",
		"CREATE TABLE checklist_items",
		"CREATE TABLE raceday_events",
		"CREATE TABLE event_checklist_item_completions",
		"UNIQUE KEY uq_checklist_sections_template_slug (template_id, slug)",
		"UNIQUE KEY uq_checklist_items_section_slug (section_id, slug)",
		"UNIQUE KEY uq_checklist_sections_template_display_order (template_id, display_order)",
		"UNIQUE KEY uq_checklist_items_section_display_order (section_id, display_order)",
		"CONSTRAINT fk_checklist_sections_template",
		"CONSTRAINT fk_checklist_items_section",
		"CONSTRAINT fk_raceday_events_template",
		"CONSTRAINT fk_event_checklist_item_completions_event",
		"CONSTRAINT fk_event_checklist_item_completions_item",
	} {
		if !strings.Contains(sql, want) {
			t.Fatalf("expected up migration to contain %q", want)
		}
	}
}

func TestChecklistStorageDownMigrationOrder(t *testing.T) {
	sql := readMigration(t, "001_create_checklist_storage.down.sql")

	assertBefore(t, sql, "DROP TABLE IF EXISTS event_checklist_item_completions", "DROP TABLE IF EXISTS raceday_events")
	assertBefore(t, sql, "DROP TABLE IF EXISTS raceday_events", "DROP TABLE IF EXISTS checklist_items")
	assertBefore(t, sql, "DROP TABLE IF EXISTS checklist_items", "DROP TABLE IF EXISTS checklist_sections")
	assertBefore(t, sql, "DROP TABLE IF EXISTS checklist_sections", "DROP TABLE IF EXISTS checklist_templates")
}

func readMigration(t *testing.T, name string) string {
	t.Helper()

	path := filepath.Join("..", "..", "migrations", name)
	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read migration %s: %v", name, err)
	}

	return string(contents)
}

func assertBefore(t *testing.T, text string, first string, second string) {
	t.Helper()

	firstIndex := strings.Index(text, first)
	if firstIndex == -1 {
		t.Fatalf("expected text to contain %q", first)
	}

	secondIndex := strings.Index(text, second)
	if secondIndex == -1 {
		t.Fatalf("expected text to contain %q", second)
	}

	if firstIndex > secondIndex {
		t.Fatalf("expected %q before %q", first, second)
	}
}
