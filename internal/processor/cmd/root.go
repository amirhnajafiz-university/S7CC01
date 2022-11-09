package cmd

import (
	"context"

	"github.com/ceit-aut/ad-registration-service/pkg/config"
	"github.com/ceit-aut/ad-registration-service/pkg/model"
	"github.com/ceit-aut/ad-registration-service/pkg/mqtt"
	"github.com/ceit-aut/ad-registration-service/pkg/storage/mongodb"
	s32 "github.com/ceit-aut/ad-registration-service/pkg/storage/s3"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson"
)

func GetCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "processor",
		Long: "starting the processor service",
		Run: func(_ *cobra.Command, _ []string) {
			main()
		},
	}
}

// main method of processor.
func main() {
	// load configs
	cfg := config.Load()

	// mongodb connection
	mongo, err := mongodb.NewConnection(cfg.Storage.Mongo)
	if err != nil {
		panic(err)
	}

	// s3 connection
	_, err = s32.NewSession(cfg.Storage.S3)
	if err != nil {
		panic(err)
	}

	// rabbitmq connection
	mq, err := mqtt.NewConnection(cfg.MQTT)
	if err != nil {
		panic(err)
	}

	events, _ := mq.Channel.Consume(
		mq.Queue,
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)

	for event := range events {
		var (
			ctx = context.Background()

			// get id from rabbitMQ
			id = string(event.Body)

			filter = bson.M{"id": id}

			c = mongo.Collection(model.AdCollection)

			ad model.Ad
		)

		value := c.FindOne(ctx, filter, nil)
		if err := value.Decode(&ad); err != nil {
			continue
		}

		// todo
		// calling imagga
		// updating ad status and category
		// sending email if valid

		if _, err := c.UpdateOne(ctx, filter, ad, nil); err != nil {
			continue
		}
	}
}
