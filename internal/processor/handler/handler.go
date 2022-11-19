package handler

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ceit-aut/ad-registration-service/pkg/enum"
	"github.com/ceit-aut/ad-registration-service/pkg/model"
	"github.com/ceit-aut/ad-registration-service/pkg/mqtt"
	"github.com/ceit-aut/ad-registration-service/pkg/service/imagga"
	"github.com/ceit-aut/ad-registration-service/pkg/service/mail"
	"github.com/ceit-aut/ad-registration-service/pkg/storage/s3"

	"github.com/aws/aws-sdk-go/aws"
	s3sdk "github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
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

		// getting the image from s3 by creating a local file
		file, err := os.Create(id)
		if err != nil {
			log.Printf("cannot create file:\n\t%s\n", id)

			continue
		}

		// creating a new downloader
		downloader := s3manager.NewDownloader(h.S3.Session)
		bytes, err := downloader.Download(
			file,
			&s3sdk.GetObjectInput{
				Bucket: aws.String(h.S3.Cfg.Bucket),
				Key:    aws.String(id),
			},
		)
		if err != nil {
			log.Printf("failed to get file from s3: %v\n", err)

			continue
		}

		log.Printf("s3 file read:\n\t%d bytes\n", bytes)

		// uploading image to imagga
		uploadResp, err := h.Imagga.Upload(id)
		if err != nil {
			log.Printf("upload file to imagga failed: %v\n", err)

			continue
		}

		log.Printf("imagga upload file:\n\t%s\n", uploadResp.Status.Type)

		// image tag process
		resp, err := h.Imagga.Process(uploadResp.Result.UploadId)
		if err != nil {
			log.Println(err)

			continue
		}

		log.Printf("imagga response:\n\t%d tags\n", len(resp.Result.Tags))

		// response check
		if len(resp.Result.Tags) > 0 {
			log.Printf("imagga response confidence: %d\n", resp.Result.Tags[0].Confidence)

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

		// remove the tmp file
		if er := os.RemoveAll(id); er != nil {
			log.Printf("failed to remove file: %s\n", er)
		}

		log.Printf("success processing {id: %s}\n", id)
	}
}
