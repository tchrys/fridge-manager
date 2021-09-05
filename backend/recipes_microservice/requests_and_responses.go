package main

type recipeResponse struct {
	Url         string   `json:"url"`
	TotalTime   string   `json:"totalTime"`
	PrepTime    string   `json:"prepTime"`
	CookTime    string   `json:"cookTime"`
	Servings    int      `json:"servings"`
	Calories    int      `json:"calories"`
	Fat         int      `json:"fat"`
	Carbs       int      `json:"carbs"`
	Fiber       int      `json:"fiber"`
	Sugar       int      `json:"sugar"`
	Protein     int      `json:"protein"`
	Ingredients []string `json:"ingredients"`
}

type mapReduceResponse struct {
	Id int `bson:"_id" json:"id"`
	Value int `bson:"value" json:"value"`
}

type urlRequest struct {
	Url string `json:"url"`
}

type propRequest struct {
	RangeType  string `json:"rangeType"`
	UpperBound int    `json:"upperBound"`
	LowerBound int    `json:"lowerBound"`
}

type mapReduceRequest struct {
	Step     string `json:"step"`
	Category string `json:"category"`
}

type ingredientsRequest struct {
	Ingredients []string `json:"ingredients"`
}
