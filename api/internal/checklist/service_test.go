package checklist

import (
	"context"
	"errors"
	"testing"
)

func TestServiceListItems(t *testing.T) {
	expected := []Item{
		{ID: "fuel", Title: "Fuel checked", Category: "Car", Done: true},
	}
	service := NewService(fakeRepository{items: expected})

	items, err := service.ListItems(context.Background())
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
	service := NewService(fakeRepository{})

	_, err := service.ListItems(context.Background())
	if !errors.Is(err, ErrUnavailable) {
		t.Fatalf("expected ErrUnavailable, got %v", err)
	}
}

func TestServiceListItemsPropagatesRepositoryError(t *testing.T) {
	expected := errors.New("database failed")
	service := NewService(fakeRepository{err: expected})

	_, err := service.ListItems(context.Background())
	if !errors.Is(err, expected) {
		t.Fatalf("expected repository error, got %v", err)
	}
}

type fakeRepository struct {
	items []Item
	err   error
}

func (repository fakeRepository) ListDefaultEventItems(context.Context) ([]Item, error) {
	return repository.items, repository.err
}
