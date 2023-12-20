package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getMongoClient() (*mongo.Client, error) {
	var err error
	var client *mongo.Client

	opts := options.Client()
	opts.ApplyURI(connString)
	opts.SetMaxPoolSize(10)

	if client, err = mongo.Connect(context.Background(), opts); err != nil {
		return client, err
	}
	client.Ping(context.Background(), nil)
	return client, err
}
