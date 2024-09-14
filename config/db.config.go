package config

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBInstance() *mongo.Client {

	// Start Server
	err := godotenv.Load(".env")

	if err != nil {
		logrus.Fatalf("Error loading .env file")
	}

	MONGODB_URL := os.Getenv("MONGODB_URL")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(MONGODB_URL))
	if err != nil {
		// error
		log.Fatal(err)

	}

	logrus.Info("Connected to MongoDB successfully. Connection")

	return mongoClient

}

var CreatedMongoClient = DBInstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	return client.Database("go-social-auth").Collection(collectionName)
	// return client.Database("test").Collection(collectionName)
}
