package services

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

var Fill = func(w http.ResponseWriter, r *http.Request) {
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
	docs := []interface{}{
		bson.D{{"title", "My Brilliant Friend"}, {"author", "Elena Ferrante"}, {"year_published", 2012}},
		bson.D{{"title", "Lucy"}, {"author", "Jamaica Kincaid"}, {"year_published", 2002}},
		bson.D{{"title", "Cat's Cradle"}, {"author", "Kurt Vonnegut Jr."}, {"year_published", 1998}},
	}
	result, err := coll.InsertMany(context.TODO(), docs)
	list_ids := result.InsertedIDs
	fmt.Printf("Documents inserted: %v\n", len(list_ids))
	for _, id := range list_ids {
		fmt.Printf("Inserted document with _id: %v\n", id)
	}

	// write back okay for good connection
	writeStatus(w, http.StatusOK, nil)
}
