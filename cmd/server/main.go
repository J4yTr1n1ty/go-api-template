package main

import (
	"log"
	"net/http"

	"github.com/J4yTr1n1ty/go-api-template/pkg/boot"
	"github.com/J4yTr1n1ty/go-api-template/pkg/models"
	"github.com/J4yTr1n1ty/go-api-template/pkg/web"
	"github.com/J4yTr1n1ty/go-api-template/pkg/web/middleware"
)

func main() {
	log.Println("Starting server...")

	boot.LoadEnv()
	boot.ConnectDB()
	err := boot.DB.AutoMigrate(&models.Example{})
	if err != nil {
		log.Fatal(err)
	}

	router := web.SetupRoutes()

	stack := middleware.CreateStack(middleware.Logging, middleware.Recovery)

	server := http.Server{
		Addr:    ":" + boot.Environment.GetEnv("PORT"),
		Handler: stack(router),
	}

	log.Println("Listening on port :" + boot.Environment.GetEnv("PORT"))
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Error starting server: " + err.Error())
	}
}
