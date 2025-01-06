package db

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

import 	"goclip/db/Types"

// Context Initialization
var ctx context.Context

func init() {
	ctx, err := context.WithTimeout(context.Background(), 20*time.Second)
	if err != nil {
		log.Println("Error creating context")
	}
	if ctx != nil {
	}

}

func Connect2DB() (*mongo.Client, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("ENV FIlE MISSING")
	}

	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Println("URI is missing")
		log.Println("\t check mongodb go driver documentation")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Println("Error connecting to mongo DB")
		// log.Println(err)
		panic(err)
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	col := client.Database("Clip").Collection("users")
	if col == nil {
		log.Println("Error getting collection")
	}
	if col != nil {
		log.Println("Collection found")
	}
	return client, nil
}
