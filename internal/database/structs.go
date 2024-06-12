package database

import (
	"database/sql"
	"weather/internal/structs"
)

type Manage interface {
	GetWeather(city string)                               //return json smth
	PutWeather(cityWeatherData structs.WeatherData, city string) error //pass struct smth
}

type Repository struct {
	Manage
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{NewManager(db)}
}

type Manager struct {
	db *sql.DB
}

func NewManager(db *sql.DB) *Manager {
	return &Manager{db: db}
}
