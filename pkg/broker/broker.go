package broker

import (
	"context"
)

type Broker interface {
	Publish(ctx context.Context, body []byte, exchange, routingKey string) error
	Close() error
}
