package main

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/sqlite"
	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io"
	"net/http"
)

type WeatherData struct {
	gorm.Model
	Time             string  `json:"time"`
	Temperature2mMax float64 `json:"temperature_2m_max"`
	Temperature2mMin float64 `json:"temperature_2m_min"`
}

type WeatherResponse struct {
	Daily struct {
		Time             []string  `json:"time"`
		Temperature2mMax []float64 `json:"temperature_2m_max"`
		Temperature2mMin []float64 `json:"temperature_2m_min"`
	} `json:"daily"`
}

type WeatherResponseFin struct {
	Daily struct {
		Time             string  `json:"time"`
		Temperature2mMax float64 `json:"temperature_2m_max"`
		Temperature2mMin float64 `json:"temperature_2m_min"`
	} `json:"daily"`
}

type Proxy struct{}

func (w *Proxy) GetWeather() (WeatherResponseFin, error) {
	resp, err := http.Get("https://api.open-meteo.com/v1/forecast?latitude=50.0614&longitude=19.9366&daily=temperature_2m_max,temperature_2m_min&timezone=Europe%2FBerlin&forecast_days=1")
	if err != nil {
		return WeatherResponseFin{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return WeatherResponseFin{}, err
	}

	var weatherResponse WeatherResponse
	err = json.Unmarshal(body, &weatherResponse)
	if err != nil {
		return WeatherResponseFin{}, err
	}

	var weatherResponseFin WeatherResponseFin
	weatherResponseFin.Daily.Time = weatherResponse.Daily.Time[0]
	weatherResponseFin.Daily.Temperature2mMax = weatherResponse.Daily.Temperature2mMax[0]
	weatherResponseFin.Daily.Temperature2mMin = weatherResponse.Daily.Temperature2mMin[0]

	return weatherResponseFin, nil
}

func main() {
	e := echo.New()
	proxy := &Proxy{}

	db, err := gorm.Open(sqlite.Open("weather.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to the database")
	}
	err = db.AutoMigrate(&WeatherData{})
	if err != nil {
		panic("failed to migrate database")
	}

	e.GET("/weather", func(c echo.Context) error {
		weatherResponse, err := proxy.GetWeather()
		if err != nil {
			return err
		}

		var weatherData WeatherData
		weatherData.Time = weatherResponse.Daily.Time
		weatherData.Temperature2mMin = weatherResponse.Daily.Temperature2mMin
		weatherData.Temperature2mMax = weatherResponse.Daily.Temperature2mMax
		if err := db.Create(&weatherData).Error; err != nil {
			return err
		}

		return c.JSON(http.StatusOK, weatherResponse)
	})

	e.GET("/db", func(c echo.Context) error {
		var weatherData []WeatherData
		result := db.Find(&weatherData)
		if result.Error != nil {
			return result.Error
		}
		return c.JSON(http.StatusOK, weatherData)
	})

	e.Logger.Fatal(e.Start(":1323"))
}
