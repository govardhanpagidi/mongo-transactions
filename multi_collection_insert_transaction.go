package main

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

func TestMultiInsertTransactionCommit(dbName string, collName ...string) {
	var err error
	var client *mongo.Client
	var collection *mongo.Collection
	var ctx = context.Background()

	if client, err = getMongoClient(); err != nil {
		log.Fatal(err)
	}

	err = client.Database(dbName).Drop(ctx)
	if err != nil {
		log.Println(err)
	}
	db := client.Database(dbName)
	db.CreateCollection(ctx, collName[0])
	db.CreateCollection(ctx, collName[1])
	db.CreateCollection(ctx, collName[2])

	var session mongo.Session
	if session, err = client.StartSession(); err != nil {
		log.Fatal(err)
	}
	if err = session.StartTransaction(); err != nil {
		log.Fatal(err)
	}

	defer func() {
		if r := recover(); r != nil {
			log.Println("recovered from panic: ", r)
		}
		session.EndSession(ctx)
		err := client.Disconnect(ctx)
		if err != nil {
			log.Fatal(err)
			return
		}
	}()

	id := primitive.NewObjectID()
	var transDoc = bson.M{"_id": id, "hometown": "GeorgeTown", "year": int32(2024)}

	if err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		collection = client.Database(dbName).Collection(collName[0])
		if _, err = collection.InsertOne(sc, transDoc); err != nil {
			log.Fatal(err)
		}

		collection2 := client.Database(dbName).Collection(collName[1])
		if _, err = collection2.InsertOne(sc, transDoc); err != nil {
			log.Fatal(err)
		}

		collection3 := client.Database(dbName).Collection(collName[2])
		if _, err = collection3.InsertOne(sc, transDoc); err == nil {
			log.Println("Abort transaction")
			return errors.New("Simulated error in transaction even there is no error in transaction")
		}

		if err = session.CommitTransaction(sc); err != nil {
			log.Fatal(err)
		}
		return nil
	}); err != nil {
		err = session.AbortTransaction(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}
	session.EndSession(ctx)

	collection = client.Database(dbName).Collection(collName[2])

	var v bson.M
	if err = collection.FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(&v); err != nil {
		log.Fatal(err)
	}
	if v["year"] != int32(2024) {
		log.Fatal("expected 2024 but got", v["year"])
	} else {
		log.Println("expected 2024 and got", v["year"])
	}

}
