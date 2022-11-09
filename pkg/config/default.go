package config

import (
	"github.com/ceit-aut/ad-registration-service/pkg/mqtt"
	"github.com/ceit-aut/ad-registration-service/pkg/service/imagga"
	"github.com/ceit-aut/ad-registration-service/pkg/service/mail"
	"github.com/ceit-aut/ad-registration-service/pkg/storage"
	"github.com/ceit-aut/ad-registration-service/pkg/storage/mongodb"
	"github.com/ceit-aut/ad-registration-service/pkg/storage/s3"
)

// Default
// loading default configs.
func Default() Config {
	return Config{
		Imagga:  imagga.Config{},
		Mailgun: mail.Config{},
		MQTT:    mqtt.Config{},
		Storage: storage.Config{
			Mongo: mongodb.Config{},
			S3:    s3.Config{},
		},
	}
}
