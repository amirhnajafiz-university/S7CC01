package test

import (
	"context"

	"github.com/ceit-aut/ad-registration-service/pkg/mqtt"

	"github.com/rabbitmq/amqp091-go"
)

// RabbitTest
// is used to test rabbitMQ cluster.
type RabbitTest struct {
	MQTT *mqtt.MQTT
}

// Publish
// over rabbitMQ queue.
func (r *RabbitTest) Publish(data []byte) error {
	// publish id over mqtt
	err := r.MQTT.Channel.PublishWithContext(
		context.Background(),
		"",
		r.MQTT.Queue,
		false,
		false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        data,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
