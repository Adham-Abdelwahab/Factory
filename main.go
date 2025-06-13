package main

import (
	f "fmt"
	"net/http"

	"Factory/internal/middleware"
	"Factory/internal/system"
	"Factory/internal/util"

	"github.com/go-chi/chi"
	chimiddle "github.com/go-chi/chi/middleware"
)

func main() {
	util.InitializeDatabase()
	defer util.Database.Close()

	var factory = chi.NewRouter()
	factory.Use(chimiddle.StripSlashes)
	factory.Use(middleware.Correlation)
	system.Initialize(factory)

	f.Println("Starting the Factory ...")

	if err := http.ListenAndServe("localhost:8080", factory); err != nil {
		f.Println("Failed to start the factory !!")
	}
}
