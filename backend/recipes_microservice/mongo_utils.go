package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type mongoConnection struct {
	client     *mongo.Client
	ctx        context.Context
	cancelFunc context.CancelFunc
}

var mgc mongoConnection

func createMongoConnection() *mongoConnection {
	client, ctx, cancel, err := connect("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}
	mgc.client = client
	mgc.ctx = ctx
	mgc.cancelFunc = cancel
	return &mgc
}

func getMongoConnection() *mongoConnection {
	return &mgc
}

func closeMongoConnection() {
	close(mgc.ctx, mgc.client, mgc.cancelFunc)
}

func close(ctx context.Context, client *mongo.Client, cancel context.CancelFunc) {
	defer cancel()
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

func connect(uri string) (*mongo.Client, context.Context, context.CancelFunc, error) {
	var cred options.Credential
	cred.Username = "admin"
	cred.Password = "fridgemanager"
	ctx, cancel := context.WithTimeout(context.Background(), 2500*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri).SetAuth(cred))
	return client, ctx, cancel, err
}

func ping(client *mongo.Client, ctx context.Context) error {
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	fmt.Println("connected successfully")
	return nil
}

func insertOne(client *mongo.Client, ctx context.Context,
	database, col string, doc interface{}) (*mongo.InsertOneResult, error) {
	collection := client.Database(database).Collection(col)
	result, err := collection.InsertOne(ctx, doc)
	return result, err
}

func insertMany(client *mongo.Client, ctx context.Context,
	database, col string, docs []interface{}) (*mongo.InsertManyResult, error) {
	collection := client.Database(database).Collection(col)
	result, err := collection.InsertMany(ctx, docs)
	return result, err
}

func query(client *mongo.Client, ctx context.Context, dataBase, col string,
	query, field interface{}) (result *mongo.Cursor, err error) {
	collection := client.Database(dataBase).Collection(col)
	result, err = collection.Find(ctx, query, options.Find().SetProjection(field))
	return
}

func updateOne(client *mongo.Client, ctx context.Context, dataBase, col string,
	filter, update interface{}) (result *mongo.UpdateResult, err error) {
	collection := client.Database(dataBase).Collection(col)
	result, err = collection.UpdateOne(ctx, filter, update)
	return
}

func updateMany(client *mongo.Client, ctx context.Context, dataBase, col string,
	filter, update interface{}) (result *mongo.UpdateResult, err error) {
	collection := client.Database(dataBase).Collection(col)
	result, err = collection.UpdateMany(ctx, filter, update)
	return
}

func deleteOne(client *mongo.Client, ctx context.Context, dataBase, col string,
	query interface{}) (result *mongo.DeleteResult, err error) {
	collection := client.Database(dataBase).Collection(col)
	result, err = collection.DeleteOne(ctx, query)
	return
}

func deleteMany(client *mongo.Client, ctx context.Context, dataBase, col string,
	query interface{}) (result *mongo.DeleteResult, err error) {
	collection := client.Database(dataBase).Collection(col)
	result, err = collection.DeleteMany(ctx, query)
	return
}
