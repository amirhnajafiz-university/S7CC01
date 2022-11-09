package handler

import (
	"context"

	"github.com/ceit-aut/ad-registration-service/pkg/model"
	"github.com/ceit-aut/ad-registration-service/pkg/mqtt"
	"github.com/ceit-aut/ad-registration-service/pkg/storage/s3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Handler
// manages to handle the processor service.
type Handler struct {
	Mongo *mongo.Database
	MQTT  *mqtt.MQTT
	S3    *s3.S3
}

func (h *Handler) Handle() {
	events, _ := h.MQTT.Channel.Consume(
		h.MQTT.Queue,
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

			c = h.Mongo.Collection(model.AdCollection)

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
