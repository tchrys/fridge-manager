package main

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
)

type Middleware func(endpoint.Endpoint) endpoint.Endpoint

func urlMiddleware() Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			defer func() {
				req := request.(urlRequest)
				fmt.Println(req.Url)

			}()
			return next(ctx, request)
			//return next
		}
	}
}
