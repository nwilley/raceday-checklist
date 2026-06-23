package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"raceday-checklist/api/internal/checklist"
)

func TestHealth(t *testing.T) {
	router := NewRouter(fakeChecklistService{})
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}

	var body map[string]string
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body["status"] != "ok" {
		t.Fatalf("expected status ok, got %q", body["status"])
	}
}

func TestChecklist(t *testing.T) {
	router := NewRouter(fakeChecklistService{
		items: []checklist.Item{
			{ID: "fuel", Title: "Fuel and fluids checked", Category: "Car", Done: true},
		},
	})
	request := httptest.NewRequest(http.MethodGet, "/api/checklist", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}

	var body struct {
		Items []checklist.Item `json:"items"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(body.Items) != 1 {
		t.Fatalf("expected 1 checklist item, got %d", len(body.Items))
	}
	if body.Items[0].ID != "fuel" {
		t.Fatalf("expected item id fuel, got %q", body.Items[0].ID)
	}
	if body.Items[0].Title != "Fuel and fluids checked" {
		t.Fatalf("expected item title, got %q", body.Items[0].Title)
	}
	if body.Items[0].Category != "Car" {
		t.Fatalf("expected item category Car, got %q", body.Items[0].Category)
	}
	if !body.Items[0].Done {
		t.Fatal("expected item done true")
	}
}

func TestChecklistUnavailable(t *testing.T) {
	router := NewRouter(fakeChecklistService{err: checklist.ErrUnavailable})
	request := httptest.NewRequest(http.MethodGet, "/api/checklist", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", response.Code)
	}

	var body map[string]string
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body["error"] != "checklist unavailable" {
		t.Fatalf("expected checklist unavailable error, got %q", body["error"])
	}
}

func TestChecklistUnexpectedError(t *testing.T) {
	router := NewRouter(fakeChecklistService{err: errors.New("database host unreachable")})
	request := httptest.NewRequest(http.MethodGet, "/api/checklist", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", response.Code)
	}

	var body map[string]string
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body["error"] != "internal server error" {
		t.Fatalf("expected generic internal server error, got %q", body["error"])
	}
}

type fakeChecklistService struct {
	items []checklist.Item
	err   error
}

func (service fakeChecklistService) ListItems(context.Context) ([]checklist.Item, error) {
	return service.items, service.err
}
