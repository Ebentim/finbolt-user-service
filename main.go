package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Ebentim/finbolt-user-service/controllers"
	"github.com/Ebentim/finbolt-user-service/services"

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

	c := services.CorsMiddleware()
	handler := c.Handler(r)

	srv := &http.Server{
		Addr:    PORT,
		Handler: handler,
	}

	// Run server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Error during server shutdown: %v", err)
	}
}
