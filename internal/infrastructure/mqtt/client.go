package mqtt

import "context"

type Client interface {
	Connect() error
	Subscribe(ctx context.Context, topic string, handler MessageHandler) error
	Disconnect()
}
