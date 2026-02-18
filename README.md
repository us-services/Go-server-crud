# Go-server-crud
Go server for crud operations with event-driven architecture using RabbitMQ.

## Features
- RESTful CRUD API for managing items
- Event-driven architecture with RabbitMQ
- Automatic event publishing for all CRUD operations
- Example consumer application for event processing
- Comprehensive unit tests

## Prerequisites
1. Go 1.24 or higher
2. RabbitMQ server (optional - for event-driven features)

## Installation
```bash
git clone this repository
cd Go-server-crud
go mod tidy
```

## Run Go server
```bash
# Without RabbitMQ (events will be disabled with warning)
go run .

# With RabbitMQ (default: localhost:5672)
go run .

# With custom RabbitMQ URL
RABBITMQ_URL="amqp://user:pass@hostname:5672/" go run .
```

## Event-Driven System

### Overview
The server now publishes events to RabbitMQ whenever items are created, updated, or deleted. This enables:
- Real-time notifications
- Event-driven microservices
- Audit logging
- Asynchronous processing

### Event Types
- `item.created` - Published when a new item is added
- `item.updated` - Published when an item is updated
- `item.deleted` - Published when an item is deleted

### Event Structure
```json
{
  "type": "item.created",
  "item": {
    "id": 1,
    "name": "Sample Item"
  },
  "timestamp": "2026-02-18T18:23:45Z"
}
```

### Running the Event Consumer
The repository includes an example consumer that listens to events:

```bash
cd examples/consumer

# With default RabbitMQ (localhost:5672)
go run .

# With custom RabbitMQ URL
RABBITMQ_URL="amqp://user:pass@hostname:5672/" go run .
```

### Setting up RabbitMQ
If you don't have RabbitMQ installed, you can run it using Docker:

```bash
docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3-management
```

Access the management UI at http://localhost:15672 (user: guest, password: guest)

## API Operations
GET Operation
```
curl -X GET http://localhost:8080/items | jq .
```

POST Operation ( Add an item )
```
curl -X POST http://localhost:8080/items -H "Content-Type: application/json" -d '{"name": "Sample Item"}'
```

PUT Operation ( Update an item )
```
curl -X PUT http://localhost:8080/items/update -H "Content-Type: application/json" -d '{"id":1,"name": "Another new Sample Item"}'
```

DELETE Operation ( Delete an item )
```
curl -X DELETE http://localhost:8080/items/delete -H "Content-Type: application/json" -d '{"id":1,"name": "Another Sample Item"}'
```

## Running Tests
```bash
# Run all tests
go test -v

# Run tests with coverage
go test -v -cover

# Run specific test
go test -v -run TestEventPublisher
```

## Architecture

### Event Publishing Flow
1. Client makes CRUD request to REST API
2. Server processes the request and updates data
3. Server publishes event to RabbitMQ queue
4. Event is stored in durable queue
5. Consumers receive and process events asynchronously

### Benefits
- **Decoupling**: Services don't need direct connections
- **Scalability**: Multiple consumers can process events
- **Reliability**: Durable queues ensure no event loss
- **Flexibility**: Add new consumers without changing the API

## Project Structure
```
Go-server-crud/
├── main.go           # Main server with CRUD endpoints
├── events.go         # RabbitMQ event publisher/consumer
├── main_test.go      # Tests for CRUD operations
├── events_test.go    # Tests for event system
└── examples/
    └── consumer/     # Example event consumer application
```

## License
See LICENSE file for details.