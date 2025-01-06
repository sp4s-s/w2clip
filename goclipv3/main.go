package main

import (
	"context"
	_ "github.com/pingstyy/wclip/db"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client             *mongo.Client
	col1               *mongo.Collection
)

func main() {
	initMongoDB()
}


func initMongoDB() {
	dbName := "clipCluster"
	uri := "mongodb+srv://pingz1:LavillaFormatDebug@clipcluster.qm8vk.mongodb.net/?retryWrites=true&w=majority&appName=clipCluster"
	clientOptions := options.Client().ApplyURI(uri)
	var err error
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	col1 = client.Database(dbName).Collection("users")
	userData := bson.M{
		"name":  "Nirman ah Doe",
		"age":   13,
		"email": "nirman.doe@example.com",
	}

	// Insert the user data
	result, err := col1.InsertOne(context.TODO(), userData)
	if err != nil {
		log.Fatal("Error inserting user:", err)
	}

	// Get the inserted ID
	insertedID := result.InsertedID

	// Find the inserted user
	var foundUser bson.M
	err = col1.FindOne(context.TODO(), bson.M{"_id": insertedID}).Decode(&foundUser)
	if err != nil {
		log.Fatal("Error finding inserted user:", err)
	}

	println("Inserted User:", foundUser)
	closeMongoDB(client)

}

func closeMongoDB(client *mongo.Client) {
	if err := client.Disconnect(context.TODO()); err != nil {
		log.Fatal(err)
	} else {
		println("MongoDB Connection Closed")
	}
}

func envUri() string {
	if err := godotenv.Load(); err != nil {
		log.Println("ENV FIlE MISSING")
	}

	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Println("URI is missing")
		log.Println("\t check mongodb go driver documentation")
	}
	return uri
}
