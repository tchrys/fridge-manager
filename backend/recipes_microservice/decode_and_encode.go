package main

import (
	"context"
	"encoding/json"
	"net/http"
)

func decodeGetByURLRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request urlRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeByPropRange(_ context.Context, r *http.Request) (interface{}, error) {
	var request propRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeForMapReduce(_ context.Context, r *http.Request) (interface{}, error) {
	var request mapReduceRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeIngredients(_ context.Context, r *http.Request) (interface{}, error) {
	var request ingredientsRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "    ")
	return enc.Encode(response)
}
