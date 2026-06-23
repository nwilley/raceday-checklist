package checklist

import (
	"context"
	"errors"
	"strings"
)

var (
	ErrInvalidClientID = errors.New("invalid race-day client id")
	ErrNotFound        = errors.New("checklist item not found")
	ErrUnavailable     = errors.New("checklist unavailable")
)

type Item struct {
	ID        string `json:"id"`
	SectionID string `json:"sectionId"`
	ItemID    string `json:"itemId"`
	Title     string `json:"title"`
	Category  string `json:"category"`
	Done      bool   `json:"done"`
}

type CompletionUpdate struct {
	SectionID string
	ItemID    string
	Done      bool
}

type Repository interface {
	ListDefaultEventItems(ctx context.Context, clientID string) ([]Item, error)
	UpdateDefaultEventItemCompletion(ctx context.Context, clientID string, update CompletionUpdate) error
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return Service{repository: repository}
}

func (service Service) ListItems(ctx context.Context, clientID string) ([]Item, error) {
	clientID = strings.TrimSpace(clientID)
	if clientID == "" {
		return nil, ErrInvalidClientID
	}

	items, err := service.repository.ListDefaultEventItems(ctx, clientID)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, ErrUnavailable
	}
	return items, nil
}

func (service Service) UpdateItemCompletion(ctx context.Context, clientID string, update CompletionUpdate) error {
	clientID = strings.TrimSpace(clientID)
	if clientID == "" {
		return ErrInvalidClientID
	}
	if strings.TrimSpace(update.SectionID) == "" || strings.TrimSpace(update.ItemID) == "" {
		return ErrNotFound
	}

	update.SectionID = strings.TrimSpace(update.SectionID)
	update.ItemID = strings.TrimSpace(update.ItemID)

	return service.repository.UpdateDefaultEventItemCompletion(ctx, clientID, update)
}

func IsInvalidClientID(err error) bool {
	return errors.Is(err, ErrInvalidClientID)
}

func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}
