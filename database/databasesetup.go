package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DBset connects to mongodb
func DBSet() *mongo.Client {
	// Create a client
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27107"))

	if err != nil {
		log.Fatal(err)
	}

	// Context for DB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	// Connect client to DB
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Test connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("failed to connect to mongodb")
		return nil
	}

	fmt. Println("successfully connected to mongodb")
	return client
}

var Client *mongo.Client = DBSet()

func UserData(client *mongo.Client, CollectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("Ecommerce").Collection(CollectionName)
	return collection

}

func ProductData(client *mongo.Client, CollectionName string) *mongo.Collection {
	var productcollection *mongo.Collection = client.Database("Ecommerce").Collection(CollectionName)
	return productcollection
}