package main

import (
	"log"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
)

func main() {
	ris := recipeInfoService{}
	createMongoConnection()

	getByURLHandler := httptransport.NewServer(
		urlMiddleware() (makeFindRecipeByURLEndpoint(ris)), decodeGetByURLRequest, encodeResponse)
	getByPropRangeHandler := httptransport.NewServer(
		makeFindRecipesByNutritionRangeEndpoint(ris), decodeByPropRange, encodeResponse)
	mapReduceByCategoryHandler := httptransport.NewServer(
		makeMapReduceByCategoryEndpoint(ris), decodeForMapReduce, encodeResponse)
	filterByIngredientsHandler := httptransport.NewServer(makeFilterByIngredientsEndpoint(ris), decodeIngredients, encodeResponse)
	getAffordableHandler := httptransport.NewServer(makeGetAffordableRecipesEndpoint(ris), decodeIngredients, encodeResponse)

	http.Handle("/by-url", getByURLHandler)
	http.Handle("/by-nutrition-prop", getByPropRangeHandler)
	http.Handle("/intervals-by-prop", mapReduceByCategoryHandler)
	http.Handle("/by-ingredients", filterByIngredientsHandler)
	http.Handle("/affordable", getAffordableHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}