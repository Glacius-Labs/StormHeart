package mock

import (
	"context"
	"errors"

	"github.com/glacius-labs/StormHeart/internal/infrastructure/mqtt"
)

type MockClient struct {
	ShouldFailConnect   bool
	ShouldFailPublish   bool
	ShouldFailSubscribe bool
	ReceivedTopic       string
	ReceivedHandler     mqtt.MessageHandler
}

func NewMockClient() *MockClient {
	return &MockClient{}
}

func (m *MockClient) Connect() error {
	if m.ShouldFailConnect {
		return errors.New("simulated connect failure")
	}
	return nil
}

func (m *MockClient) Publish(ctx context.Context, topic string, payload []byte) error {
	if m.ShouldFailPublish {
		return errors.New("simulated publish failure")
	}
	m.ReceivedTopic = topic
	return nil
}

func (m *MockClient) Subscribe(ctx context.Context, topic string, handler mqtt.MessageHandler) error {
	if m.ShouldFailSubscribe {
		return errors.New("simulated subscribe failure")
	}
	m.ReceivedTopic = topic
	m.ReceivedHandler = handler
	return nil
}

func (m *MockClient) Disconnect() {
	// No-op
}
