package mqtt

import "context"

type MessageHandler func(ctx context.Context, payload []byte)
