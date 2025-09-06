package main

import (
	"testing"

	"github.com/Shopify/sarama"
	"github.com/stretchr/testify/assert"
)

// TestConsumer_processMessage tests the message processing function
func TestConsumer_processMessage(t *testing.T) {
	consumer := &Consumer{}

	message := &sarama.ConsumerMessage{
		Topic:     "test-topic",
		Partition: 0,
		Offset:    123,
		Key:       []byte("test-key"),
		Value:     []byte("test-value"),
	}

	// This test mainly verifies that processMessage doesn't panic
	assert.NotPanics(t, func() {
		consumer.processMessage(message)
	})
}

// TestConsumer_Cleanup tests the Cleanup method
func TestConsumer_Cleanup(t *testing.T) {
	consumer := &Consumer{}

	// Cleanup should always return nil
	err := consumer.Cleanup(nil)
	assert.NoError(t, err)
}

// TestExampleFunctions tests the example functions don't panic
func TestExampleFunctions(t *testing.T) {
	assert.NotPanics(t, func() {
		ExampleConsumerWithFilters()
	})

	assert.NotPanics(t, func() {
		ExampleConsumerWithRetry()
	})
}
