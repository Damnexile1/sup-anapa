package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type WeatherService struct {
	apiKey string
}

type WeatherData struct {
	AirTemp     float64
	WaterTemp   float64
	WindSpeed   float64
	CloudCover  int
	Description string
}

func NewWeatherService(apiKey string) *WeatherService {
	return &WeatherService{apiKey: apiKey}
}

func (s *WeatherService) GetWeather(date time.Time, lat, lon float64) (*WeatherData, error) {
	// TODO: Implement actual OpenWeatherMap API call
	// For now, return mock data
	return &WeatherData{
		AirTemp:     24.0,
		WaterTemp:   20.0,
		WindSpeed:   3.5,
		CloudCover:  20,
		Description: "Ясно",
	}, nil
}

func (s *WeatherService) GetWeatherFromAPI(date time.Time, lat, lon float64) (*WeatherData, error) {
	url := fmt.Sprintf(
		"https://api.openweathermap.org/data/2.5/forecast?lat=%f&lon=%f&appid=%s&units=metric&lang=ru",
		lat, lon, s.apiKey,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	// Parse response and extract weather data
	// This is a simplified version
	return &WeatherData{
		AirTemp:     24.0,
		WaterTemp:   20.0,
		WindSpeed:   3.5,
		CloudCover:  20,
		Description: "Ясно",
	}, nil
}
