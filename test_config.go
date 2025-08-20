package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Shopify/sarama"
)

// TestConfig holds configuration for tests
type TestConfig struct {
	Brokers []string
	GroupID string
	Topics  []string
}

// DefaultTestConfig returns a default test configuration
func DefaultTestConfig() *TestConfig {
	return &TestConfig{
		Brokers: []string{"localhost:9092"},
		GroupID: "test-consumer-group",
		Topics:  []string{"test-topic"},
	}
}

// CreateTestConsumer creates a consumer for testing purposes
func CreateTestConsumer(config *TestConfig) *Consumer {
	return &Consumer{
		topics: config.Topics,
		ready:  make(chan bool),
	}
}

// CreateTestMessage creates a test message with the given parameters
func CreateTestMessage(topic string, partition int32, offset int64, key, value string) *sarama.ConsumerMessage {
	return &sarama.ConsumerMessage{
		Topic:     topic,
		Partition: partition,
		Offset:    offset,
		Key:       []byte(key),
		Value:     []byte(value),
	}
}

// CreateTestContext creates a test context with timeout
func CreateTestContext(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}

// CreateTestContextWithCancel creates a test context with cancellation
func CreateTestContextWithCancel() (context.Context, context.CancelFunc) {
	return context.WithCancel(context.Background())
}

// WaitForChannel waits for a channel to receive a value or timeout
func WaitForChannel(ch <-chan interface{}, timeout time.Duration) bool {
	select {
	case <-ch:
		return true
	case <-time.After(timeout):
		return false
	}
}

// WaitForBoolChannel waits for a bool channel to receive a value or timeout
func WaitForBoolChannel(ch <-chan bool, timeout time.Duration) bool {
	select {
	case <-ch:
		return true
	case <-time.After(timeout):
		return false
	}
}

// TestSetup performs common test setup
func TestSetup(t *testing.T) *TestConfig {
	t.Helper()

	config := DefaultTestConfig()

	// Add any additional test setup here
	// For example, you might want to create test topics in Kafka
	// or set up test data

	return config
}

// TestTeardown performs common test cleanup
func TestTeardown(t *testing.T, config *TestConfig) {
	t.Helper()

	// Add any additional test cleanup here
	// For example, you might want to delete test topics
	// or clean up test data
}

// BenchmarkConfig holds configuration for benchmarks
type BenchmarkConfig struct {
	MessageCount int
	MessageSize  int
}

// DefaultBenchmarkConfig returns a default benchmark configuration
func DefaultBenchmarkConfig() *BenchmarkConfig {
	return &BenchmarkConfig{
		MessageCount: 1000,
		MessageSize:  1024,
	}
}

// CreateBenchmarkMessages creates a slice of messages for benchmarking
func CreateBenchmarkMessages(config *BenchmarkConfig) []*sarama.ConsumerMessage {
	messages := make([]*sarama.ConsumerMessage, config.MessageCount)

	for i := 0; i < config.MessageCount; i++ {
		key := []byte(fmt.Sprintf("key-%d", i))
		value := make([]byte, config.MessageSize)
		for j := 0; j < config.MessageSize; j++ {
			value[j] = byte(i + j)
		}

		messages[i] = &sarama.ConsumerMessage{
			Topic:     "benchmark-topic",
			Partition: 0,
			Offset:    int64(i),
			Key:       key,
			Value:     value,
		}
	}

	return messages
}
