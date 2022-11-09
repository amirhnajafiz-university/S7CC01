package storage

import (
	"github.com/ceit-aut/ad-registration-service/pkg/storage/mongodb"
	"github.com/ceit-aut/ad-registration-service/pkg/storage/s3"
)

type Config struct {
	Mongo mongodb.Config `koanf:"mongodb"`
	S3    s3.Config      `koanf:"amazon"`
}
