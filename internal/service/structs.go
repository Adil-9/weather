package service

import (
	"weather/internal/database"
)

type Service struct {
	Requester
}

func NewService(repo *database.Repository) *Service {
	return &Service{Requester: NewRequest(repo.Manage)}
}

type Requester interface {
	RequestWeather()
	PutWeatherInDatabase(city string) error
}

type Request struct {
	repo database.Manage
}

func NewRequest(repo database.Manage) *Request {
	return &Request{repo: repo}
}

