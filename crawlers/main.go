package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// this method closes mongoDB connection and cancel context
func close(client *mongo.Client, ctx context.Context, cancel context.CancelFunc) {
	defer cancel()
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

// mongo.Client will be used for further database operation.
// context.Context will be used set deadlines for process.
// context.CancelFunc will be used to cancel context and
// resource associtated with it.
func connect(uri string) (*mongo.Client, context.Context, context.CancelFunc, error) {
	var cred options.Credential
	cred.Username = "admin"
	cred.Password = "fridgemanager"
	ctx, cancel := context.WithTimeout(context.Background(), 2500 * time.Second)
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

type recipe struct {
	url string
	totalTime string
	prepTime  string
	cookTime  string
	servings int
	calories int
	fat int
	carbs int
	fiber int
	sugar int
	protein int
	ingredients []string
}

func createRecipeFromURL(tastyURL string) recipe {
	var myRecipe recipe = recipe{}
	myRecipe.url = tastyURL
	myRecipe.ingredients = make([]string, 0)
	c := colly.NewCollector(
		colly.AllowedDomains("tasty.co", "www.tasty.co"),
	)

	c.OnHTML("div.recipe-time", func(e *colly.HTMLElement) {
		switch recipeDetail := e.ChildText("h5"); recipeDetail {
		case "Total Time":
			myRecipe.totalTime = e.ChildText("p:nth-child(3)")
		case "Prep Time":
			myRecipe.prepTime = e.ChildText("p:nth-child(3)")
		case "Cook Time":
			myRecipe.cookTime = e.ChildText("p:nth-child(3)")
		}
	})
	c.OnHTML("p.servings-display", func(e *colly.HTMLElement) {
		myRecipe.servings, _ = strconv.Atoi(strings.Split(e.Text, " ")[1])
	})
	c.OnHTML("li.ingredient", func(e *colly.HTMLElement) {
		myRecipe.ingredients = append(myRecipe.ingredients, e.Text)
	})
	c.OnHTML("div.nutrition-details", func(e * colly.HTMLElement) {
		var nutritionValues []string = strings.Split(e.ChildText("ul > span"), " ")
		myRecipe.calories, _ = strconv.Atoi(nutritionValues[0])
		myRecipe.fat, _ = strconv.Atoi(nutritionValues[1][: len(nutritionValues[1]) - 1])
		myRecipe.carbs, _ = strconv.Atoi(nutritionValues[2][: len(nutritionValues[2]) - 1])
		myRecipe.fiber, _ = strconv.Atoi(nutritionValues[3][: len(nutritionValues[3]) - 1])
		myRecipe.sugar, _ = strconv.Atoi(nutritionValues[4][: len(nutritionValues[4]) - 1])
		myRecipe.protein, _ = strconv.Atoi(nutritionValues[5][: len(nutritionValues[5]) - 1])
	})
	c.Visit(tastyURL)
	return myRecipe
}

func getAllRecipesFromInputFile(ctx context.Context, client *mongo.Client) {
	file, err := os.Open("recipes_links.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		crtRecipe := createRecipeFromURL(scanner.Text())
		// fmt.Println(scanner.Text())
		// fmt.Println(crtRecipe)
		go insertRecipeIntoDB(crtRecipe, ctx, client)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func insertRecipeIntoDB(crtRecipe recipe, ctx context.Context, client *mongo.Client) {
	// doc, err := toDoc(crtRecipe)
	var doc interface{}
	doc = bson.D {
		{"url", crtRecipe.url},
		{"totalTime", crtRecipe.totalTime},
		{"cookTime", crtRecipe.cookTime},
		{"prepTime", crtRecipe.prepTime},
		{"servings", crtRecipe.servings},
		{"calories", crtRecipe.calories},
		{"fat", crtRecipe.fat},
		{"carbs", crtRecipe.carbs},
		{"fiber", crtRecipe.fiber},
		{"sugar", crtRecipe.sugar},
		{"protein", crtRecipe.protein},
		{"ingredients", crtRecipe.ingredients},
	}
	insertOneResult, err := insertOne(client, ctx, "mydb", "recipes", doc)
	if err != nil {
		//panic(err)
		fmt.Println(err)
	} else {
		fmt.Println(insertOneResult.InsertedID)
	}
}

func toDoc(v interface{}) (doc *bson.D, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}
	err = bson.Unmarshal(data, &doc)
	return
}

func main() {
	client, ctx, cancel, err := connect("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}
	defer close(client, ctx, cancel)
	getAllRecipesFromInputFile(ctx, client)
	
	// ping(client, ctx)

	// insert######################################################
	// var document interface{}
	// document = bson.D {
	// 	{"url", "https://tasyt.ksbfskfb"},
	// 	{"rollNo", 175},
	// 	{"maths", 80},
	// }
	// insertOneResult, err := insertOne(client, ctx, "mydb", "recipes", document)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("Result of insert one")
	// fmt.Println(insertOneResult.InsertedID)


	// find##########################################################
	// var filter, option interface{}
	// filter = bson.D {
	// 	{"maths", bson.D{{"$gt", 70}}},
	// }
	// option = bson.D {{"_id", 0}}
	// cursor, err := query(client, ctx, "mydb", "recipes", filter, option)
	// if err != nil {
	// 	panic(err)
	// }
	// var results []bson.D
	// if err := cursor.All(ctx, &results); err != nil {
	// 	panic(err)
	// }
	// fmt.Println("Query Result")
	// for _, doc := range results {
	// 	fmt.Println(doc)
	// }



	//updateOne#########################################################
	// filter := bson.D {
	// 	{"maths", bson.D{{"$lt", 100}}},
	// }
	// update := bson.D {
	// 	{"$set", bson.D {
	// 		{"maths", 100},
	// 	}},
	// }
	// result, err := UpdateOne(client, ctx, "mydb", "recipes")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(result.ModifiedCount)

	//deleteOne#########################################################
	// query := bson.D {
	// 	{"maths", bson.D{{"$gt", 60}}},
	// }
	// result, err := deleteOne(client, ctx, "mydb", "recipes", query)
	// fmt.Println(result.DeletedCount)
}
