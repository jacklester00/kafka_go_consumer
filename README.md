# Kafka Go Consumer Example

A simple, clean example of building a Kafka consumer in Go using the Sarama library with Docker.

## Prerequisites

- [Go](https://golang.org/dl/) 1.21+ 
- [Docker](https://docs.docker.com/get-docker/) and Docker Compose

## Project Structure

```
kafka_go_consumer/
├── README.md              # This file
├── docker-compose.yml     # Kafka cluster setup
├── example_consumer.go    # Main consumer implementation
├── example_consumer_test.go # Simple unit tests
├── integration_test.go    # Integration tests (requires Kafka)
├── .github/workflows/test.yml # CI/CD pipeline
├── .golangci.yml         # Code linting rules
├── Makefile              # Development commands
├── go.mod                # Go dependencies
└── .gitignore            # Git ignore rules
```

## Quick Start

### 1. Start Kafka

```bash
# Start Kafka cluster
make docker-up

# Or manually with docker-compose
docker-compose up -d
```

### 2. Run the Consumer

```bash
# Install dependencies and run
go mod tidy
go run example_consumer.go
```

### 3. Send Test Messages

```bash
# In another terminal, send a test message
docker exec -it kafka kafka-console-producer --bootstrap-server localhost:9092 --topic example-topic
> Hello, Kafka!
> ^D
```

Your consumer should display the received message!

## Configuration

### Consumer Settings

The consumer is configured to:
- Connect to `localhost:9092`
- Join consumer group `example-consumer-group`
- Consume from topic `example-topic`
- Start from the oldest offset
- Auto-commit offsets every second

### Customization

Modify these values in `example_consumer.go`:

```go
brokers := []string{"localhost:9092"}  // Kafka brokers
groupID := "example-consumer-group"    // Consumer group ID
topics := []string{"example-topic"}    // Topics to consume
```

## Development Commands

```bash
# Docker operations
make docker-up          # Start Kafka cluster
make docker-down        # Stop Kafka cluster
make docker-logs        # View Kafka logs

# Development
make build              # Build the application
make run                # Build and run consumer
make test               # Run unit tests
make integration-test   # Run integration tests (requires Kafka)

# Code quality
make fmt                # Format code
make lint               # Run linter
make check              # Run all checks
```

## Testing

### Unit Tests
```bash
go test -v                    # Run unit tests
make test                     # Run via Makefile
```

### Integration Tests
```bash
make integration-test         # Starts Kafka, runs tests, cleans up
go test -tags=integration -v  # Run manually (requires Kafka running)
```

## Troubleshooting

**Kafka won't start**: Check `docker-compose logs kafka`

**Consumer can't connect**: Verify Kafka is healthy with `docker-compose ps`

**No messages received**: Check topic exists with:
```bash
docker exec -it kafka kafka-topics --bootstrap-server localhost:9092 --list
```

**Reset everything**:
```bash
docker-compose down -v  # Removes all data
docker-compose up -d    # Fresh start
```

## Resources

- [Apache Kafka Documentation](https://kafka.apache.org/documentation/)
- [Sarama Go Client](https://github.com/Shopify/sarama)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
