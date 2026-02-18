# Event-Driven System Implementation Guide

## Overview
This document explains the event-driven architecture implementation using RabbitMQ in the Go CRUD server.

## Architecture

### Components

#### 1. Event Publisher (`EventPublisher` in events.go)
- Publishes events to RabbitMQ when CRUD operations occur
- Maintains a connection to RabbitMQ
- Declares a durable queue named "item_events"
- Publishes messages with persistent delivery mode

#### 2. Event Consumer (`EventConsumer` in events.go)
- Subscribes to events from RabbitMQ
- Processes events with a custom handler function
- Acknowledges or rejects messages based on processing success
- Can be scaled horizontally (multiple consumers)

#### 3. Event Types
- `item.created` - Published when a new item is added
- `item.updated` - Published when an item is updated
- `item.deleted` - Published when an item is deleted

### Event Flow

```
Client Request
    ↓
REST API Endpoint (POST/PUT/DELETE)
    ↓
CRUD Operation (update in-memory data)
    ↓
Publish Event to RabbitMQ
    ↓
Event stored in durable queue
    ↓
Consumer receives and processes event
    ↓
Event acknowledged/rejected
```

## Implementation Details

### Main Server Changes

The main server (`main.go`) was updated to:
1. Initialize an `EventPublisher` on startup
2. Publish events after successful CRUD operations
3. Continue functioning even if RabbitMQ is unavailable

#### Example: Add Item with Event Publishing
```go
func addItem(w http.ResponseWriter, r *http.Request) {
    // ... existing code ...
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
    // ... rest of code ...
}
```

### Event Structure

```go
type ItemEvent struct {
    Type      EventType `json:"type"`      // Event type
    Item      Item      `json:"item"`      // Item data
    Timestamp time.Time `json:"timestamp"` // When event occurred
}
```

### Configuration

The system uses environment variables for configuration:

- `RABBITMQ_URL`: RabbitMQ connection string (default: `amqp://guest:guest@localhost:5672/`)

## Usage Examples

### Starting the Server

```bash
# Default configuration
go run .

# Custom RabbitMQ
RABBITMQ_URL="amqp://user:pass@rabbitmq-host:5672/" go run .
```

### Starting a Consumer

```bash
cd examples/consumer
go run .
```

### Consumer Output Example

```
Event consumer started. Waiting for events... Press Ctrl+C to exit.
[CREATED] Item ID: 1, Name: Sample Item at 2026-02-18 18:23:45
[UPDATED] Item ID: 1, Name: Updated Item at 2026-02-18 18:24:10
[DELETED] Item ID: 1, Name: Updated Item at 2026-02-18 18:24:30
```

## Benefits

### 1. Decoupling
- Services don't need direct connections
- Changes to consumers don't affect the API
- New consumers can be added without modifying the server

### 2. Scalability
- Multiple consumers can process events in parallel
- Load can be distributed across consumer instances
- Event processing doesn't block API responses

### 3. Reliability
- Durable queues survive RabbitMQ restarts
- Messages are persisted to disk
- Failed messages can be requeued
- Server continues working if RabbitMQ is down

### 4. Flexibility
- Add new event types easily
- Implement different processing logic per consumer
- Filter events based on type or other criteria

## Testing

The implementation includes comprehensive unit tests:

### Test Coverage
- Event structure validation
- Event serialization/deserialization
- Event type definitions
- Publisher behavior
- Consumer behavior
- Error handling

### Running Tests
```bash
# Run all tests
go test -v

# Run with coverage
go test -v -cover

# Run specific test
go test -v -run TestEventPublisher
```

## Error Handling

### Server-Side
- If RabbitMQ is unavailable at startup, server logs a warning and continues
- If event publishing fails, error is logged but API request succeeds
- Events are fire-and-forget from the API perspective

### Consumer-Side
- If event unmarshaling fails, message is rejected
- If handler returns error, message is requeued
- If handler succeeds, message is acknowledged

## Best Practices

1. **Always close connections**: Use `defer publisher.Close()` and `defer consumer.Close()`

2. **Handle errors gracefully**: Don't let event failures break API functionality

3. **Use durable queues**: Ensure messages survive restarts

4. **Acknowledge messages properly**: 
   - ACK on success
   - NACK with requeue on transient errors
   - NACK without requeue on permanent errors

5. **Monitor queue depth**: High queue depth may indicate processing bottleneck

## Future Enhancements

Potential improvements for the event-driven system:

1. **Dead Letter Queue**: For messages that fail repeatedly
2. **Event Filtering**: Allow consumers to subscribe to specific event types
3. **Event Replay**: Store events for replay/debugging
4. **Metrics**: Track event publishing/consumption rates
5. **Multiple Exchanges**: Route different event types to different queues
6. **Schema Versioning**: Handle event schema evolution

## Troubleshooting

### Issue: Server fails to start
**Solution**: Check if RabbitMQ is running and accessible. Server will continue with warning if RabbitMQ is unavailable.

### Issue: Events not being consumed
**Solution**: 
- Verify consumer is running
- Check RabbitMQ connection
- Verify queue name matches ("item_events")

### Issue: Messages stuck in queue
**Solution**:
- Check consumer logs for errors
- Verify handler logic is working
- Check if messages are being rejected

### Issue: Duplicate event processing
**Solution**: Ensure consumers are acknowledging messages properly. Check for consumer crashes before ACK.

## Security Considerations

1. **Authentication**: Use secure credentials in `RABBITMQ_URL`
2. **TLS**: Use `amqps://` for encrypted connections
3. **Input Validation**: Validate event data in consumers
4. **Rate Limiting**: Protect consumers from event floods

## Conclusion

The event-driven architecture provides a scalable, reliable, and flexible way to handle CRUD operations. The implementation maintains backward compatibility while adding powerful event processing capabilities.
