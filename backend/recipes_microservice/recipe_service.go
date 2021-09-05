package main

import (
	"fmt"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
)

type RecipeInfoService interface {
	FindRecipeByURL(url string) (recipeResponse, error)
	FindRecipesByNutritionRange(rangeType string, lowerBound, upperBound int) ([]recipeResponse, error)
	MapReduceByCategory(step, category string) ([]mapReduceResponse, error)
	FilterByIngredients(ingredients []string) ([]recipeResponse, error)
	GetAffordableRecipes(ingredients []string) ([]recipeResponse, error)
}

var DbNotFound = errors.New("Query returned 0 results")

type recipeInfoService struct {}

func (recipeInfoService) FindRecipeByURL(url string) (recipeResponse, error) {
	var filter, option interface{}
	filter = bson.D{{"url", url}}
	option = bson.D{{"_id", 0}}
	results := queryCollectionAndReturnResult("recipes", filter, option)
	if len(results) > 0 {
		var recipe recipeResponse
		bsonBytes, _ := bson.Marshal(results[0])
		bson.Unmarshal(bsonBytes, &recipe)
		return recipe, nil
	}
	return recipeResponse{}, DbNotFound
}

func (recipeInfoService) FindRecipesByNutritionRange(rangeType string, lowerBound, upperBound int) ([]recipeResponse, error) {
	var filter, option interface{}
	conditions := [2]bson.D {
		bson.D{{rangeType, bson.D{{"$gt", lowerBound}}}},
		bson.D{{rangeType, bson.D{{"$lte", upperBound}}}},
	}
	filter = bson.D{{"$and", conditions}}
	option = bson.D{{"_id", 0}}
	results := queryCollectionAndReturnResult("recipes", filter, option)
	recipes := convertRecipeBSONToSlice(results)
	return recipes, nil
}

func (recipeInfoService) MapReduceByCategory(step, category string) ([]mapReduceResponse, error) {
	var filter interface{} = bson.D{{category, bson.D{{"$gt", 0}}}}
	mapFn := fmt.Sprintf(" function() { emit(Math.floor(this.%s / %s), 1); }", category, step)
	reduceFn := " function(key, arr) { return Array.sum(arr); } "
	mapReduceByCategoryAndStep(mgc.client, mgc.ctx, "mydb", "recipes", filter, mapFn, reduceFn, step, category)

	results := queryCollectionAndReturnResult(category + step, bson.M{}, bson.M{})
	aggRes := make([]mapReduceResponse, 0)
	for _, v := range results {
		var crtAgg mapReduceResponse
		bsonBytes, _ := bson.Marshal(v)
		bson.Unmarshal(bsonBytes, &crtAgg)
		aggRes = append(aggRes, crtAgg)
	}
	return aggRes, nil
}

func (recipeInfoService) FilterByIngredients(ingredients []string) ([]recipeResponse, error) {
	var filter, option interface{}
	conditions := make([]bson.D, 0)
	for _, ingredient := range ingredients {
		conditions = append(conditions, bson.D{{"ingredients", bson.D{{"$regex", ingredient}} }})
	}
	filter = bson.D{{"$and", conditions}}
	option = bson.D{{"_id", 0}}
	results := queryCollectionAndReturnResult("recipes", filter, option)
	recipes := convertRecipeBSONToSlice(results)
	return recipes, nil
}

func (recipeInfoService) GetAffordableRecipes(ingredients []string) ([]recipeResponse, error) {
	var filter, option interface{}
	regex := ""
	for i, v := range ingredients {
		regex += v
		if i != len(ingredients) - 1 {
			regex += "|"
		}
	}
	conditions := [2]bson.D {
		bson.D{{"ingredients.1", bson.D{{"$exists", "true"}}}},
		bson.D{{"ingredients",
			bson.D{{"$not",
				bson.D{{"$elemMatch",
					bson.D{{"$not", bson.D{{"$regex", regex}}  }},
				}},
			}},
		}},
	}
	filter = bson.D{{"$and", conditions}}
	option = bson.D{{"_id", 0}}
	results := queryCollectionAndReturnResult("recipes", filter, option)
	recipes := convertRecipeBSONToSlice(results)
	return recipes, nil
}