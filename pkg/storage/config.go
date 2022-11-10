package storage

import (
	"github.com/ceit-aut/ad-registration-service/pkg/storage/mongodb"
	"github.com/ceit-aut/ad-registration-service/pkg/storage/s3"
)

// Config
// contains data for database connections.
type Config struct {
	Mongo mongodb.Config `koanf:"mongodb"`
	S3    s3.Config      `koanf:"s3"`
}
