package application

import (
	"encoding/json"
	"fmt"
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
		fmt.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

}
