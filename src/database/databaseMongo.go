package database

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DataBaseMongo struct {
	client *mongo.Client
}

func (d *DataBaseMongo) Connect() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URL"))

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	d.client = client
	return client, nil
}
