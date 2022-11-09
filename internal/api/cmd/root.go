package cmd

import (
	"github.com/ceit-aut/ad-registration-service/internal/api/http"
	"github.com/ceit-aut/ad-registration-service/pkg/config"
	"github.com/ceit-aut/ad-registration-service/pkg/mqtt"
	"github.com/ceit-aut/ad-registration-service/pkg/storage/mongodb"
	s32 "github.com/ceit-aut/ad-registration-service/pkg/storage/s3"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "api",
		Long: "command for starting the api service",
		Run: func(_ *cobra.Command, args []string) {
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
	s3, err := s32.NewSession(cfg.Storage.S3)
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
		S3:    s3,
	}

	// creating a new fiber app
	app := fiber.New()

	// declaring endpoints
	app.Get("/{id}", h.HandleGetRequests)
	app.Post("/", h.HandlePostRequests)

	// starting fiber
	if err := app.Listen(":5050"); err != nil {
		panic(err)
	}
}
