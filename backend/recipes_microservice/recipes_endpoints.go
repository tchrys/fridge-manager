package main

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

func makeFindRecipeByURLEndpoint(ris RecipeInfoService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(urlRequest)
		recipe, err := ris.FindRecipeByURL(req.Url)
		if err != nil {
			return nil, DbNotFound
		}
		return recipe, nil
	}
}

func makeFindRecipesByNutritionRangeEndpoint(ris RecipeInfoService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(propRequest)
		recipes, _ := ris.FindRecipesByNutritionRange(req.RangeType, req.LowerBound, req.UpperBound)
		return recipes, nil
	}
}

func makeMapReduceByCategoryEndpoint(ris RecipeInfoService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(mapReduceRequest)
		result, _ := ris.MapReduceByCategory(req.Step, req.Category)
		return result, nil
	}
}

func makeFilterByIngredientsEndpoint(ris RecipeInfoService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(ingredientsRequest)
		result, _ := ris.FilterByIngredients(req.Ingredients)
		return result, nil
	}
}

func makeGetAffordableRecipesEndpoint(ris RecipeInfoService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(ingredientsRequest)
		result, _ := ris.GetAffordableRecipes(req.Ingredients)
		return result, nil
	}
}
