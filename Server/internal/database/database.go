package database

import (
	"context"
	"log"
	"time"

	"github.com/ykkalexx/recommendation-system/internal/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func Connect(uri string) (*mongo.Client, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    clientOptions := options.Client().ApplyURI(uri)
    var err error
    client, err = mongo.Connect(ctx, clientOptions)
    if err != nil {
        return nil, utils.NewAppError(500, "Failed to connect to MongoDB", err)
    }

    // ping db to check if connection is successful
    err = client.Ping(ctx, nil)
    if err != nil {
        return nil, utils.NewAppError(500, "Failed to ping MongoDB", err)
    }

    utils.InfoLogger.Println("Connected to MongoDB!")
    return client, nil
}

func GetCollection(database, collection string) (*mongo.Collection, error) {
    if client == nil {
        return nil, utils.NewAppError(500, "database client is not initialized", nil)
    }
    return client.Database(database).Collection(collection), nil
}

func Disconnect(client *mongo.Client) error {
    if client == nil {
        return nil
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    if err := client.Disconnect(ctx); err != nil {
        log.Printf("Error disconnecting from MongoDB: %v", err)
        return err
    }

    utils.InfoLogger.Println("Disconnected from MongoDB!")
    return nil
}


