# Kafka Go Consumer Example

This project demonstrates how to create a Kafka consumer in Go using the Sarama library, with a local Kafka cluster running on Docker using KRaft.

## Prerequisites

- [Go](https://golang.org/dl/) 1.25.0 or later
- [Docker](https://docs.docker.com/get-docker/) and Docker Compose
- [Git](https://git-scm.com/) (optional, for cloning)

## Project Structure

```
kafka_go_consumer/
├── README.md              # This file
├── docker-compose.yml     # Kafka cluster configuration
├── example_consumer.go    # Go Kafka consumer example
├── go.mod                 # Go module dependencies
└── go.sum                 # Go module checksums (generated)
```

## Quick Start

### 1. Start Kafka Cluster

Start the Kafka cluster using Docker Compose:

```bash
# Start Kafka only
docker-compose up -d

# Or start with Kafka UI (optional)
docker-compose --profile ui up -d
```

**Note**: The first startup may take a few minutes as Kafka formats its storage and initializes.

### 2. Verify Kafka is Running

Check if Kafka is healthy:

```bash
docker-compose ps
```

You should see Kafka running and healthy. If you started with the UI profile, you can also visit `http://localhost:8080` to see the Kafka UI.

### 3. Install Go Dependencies

```bash
go mod tidy
```

### 4. Run the Consumer

```bash
go run example_consumer.go
```

The consumer will start and wait for messages on the `example-topic`.

## Detailed Setup

### Kafka Configuration

The `docker-compose.yml` file sets up:

- **Kafka 3.6.1** with KRaft consensus (no Zookeeper)
- **Single-node cluster** for development
- **Port 9092** exposed for external connections
- **Persistent storage** using Docker volumes
- **Health checks** to ensure Kafka is ready

### Consumer Configuration

The Go consumer is configured to:

- Connect to `localhost:9092`
- Join consumer group `example-consumer-group`
- Consume from topic `example-topic`
- Start from the oldest offset
- Auto-commit offsets every second

## Testing the Setup

### 1. Create a Test Topic (Optional)

Topics are auto-created by default, but you can manually create one:

```bash
# Connect to Kafka container
docker exec -it kafka bash

# Create a topic
kafka-topics --bootstrap-server localhost:9093 --create --topic example-topic --partitions 1 --replication-factor 1

# List topics
kafka-topics --bootstrap-server localhost:9093 --list
```

### 2. Send Test Messages

```bash
# Connect to Kafka container
docker exec -it kafka bash

# Send a test message
echo "Hello, Kafka!" | kafka-console-producer --bootstrap-server localhost:9093 --topic example-topic

# Or send multiple messages
kafka-console-producer --bootstrap-server localhost:9093 --topic example-topic
> Message 1
> Message 2
> Message 3
> ^D
```

### 3. Verify Consumer Receives Messages

Your Go consumer should now display the received messages:

```
2024/01/XX XX:XX:XX Consumer is ready, starting to consume messages...
2024/01/XX XX:XX:XX Message topic:"example-topic" partition:0 offset:0
2024/01/XX XX:XX:XX Message value: Hello, Kafka!
2024/01/XX XX:XX:XX Message key: 
```

## Customization

### Change Kafka Configuration

Edit the environment variables in `docker-compose.yml`:

```yaml
environment:
  KAFKA_AUTO_CREATE_TOPICS_ENABLE: true  # Enable/disable auto topic creation
  KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1  # Replication factor for internal topics
```

### Change Consumer Configuration

Modify the configuration in `example_consumer.go`:

```go
// Configuration
brokers := []string{"localhost:9092"}  // Change broker addresses
groupID := "example-consumer-group"    // Change consumer group
topics := []string{"example-topic"}    // Change topic names
```

### Add More Brokers

For a multi-broker setup, add more services to `docker-compose.yml`:

```yaml
kafka-2:
  image: apache/kafka:3.6.1
  # ... similar configuration with different NODE_ID
```

## Troubleshooting

### Common Issues

#### 1. Kafka Won't Start

**Problem**: Kafka container exits immediately
**Solution**: Check logs for errors:
```bash
docker-compose logs kafka
```

**Common causes**:
- Port conflicts (check if 9092 is already in use)
- Insufficient memory (Kafka needs at least 1GB)
- Storage permission issues

#### 2. Consumer Can't Connect

**Problem**: Consumer fails to connect to Kafka
**Solution**: Verify Kafka is healthy:
```bash
docker-compose ps
docker-compose logs kafka
```

**Common causes**:
- Kafka not fully started (wait for health check)
- Wrong broker address
- Network issues between host and container

#### 3. No Messages Received

**Problem**: Consumer runs but doesn't receive messages
**Solution**: Check if messages are being sent to the correct topic:
```bash
# Check topic exists
docker exec -it kafka kafka-topics --bootstrap-server localhost:9093 --list

# Check message count
docker exec -it kafka kafka-run-class kafka.tools.GetOffsetShell --bootstrap-server localhost:9093 --topic example-topic
```

### Reset Everything

If you need to start fresh:

```bash
# Stop and remove containers
docker-compose down

# Remove volumes (WARNING: This deletes all data)
docker-compose down -v

# Start again
docker-compose up -d
```

## Development Workflow

### 1. Start Development Environment

```bash
# Start Kafka
docker-compose up -d

# Wait for Kafka to be healthy
docker-compose ps
```

### 2. Develop and Test

```bash
# Run consumer
go run example_consumer.go

# In another terminal, send test messages
docker exec -it kafka kafka-console-producer --bootstrap-server localhost:9093 --topic example-topic
```

### 3. Iterate

- Modify `example_consumer.go`
- Test changes
- Repeat

## Production Considerations

This setup is for **development only**. For production:

- Use multiple Kafka brokers
- Configure proper security (SASL/SSL)
- Set appropriate replication factors
- Use external monitoring and alerting
- Configure proper retention policies
- Use dedicated storage volumes

## Additional Resources

- [Apache Kafka Documentation](https://kafka.apache.org/documentation/)
- [Sarama Go Client](https://github.com/Shopify/sarama)
- [KRaft Mode](https://kafka.apache.org/documentation/#kraft)
- [Docker Compose Documentation](https://docs.docker.com/compose/)

## License

This project is provided as an example and is open source.
