package http

import (
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/ceit-aut/ad-registration-service/pkg/enum"
	"github.com/ceit-aut/ad-registration-service/pkg/model"
	"github.com/ceit-aut/ad-registration-service/pkg/mqtt"
	"github.com/ceit-aut/ad-registration-service/pkg/storage/s3"
	"github.com/gofiber/fiber/v2"
	"github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Handler
// manages to handle http endpoints.
type Handler struct {
	Mongo *mongo.Database
	MQTT  *mqtt.MQTT
	S3    *s3.S3
}

// HandleGetRequests
// return a json response of user ad.
func (h *Handler) HandleGetRequests(ctx *fiber.Ctx) error {
	var (
		id = ctx.FormValue("id")

		filter = bson.M{"_id": id}
		c      = h.Mongo.Collection(model.AdCollection, nil)

		ad model.Ad
	)

	value := c.FindOne(ctx.Context(), filter, nil)
	if err := value.Decode(&ad); err != nil {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	switch ad.State {
	case enum.PendingState:
		return ctx.SendString(enum.PendingState)
	case enum.RejectState:
		return ctx.SendString(enum.RejectState)
	}

	return ctx.JSON(ad)
}

// HandlePostRequests
// get ad request and save it into mongodb and s3.
// after that send the id over rabbitMQ.
func (h *Handler) HandlePostRequests(ctx *fiber.Ctx) error {
	var (
		c = h.Mongo.Collection(model.AdCollection, nil)

		email       = ctx.FormValue("email")
		description = ctx.FormValue("description")
	)

	image, err := ctx.FormFile("image")
	if err != nil {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	file, err := image.Open()
	if err != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	uploader := s3manager.NewUploader(h.S3.Session)
	up, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(h.S3.Bucket),
		Key:    aws.String("file"),
		Body:   file,
	})
	if err != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	ad := model.Ad{
		Email:       email,
		Description: description,
		State:       enum.PendingState,
		Category:    "",
		Image:       up.UploadID,
	}

	id, err := c.InsertOne(ctx.Context(), ad, nil)
	if err != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	err = h.MQTT.Channel.PublishWithContext(
		ctx.Context(),
		"",
		h.MQTT.Queue,
		false,
		false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(id.InsertedID.(string)),
		},
	)
	if err != nil {
		return ctx.SendStatus(http.StatusBadGateway)
	}

	return ctx.SendString(id.InsertedID.(string))
}
