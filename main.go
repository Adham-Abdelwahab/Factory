package main

import (
	"fmt"
	"net/http"
	"os"

	"Factory/internal/handlers"
	"Factory/internal/middleware"

	"github.com/go-chi/chi"
	chimiddle "github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetReportCaller(true)
	logrus.SetOutput(os.Stdout)

	var factory = chi.NewRouter()
	factory.Use(chimiddle.StripSlashes)
	factory.Use(middleware.Correlation)

	handlers.GridHandler(factory)
	handlers.RailwayHandler(factory)
	handlers.ResourceHandler(factory)

	fmt.Println("Starting the Factory ...")

	err := http.ListenAndServe("localhost:8080", factory)
	if err != nil {
		fmt.Println("Failed to start the factory.")
		logrus.Error(err)
	}
}
