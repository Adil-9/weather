package application

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"weather/internal/structs"
)

func (h *Handler) HandleWeather(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.HandleGet(w, r)
	case http.MethodPut:
		h.HandlePut(w, r)
	default:
		w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
	}
}

func (h *Handler) HandleGet(w http.ResponseWriter, r *http.Request) {
	var city structs.City
	err := json.NewDecoder(r.Body).Decode(&city)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if city.Name == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	weatherData, err := h.service.RequestWeather(city.Name)
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	marshalled, err := json.MarshalIndent(weatherData, "", "\t")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	w.Write(marshalled)
}

func (h *Handler) HandlePut(w http.ResponseWriter, r *http.Request) {
	var city structs.City
	err := json.NewDecoder(r.Body).Decode(&city)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.service.PutWeatherInDatabase(city.Name)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
