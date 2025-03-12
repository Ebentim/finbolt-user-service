package services

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// UserAggregationOptions contains options for the user aggregation pipeline
type UserAggregationOptions struct {
	UID   string
	Skip  int
	Limit int
}

func CreateUserAggregationPipeline(opts UserAggregationOptions) mongo.Pipeline {
	pipeline := mongo.Pipeline{
		// Add pagination early
		{{Key: "$skip", Value: opts.Skip}},
		{{Key: "$limit", Value: opts.Limit}},

		// Lookups
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "user_accounts"},
			{Key: "localField", Value: "uid"},
			{Key: "foreignField", Value: "uid"},
			{Key: "as", Value: "account"},
		}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "user_subscriptions"},
			{Key: "localField", Value: "uid"},
			{Key: "foreignField", Value: "uid"},
			{Key: "as", Value: "subscription"},
		}}},

		// Unwind arrays (with preservation)
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$account"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$subscription"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},

		// Simplified projection
		{{Key: "$project", Value: bson.D{
			{Key: "us_pro", Value: "$$ROOT"},
			{Key: "us_acc", Value: "$account"},
			{Key: "us_sub", Value: "$subscription"},
		}}},
	}

	// Add UID filter if specified
	if opts.UID != "" {
		matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "uid", Value: opts.UID}}}}
		pipeline = append([]bson.D{matchStage}, pipeline...)
	}

	return pipeline
}

func FetchDocument(ctx context.Context, collection *mongo.Collection, uid string, result any) error {
	return collection.FindOne(ctx, bson.M{"uid": uid}).Decode(result)
}
