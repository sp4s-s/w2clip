package db

type User struct {
	Name string   `bson:"name"`
	ID   int      `bson:"id"`
	URI  string   `bson:"uri"`
	Msgs []string `bson:"msgs"`
}

type Message struct {
	ID        int    `bson:"id"`
	Data      string `bson:"data"`
	CreatedAt string `bson:"createdAt"`
}
