package handler

import (
	"context"
	"log"

	"github.com/ceit-aut/ad-registration-service/pkg/model"
	"github.com/ceit-aut/ad-registration-service/pkg/mqtt"
	"github.com/ceit-aut/ad-registration-service/pkg/service/mail"
	"github.com/ceit-aut/ad-registration-service/pkg/storage/s3"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Handler
// manages to handle the processor service.
type Handler struct {
	Mongo *mongo.Database
	Mail  *mail.Mailgun
	MQTT  *mqtt.MQTT
	S3    *s3.S3
}

// Handle
// listens over rabbitMQ and processes
// the input images.
func (h *Handler) Handle() {
	// creating a consumer for rabbitMQ
	events, _ := h.MQTT.Channel.Consume(
		h.MQTT.Queue,
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)

	// listen over rabbitMQ events
	for event := range events {
		var (
			// creating a new context
			ctx = context.Background()
			// get id from rabbitMQ
			id = string(event.Body)
			// mongodb filter
			filter = bson.M{"id": id}
			// connecting to mongodb collection
			c = h.Mongo.Collection(model.AdCollection)
			// creating a new ad model
			ad model.Ad
		)

		// finding the ad
		value := c.FindOne(ctx, filter, nil)
		if err := value.Decode(&ad); err != nil {
			log.Println(err)

			continue
		}

		// todo
		// calling imagga
		// updating ad status and category
		// sending email if valid

		if _, err := c.UpdateOne(ctx, filter, ad, nil); err != nil {
			log.Println(err)

			continue
		}

		log.Printf("success processing {id: %s}", id)
	}
}
