package main

import (
	"fmt"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func convertRecipeBSONToSlice(bsonResult []bson.D) []recipeResponse {
	resultSlice := make([]recipeResponse, 0)
	var crtObject recipeResponse
	for _, v := range bsonResult {
		bsonBytes, _ := bson.Marshal(v)
		bson.Unmarshal(bsonBytes, &crtObject)
		resultSlice = append(resultSlice, crtObject)
	}
	return resultSlice
}

func queryCollectionAndReturnResult(collection string, filter, option interface{}) []bson.D {
	cursor, err := query(mgc.client, mgc.ctx, "mydb", collection, filter, option)
	if err != nil {
		panic(err)
	}
	var results []bson.D
	if err := cursor.All(mgc.ctx, &results); err != nil {
		panic(err)
	}
	return results
}

func mapReduceByCategoryAndStep(client *mongo.Client, ctx context.Context, dataBase, col string,
	query interface{}, mapFn, reduceFn string, step, category string) {
	fmt.Println(category + step)
	fmt.Println(mapFn)
	fmt.Println(reduceFn)
	par := bson.D{
		{"mapreduce", col},
		{"map", mapFn},
		{"reduce", reduceFn},
		{"out", category + step},
		{"query", query},
	}
	_ = client.Database(dataBase).RunCommand(nil, par)
}


