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
		return nil, fmt.Errorf("failed to create watermill command router: %w", err)
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

	r.router.AddNoPublisherHandler(
		h.Name(),
		topic,
		r.sub,
		func(msg *message.Message) error {
			cmd, err := r.decodeCommand(cmdType, msg.Payload)
			if err != nil {
				return fmt.Errorf("failed to decode command: %w", err)
			}

			return h.Handle(msg.Context(), cmd)
		},
	)

	return nil
}

func (r *WatermillCommandRouter) Publish(cmd command.Command) error {
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
