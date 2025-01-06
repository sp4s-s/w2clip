import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var db *mongo.Database
var usersCollection *mongo.Collection

type Msg struct {
	MsgID string    `bson:"msgId"`
	Msg   string    `bson:"msg"`
	Time  time.Time `bson:"time"`
}

type User struct {
	Name string `bson:"name"`
	ID   string `bson:"id"`
	URI  string `bson:"uri"`
	Msgs []Msg  `bson:"msg"`
}

func connectDB() {
	var err error
	client, err = mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	db = client.Database("socialApp")
	usersCollection = db.Collection("users")
}

func createPost(userID, newMsg string) {
	msg := Msg{
		MsgID: fmt.Sprintf("%d", time.Now().UnixNano()),
		Msg:   newMsg,
		Time:  time.Now(),
	}

	filter := bson.M{"id": userID}
	update := bson.M{
		"$push": bson.M{"msg": msg},
	}

	_, err := usersCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
}

func deletePost(userID, msgID string) {
	filter := bson.M{"id": userID, "msg.msgId": msgID}
	update := bson.M{
		"$pull": bson.M{"msg": bson.M{"msgId": msgID}},
	}

	_, err := usersCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
}

func deleteAllPosts(userID string) {
	filter := bson.M{"id": userID}
	update := bson.M{
		"$set": bson.M{"msg": []Msg{}},
	}

	_, err := usersCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
}

func readUserPosts(userID string, startID, endID int) {
	filter := bson.M{"id": userID}
	var user User
	err := usersCollection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Posts for user %s (ID: %s):\n", user.Name, user.ID)
	for i := startID; i <= endID && i < len(user.Msgs); i++ {
		fmt.Printf("%s: %s\n", user.Msgs[i].Time.Format(time.RFC3339), user.Msgs[i].Msg)
	}
}

func listAllUsers() {
	cursor, err := usersCollection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var user User
		err := cursor.Decode(&user)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("User: %s (ID: %s)\n", user.Name, user.ID)
	}
}

func getUserMsgs(userID string) {
	filter := bson.M{"id": userID}
	var user User
	err := usersCollection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Messages for user %s (ID: %s):\n", user.Name, user.ID)
	for _, msg := range user.Msgs {
		fmt.Printf("Message ID: %s, Time: %s, Msg: %s\n", msg.MsgID, msg.Time.Format(time.RFC3339), msg.Msg)
	}
}

func main() {
	connectDB()
	defer client.Disconnect(context.Background())

	// Example Usage
	createPost("user1", "New message here!")
	createPost("user2", "Another message here!")

	deletePost("user1", "1632968483")

	readUserPosts("user1", 0, 1)
	listAllUsers()
	getUserMsgs("user1")
	deleteAllPosts("user1")
}
