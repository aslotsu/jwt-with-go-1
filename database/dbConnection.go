package database

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

func DbInstance() *mongo.Client {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err.Error())
	}

	uri := os.Getenv("MONGODB_URI")
	serverApiOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(uri).SetServerAPIOptions(serverApiOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	fmt.Println("Looks like it connected successfully")
	return client
}

var Client = DbInstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection = client.Database("users").Collection(collectionName)
	return collection
}
