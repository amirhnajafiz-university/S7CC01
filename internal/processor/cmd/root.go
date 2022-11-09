package cmd

import (
	"github.com/ceit-aut/ad-registration-service/internal/processor/handler"
	"github.com/ceit-aut/ad-registration-service/pkg/config"
	"github.com/ceit-aut/ad-registration-service/pkg/mqtt"
	"github.com/ceit-aut/ad-registration-service/pkg/storage/mongodb"
	"github.com/ceit-aut/ad-registration-service/pkg/storage/s3"

	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "processor",
		Long: "starting the processor service",
		Run: func(_ *cobra.Command, _ []string) {
			main()
		},
	}
}

// main method of processor.
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

	// creating a new handler
	h := handler.Handler{
		Mongo: mongo,
		MQTT:  mq,
		S3:    s,
	}

	// start processing
	h.Handle()
}
