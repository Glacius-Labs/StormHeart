package mqtt

import "context"

type Client interface {
	Connect() error
	Subscribe(ctx context.Context, topic string, callback MessageHandler) error
	Disconnect()
}

type MessageHandler func(ctx context.Context, payload []byte)
