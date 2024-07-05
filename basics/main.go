package main

import (
	"context"
	"log"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
)

type HelloResponse struct {
	Message string `json:"message"`
}

func main() {
	router := humachi.New()

	api := huma.NewAPI(huma.DefaultConfig("My API", "1.0.0"), router)

	api.Get("/hello", func(ctx context.Context) (*HelloResponse, error) {
		return &HelloResponse{Message: "Hello, World!"}, nil
	})

	log.Fatal(api.Listen(":8080"))
}