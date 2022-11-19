package handler

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ceit-aut/ad-registration-service/pkg/enum"
	"github.com/ceit-aut/ad-registration-service/pkg/model"
	"github.com/ceit-aut/ad-registration-service/pkg/mqtt"
	"github.com/ceit-aut/ad-registration-service/pkg/service/imagga"
	"github.com/ceit-aut/ad-registration-service/pkg/service/mail"
	"github.com/ceit-aut/ad-registration-service/pkg/storage/s3"

	"github.com/aws/aws-sdk-go/aws"
	s3sdk "github.com/aws/aws-sdk-go/service/s3"
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
	log.Printf("start listening on queue: %s\n", h.MQTT.Queue)

	// creating a consumer for rabbitMQ
	events, err := h.MQTT.Channel.Consume(
		h.MQTT.Queue,
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		log.Printf("failed to consume messages: %v\n", err)

		return
	}

	log.Println("processor started ...")

	// listen over rabbitMQ events
	for event := range events {
		var (
			// creating a new context
			ctx = context.Background()
			// get id from rabbitMQ
			id = string(event.Body)
			// mongodb filter
			filter = bson.D{{"id", id}}
			// connecting to mongodb collection
			c = h.Mongo.Collection(model.AdCollection)
			// creating a new ad model
			ad model.Ad
		)

		log.Printf("receive id:\n\t%s\n", id)

		// finding the ad
		value := c.FindOne(ctx, filter, nil)
		if err := value.Decode(&ad); err != nil {
			log.Println(err)

			continue
		}

		log.Println("mongodb get by id succeed")

		// getting the image from s3
		svc := s3sdk.New(h.S3.Session, &aws.Config{
			Region:   aws.String(h.S3.Cfg.Region),
			Endpoint: aws.String(h.S3.Cfg.Endpoint),
		})

		req, _ := svc.GetObjectRequest(&s3sdk.GetObjectInput{
			Bucket: aws.String(h.S3.Cfg.Bucket),
			Key:    aws.String(ad.Id),
		})

		urlStr, err := req.Presign(15 * time.Minute)
		if err != nil {
			log.Println(err)

			continue
		}

		log.Printf("image url from s3:\n\t%s\n", urlStr)

		// image tag process
		resp, err := h.Imagga.Process(urlStr)
		if err != nil {
			log.Println(err)

			continue
		}

		log.Printf("imagga response:\n\t%d tags\n", len(resp.Result.Tags))

		// response check
		if len(resp.Result.Tags) > 0 {
			if validateTag(resp) {
				log.Printf("imagga response confidence: %f\n", resp.Result.Tags[0].Confidence)

				if resp.Result.Tags[0].Confidence > 50 {
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

							return
						}

						log.Printf("email send {id: %s}\n", id)
					}()
				} else {
					ad.State = enum.RejectState
				}
			} else {
				ad.State = enum.RejectState
			}
		} else {
			ad.State = enum.RejectState
		}

		log.Printf("imagga api call succeed")

		// update filter for mongodb
		update := bson.D{
			{
				"$set",
				bson.D{
					{"state", ad.State},
					{"category", ad.Category},
				},
			},
		}

		// update mongodb
		if _, err := c.UpdateOne(ctx, filter, update, nil); err != nil {
			log.Println(err)

			continue
		}

		log.Printf("success processing:\n\t{id: %s}\n", id)
	}
}

// validateTag
// will check the vehicle validation of that image.
func validateTag(response *imagga.Response) bool {
	for _, item := range response.Result.Tags {
		if item.Tag.En == "vehicle" {
			return true
		}
	}

	return false
}
