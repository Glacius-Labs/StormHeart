package router

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/glacius-labs/StormHeart/internal/application/command"
	"github.com/glacius-labs/StormHeart/internal/application/handler"
)

type WatermillCommandRouter struct {
	pub      message.Publisher
	sub      message.Subscriber
	log      watermill.LoggerAdapter
	handlers map[command.CommandType]handler.Handler
	router   *message.Router
}

func NewWatermillCommandRouter(
	pub message.Publisher,
	sub message.Subscriber,
	log watermill.LoggerAdapter,
) (*WatermillCommandRouter, error) {
	r, err := message.NewRouter(message.RouterConfig{}, log)
	if err != nil {
		return nil, fmt.Errorf("failed to create watermill router: %w", err)
	}

	return &WatermillCommandRouter{
		pub:      pub,
		sub:      sub,
		log:      log,
		handlers: make(map[command.CommandType]handler.Handler),
		router:   r,
	}, nil
}

func (r *WatermillCommandRouter) RegisterHandler(h handler.Handler) error {
	cmdType := h.CommandType()
	topic := string(cmdType)

	if _, exists := r.handlers[cmdType]; exists {
		return fmt.Errorf("handler for command type %s already registered", cmdType)
	}

	r.handlers[cmdType] = h

	r.router.AddHandler(
		topic+"-handler",
		topic,
		r.sub,
		"", // no output topic
		r.pub,
		func(msg *message.Message) ([]*message.Message, error) {
			ctx := context.Background()

			cmd, err := r.decodeCommand(cmdType, msg.Payload)
			if err != nil {
				r.log.Error("failed to decode command", err, nil)
				return nil, err
			}

			return nilOrError(h.Handle(ctx, cmd))
		},
	)

	return nil
}

func (r *WatermillCommandRouter) Publish(ctx context.Context, cmd command.Command) error {
	data, err := json.Marshal(cmd)
	if err != nil {
		return fmt.Errorf("failed to marshal command: %w", err)
	}

	topic := string(cmd.CommandType())

	msg := message.NewMessage(watermill.NewUUID(), data)

	return r.pub.Publish(topic, msg)
}

func (r *WatermillCommandRouter) Start() error {
	return r.router.Run(context.Background())
}

func (r *WatermillCommandRouter) Close() error {
	return r.router.Close()
}

func nilOrError(err error) ([]*message.Message, error) {
	if err != nil {
		return nil, err
	}
	return nil, nil
}
