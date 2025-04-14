package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	fmt.Println("started")
	router := chi.NewRouter()

	router.Use(middleware.Logger)

}
