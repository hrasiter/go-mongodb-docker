package services

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Connect returns a successfully connected mongo client
func Connect(ctx context.Context) (*mongo.Client, error) {
	url := os.Getenv("MONGODB_URI")
	log.Printf("connect to mongodb at: %v", url)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		msg := fmt.Sprintf("invalid mongodb connection: %v", err)
		log.Println(msg)
		return nil, fmt.Errorf(msg)
	}

	// test connection
	err = client.Ping(ctx, readpref.Primary())

	if err != nil {
		log.Printf("connection test to mongodb client failed: %v", err)
		return nil, fmt.Errorf(fmt.Sprintf("connection test to mongodb client failed: %v", err))
	}

	log.Println("connected successfully to mongodb")

	return client, nil
}
