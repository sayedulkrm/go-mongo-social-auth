package config

import "go.mongodb.org/mongo-driver/mongo"

var UserCollection *mongo.Collection = OpenCollection(CreatedMongoClient, "user")
