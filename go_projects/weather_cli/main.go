package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Location struct {
	Name    string `json:"name"`
	Country string `json:"country"`
}

type Current struct {
	Temperature int `json:"temperature"`
	WindSpeed   int `json:"wind_speed"`
	Humidity    int `json:"humidity"`
}

type WeatherResponse struct {
	Location Location `json:"location"`
	Current  Current  `json:"current"`
}

func getWeather(city, apiKey string) (*WeatherResponse, error) {

	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s", apiKey, city)

	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	//Print the response body for debugging
	// body, _ := io.ReadAll(resp.Body)
	// fmt.Println("Response Body:", string(body))

	var weatherResponse WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherResponse); err != nil {
		return nil, err
	}

	return &weatherResponse, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a city name.")
		return
	}

	city := os.Args[1]
	apiKey := "50d4cc9d4ac9484484195735242308"

	weather, err := getWeather(city, apiKey)
	if err != nil {
		fmt.Printf("Error fetching weather data: %v\n", err)
		return
	}

	fmt.Printf("Weather in %s, %s:\n", weather.Location.Name, weather.Location.Country)
	fmt.Printf("Temperature: %d°C\n", weather.Current.Temperature)
	fmt.Printf("Wind Speed: %d km/h\n", weather.Current.WindSpeed)
	fmt.Printf("Humidity: %d%%\n", weather.Current.Humidity)
}
