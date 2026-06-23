package checklist

import (
	"context"
	"errors"
	"testing"
)

func TestServiceListItems(t *testing.T) {
	expected := []Item{
		{ID: "set-droop", SectionID: "pre-practice", ItemID: "set-droop", Title: "Set Droop", Category: "Pre-practice", Done: true},
	}
	service := NewService(&fakeRepository{items: expected})

	items, err := service.ListItems(context.Background(), "client-1")
	if err != nil {
		t.Fatalf("list items: %v", err)
	}

	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	if items[0] != expected[0] {
		t.Fatalf("expected item %#v, got %#v", expected[0], items[0])
	}
}

func TestServiceListItemsUnavailableWhenRepositoryReturnsNoItems(t *testing.T) {
	service := NewService(&fakeRepository{})

	_, err := service.ListItems(context.Background(), "client-1")
	if !errors.Is(err, ErrUnavailable) {
		t.Fatalf("expected ErrUnavailable, got %v", err)
	}
}

func TestServiceListItemsRequiresClientID(t *testing.T) {
	service := NewService(&fakeRepository{})

	_, err := service.ListItems(context.Background(), "")
	if !errors.Is(err, ErrInvalidClientID) {
		t.Fatalf("expected ErrInvalidClientID, got %v", err)
	}
}

func TestServiceListItemsPropagatesRepositoryError(t *testing.T) {
	expected := errors.New("database failed")
	service := NewService(&fakeRepository{err: expected})

	_, err := service.ListItems(context.Background(), "client-1")
	if !errors.Is(err, expected) {
		t.Fatalf("expected repository error, got %v", err)
	}
}

func TestServiceUpdateItemCompletion(t *testing.T) {
	repository := &fakeRepository{}
	service := NewService(repository)

	err := service.UpdateItemCompletion(context.Background(), "client-1", CompletionUpdate{
		SectionID: "pre-practice",
		ItemID:    "set-droop",
		Done:      true,
	})
	if err != nil {
		t.Fatalf("update item completion: %v", err)
	}

	if repository.clientID != "client-1" {
		t.Fatalf("expected client id client-1, got %q", repository.clientID)
	}
	if repository.update.SectionID != "pre-practice" || repository.update.ItemID != "set-droop" || !repository.update.Done {
		t.Fatalf("unexpected update: %#v", repository.update)
	}
}

func TestServiceUpdateItemCompletionRequiresClientID(t *testing.T) {
	service := NewService(&fakeRepository{})

	err := service.UpdateItemCompletion(context.Background(), "", CompletionUpdate{SectionID: "pre-practice", ItemID: "set-droop"})
	if !errors.Is(err, ErrInvalidClientID) {
		t.Fatalf("expected ErrInvalidClientID, got %v", err)
	}
}

type fakeRepository struct {
	items    []Item
	err      error
	clientID string
	update   CompletionUpdate
}

func (repository fakeRepository) ListDefaultEventItems(_ context.Context, _ string) ([]Item, error) {
	return repository.items, repository.err
}

func (repository *fakeRepository) UpdateDefaultEventItemCompletion(_ context.Context, clientID string, update CompletionUpdate) error {
	repository.clientID = clientID
	repository.update = update
	return repository.err
}
