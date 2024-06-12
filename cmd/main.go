package main

import (
	"fmt"
	"log"
	"net/http"
	"weather/internal/application"
	"weather/internal/database"
	"weather/internal/service"
)

func main() {
	handler := createHandler()

	http.HandleFunc("/weather", handler.HandleWeather)

	fmt.Println("server starting on localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err.Error())
	}
}

func createHandler() *application.Handler {	//this makes it more complicated but architechture ¯\_(ツ)_/¯
	db, err := database.ConnectPSQL()
	if err != nil {
		log.Fatal(err.Error())
	}

	err = database.CreateTables(db)
	if err != nil {
		log.Fatal(err.Error())
	}

	repo := database.NewRepository(db)
	service := service.NewService(repo)
	handler := application.NewHandler(service)
	return handler
}
