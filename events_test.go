package main

import (
	"encoding/json"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// MockAMQPChannel is a mock implementation of AMQP channel for testing
type MockAMQPChannel struct {
	publishedMessages []amqp.Publishing
	queueName         string
	consumeHandler    func(amqp.Delivery)
	shouldFail        bool
}

func (m *MockAMQPChannel) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	if m.shouldFail {
		return amqp.ErrClosed
	}
	m.publishedMessages = append(m.publishedMessages, msg)
	return nil
}

func (m *MockAMQPChannel) QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error) {
	if m.shouldFail {
		return amqp.Queue{}, amqp.ErrClosed
	}
	m.queueName = name
	return amqp.Queue{Name: name}, nil
}

func (m *MockAMQPChannel) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	if m.shouldFail {
		return nil, amqp.ErrClosed
	}
	deliveryChan := make(chan amqp.Delivery, 10)
	return deliveryChan, nil
}

func (m *MockAMQPChannel) Close() error {
	return nil
}

// TestEventPublisher tests the event publisher
func TestEventPublisher(t *testing.T) {
	// This test validates the event structure
	t.Run("ValidateEventStructure", func(t *testing.T) {
		event := ItemEvent{
			Type:      EventItemCreated,
			Item:      Item{ID: 1, Name: "Test Item"},
			Timestamp: time.Now(),
		}

		// Validate event can be marshaled
		data, err := json.Marshal(event)
		if err != nil {
			t.Fatalf("Failed to marshal event: %v", err)
		}

		// Validate event can be unmarshaled
		var unmarshaled ItemEvent
		if err := json.Unmarshal(data, &unmarshaled); err != nil {
			t.Fatalf("Failed to unmarshal event: %v", err)
		}

		if unmarshaled.Type != EventItemCreated {
			t.Errorf("Expected event type %s, got %s", EventItemCreated, unmarshaled.Type)
		}
		if unmarshaled.Item.ID != 1 || unmarshaled.Item.Name != "Test Item" {
			t.Errorf("Item data mismatch: %+v", unmarshaled.Item)
		}
	})

	t.Run("EventTypesAreDefined", func(t *testing.T) {
		if EventItemCreated == "" {
			t.Error("EventItemCreated is not defined")
		}
		if EventItemUpdated == "" {
			t.Error("EventItemUpdated is not defined")
		}
		if EventItemDeleted == "" {
			t.Error("EventItemDeleted is not defined")
		}

		// Validate they are different
		if EventItemCreated == EventItemUpdated || EventItemCreated == EventItemDeleted || EventItemUpdated == EventItemDeleted {
			t.Error("Event types should be unique")
		}
	})
}

// TestEventPublisherPublish tests the Publish method
func TestEventPublisherPublish(t *testing.T) {
	t.Run("PublishWithNilChannel", func(t *testing.T) {
		publisher := &EventPublisher{
			channel: nil,
		}

		event := ItemEvent{
			Type:      EventItemCreated,
			Item:      Item{ID: 1, Name: "Test"},
			Timestamp: time.Now(),
		}

		err := publisher.Publish(event)
		if err == nil {
			t.Error("Expected error when publishing with nil channel")
		}
	})
}

// TestEventConsumer tests the event consumer
func TestEventConsumer(t *testing.T) {
	t.Run("ConsumeWithNilChannel", func(t *testing.T) {
		consumer := &EventConsumer{
			channel: nil,
		}

		err := consumer.Consume(func(event ItemEvent) error {
			return nil
		})

		if err == nil {
			t.Error("Expected error when consuming with nil channel")
		}
	})
}

// TestItemEvent tests the ItemEvent structure
func TestItemEvent(t *testing.T) {
	t.Run("CreateItemEvent", func(t *testing.T) {
		item := Item{ID: 1, Name: "Test Item"}
		timestamp := time.Now()

		event := ItemEvent{
			Type:      EventItemCreated,
			Item:      item,
			Timestamp: timestamp,
		}

		if event.Type != EventItemCreated {
			t.Errorf("Expected type %s, got %s", EventItemCreated, event.Type)
		}
		if event.Item.ID != 1 {
			t.Errorf("Expected item ID 1, got %d", event.Item.ID)
		}
		if event.Item.Name != "Test Item" {
			t.Errorf("Expected item name 'Test Item', got %s", event.Item.Name)
		}
	})

	t.Run("UpdateItemEvent", func(t *testing.T) {
		item := Item{ID: 2, Name: "Updated Item"}
		event := ItemEvent{
			Type:      EventItemUpdated,
			Item:      item,
			Timestamp: time.Now(),
		}

		if event.Type != EventItemUpdated {
			t.Errorf("Expected type %s, got %s", EventItemUpdated, event.Type)
		}
		if event.Item.ID != 2 {
			t.Errorf("Expected item ID 2, got %d", event.Item.ID)
		}
	})

	t.Run("DeleteItemEvent", func(t *testing.T) {
		item := Item{ID: 3, Name: "Deleted Item"}
		event := ItemEvent{
			Type:      EventItemDeleted,
			Item:      item,
			Timestamp: time.Now(),
		}

		if event.Type != EventItemDeleted {
			t.Errorf("Expected type %s, got %s", EventItemDeleted, event.Type)
		}
		if event.Item.ID != 3 {
			t.Errorf("Expected item ID 3, got %d", event.Item.ID)
		}
	})
}

// TestEventSerialization tests JSON serialization of events
func TestEventSerialization(t *testing.T) {
	t.Run("SerializeAndDeserialize", func(t *testing.T) {
		originalEvent := ItemEvent{
			Type: EventItemCreated,
			Item: Item{
				ID:   123,
				Name: "Serialization Test",
			},
			Timestamp: time.Now().UTC().Truncate(time.Second),
		}

		// Serialize
		data, err := json.Marshal(originalEvent)
		if err != nil {
			t.Fatalf("Failed to serialize event: %v", err)
		}

		// Deserialize
		var deserializedEvent ItemEvent
		if err := json.Unmarshal(data, &deserializedEvent); err != nil {
			t.Fatalf("Failed to deserialize event: %v", err)
		}

		// Verify
		if deserializedEvent.Type != originalEvent.Type {
			t.Errorf("Type mismatch: expected %s, got %s", originalEvent.Type, deserializedEvent.Type)
		}
		if deserializedEvent.Item.ID != originalEvent.Item.ID {
			t.Errorf("Item ID mismatch: expected %d, got %d", originalEvent.Item.ID, deserializedEvent.Item.ID)
		}
		if deserializedEvent.Item.Name != originalEvent.Item.Name {
			t.Errorf("Item name mismatch: expected %s, got %s", originalEvent.Item.Name, deserializedEvent.Item.Name)
		}
	})

	t.Run("SerializeMultipleEvents", func(t *testing.T) {
		events := []ItemEvent{
			{Type: EventItemCreated, Item: Item{ID: 1, Name: "Event 1"}, Timestamp: time.Now()},
			{Type: EventItemUpdated, Item: Item{ID: 2, Name: "Event 2"}, Timestamp: time.Now()},
			{Type: EventItemDeleted, Item: Item{ID: 3, Name: "Event 3"}, Timestamp: time.Now()},
		}

		for i, event := range events {
			data, err := json.Marshal(event)
			if err != nil {
				t.Errorf("Failed to serialize event %d: %v", i, err)
			}

			var unmarshaled ItemEvent
			if err := json.Unmarshal(data, &unmarshaled); err != nil {
				t.Errorf("Failed to deserialize event %d: %v", i, err)
			}

			if unmarshaled.Type != event.Type {
				t.Errorf("Event %d type mismatch", i)
			}
		}
	})
}

// TestEventPublisherClose tests closing the publisher
func TestEventPublisherClose(t *testing.T) {
	t.Run("CloseWithNilChannelAndConnection", func(t *testing.T) {
		publisher := &EventPublisher{}
		err := publisher.Close()
		if err != nil {
			t.Errorf("Expected no error when closing nil publisher, got %v", err)
		}
	})
}

// TestEventConsumerClose tests closing the consumer
func TestEventConsumerClose(t *testing.T) {
	t.Run("CloseWithNilChannelAndConnection", func(t *testing.T) {
		consumer := &EventConsumer{}
		err := consumer.Close()
		if err != nil {
			t.Errorf("Expected no error when closing nil consumer, got %v", err)
		}
	})
}
