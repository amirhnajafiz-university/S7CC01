package http

import (
	"net/http"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/ceit-aut/ad-registration-service/firstAPI/model"
	"github.com/ceit-aut/ad-registration-service/firstAPI/port/mqtt"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	pendingState = "pending"
	rejectState  = "rejected"
	acceptState  = "accepted"
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
	case pendingState:
		return ctx.SendString(pendingState)
	case rejectState:
		return ctx.SendString(rejectState)
	}

	return ctx.JSON(ad)
}

// HandlePostRequests
// get ad request and save it into mongodb and s3.
// after that send the id over rabbitMQ.
func (h *Handler) HandlePostRequests(ctx *fiber.Ctx) {

}
