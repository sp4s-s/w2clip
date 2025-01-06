package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetAll(collection *mongo.Collection, result interface{}) error {
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return err
	}
	if err = cursor.All(context.TODO(), result); err != nil {
		return err
	}
	return nil
}

func GetInRange(collection *mongo.Collection, result interface{}, startID, endID int) error {
	filter := bson.M{"id": bson.M{"$gte": startID, "$lte": endID}}
	return Get(collection, result, filter)
}

func Get(collection *mongo.Collection, result interface{}, filter interface{}) error {
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return err
	}
	if err = cursor.All(context.TODO(), result); err != nil {
		return err
	}
	return nil
}

func Add(collection *mongo.Collection, document interface{}) error {
	_, err := collection.InsertOne(context.TODO(), document)
	return err
}

func DelAll(collection *mongo.Collection, filter interface{}) error {
	_, err := collection.DeleteMany(context.TODO(), filter)
	return err
}

func DelOne(collection *mongo.Collection, filter interface{}) error {
	_, err := collection.DeleteOne(context.TODO(), filter)
	return err
}

func Update(collection *mongo.Collection, filter, update interface{}) error {
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	return err
}
