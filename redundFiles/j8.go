package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB client and collections
var client *mongo.Client
var usersCollection *mongo.Collection
var messagesCollection *mongo.Collection

// Initialize MongoDB connection
func connectMongoDB(uri, dbName string) {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}
	usersCollection = client.Database(dbName).Collection("users")
	messagesCollection = client.Database(dbName).Collection("messages")
	fmt.Println("Connected to MongoDB!")
}

// Disconnect MongoDB
func disconnectMongoDB() {
	if err := client.Disconnect(context.TODO()); err != nil {
		log.Fatalf("Error disconnecting MongoDB: %v", err)
	}
	fmt.Println("Disconnected from MongoDB!")
}

// User represents a document in the "users" collection
type User struct {
	ID        int       `bson:"id"`
	Name      string    `bson:"name"`
	StringVal string    `bson:"string_val"`
	CreatedAt time.Time `bson:"created_at"`
}

// Message represents a document in the "messages" collection
type Message struct {
	ID        int       `bson:"id"`
	UserID    int       `bson:"user_id"`
	Title     string    `bson:"title"`
	CreatedAt time.Time `bson:"created_at"`
}

// Get all users
func getAllUsers() ([]User, error) {
	cursor, err := usersCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	var users []User
	if err = cursor.All(context.TODO(), &users); err != nil {
		return nil, err
	}
	return users, nil
}

// Get all messages
func getAllMessages() ([]Message, error) {
	cursor, err := messagesCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	var messages []Message
	if err = cursor.All(context.TODO(), &messages); err != nil {
		return nil, err
	}
	return messages, nil
}

// Get users within a range
func getUsersInRange(startID, endID int) ([]User, error) {
	filter := bson.M{"id": bson.M{"$gte": startID, "$lte": endID}}
	cursor, err := usersCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	var users []User
	if err = cursor.All(context.TODO(), &users); err != nil {
		return nil, err
	}
	return users, nil
}

// Add a new user
func addUser(name, stringVal string) error {
	user := User{
		ID:        int(time.Now().UnixNano() % 1e9), // Example unique ID
		Name:      name,
		StringVal: stringVal,
		CreatedAt: time.Now(),
	}
	_, err := usersCollection.InsertOne(context.TODO(), user)
	return err
}

// Add a message for a user
func addMessage(userID int, title string) error {
	message := Message{
		ID:        int(time.Now().UnixNano() % 1e9), // Example unique ID
		UserID:    userID,
		Title:     title,
		CreatedAt: time.Now(),
	}
	_, err := messagesCollection.InsertOne(context.TODO(), message)
	return err
}

// Delete all messages for a user
func deleteMessagesByUser(userID int) error {
	_, err := messagesCollection.DeleteMany(context.TODO(), bson.M{"user_id": userID})
	return err
}

// Delete a user and their messages
func deleteUser(userID int) error {
	// Delete user
	_, err := usersCollection.DeleteOne(context.TODO(), bson.M{"id": userID})
	if err != nil {
		return err
	}
	// Delete user's messages
	return deleteMessagesByUser(userID)
}

// Update a user's string value
func updateUserString(userID int, newString string) error {
	filter := bson.M{"id": userID}
	update := bson.M{"$set": bson.M{"string_val": newString}}
	_, err := usersCollection.UpdateOne(context.TODO(), filter, update)
	return err
}

// Update a message title
func updateMessageTitle(messageID int, newTitle string) error {
	filter := bson.M{"id": messageID}
	update := bson.M{"$set": bson.M{"title": newTitle}}
	_, err := messagesCollection.UpdateOne(context.TODO(), filter, update)
	return err
}

// Main function
func main() {
	// Replace with your MongoDB URI and database name
	connectMongoDB("mongodb://localhost:27017", "testdb")
	defer disconnectMongoDB()

	// Example operations
	if err := addUser("John Doe", "example-string"); err != nil {
		log.Fatalf("Error adding user: %v", err)
	}
	users, _ := getAllUsers()
	fmt.Println("Users:", users)
}