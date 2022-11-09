package handler

import (
	"context"
	"fmt"
	"log"

	"github.com/ceit-aut/ad-registration-service/pkg/enum"
	"github.com/ceit-aut/ad-registration-service/pkg/model"
	"github.com/ceit-aut/ad-registration-service/pkg/mqtt"
	"github.com/ceit-aut/ad-registration-service/pkg/service/imagga"
	"github.com/ceit-aut/ad-registration-service/pkg/service/mail"
	"github.com/ceit-aut/ad-registration-service/pkg/storage/s3"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Handler
// manages to handle the processor service.
type Handler struct {
	Imagga *imagga.Imagga
	Mongo  *mongo.Database
	Mail   *mail.Mailgun
	MQTT   *mqtt.MQTT
	S3     *s3.S3
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

		// image tag
		resp, err := h.Imagga.Process("")
		if err != nil {
			log.Println(err)

			continue
		}

		// response check
		if len(resp.Result.Tags) > 0 {
			if resp.Result.Tags[0].Confidence == 100 {
				ad.Category = resp.Result.Tags[0].Tag.En
				ad.State = enum.AcceptState

				// send email
				go func() {
					msg := fmt.Sprintf(
						"Dear '%s', your ad about \"%s\", registered sucessfully.",
						ad.Email,
						ad.Description,
					)
					if err := h.Mail.Send(msg, "ad-status", ad.Email); err != nil {
						log.Println(err)
					}
				}()
			} else {
				ad.State = enum.RejectState
			}
		} else {
			ad.State = enum.RejectState
		}

		// update mongodb
		if _, err := c.UpdateOne(ctx, filter, ad, nil); err != nil {
			log.Println(err)

			continue
		}

		log.Printf("success processing {id: %s}", id)
	}
}
