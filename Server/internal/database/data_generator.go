package database

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/ykkalexx/recommendation-system/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
    numOfUsers      = 50
    numOfItems      = 10
    numOfBehaviours = 150
)

var actions = []string{"view", "like", "comment", "share", "purchase"}

func GenerateFakeData(ctx context.Context, db *mongo.Database) error {
    // generate users
    users := make([]interface{}, numOfUsers)
    for i := 0; i < numOfUsers; i++ {
        users[i] = bson.M{"_id": fmt.Sprintf("user%d", i+1)}
    }
    _, err := db.Collection("users").InsertMany(ctx, users)
    if err != nil {
        return err
    }

    // generate Items
    items := make([]interface{}, numOfItems)
    for i := 0; i < numOfItems; i++ {
        items[i] = bson.M{"_id": fmt.Sprintf("item%d", i+1)}
    }
    _, err = db.Collection("items").InsertMany(ctx, items)
    if err != nil {
        return err
    }

    // generate behaviours
    behaviors := make([]interface{}, numOfBehaviours)
    for i := 0; i < numOfBehaviours; i++ {
        behaviors[i] = models.UserBehavior{
            UserID:    fmt.Sprintf("user%d", rand.Intn(numOfUsers)+1),
            ItemID:    fmt.Sprintf("item%d", rand.Intn(numOfItems)+1),
            Action:    actions[rand.Intn(len(actions))],
            Timestamp: time.Now().Add(-time.Duration(rand.Intn(7*24)) * time.Hour).Unix(),
        }
    }
    _, err = db.Collection("behaviors").InsertMany(ctx, behaviors)
    if err != nil {
        return err
    }

    return nil
}