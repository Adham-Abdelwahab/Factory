package main

import (
	"fmt"
	"net/http"
	"os"

	"Factory/internal/middleware"
	"Factory/internal/system"
	"Factory/internal/util"

	"github.com/go-chi/chi"
	chimiddle "github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetReportCaller(true)
	logrus.SetOutput(os.Stdout)

	util.InitializeDatabase()
	defer util.Database.Close()

	var factory = chi.NewRouter()
	factory.Use(chimiddle.StripSlashes)
	factory.Use(middleware.Correlation)
	system.Initialize(factory)

	fmt.Println("Starting the Factory ...")

	err := http.ListenAndServe("localhost:8080", factory)
	if err != nil {
		fmt.Println("Failed to start the factory !")
		logrus.Error(err)
	}
}
