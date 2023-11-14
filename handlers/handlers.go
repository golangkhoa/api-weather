package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Weather struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct {
		TempC     float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`
	Forecast struct {
		Forecastday []struct {
			Hour []struct {
				TimeEpoch int64   `json:"time_epoch"`
				TempC     float64 `json:"temp_c"`
				Condition struct {
					Text string `json:"text"`
				} `json:"condition"`
				ChanceOfRain float64 `json:"chance_of_rain"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func weather(q string, w http.ResponseWriter, r *http.Request) (string, error) {
	if q == "" {
		q = "metar:VVTS"
	}
	res, err := http.Get("http://api.weatherapi.com/v1/forecast.json?key=5c4ea4e5c9284c3d813123717231011&q=" + q + "&days=1&aqi=no&alerts=no")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		panic("Weather API not available")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Sprint(err), err
	}
	var weather Weather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		return fmt.Sprint(err), err
	}

	location, current, hours := weather.Location, weather.Current, weather.Forecast.Forecastday[0].Hour
	now := fmt.Sprintf("%s, %s: %.1f°C, %s\n", location.Name, location.Country, current.TempC, current.Condition.Text)
	w.Write([]byte(now))

	for _, hour := range hours {
		date := time.Unix(hour.TimeEpoch, 0)

		if date.Before(time.Now()) {
			continue
		}

		message := fmt.Sprintf("%s - %.1f°C, %.0f%%, %s\n", date.Format("15:04"), hour.TempC, hour.ChanceOfRain, hour.Condition.Text)
		w.Write([]byte(message))
	}

	return now, nil
}

func weatherJson(q string, w http.ResponseWriter, r *http.Request) error {
	if q == "" {
		q = "metar:VVTS"
	}
	res, err := http.Get("http://api.weatherapi.com/v1/forecast.json?key=5c4ea4e5c9284c3d813123717231011&q=" + q + "&days=1&aqi=no&alerts=no")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		panic("Weather API not available")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var weather Weather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		return err
	}

	out, _ := json.MarshalIndent(weather, "", "    ")
	w.Write(out)

	return nil
}

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	weather("", w, r)
}

func PlaceHandler(q string, w http.ResponseWriter, r *http.Request) {
	weather(q, w, r)
}

func JsonHandler(q string, w http.ResponseWriter, r *http.Request) {
	weatherJson(q, w, r)
}