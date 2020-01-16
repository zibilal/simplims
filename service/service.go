package service

import "context"

// Service is meant to be a layer that holds all the business logic
// All service types must contains Serve function

type Service interface {
	Serve(ctx context.Context, input interface{}, output interface{}, status chan<- error)
}
