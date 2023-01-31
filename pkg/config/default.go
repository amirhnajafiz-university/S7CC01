package config

import (
	"github.com/ceit-aut/S7CC01/pkg/mqtt"
	"github.com/ceit-aut/S7CC01/pkg/service/imagga"
	"github.com/ceit-aut/S7CC01/pkg/service/mail"
	"github.com/ceit-aut/S7CC01/pkg/storage"
	"github.com/ceit-aut/S7CC01/pkg/storage/mongodb"
	"github.com/ceit-aut/S7CC01/pkg/storage/s3"
)

// Default
// loading default configs.
func Default() Config {
	return Config{
		Imagga: imagga.Config{
			ApiKey:    "",
			ApiSecret: "",
		},
		Mailgun: mail.Config{
			Domain: "",
			APIKEY: "",
		},
		MQTT: mqtt.Config{
			Queue: "",
			URI:   "",
		},
		Storage: storage.Config{
			Mongo: mongodb.Config{
				Database: "",
				URI:      "",
			},
			S3: s3.Config{
				AccessKeyID:     "",
				SecretAccessKey: "",
				Region:          "",
				Bucket:          "",
				Endpoint:        "",
			},
		},
	}
}
