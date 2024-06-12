package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"weather/internal/application"
	"weather/internal/database"
	"weather/internal/service"
	"weather/internal/structs"
)

func TestGet(t *testing.T) {
	test := createHandler()

	item := structs.City{Name: "London"}
	jsonData, err := json.Marshal(item)
	if err != nil {
		t.Fatalf("Failed to marshal item: %v", err)
	}

	req, err := http.NewRequest("GET", "/weather", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(test.HandleWeather)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	var returnedItem structs.WeatherData
	err = json.Unmarshal(rr.Body.Bytes(), &returnedItem)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	fmt.Println(returnedItem)	//data always changes, so i just printed it out
}

func TestPut(t *testing.T) {	//in method put i did not send anything
	test := createHandler()

	item := structs.City{Name: "London"}
	jsonData, err := json.Marshal(item)
	if err != nil {
		t.Fatalf("Failed to marshal item: %v", err)
	}

	req, err := http.NewRequest("PUT", "/weather", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(test.HandleWeather)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}
}

func createHandler() *application.Handler {
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
