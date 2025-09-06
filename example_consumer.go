package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
)

// Consumer represents a Kafka consumer
type Consumer struct {
	consumer sarama.ConsumerGroup
	topics   []string
	ready    chan bool
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(brokers []string, groupID string, topics []string) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_8_1_0 // Use appropriate Kafka version
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second

	consumer, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	return &Consumer{
		consumer: consumer,
		topics:   topics,
		ready:    make(chan bool),
	}, nil
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (c *Consumer) Setup(session sarama.ConsumerGroupSession) error {
	log.Printf("Consumer group session setup - Member ID: %s", session.MemberID())
	log.Printf("Consumer group session setup - Claims: %v", session.Claims())
	// Mark the consumer as ready
	close(c.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages()
func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE: Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/main/consumer_group.go#L27-L29
	for {
		select {
		case message := <-claim.Messages():
			if message == nil {
				return nil
			}

			// Process the message
			c.processMessage(message)

			// Mark message as processed
			session.MarkMessage(message, "")

		case <-session.Context().Done():
			return nil
		}
	}
}

// processMessage handles individual Kafka messages
func (c *Consumer) processMessage(message *sarama.ConsumerMessage) {
	log.Printf("Message topic:%q partition:%d offset:%d\n", message.Topic, message.Partition, message.Offset)
	log.Printf("Message value: %s\n", string(message.Value))
	log.Printf("Message key: %s\n", string(message.Key))

	// Add your custom message processing logic here
	// For example:
	// - Parse JSON payload
	// - Validate data
	// - Store in database
	// - Send to another service
	// - etc.
}

// Start begins consuming messages
func (c *Consumer) Start(ctx context.Context) error {
	log.Println("Starting consumer, waiting for setup...")
	// Start consuming in a loop
	for {
		log.Printf("Calling consumer.Consume() for topics: %v", c.topics)
		err := c.consumer.Consume(ctx, c.topics, c)
		if err != nil {
			log.Printf("Error from consumer: %v", err)
		}

		// Check if context was cancelled
		if ctx.Err() != nil {
			return ctx.Err()
		}

		// Wait until the consumer has been set up
		select {
		case <-c.ready:
			log.Println("Consumer is ready, continuing to consume messages...")
		default:
			// ready channel already closed, continue
		}
	}
}

// Close closes the consumer
func (c *Consumer) Close() error {
	return c.consumer.Close()
}

func main() {
	// Configuration
	brokers := []string{"localhost:9092"} // Change to your Kafka broker addresses
	groupID := "example-consumer-group"
	topics := []string{"example-topic"} // Change to your topic names

	// Create consumer
	consumer, err := NewConsumer(brokers, groupID, topics)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer consumer.Close()

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup
	wg.Add(1)

	// Start consumer in goroutine
	go func() {
		defer wg.Done()
		if err := consumer.Start(ctx); err != nil {
			log.Printf("Consumer error: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-signals
	log.Println("Received shutdown signal, stopping consumer...")

	// Cancel context to stop consumer
	cancel()

	// Wait for consumer to finish
	wg.Wait()
	log.Println("Consumer stopped gracefully")
}

// Example usage functions for different scenarios

// ExampleSimpleConsumer shows a simpler consumer without consumer groups
func ExampleSimpleConsumer() {
	config := sarama.NewConfig()
	config.Version = sarama.V2_8_1_0
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, config)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition("example-topic", 0, sarama.OffsetOldest)
	if err != nil {
		log.Printf("Failed to create partition consumer: %v", err)
		return
	}
	defer partitionConsumer.Close()

	for message := range partitionConsumer.Messages() {
		log.Printf("Received message: %s", string(message.Value))
	}
}

// ExampleConsumerWithFilters shows how to filter messages
func ExampleConsumerWithFilters() {
	// This would be implemented with custom filtering logic
	// based on message headers, keys, or values
	log.Println("Example consumer with filters - implement custom filtering logic")
}

// ExampleConsumerWithRetry shows how to implement retry logic
func ExampleConsumerWithRetry() {
	// This would be implemented with exponential backoff
	// and dead letter queue handling
	log.Println("Example consumer with retry - implement retry logic with backoff")
}
