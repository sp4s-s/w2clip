package db

import "time"

type User struct {
	ID        int       `bson:"id"`
	Name      string    `bson:"name"`
	StringVal string    `bson:"string_val"`
	CreatedAt time.Time `bson:"created_at"`
}

type Message struct {
	ID        int       `bson:"id"`
	UserID    int       `bson:"user_id"`
	Title     string    `bson:"title"`
	CreatedAt time.Time `bson:"created_at"`
}