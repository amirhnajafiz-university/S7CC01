package http

import (
	"net/http"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/ceit-aut/ad-registration-service/pkg/enum"
	"github.com/ceit-aut/ad-registration-service/pkg/model"
	"github.com/ceit-aut/ad-registration-service/pkg/mqtt"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
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
func (h *Handler) HandlePostRequests(ctx *fiber.Ctx) {

}
