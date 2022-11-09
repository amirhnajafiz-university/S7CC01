package http

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/ceit-aut/ad-registration-service/firstAPI/port/mqtt"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

// Handler
// manages to handle http endpoints.
type Handler struct {
	Mongo *mongo.Database
	MQTT  *mqtt.MQTT
	S3    *session.Session
}

// HandleGetRequests
// return a json response of user ad.
func (h *Handler) HandleGetRequests(ctx *fiber.Ctx) {

}

// HandlePostRequests
// get ad request and save it into mongodb and s3.
// after that send the id over rabbitMQ.
func (h *Handler) HandlePostRequests(ctx *fiber.Ctx) {

}
