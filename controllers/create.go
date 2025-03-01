package controllers

import (
	"net/http"
	"sync"

	"github.com/Ebentim/finbolt-user-service/lib"
	"github.com/Ebentim/finbolt-user-service/services"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

func CreateUserProfile(db *mongo.Database) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := lib.CreateCollectionIfNotExist(db, "user_profile", &wg); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			if err := lib.CreateCollectionIfNotExist(db, "user_subscription", &wg); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			if err := lib.CreateCollectionIfNotExist(db, "user_account", &wg); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}()

		wg.Wait()

		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := services.CreateUserProfile(db, r); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}

			if err := services.CreateSubscriptionProfile(db, r); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}

			if err := services.CreateUserAccount(db, r); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
		}()
	}
}
