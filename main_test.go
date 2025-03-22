package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetItems(t *testing.T) {
	// Arrange
	items = []Item{{ID: 1, Name: "Test Item"}}
	_, err := http.NewRequest(http.MethodGet, "/items", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}
	rec := httptest.NewRecorder()

	// Act
	getItems(rec)

	// Assert
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status OK; got %v", rec.Code)
	}
	var got []Item
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("Could not decode response: %v", err)
	}
	if len(got) != 1 || got[0].Name != "Test Item" {
		t.Errorf("Unexpected response: %v", got)
	}
}

func TestAddItem(t *testing.T) {
	// Arrange
	items = []Item{}
	nextID = 1
	newItem := Item{Name: "New Item"}
	body, err := json.Marshal(newItem)
	if err != nil {
		t.Fatalf("Could not marshal item: %v", err)
	}
	req, err := http.NewRequest(http.MethodPost, "/items", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// Act
	addItem(rec, req)

	// Assert
	if rec.Code != http.StatusCreated {
		t.Errorf("Expected status Created; got %v", rec.Code)
	}
	var got Item
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("Could not decode response: %v", err)
	}
	if got.ID != 1 || got.Name != "New Item" {
		t.Errorf("Unexpected response: %v", got)
	}
	if len(items) != 1 || items[0].Name != "New Item" {
		t.Errorf("Item was not added correctly: %v", items)
	}
}

func TestHandleItems(t *testing.T) {
	// Test GET method
	t.Run("GET /items", func(t *testing.T) {
		items = []Item{{ID: 1, Name: "Test Item"}}
		req, err := http.NewRequest(http.MethodGet, "/items", nil)
		if err != nil {
			t.Fatalf("Could not create request: %v", err)
		}
		rec := httptest.NewRecorder()

		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				getItems(w)
			case http.MethodPost:
				addItem(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		}).ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status OK; got %v", rec.Code)
		}
	})

	// Test POST method
	t.Run("POST /items", func(t *testing.T) {
		items = []Item{}
		nextID = 1
		newItem := Item{Name: "New Item"}
		body, err := json.Marshal(newItem)
		if err != nil {
			t.Fatalf("Could not marshal item: %v", err)
		}
		req, err := http.NewRequest(http.MethodPost, "/items", bytes.NewReader(body))
		if err != nil {
			t.Fatalf("Could not create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				getItems(w)
			case http.MethodPost:
				addItem(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		}).ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Errorf("Expected status Created; got %v", rec.Code)
		}
	})
}
