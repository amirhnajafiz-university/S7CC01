package cmd

import (
	"github.com/ceit-aut/ad-registration-service/internal/api/port/http"
	"github.com/ceit-aut/ad-registration-service/pkg/config"
	"github.com/ceit-aut/ad-registration-service/pkg/mqtt"
	"github.com/ceit-aut/ad-registration-service/pkg/storage/mongodb"
	"github.com/ceit-aut/ad-registration-service/pkg/storage/s3"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "api",
		Long: "starting the api service",
		Run: func(_ *cobra.Command, _ []string) {
			main()
		},
	}
}

// main method of api service.
func main() {
	// load configs
	cfg := config.Load()

	// mongodb connection
	mongo, err := mongodb.NewConnection(cfg.Storage.Mongo)
	if err != nil {
		panic(err)
	}

	// s3 connection
	s, err := s3.NewSession(cfg.Storage.S3)
	if err != nil {
		panic(err)
	}

	// rabbitmq connection
	mq, err := mqtt.NewConnection(cfg.MQTT)
	if err != nil {
		panic(err)
	}

	// creating a handler
	h := http.Handler{
		Mongo: mongo,
		MQTT:  mq,
		S3:    s,
	}

	// creating a new fiber app
	app := fiber.New()

	// declaring endpoints
	app.Get("api/{id}", h.HandleGetRequests)
	app.Post("api/", h.HandlePostRequests)

	// starting fiber
	if er := app.Listen(":5050"); er != nil {
		panic(er)
	}
}
