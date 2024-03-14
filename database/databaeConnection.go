package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBinstance() *mongo.Client {
	MongoDb := "mongodb://localhost:27017"
	fmt.Print(MongoDb)

	// creates a new MongoDB client using the connection options specified in
	client, err := mongo.NewClient(options.Client().ApplyURI(MongoDb))

	if err != nil {
		log.Fatal(err)
	}

	// It sets up a context with a timeout of 10 seconds and attempts to connect to the MongoDB server.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("connected to mongodb")
	return client
}

var Client *mongo.Client = DBinstance()

// This function opens a specific collection within the MongoDB database.
func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {

	//    Here's how it works:
	// It accesses the "restaurant" database using the client.Database("restaurant") method.
	// It then accesses the specified collection within the database using Collection(collectionName).
	// Finally, it returns a pointer to the MongoDB collection.

	var collection *mongo.Collection = client.Database("restaurant").Collection(collectionName)
	return collection
}
