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
	req, err := http.NewRequest(http.MethodPost, "/items/add", bytes.NewReader(body))
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
func TestUpdateItem(t *testing.T) {
	// Arrange
	items = []Item{{ID: 1, Name: "Old Item"}}
	updatedItem := Item{ID: 1, Name: "Updated Item"}
	body, err := json.Marshal(updatedItem)
	if err != nil {
		t.Fatalf("Could not marshal item: %v", err)
	}
	req, err := http.NewRequest(http.MethodPut, "/items/update", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// Act
	updateItem(rec, req)

	// Assert
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status OK; got %v", rec.Code)
	}
	var got Item
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("Could not decode response: %v", err)
	}
	if got.ID != 1 || got.Name != "Updated Item" {
		t.Errorf("Unexpected response: %v", got)
	}
	if len(items) != 1 || items[0].Name != "Updated Item" {
		t.Errorf("Item was not updated correctly: %v", items)
	}
}

func TestDeleteItem(t *testing.T) {
	// Arrange
	items = []Item{{ID: 1, Name: "Item to Delete"}}
	itemToDelete := Item{ID: 1}
	body, err := json.Marshal(itemToDelete)
	if err != nil {
		t.Fatalf("Could not marshal item: %v", err)
	}
	req, err := http.NewRequest(http.MethodDelete, "/items/delete", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// Act
	deleteItem(rec, req)

	// Assert
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status OK; got %v", rec.Code)
	}
	var got Item
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("Could not decode response: %v", err)
	}
	if got.ID != 1 || got.Name != "Item to Delete" {
		t.Errorf("Unexpected response: %v", got)
	}
	if len(items) != 0 {
		t.Errorf("Item was not deleted correctly: %v", items)
	}
}
