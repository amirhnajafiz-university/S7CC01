package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewConnection
// opens a new connection to mongodb database.
func NewConnection(cfg Config) (*mongo.Database, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.URI))
	if err != nil {
		return nil, fmt.Errorf("mongoDB connection failed: %w", err)
	}

	return client.Database(cfg.Database), nil
}
