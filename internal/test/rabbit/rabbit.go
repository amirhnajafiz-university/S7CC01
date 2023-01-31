package rabbit

import (
	"context"
	"log"

	"github.com/ceit-aut/S7CC01/pkg/config"
	"github.com/ceit-aut/S7CC01/pkg/mqtt"

	"github.com/rabbitmq/amqp091-go"
	"github.com/spf13/cobra"
)

// rabbit
// is used to test rabbitMQ cluster.
type rabbit struct {
	MQTT *mqtt.MQTT
}

// publish
// over rabbitMQ queue.
func (r *rabbit) publish(data []byte) error {
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

// GetCommand
// returns the cobra command.
func GetCommand() *cobra.Command {
	return &cobra.Command{
		Use: "rabbit",
		Run: func(_ *cobra.Command, args []string) {
			main(args[0])
		},
	}
}

// main
// start rabbit test.
func main(id string) {
	// load configs
	cfg := config.Load()

	// mqtt connection
	mq, err := mqtt.NewConnection(cfg.MQTT)
	if err != nil {
		panic(err)
	}

	// creating a new rabbit
	r := rabbit{
		MQTT: mq,
	}

	if err := r.publish([]byte(id)); err != nil {
		panic(err)
	}

	log.Println("succeed")
}
