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

var client *mongo.Client = DBSet()

func Userdata(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = (*mongo.Collection)(client.Database("Ecommerce").Collection(collectionName))
	return collection

}

func productData(client *mongo.Client, collectionName string) *mongo.Collection {
	var productCollection *mongo.Collection = (*mongo.Collection)(client.Database("Ecommerce").Collection(collectionName))
	return productCollection

}