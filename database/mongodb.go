package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func Connect(uri string) *mongo.Database {
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("MongoDB connection failed:", err)
	}

	log.Println("MongoDB connected!")

	return client.Database("course_selling_app")
}
