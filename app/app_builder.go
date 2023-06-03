package app

import (
	"errors"

	"github.com/arael34/qlogger/types"
	"go.mongodb.org/mongo-driver/mongo"
)

type AppBuilder struct {
	client *mongo.Client
	logger *types.QLogger
}

func NewAppBuilder() *AppBuilder {
	return &AppBuilder{}
}

func (ab *AppBuilder) WithClient(client *mongo.Client) *AppBuilder {
	ab.client = client
	return ab
}

func (ab *AppBuilder) WithLogger(logger *types.QLogger) *AppBuilder {
	ab.logger = logger
	return ab
}

func (ab *AppBuilder) Build() (*App, error) {
	if ab.client == nil ||
		ab.logger == nil {
		return nil, errors.New("failed to build app")
	}

	return &App{
		client: ab.client,
		logger: ab.logger,
	}, nil
}
