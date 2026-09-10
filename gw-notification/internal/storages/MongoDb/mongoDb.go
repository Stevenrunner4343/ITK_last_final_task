package storages

import "go.mongodb.org/mongo-driver/v2/mongo"

type MongoDb struct {
	client *mongo.Client
	db     *mongo.Database
}

type MongoConfig struct {
}
