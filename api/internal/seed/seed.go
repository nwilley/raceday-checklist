package seed

import (
	"context"
	"database/sql"
	"fmt"
)

const (
	defaultTemplateSlug = "default-raceday-checklist"
	defaultTemplateName = "Default Race Day Checklist"
	defaultEventSlug    = "default-raceday"
	defaultEventName    = "Default Race Day"
	defaultEventDate    = "2026-01-01"
)

type executor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type Template struct {
	Slug        string
	Name        string
	Description string
	Sections    []Section
}

type Section struct {
	Slug         string
	Title        string
	DisplayOrder int
	Items        []Item
}

type Item struct {
	Slug         string
	Title        string
	DisplayOrder int
}

type Event struct {
	Slug      string
	Name      string
	EventDate string
}

type Plan struct {
	Template Template
	Event    Event
}

var defaultPlan = Plan{
	Template: Template{
		Slug:        defaultTemplateSlug,
		Name:        defaultTemplateName,
		Description: "Default race-day maintenance checklist.",
		Sections: []Section{
			{
				Slug:         "pre-practice",
				Title:        "Pre-practice",
				DisplayOrder: 1,
				Items: []Item{
					{Slug: "set-droop", Title: "Set Droop", DisplayOrder: 1},
					{Slug: "set-camber", Title: "Set Camber", DisplayOrder: 2},
					{Slug: "set-ride-height", Title: "Set Ride Height", DisplayOrder: 3},
				},
			},
			{
				Slug:         "pre-qualifying",
				Title:        "Pre-qualifying",
				DisplayOrder: 2,
				Items: []Item{
					{Slug: "swap-diffs", Title: "Swap Diffs", DisplayOrder: 1},
					{Slug: "bleed-shocks", Title: "Bleed Shocks", DisplayOrder: 2},
					{Slug: "swap-hinge-pins", Title: "Swap Hinge Pins", DisplayOrder: 3},
					{Slug: "bolt-check", Title: "Bolt Check", DisplayOrder: 4},
					{Slug: "clean-car", Title: "Clean Car", DisplayOrder: 5},
					{Slug: "set-droop", Title: "Set Droop", DisplayOrder: 6},
					{Slug: "set-camber", Title: "Set Camber", DisplayOrder: 7},
					{Slug: "set-ride-height", Title: "Set Ride Height", DisplayOrder: 8},
				},
			},
			{
				Slug:         "mid-qualifying-maintenance",
				Title:        "Mid-qualifying maintenance",
				DisplayOrder: 3,
				Items: []Item{
					{Slug: "rebuild-first-diff-set", Title: "Rebuild first diff set", DisplayOrder: 1},
					{Slug: "polish-first-hinge-pin-set", Title: "Polish first hinge pin set", DisplayOrder: 2},
					{Slug: "check-droop-before-q2", Title: "Check droop before Q2", DisplayOrder: 3},
					{Slug: "check-camber-before-q2", Title: "Check camber before Q2", DisplayOrder: 4},
					{Slug: "check-ride-height-before-q2", Title: "Check ride height before Q2", DisplayOrder: 5},
				},
			},
			{
				Slug:         "pre-main",
				Title:        "Pre-main",
				DisplayOrder: 4,
				Items: []Item{
					{Slug: "swap-in-fresh-diffs", Title: "Swap in fresh diffs", DisplayOrder: 1},
					{Slug: "swap-in-fresh-hinge-pins", Title: "Swap in fresh hinge pins", DisplayOrder: 2},
					{Slug: "bleed-shocks", Title: "Bleed Shocks", DisplayOrder: 3},
					{Slug: "bolt-check", Title: "Bolt Check", DisplayOrder: 4},
					{Slug: "clean-car", Title: "Clean Car", DisplayOrder: 5},
					{Slug: "set-droop", Title: "Set Droop", DisplayOrder: 6},
					{Slug: "set-camber", Title: "Set Camber", DisplayOrder: 7},
					{Slug: "set-ride-height", Title: "Set Ride Height", DisplayOrder: 8},
				},
			},
		},
	},
	Event: Event{
		Slug:      defaultEventSlug,
		Name:      defaultEventName,
		EventDate: defaultEventDate,
	},
}

func Run(ctx context.Context, db *sql.DB) error {
	return run(ctx, db, defaultPlan)
}

func run(ctx context.Context, executor executor, plan Plan) error {
	templateID, err := upsertTemplate(ctx, executor, plan.Template)
	if err != nil {
		return fmt.Errorf("seed template: %w", err)
	}

	for _, section := range plan.Template.Sections {
		sectionID, err := upsertSection(ctx, executor, templateID, section)
		if err != nil {
			return fmt.Errorf("seed section %s: %w", section.Slug, err)
		}

		for _, item := range section.Items {
			if _, err := upsertItem(ctx, executor, sectionID, item); err != nil {
				return fmt.Errorf("seed item %s/%s: %w", section.Slug, item.Slug, err)
			}
		}
	}

	if _, err := upsertEvent(ctx, executor, templateID, plan.Event); err != nil {
		return fmt.Errorf("seed event: %w", err)
	}

	return nil
}

func upsertTemplate(ctx context.Context, executor executor, template Template) (int64, error) {
	return execReturningID(ctx, executor, `
INSERT INTO checklist_templates (slug, name, description)
VALUES (?, ?, ?)
ON DUPLICATE KEY UPDATE
  id = LAST_INSERT_ID(id),
  name = VALUES(name),
  description = VALUES(description)`,
		template.Slug,
		template.Name,
		template.Description,
	)
}

func upsertSection(ctx context.Context, executor executor, templateID int64, section Section) (int64, error) {
	return execReturningID(ctx, executor, `
INSERT INTO checklist_sections (template_id, slug, title, display_order)
VALUES (?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
  id = LAST_INSERT_ID(id),
  title = VALUES(title),
  display_order = VALUES(display_order)`,
		templateID,
		section.Slug,
		section.Title,
		section.DisplayOrder,
	)
}

func upsertItem(ctx context.Context, executor executor, sectionID int64, item Item) (int64, error) {
	return execReturningID(ctx, executor, `
INSERT INTO checklist_items (section_id, slug, title, display_order, default_done)
VALUES (?, ?, ?, ?, 0)
ON DUPLICATE KEY UPDATE
  id = LAST_INSERT_ID(id),
  title = VALUES(title),
  display_order = VALUES(display_order),
  default_done = VALUES(default_done)`,
		sectionID,
		item.Slug,
		item.Title,
		item.DisplayOrder,
	)
}

func upsertEvent(ctx context.Context, executor executor, templateID int64, event Event) (int64, error) {
	return execReturningID(ctx, executor, `
INSERT INTO raceday_events (template_id, slug, name, event_date)
VALUES (?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
  id = LAST_INSERT_ID(id),
  template_id = VALUES(template_id),
  name = VALUES(name),
  event_date = VALUES(event_date)`,
		templateID,
		event.Slug,
		event.Name,
		event.EventDate,
	)
}

func execReturningID(ctx context.Context, executor executor, query string, args ...any) (int64, error) {
	result, err := executor.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}
