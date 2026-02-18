package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// Example consumer application that listens to item events
func main() {
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		rabbitMQURL = "amqp://guest:guest@localhost:5672/"
	}

	consumer, err := NewEventConsumer(rabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to create event consumer: %v", err)
	}
	defer consumer.Close()

	// Define event handler
	handler := func(event ItemEvent) error {
		switch event.Type {
		case EventItemCreated:
			fmt.Printf("[CREATED] Item ID: %d, Name: %s at %s\n",
				event.Item.ID, event.Item.Name, event.Timestamp.Format("2006-01-02 15:04:05"))
		case EventItemUpdated:
			fmt.Printf("[UPDATED] Item ID: %d, Name: %s at %s\n",
				event.Item.ID, event.Item.Name, event.Timestamp.Format("2006-01-02 15:04:05"))
		case EventItemDeleted:
			fmt.Printf("[DELETED] Item ID: %d, Name: %s at %s\n",
				event.Item.ID, event.Item.Name, event.Timestamp.Format("2006-01-02 15:04:05"))
		default:
			fmt.Printf("[UNKNOWN] Event type: %s\n", event.Type)
		}
		return nil
	}

	// Start consuming
	if err := consumer.Consume(handler); err != nil {
		log.Fatalf("Failed to start consuming: %v", err)
	}

	fmt.Println("Event consumer started. Waiting for events... Press Ctrl+C to exit.")

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\nShutting down consumer...")
}
