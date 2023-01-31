package storage

import (
	"github.com/ceit-aut/S7CC01/pkg/storage/mongodb"
	"github.com/ceit-aut/S7CC01/pkg/storage/s3"
)

// Config
// contains data for database connections.
type Config struct {
	Mongo mongodb.Config `koanf:"mongodb"`
	S3    s3.Config      `koanf:"s3"`
}
