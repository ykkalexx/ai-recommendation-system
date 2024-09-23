package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func Connect(uri string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    clientOptions := options.Client().ApplyURI(uri)
    var err error
    client, err = mongo.Connect(ctx, clientOptions)
    if err != nil {
        return err
    }

    // ping db to check if connection is successful
    err = client.Ping(ctx, nil)
    if err != nil {
        return err
    }

    log.Println("Connected to MongoDB")
    return nil
}

func GetCollection(database, collection string) *mongo.Collection {
    return client.Database(database).Collection(collection)
}

func Disconnect() error {
    if client == nil {
        return nil
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    if err := client.Disconnect(ctx); err != nil {
        log.Printf("Error disconnecting from MongoDB: %v", err)
        return err
    }

    log.Println("Disconnected from MongoDB")
    return nil
}