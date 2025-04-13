package mqtt

import (
	"context"

	paho "github.com/eclipse/paho.mqtt.golang"
)

type PahoClient struct {
	client paho.Client
}

func NewPahoClient(id, brokerURL string) *PahoClient {
	if id == "" {
		panic("id cannot be empty")
	}

	if brokerURL == "" {
		panic("broker URL cannot be empty")
	}

	opts := paho.NewClientOptions().
		AddBroker(brokerURL).
		SetClientID("stormlink").
		SetAutoReconnect(true).
		SetCleanSession(false)

	client := paho.NewClient(opts)

	return &PahoClient{client: client}
}

func (p *PahoClient) Connect() error {
	token := p.client.Connect()
	token.Wait()
	return token.Error()
}

func (p *PahoClient) Subscribe(ctx context.Context, topic string, handler MessageHandler) error {
	token := p.client.Subscribe(topic, 1, func(client paho.Client, msg paho.Message) {
		handler(ctx, msg.Payload())
	})
	token.Wait()

	return token.Error()
}

func (p *PahoClient) Disconnect() {
	p.client.Disconnect(250)
}
