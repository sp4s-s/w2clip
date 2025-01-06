package db

// import . "goclip/db"
// import . "db/init"

import (
	"context"
	"log"
	"os"
	"time"

	// "encoding/json"
	// "github.com/99designs/gqlgen/client"
	// "github.com/99designs/gqlgen/client"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


// Context Initialization
var _ctx context.Context
func init() {
	_ctx, err := context.WithTimeout(context.Background(), 20*time.Second)
	if err != nil {
		log.Println("Error creating context")
	}
	if _ctx != nil {}

}

type usr struct {
	name string   `bson:"name"`
	id   int      `bson:"id"`
	uri  string   `bson:"uri"`
	msgs []string `bson:"msgs"`
}

type msg struct {
	id        int    `bson:"id"`
	data      string `bson:"data"`
	createdAt string `bson:"createdAt"`
}

func connect2DB() (*mongo.Client, error) {
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

func CreateUser(client *mongo.Client, db string, cl string, usrName string, dburi string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	col := client.Database(db).Collection(cl)

	var _id struct {
		id int `bson:"id"`
	}
	opts := options.FindOne().SetSort(bson.D{{Key: "id", Value: -1}})
	err := col.FindOne(ctx, bson.M{}, opts).Decode(&_id)
	if err != nil {
		_id.id = 0
	}
	id := _id.id + 1
	_, err = col.InsertOne(ctx, bson.M{"id": id, "userName": usrName, "uri": dburi})
	if err != nil {
		return 0, err
	}
	return id, nil

}

// log.Println("Error getting last id \n Error User Creation")
// panic(err)

func UpdateUser(client *mongo.Client, db string, cl string, id int, dturi string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	col := client.Database(db).Collection(cl)
	filter := bson.M{"id": id}

	updatedData := bson.M{"$set": bson.M{"uri": dturi}}

	_, err := col.UpdateOne(ctx, filter, updatedData)
	if err != nil {
		return err
	}
	return nil

}

func DelusrData(client *mongo.Client, db string, cl string, id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	col := client.Database(db).Collection(cl)
	filter := bson.M{"id": id}

	_, err := col.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func DeleteUserAllData(client *mongo.Client, db string, cl string, id_from int, id_2 int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	col := client.Database(db).Collection(cl)
	filter := bson.M{"id": bson.M{"$gte": id_from, "$lte": id_2}}

	_, err := col.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func DeleteUser(client *mongo.Client, db string, cl string, id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	col := client.Database(db).Collection(cl)
	filter := bson.M{"id": id}

	_, err := col.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func GetTotalUser(client *mongo.Client, db string, cl string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	col := client.Database(db).Collection(cl)

	count, err := col.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func PostData(client *mongo.Client, db string, cl string, id int, data string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	col := client.Database(db).Collection(cl)
	filter := bson.M{"id": id}

	updatedData := bson.M{"$set": bson.M{"data": data}}

	_, err := col.UpdateOne(ctx, filter, updatedData)
	if err != nil {
		return err
	}
	return nil

}
