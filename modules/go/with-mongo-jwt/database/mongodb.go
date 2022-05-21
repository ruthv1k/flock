package mongoconnect

import (
	"context"
	"log"
	"time"

	"github.com/ruthv1k/flock/modules/go/with-mongo-jwt/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var client *mongo.Client

// Starts a mongodb connection, returns a mongo client if valid URI is provided
//
// InitializeDbConnection(MONGO_URI)
func InitializeDbConnection(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))

	if err != nil {
		log.Println("Error connecting to db. Err: ", err.Error())
		return client, err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Println("Unable to ping to primary. Err: ", err.Error())
		return client, err
	}

	log.Println("Connected to db.")

	return client, nil
}

// Disconnects the connection with db
func DisconnectDb() {
	if client == nil {
		log.Fatal("Error accessing db client, db client is nil")
		return
	}

	if err := client.Disconnect(context.TODO()); err != nil {
		log.Fatal(err.Error())
		return
	}

	log.Println("Disconnecting from db.")
}

// Gets mongo client to access databases and collections
//
// GetClient()
func GetClient() *mongo.Client {
	if client != nil {
		return client
	}
	MONGO_URI, _ := utils.GetEnv("MONGODB_URI")

	client, dbErr := InitializeDbConnection(MONGO_URI)

	if dbErr != nil {
		log.Fatal(dbErr.Error())
	}

	return client
}

// Gets a specific database, returns a mongo database
//
// GetDatabase(databaseName)
func GetDatabase(dbName string) *mongo.Database {
	if client != nil {
		return client.Database(dbName)
	}

	return GetClient().Database(dbName)
}

// Gets a mongo collection, mongo context and cancel function from specified DB
//
// GetCollection(collectionName)
func GetCollection(collectionName string) (*mongo.Collection, context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	db := GetDatabase("go_auth")
	postsCollection := db.Collection(collectionName)

	return postsCollection, ctx, cancel
}
