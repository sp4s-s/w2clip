package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect2DB(mongoUri string) (*mongo.Client, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUri))
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

	return client, nil
}

func CloseMongoDB(client *mongo.Client) {
	if err := client.Disconnect(context.TODO()); err != nil {
		log.Fatal(err)
	}
}

func ConnectionPractice(client *mongo.Client) {
	err := client.Ping(context.TODO(), nil)
	if err != nil {
		println("Failed to ping MongoDB: %w", err)
			panic(err)
	}

	// Change to db config
	col := client.Database("Clip").Collection("users")
	if col == nil {
		log.Println("Error getting collection")
	}
	if col != nil {
		log.Println("Collection found")
	}
}
