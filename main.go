package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/Ebentim/finbolt-user-service/controllers"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error Loading ENV file")
	}

	MONGO_URI := os.Getenv("MONGO_URI")
	PORT := os.Getenv("PORT")

	if PORT == "" {
		PORT = ":5000"
	} else if PORT[0] != ':' {
		PORT = ":" + PORT
	}

	clientOptions := options.Client().ApplyURI(MONGO_URI)
	client, err := mongo.Connect(clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")

	db := client.Database("test")

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	r := mux.NewRouter()
	r.HandleFunc("/api/v1/create_userdata", controllers.CreateUserProfile(db)).Methods("POST")
	r.HandleFunc("/", controllers.ListAllUsers(db)).Methods("GET")
	r.HandleFunc("/api/v1/get_user", controllers.LoginUser(db)).Methods("POST")

	c := CorsMiddleware()
	handler := c.Handler(r)
	http.ListenAndServe(PORT, handler)
}
