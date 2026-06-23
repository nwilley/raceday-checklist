package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"raceday-checklist/api/internal/checklist"
)

func TestHealth(t *testing.T) {
	router := NewRouter(&fakeChecklistService{})
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
	service := &fakeChecklistService{
		items: []checklist.Item{
			{ID: "set-droop", SectionID: "pre-practice", ItemID: "set-droop", Title: "Set Droop", Category: "Pre-practice", Done: true},
		},
	}
	router := NewRouter(service)
	request := httptest.NewRequest(http.MethodGet, "/api/checklist", nil)
	request.Header.Set(raceDayClientHeader, "client-1")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}
	if service.listClientID != "client-1" {
		t.Fatalf("expected client id client-1, got %q", service.listClientID)
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
	if body.Items[0].SectionID != "pre-practice" {
		t.Fatalf("expected section id pre-practice, got %q", body.Items[0].SectionID)
	}
	if body.Items[0].ItemID != "set-droop" {
		t.Fatalf("expected item id set-droop, got %q", body.Items[0].ItemID)
	}
}

func TestChecklistMissingClientID(t *testing.T) {
	router := NewRouter(&fakeChecklistService{listErr: checklist.ErrInvalidClientID})
	request := httptest.NewRequest(http.MethodGet, "/api/checklist", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", response.Code)
	}
}

func TestChecklistUnavailable(t *testing.T) {
	router := NewRouter(&fakeChecklistService{listErr: checklist.ErrUnavailable})
	request := httptest.NewRequest(http.MethodGet, "/api/checklist", nil)
	request.Header.Set(raceDayClientHeader, "client-1")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", response.Code)
	}
}

func TestChecklistUnexpectedError(t *testing.T) {
	router := NewRouter(&fakeChecklistService{listErr: errors.New("database host unreachable")})
	request := httptest.NewRequest(http.MethodGet, "/api/checklist", nil)
	request.Header.Set(raceDayClientHeader, "client-1")
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

func TestChecklistCompletionUpdate(t *testing.T) {
	service := &fakeChecklistService{}
	router := NewRouter(service)
	request := httptest.NewRequest(http.MethodPatch, "/api/checklist/items/pre-practice/set-droop", strings.NewReader(`{"done":true}`))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set(raceDayClientHeader, "client-1")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}
	if service.updateClientID != "client-1" {
		t.Fatalf("expected client id client-1, got %q", service.updateClientID)
	}
	if service.update.SectionID != "pre-practice" || service.update.ItemID != "set-droop" || !service.update.Done {
		t.Fatalf("unexpected update: %#v", service.update)
	}
}

func TestChecklistCompletionUpdateNotFound(t *testing.T) {
	router := NewRouter(&fakeChecklistService{updateErr: checklist.ErrNotFound})
	request := httptest.NewRequest(http.MethodPatch, "/api/checklist/items/bad/missing", strings.NewReader(`{"done":true}`))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set(raceDayClientHeader, "client-1")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", response.Code)
	}
}

func TestChecklistCompletionUpdateInvalidJSON(t *testing.T) {
	router := NewRouter(&fakeChecklistService{})
	request := httptest.NewRequest(http.MethodPatch, "/api/checklist/items/pre-practice/set-droop", strings.NewReader(`{`))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set(raceDayClientHeader, "client-1")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", response.Code)
	}
}

type fakeChecklistService struct {
	items          []checklist.Item
	listErr        error
	updateErr      error
	listClientID   string
	updateClientID string
	update         checklist.CompletionUpdate
}

func (service *fakeChecklistService) ListItems(_ context.Context, clientID string) ([]checklist.Item, error) {
	service.listClientID = clientID
	return service.items, service.listErr
}

func (service *fakeChecklistService) UpdateItemCompletion(_ context.Context, clientID string, update checklist.CompletionUpdate) error {
	service.updateClientID = clientID
	service.update = update
	return service.updateErr
}
