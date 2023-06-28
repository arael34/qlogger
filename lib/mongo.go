package lib

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type mongoWrapper struct {
	collection *mongo.Collection
}

func (m mongoWrapper) Write(ctx context.Context, log LogSchema) error {
	ctx, cancel := context.WithTimeout(
		ctx,
		time.Duration(15*time.Second),
	)
	defer cancel()

	_, err := m.collection.InsertOne(ctx, log)
	return err
}

func NewMongoDatabase(collection *mongo.Collection) Database {
	return mongoWrapper{collection: collection}
}
