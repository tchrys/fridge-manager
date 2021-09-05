//package user_microservice
//
//import (
//	"net/http"
//
//	"fridge-manager/repositories"
//
//	"github.com/gin-gonic/gin"
//)
//
//type urlPayload struct {
//	Url string `json:"url"`
//}
//
//type caloriesPayload struct {
//	RangeType string `json:"rangeType"`
//	UpperBound int `json:"upperBound"`
//	LowerBound int `json:"lowerBound"`
//}
//
//type mapReducePayload struct {
//	Step string `json:"step"`
//	Category string `json:"category"`
//}
//
//type ingredientsPayload struct {
//	Ingredients []string `json:"ingredients"`
//}
//
//func getByURL(c *gin.Context) {
//	var payload urlPayload
//	if err := c.BindJSON(&payload); err != nil {
//		return
//	}
//	recipe := repositories.FindRecipeByURL(payload.Url)
//	c.IndentedJSON(http.StatusOK, recipe)
//}
//
//func getByCaloriesRange(c *gin.Context) {
//	var payload caloriesPayload
//	if err := c.BindJSON(&payload); err != nil {
//		return
//	}
//	recipes := repositories.FindRecipesByNutritionRange(payload.LowerBound, payload.UpperBound, payload.RangeType)
//	c.IndentedJSON(http.StatusOK, recipes)
//}
//
//func mapReduceByCategory(c *gin.Context) {
//	var payload mapReducePayload
//	if err := c.BindJSON(&payload); err != nil {
//		return
//	}
//	res := repositories.MapReduceByCategory(payload.Step, payload.Category)
//	c.IndentedJSON(http.StatusOK, res)
//}
//
//func filterByIngredients(c *gin.Context) {
//	var payload ingredientsPayload
//	if err := c.BindJSON(&payload); err != nil {
//		return
//	}
//	res := repositories.FilterByIngredients(payload.Ingredients)
//	c.IndentedJSON(http.StatusOK, res)
//}
//
//func getAffordableRecipes(c *gin.Context) {
//	var payload ingredientsPayload
//	if err := c.BindJSON(&payload); err != nil {
//		return
//	}
//	res := repositories.GetAffordableRecipes(payload.Ingredients)
//	c.IndentedJSON(http.StatusOK, res)
//}
//
//func RecipesRoutes(route *gin.Engine) {
//	recipesRoutes := route.Group("/recipes")
//	{
//		recipesRoutes.GET("/by-url", getByURL)
//		recipesRoutes.GET("/by-nutrition-prop", getByCaloriesRange)
//		recipesRoutes.GET("/intervals-by-prop", mapReduceByCategory)
//		recipesRoutes.GET("/by-ingredients", filterByIngredients)
//		recipesRoutes.GET("/affordable", getAffordableRecipes)
//	}
//}
