package mongodb

import (
	"Rsbot_only/internal/config"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoClient() (client *mongo.Client, err error) {
	uri := config.Instance.Mongo
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	return client, err
}
