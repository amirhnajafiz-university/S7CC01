package http

import (
	"log"
	"net/http"

	"github.com/ceit-aut/ad-registration-service/pkg/enum"
	"github.com/ceit-aut/ad-registration-service/pkg/model"
	"github.com/ceit-aut/ad-registration-service/pkg/mqtt"
	"github.com/ceit-aut/ad-registration-service/pkg/storage/s3"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
		// get ad id from form request
		id = ctx.FormValue("id")
		// creating a filter for mongodb
		filter = bson.M{"id": id}
		// connecting to mongodb collection
		c = h.Mongo.Collection(model.AdCollection, nil)
		// creating an add model
		ad model.Ad
	)

	// find the object in mongodb
	value := c.FindOne(ctx.Context(), filter, nil)
	if err := value.Decode(&ad); err != nil {
		log.Println(err)

		return ctx.SendStatus(http.StatusBadRequest)
	}

	// switch case on ad state
	switch ad.State {
	case enum.PendingState:
		return ctx.SendString(enum.PendingState)
	case enum.RejectState:
		return ctx.SendString(enum.RejectState)
	default:
		return ctx.JSON(ad)
	}
}

// HandlePostRequests
// get ad request and save it into mongodb and s3.
// after that send the id over rabbitMQ.
func (h *Handler) HandlePostRequests(ctx *fiber.Ctx) error {
	var (
		// get email and description from form value
		email       = ctx.FormValue("email")
		description = ctx.FormValue("description")
		// creating a new unique id
		uid = uuid.New().String()
		// connecting to mongodb collection
		c = h.Mongo.Collection(model.AdCollection, nil)
	)

	// get image file
	image, err := ctx.FormFile("image")
	if err != nil {
		log.Println(err)

		return ctx.SendStatus(http.StatusBadRequest)
	}

	// open image file
	file, err := image.Open()
	if err != nil {
		log.Println(err)

		return ctx.SendStatus(http.StatusInternalServerError)
	}

	// creating a new uploader
	uploader := s3manager.NewUploader(h.S3.Session)
	// upload image into s3 database
	up, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(h.S3.Bucket),
		Key:    aws.String(uid),
		Body:   file,
	})
	if err != nil {
		log.Println(err)

		return ctx.SendStatus(http.StatusInternalServerError)
	}

	// filling the ad model
	ad := model.Ad{
		Id:          uid,
		Email:       email,
		Description: description,
		State:       enum.PendingState,
		Category:    "",
		Image:       up.UploadID,
	}

	// insert ad into mongodb
	if _, err = c.InsertOne(ctx.Context(), ad, nil); err != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	// publish id over mqtt
	err = h.MQTT.Channel.PublishWithContext(
		ctx.Context(),
		"",
		h.MQTT.Queue,
		false,
		false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(uid),
		},
	)
	if err != nil {
		log.Println(err)

		return ctx.SendStatus(http.StatusBadGateway)
	}

	return ctx.SendString(uid)
}
