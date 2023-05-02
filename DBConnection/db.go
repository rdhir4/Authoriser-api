package dbConnection

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"


)



// Connection URI
const val = "Trump@123"
const uri string = "mongodb+srv://rdhir4:Trump%40123@usercluster.xy48btj.mongodb.net/?retryWrites=true&w=majority"

var connectionPool = make(map[string]*mongo.Client)
var name = "Raghav"

func DbConnection() *mongo.Client {
	// Create a new client and connect to the server
	var client *mongo.Client
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		panic(err)
	}
	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	return client
}
func GetPool() *mongo.Client {
	var client *mongo.Client

	if _, exist := connectionPool[name]; exist {
		client = connectionPool[name]
	} else {
		client = DbConnection()
		connectionPool[name] = client
	}
	return client
}
