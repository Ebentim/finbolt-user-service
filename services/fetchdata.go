package services

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func FetchDocument(ctx context.Context, collection *mongo.Collection, uid string, result interface{}) error {
	return collection.FindOne(ctx, bson.M{"uid": uid}).Decode(result)
}
