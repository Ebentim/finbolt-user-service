package lib

import (
	"context"
	"fmt"
	"slices"
	"sync"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func CreateCollectionIfNotExist(db *mongo.Database, collectionName string, wg *sync.WaitGroup) error {
	defer wg.Done()

	wg.Add(1)
	collections, err := db.ListCollectionNames(context.TODO(), bson.D{})
	if err != nil {
		return err
	}
	wg.Wait()

	exists := false

	exists = slices.Contains(collections, collectionName)

	if !exists {
		err = db.CreateCollection(context.TODO(), collectionName)
		if err != nil {
			return fmt.Errorf("collection %s could not be created", collectionName)
		}
	}

	return nil
}
