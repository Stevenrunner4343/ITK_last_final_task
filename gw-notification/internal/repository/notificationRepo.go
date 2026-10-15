package repository

import (
	"context"

	"gw-notification/internal/storages/model"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type NotificationRepo struct {
	collection *mongo.Collection
}

func NewNotificationRepo(client *mongo.Client, dbName, collName string) *NotificationRepo {
	return &NotificationRepo{
		collection: client.Database(dbName).Collection(collName),
	}
}

func (r *NotificationRepo) SendToMongo(ctx context.Context, d model.Data) error {
	_, err := r.collection.InsertOne(ctx, d)
	return err
}
