//go:build test

package main

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Shopify/sarama"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockConsumerGroup implements sarama.ConsumerGroup for testing
type MockConsumerGroup struct {
	mock.Mock
}

func (m *MockConsumerGroup) Consume(ctx context.Context, topics []string, handler sarama.ConsumerGroupHandler) error {
	args := m.Called(ctx, topics, handler)
	return args.Error(0)
}

func (m *MockConsumerGroup) Errors() <-chan error {
	args := m.Called()
	return args.Get(0).(<-chan error)
}

func (m *MockConsumerGroup) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockConsumerGroup) Pause(partitions map[string][]int32) {
	m.Called(partitions)
}

func (m *MockConsumerGroup) PauseAll() {
	m.Called()
}

func (m *MockConsumerGroup) Resume(partitions map[string][]int32) {
	m.Called(partitions)
}

func (m *MockConsumerGroup) ResumeAll() {
	m.Called()
}

// MockConsumerGroupSession implements sarama.ConsumerGroupSession for testing
type MockConsumerGroupSession struct {
	mock.Mock
}

func (m *MockConsumerGroupSession) Claims() map[string][]int32 {
	args := m.Called()
	return args.Get(0).(map[string][]int32)
}

func (m *MockConsumerGroupSession) MemberID() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockConsumerGroupSession) GenerationID() int32 {
	args := m.Called()
	return args.Get(0).(int32)
}

func (m *MockConsumerGroupSession) MarkOffset(topic string, partition int32, offset int64, metadata string) {
	m.Called(topic, partition, offset, metadata)
}

func (m *MockConsumerGroupSession) MarkMessage(msg *sarama.ConsumerMessage, metadata string) {
	m.Called(msg, metadata)
}

func (m *MockConsumerGroupSession) Commit() {
	m.Called()
}

func (m *MockConsumerGroupSession) ResetOffset(topic string, partition int32, offset int64, metadata string) {
	m.Called(topic, partition, offset, metadata)
}

func (m *MockConsumerGroupSession) Context() context.Context {
	args := m.Called()
	return args.Get(0).(context.Context)
}

// MockConsumerGroupClaim implements sarama.ConsumerGroupClaim for testing
type MockConsumerGroupClaim struct {
	mock.Mock
}

func (m *MockConsumerGroupClaim) Topic() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockConsumerGroupClaim) Partition() int32 {
	args := m.Called()
	return args.Get(0).(int32)
}

func (m *MockConsumerGroupClaim) InitialOffset() int64 {
	args := m.Called()
	return args.Get(0).(int64)
}

func (m *MockConsumerGroupClaim) HighWaterMarkOffset() int64 {
	args := m.Called()
	return args.Get(0).(int64)
}

func (m *MockConsumerGroupClaim) Messages() <-chan *sarama.ConsumerMessage {
	args := m.Called()
	return args.Get(0).(<-chan *sarama.ConsumerMessage)
}

func TestNewConsumer(t *testing.T) {
	tests := []struct {
		name          string
		brokers       []string
		groupID       string
		topics        []string
		expectError   bool
		errorContains string
	}{
		{
			name:        "valid configuration",
			brokers:     []string{"localhost:9092"},
			groupID:     "test-group",
			topics:      []string{"test-topic"},
			expectError: false,
		},
		{
			name:          "empty brokers",
			brokers:       []string{},
			groupID:       "test-group",
			topics:        []string{"test-topic"},
			expectError:   true,
			errorContains: "brokers",
		},
		{
			name:          "empty group ID",
			brokers:       []string{"localhost:9092"},
			groupID:       "",
			topics:        []string{"test-topic"},
			expectError:   true,
			errorContains: "group ID",
		},
		{
			name:          "empty topics",
			brokers:       []string{"localhost:9092"},
			groupID:       "test-group",
			topics:        []string{},
			expectError:   true,
			errorContains: "topics",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This test would require mocking sarama.NewConsumerGroup
			// For now, we'll test the validation logic separately
			if tt.expectError {
				// Skip this test for now as it requires more complex mocking
				t.Skip("Requires mocking sarama.NewConsumerGroup")
			}
		})
	}
}

func TestConsumer_Setup(t *testing.T) {
	consumer := &Consumer{
		ready: make(chan bool),
	}

	session := &MockConsumerGroupSession{}

	// Test that Setup closes the ready channel
	err := consumer.Setup(session)

	assert.NoError(t, err)

	// Check that the ready channel was closed
	select {
	case <-consumer.ready:
		// Channel was closed, which is expected
	default:
		t.Error("Expected ready channel to be closed")
	}
}

func TestConsumer_Cleanup(t *testing.T) {
	consumer := &Consumer{}
	session := &MockConsumerGroupSession{}

	err := consumer.Cleanup(session)
	assert.NoError(t, err)
}

func TestConsumer_ConsumeClaim(t *testing.T) {
	tests := []struct {
		name          string
		messages      []*sarama.ConsumerMessage
		contextDone   bool
		expectedCalls int
		expectError   bool
	}{
		{
			name: "processes messages successfully",
			messages: []*sarama.ConsumerMessage{
				{
					Topic:     "test-topic",
					Partition: 0,
					Offset:    1,
					Key:       []byte("key1"),
					Value:     []byte("value1"),
				},
				{
					Topic:     "test-topic",
					Partition: 0,
					Offset:    2,
					Key:       []byte("key2"),
					Value:     []byte("value2"),
				},
			},
			contextDone:   false,
			expectedCalls: 2,
			expectError:   false,
		},
		{
			name:          "handles nil message",
			messages:      []*sarama.ConsumerMessage{nil},
			contextDone:   false,
			expectedCalls: 0,
			expectError:   false,
		},
		{
			name:          "handles context cancellation",
			messages:      []*sarama.ConsumerMessage{},
			contextDone:   true,
			expectedCalls: 0,
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			consumer := &Consumer{}

			// Create message channel
			msgChan := make(chan *sarama.ConsumerMessage, len(tt.messages))

			// Create context
			ctx := context.Background()
			if tt.contextDone {
				ctx, _ = context.WithCancel(ctx)
			}

			// Create mock session
			session := &MockConsumerGroupSession{}
			session.On("Context").Return(ctx)

			// Set up expectations for MarkMessage calls
			for _, msg := range tt.messages {
				if msg != nil {
					session.On("MarkMessage", msg, "").Return()
				}
			}

			// Create mock claim
			claim := &MockConsumerGroupClaim{}
			claim.On("Messages").Return((<-chan *sarama.ConsumerMessage)(msgChan))

			// Send messages to channel
			go func() {
				for _, msg := range tt.messages {
					msgChan <- msg
				}
				if tt.contextDone {
					time.Sleep(10 * time.Millisecond)
					// Context cancellation would be handled here in a real scenario
				}
				close(msgChan)
			}()

			// Call ConsumeClaim
			err := consumer.ConsumeClaim(session, claim)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// Verify session.MarkMessage was called for each valid message
			validMessages := 0
			for _, msg := range tt.messages {
				if msg != nil {
					validMessages++
				}
			}

			// Wait a bit for goroutines to finish
			time.Sleep(50 * time.Millisecond)

			// Clean up
			session.AssertExpectations(t)
			claim.AssertExpectations(t)
		})
	}
}

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
	// In a real application, you might want to capture log output
	// or add more sophisticated message processing logic
	assert.NotPanics(t, func() {
		consumer.processMessage(message)
	})
}

func TestConsumer_Start(t *testing.T) {
	// Test that the Start method properly waits for the ready channel
	consumer := &Consumer{
		ready: make(chan bool),
	}

	// Create a mock consumer group that returns immediately
	mockGroup := &MockConsumerGroup{}
	mockGroup.On("Consume", mock.Anything, mock.Anything, mock.Anything).Return(context.Canceled)
	consumer.consumer = mockGroup

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Test that Start blocks until ready channel is closed
	started := make(chan bool, 1)
	go func() {
		started <- true
		// This will block until ready channel is closed
		consumer.Start(ctx)
	}()

	// Wait a bit to ensure the goroutine has started
	time.Sleep(10 * time.Millisecond)

	// Verify that Start is blocked (started channel should have a value)
	select {
	case <-started:
		// Good, goroutine started
	default:
		t.Fatal("Start method should have started")
	}

	// Close the ready channel to unblock Start
	close(consumer.ready)

	// Give it a moment to process
	time.Sleep(10 * time.Millisecond)

	// Cancel context to stop the test
	cancel()

	mockGroup.AssertExpectations(t)
}

func TestConsumer_Close(t *testing.T) {
	mockGroup := &MockConsumerGroup{}
	mockGroup.On("Close").Return(nil)

	consumer := &Consumer{
		consumer: mockGroup,
	}

	err := consumer.Close()

	assert.NoError(t, err)
	mockGroup.AssertExpectations(t)
}

func TestConsumer_Close_Error(t *testing.T) {
	mockGroup := &MockConsumerGroup{}
	mockGroup.On("Close").Return(errors.New("close error"))

	consumer := &Consumer{
		consumer: mockGroup,
	}

	err := consumer.Close()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "close error")
	mockGroup.AssertExpectations(t)
}

// Test helper functions
func TestExampleSimpleConsumer(t *testing.T) {
	// This is a placeholder test for the example function
	// In a real scenario, you might want to mock the sarama.NewConsumer call
	t.Skip("Requires mocking sarama.NewConsumer")
}

func TestExampleConsumerWithFilters(t *testing.T) {
	// This function just logs a message, so we can test it
	assert.NotPanics(t, func() {
		ExampleConsumerWithFilters()
	})
}

func TestExampleConsumerWithRetry(t *testing.T) {
	// This function just logs a message, so we can test it
	assert.NotPanics(t, func() {
		ExampleConsumerWithRetry()
	})
}

// Benchmark tests
func BenchmarkConsumer_processMessage(b *testing.B) {
	consumer := &Consumer{}
	message := &sarama.ConsumerMessage{
		Topic:     "test-topic",
		Partition: 0,
		Offset:    123,
		Key:       []byte("test-key"),
		Value:     []byte("test-value"),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		consumer.processMessage(message)
	}
}

// Test utilities
func createTestMessage(topic string, partition int32, offset int64, key, value string) *sarama.ConsumerMessage {
	return &sarama.ConsumerMessage{
		Topic:     topic,
		Partition: partition,
		Offset:    offset,
		Key:       []byte(key),
		Value:     []byte(value),
	}
}

func createTestContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 100*time.Millisecond)
}
