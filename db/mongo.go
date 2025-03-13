package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

func ConnectToDb(db_uri string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(db_uri)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
		return nil, err
	}

	// Check connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Could not ping MongoDB:", err)
		return nil, err
	}

	fmt.Println("Connected to MongoDB!")

	// Store the client in a global variable
	MongoClient = client

	return client, nil
}

// Disconnect MongoDB when the app shuts down
func DisconnectMongoDB() {
	if MongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := MongoClient.Disconnect(ctx)
		if err != nil {
			log.Fatal("Error disconnecting from MongoDB:", err)
		}
		fmt.Println("Disconnected from MongoDB")
	}
}
