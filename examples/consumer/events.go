package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// EventType represents the type of event
type EventType string

// Item represents an item in the CRUD system
type Item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

const (
	EventItemCreated EventType = "item.created"
	EventItemUpdated EventType = "item.updated"
	EventItemDeleted EventType = "item.deleted"
)

// ItemEvent represents an event related to an item
type ItemEvent struct {
	Type      EventType `json:"type"`
	Item      Item      `json:"item"`
	Timestamp time.Time `json:"timestamp"`
}

// EventPublisher handles publishing events to RabbitMQ
type EventPublisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

// NewEventPublisher creates a new event publisher
func NewEventPublisher(amqpURL string) (*EventPublisher, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	q, err := ch.QueueDeclare(
		"item_events", // name
		true,          // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	return &EventPublisher{
		conn:    conn,
		channel: ch,
		queue:   q,
	}, nil
}

// Publish publishes an event to RabbitMQ
func (ep *EventPublisher) Publish(event ItemEvent) error {
	if ep.channel == nil {
		return fmt.Errorf("channel is not initialized")
	}

	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	err = ep.channel.Publish(
		"",           // exchange
		ep.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	log.Printf("Published event: %s for item ID: %d", event.Type, event.Item.ID)
	return nil
}

// Close closes the connection and channel
func (ep *EventPublisher) Close() error {
	if ep.channel != nil {
		if err := ep.channel.Close(); err != nil {
			return err
		}
	}
	if ep.conn != nil {
		if err := ep.conn.Close(); err != nil {
			return err
		}
	}
	return nil
}

// EventConsumer handles consuming events from RabbitMQ
type EventConsumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

// NewEventConsumer creates a new event consumer
func NewEventConsumer(amqpURL string) (*EventConsumer, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	q, err := ch.QueueDeclare(
		"item_events", // name
		true,          // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	return &EventConsumer{
		conn:    conn,
		channel: ch,
		queue:   q,
	}, nil
}

// Consume starts consuming events from RabbitMQ
func (ec *EventConsumer) Consume(handler func(ItemEvent) error) error {
	if ec.channel == nil {
		return fmt.Errorf("channel is not initialized")
	}

	msgs, err := ec.channel.Consume(
		ec.queue.Name, // queue
		"",            // consumer
		false,         // auto-ack
		false,         // exclusive
		false,         // no-local
		false,         // no-wait
		nil,           // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	go func() {
		for d := range msgs {
			var event ItemEvent
			if err := json.Unmarshal(d.Body, &event); err != nil {
				log.Printf("Error unmarshaling event: %v", err)
				d.Nack(false, false) // reject message
				continue
			}

			if err := handler(event); err != nil {
				log.Printf("Error handling event: %v", err)
				d.Nack(false, true) // requeue message
			} else {
				d.Ack(false) // acknowledge message
				log.Printf("Processed event: %s for item ID: %d", event.Type, event.Item.ID)
			}
		}
	}()

	log.Printf("Consumer started, waiting for events...")
	return nil
}

// Close closes the connection and channel
func (ec *EventConsumer) Close() error {
	if ec.channel != nil {
		if err := ec.channel.Close(); err != nil {
			return err
		}
	}
	if ec.conn != nil {
		if err := ec.conn.Close(); err != nil {
			return err
		}
	}
	return nil
}
