package graph

import "src/http/pkg/services"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	service services.ServicesInterface
}

func NewResolver(store services.ServicesInterface) *Resolver {
	return &Resolver{
		service: store,
	}
}
