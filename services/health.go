package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Connection tests the mongodb connection health
var Health = func(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := Connect(ctx)
	if err != nil {
		writeStatus(w, http.StatusBadRequest, err)
		return
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Printf("error disconnecting client: %v", err)
		}

		cancel()
	}()

	coll := client.Database("myDB").Collection("favorite_books")
	title := "Lucy"
	var result bson.M
	err = coll.FindOne(context.TODO(), bson.D{{"title", title}}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		fmt.Printf("No document was found with the title %s\n", title)
		return
	}
	if err != nil {
		panic(err)
	}
	jsonData, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", jsonData)

	// write back okay for good connection
	writeStatus(w, http.StatusOK, nil)
}

func writeStatus(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)
	w.Write([]byte("505 - something connected!"))
}
