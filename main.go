package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type Item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var (
	items          []Item
	mutex          sync.Mutex
	nextID         int = 1
	eventPublisher *EventPublisher
)

func getItems(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	mutex.Lock()
	defer mutex.Unlock()
	json.NewEncoder(w).Encode(items)
}

func addItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newItem Item
	if err := json.NewDecoder(r.Body).Decode(&newItem); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	mutex.Lock()
	defer mutex.Unlock()
	newItem.ID = nextID
	nextID++
	items = append(items, newItem)
	
	// Publish event
	if eventPublisher != nil {
		event := ItemEvent{
			Type:      EventItemCreated,
			Item:      newItem,
			Timestamp: time.Now(),
		}
		if err := eventPublisher.Publish(event); err != nil {
			log.Printf("Failed to publish event: %v", err)
		}
	}
	
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newItem)
}

func updateItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var updatedItem Item
	if err := json.NewDecoder(r.Body).Decode(&updatedItem); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	mutex.Lock()
	defer mutex.Unlock()
	for i, item := range items {
		if item.ID == updatedItem.ID {
			items[i] = updatedItem
			
			// Publish event
			if eventPublisher != nil {
				event := ItemEvent{
					Type:      EventItemUpdated,
					Item:      updatedItem,
					Timestamp: time.Now(),
				}
				if err := eventPublisher.Publish(event); err != nil {
					log.Printf("Failed to publish event: %v", err)
				}
			}
			
			json.NewEncoder(w).Encode(updatedItem)
			return
		}
	}
	http.Error(w, "Item not found", http.StatusNotFound)
}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var itemToDelete Item
	if err := json.NewDecoder(r.Body).Decode(&itemToDelete); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	mutex.Lock()
	defer mutex.Unlock()
	for i, item := range items {
		if item.ID == itemToDelete.ID {
			items = append(items[:i], items[i+1:]...)
			
			// Publish event
			if eventPublisher != nil {
				event := ItemEvent{
					Type:      EventItemDeleted,
					Item:      item,
					Timestamp: time.Now(),
				}
				if err := eventPublisher.Publish(event); err != nil {
					log.Printf("Failed to publish event: %v", err)
				}
			}
			
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	http.Error(w, "Item not found", http.StatusNotFound)
}

func main() {
	// Initialize event publisher if RabbitMQ URL is provided
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		rabbitMQURL = "amqp://guest:guest@localhost:5672/" // default
	}
	
	var err error
	eventPublisher, err = NewEventPublisher(rabbitMQURL)
	if err != nil {
		log.Printf("Warning: Failed to initialize event publisher: %v", err)
		log.Println("Server will continue without event publishing")
	} else {
		defer eventPublisher.Close()
		log.Println("Event publisher initialized successfully")
	}

	http.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getItems(w)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/items/add", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			addItem(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/items/update", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			updateItem(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/items/delete", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			deleteItem(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
