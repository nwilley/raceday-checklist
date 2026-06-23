package checklist

import (
	"context"
	"errors"
)

var ErrUnavailable = errors.New("checklist unavailable")

type Item struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Category string `json:"category"`
	Done     bool   `json:"done"`
}

type Repository interface {
	ListDefaultEventItems(ctx context.Context) ([]Item, error)
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return Service{repository: repository}
}

func (service Service) ListItems(ctx context.Context) ([]Item, error) {
	items, err := service.repository.ListDefaultEventItems(ctx)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, ErrUnavailable
	}
	return items, nil
}
