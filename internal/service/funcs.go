package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"weather/internal/structs"
)

const CityLocationAPI = "http://api.openweathermap.org/geo/1.0/direct?q=%s&limit=1&appid=0c7c9c223884e6ee4768d39df765340c" //left api keys on purpose
const WeatherLatLonAPI = "https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&appid=0c7c9c223884e6ee4768d39df765340c"

func (r *Request) RequestWeather(city string) (structs.WeatherData, error) {
	ret, err :=	r.repo.GetWeather(city)
	if err != nil {
		return ret, err
	}
	return ret, nil
}

func (r *Request) PutWeatherInDatabase(city string) error {
	requestLatLonURL := fmt.Sprintf(CityLocationAPI, city)

	cityData, err := cityDataRequest(requestLatLonURL)
	if err != nil {
		return err
	}

	requestWeatherDataURL := fmt.Sprintf(WeatherLatLonAPI, cityData[0].Lat, cityData[0].Lon)
	cityWeatherData, err := cityWeatherRequest(requestWeatherDataURL)
	if err != nil {
		return err
	}

	err = r.repo.PutWeather(cityWeatherData, city)
	if err != nil {
		return err
	}
	return nil
}

func cityDataRequest(api string) (structs.CityData, error) {
	var cityData structs.CityData

	resp, err := http.Get(api)
	if err != nil {
		return cityData, err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&cityData)
	if err != nil {
		return cityData, err
	}

	return cityData, nil
}

func cityWeatherRequest(api string) (structs.WeatherData, error) {
	var cityDataWeather structs.WeatherData

	resp, err := http.Get(api)
	if err != nil {
		return cityDataWeather, err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&cityDataWeather)
	if err != nil {
		return cityDataWeather, err
	}

	return cityDataWeather, nil
}
