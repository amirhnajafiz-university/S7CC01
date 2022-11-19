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
	s3Sdk "github.com/aws/aws-sdk-go/service/s3"
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

	log.Println("processor started ...")

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

		log.Printf("receive id: %s\n", id)

		// finding the ad
		value := c.FindOne(ctx, filter, nil)
		if err := value.Decode(&ad); err != nil {
			log.Println(err)

			continue
		}

		log.Println("mongodb get by id succeed")

		// getting the image from s3
		svc := s3Sdk.New(h.S3.Session, &aws.Config{
			Region:   aws.String(h.S3.Cfg.Region),
			Endpoint: aws.String(h.S3.Cfg.Endpoint),
		})

		req, _ := svc.GetObjectRequest(&s3Sdk.GetObjectInput{
			Bucket: aws.String(h.S3.Cfg.Bucket),
			Key:    aws.String(ad.Id),
		})

		urlStr, err := req.Presign(15 * time.Minute)
		if err != nil {
			log.Println(err)

			continue
		}

		log.Println("s3 link generated")

		// image tag
		resp, err := h.Imagga.Process(urlStr)
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

					log.Printf("email send {id: %s}\n", id)
				}()
			} else {
				ad.State = enum.RejectState
			}
		} else {
			ad.State = enum.RejectState
		}

		log.Println("imagga api call succeed")

		// update mongodb
		if _, err := c.UpdateOne(ctx, filter, ad, nil); err != nil {
			log.Println(err)

			continue
		}

		log.Printf("success processing {id: %s}\n", id)
	}
}
