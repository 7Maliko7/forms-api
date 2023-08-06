package rabbitmq

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Broker struct {
	channel *amqp.Channel
	conn    *amqp.Connection
}

func NewBroker(url string) (*Broker, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &Broker{
		channel: ch,
		conn:    conn,
	}, nil
}

func (b *Broker) Close() error {
	err := b.conn.Close()
	if err != nil {
		return err
	}
	err = b.channel.Close()
	if err != nil {
		return err
	}
	return nil
}

func (b *Broker) Publish(ctx context.Context, body []byte, exchange, routingKey string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := b.channel.PublishWithContext(ctx,
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	if err != nil {
		return err
	}
	return nil
}
