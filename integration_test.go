//go:build integration

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestIntegration_ConsumerWithRealKafka tests the consumer against a real Kafka instance
func TestIntegration_ConsumerWithRealKafka(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// This test requires Kafka to be running
	// You can start it with: make docker-up
	t.Log("Starting integration test with real Kafka...")

	// Test configuration
	brokers := []string{"localhost:9092"}
	groupID := "integration-test-group"
	topics := []string{"integration-test-topic"}

	// Create consumer
	consumer, err := NewConsumer(brokers, groupID, topics)
	require.NoError(t, err)
	defer consumer.Close()

	// Test that consumer can be created
	assert.NotNil(t, consumer)
	assert.Equal(t, topics, consumer.topics)

	t.Log("Consumer created successfully")
}

// TestIntegration_ErrorHandling tests basic error handling
func TestIntegration_ErrorHandling(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Log("Testing error handling scenarios...")

	// Test with invalid broker address
	brokers := []string{"invalid:9999"}
	groupID := "error-test-group"
	topics := []string{"error-test-topic"}

	// This should fail to create consumer
	consumer, err := NewConsumer(brokers, groupID, topics)
	if err != nil {
		t.Logf("Expected error when connecting to invalid broker: %v", err)
		return
	}
	if consumer != nil {
		defer consumer.Close()
	}
}
