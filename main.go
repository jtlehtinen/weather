package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const BASE_URL = "https://api.openweathermap.org/data/2.5/weather"

type options struct {
	apiKey  string
	units   string
	verbose bool
	city    string
}

func exitWithError(errorMessage string) {
	fmt.Fprintf(os.Stderr, "ERROR: %s\n", errorMessage)
	os.Exit(1)
}

type Weather struct {
	CityName    string
	TimeZone    int
	Visibility  float64
	Temperature float64
	Pressure    float64
	Humidity    float64
	WindSpeed   float64
	WindDegrees float64
	Conditions  string
	Icon        string
}

func makeRequestURL(cityName, units, apiKey string) string {
	cityName = url.QueryEscape(cityName)
	apiKey = url.QueryEscape(apiKey)
	return fmt.Sprintf("%s?q=%s&units=%s&appid=%s", BASE_URL, cityName, units, apiKey)
}

func fetchWeather(apiKey, cityName, units string) (*Weather, error) {
	u := makeRequestURL(cityName, units, apiKey)

	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request status %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	// API docs: https://openweathermap.org/current
	type response struct {
		Weather []struct {
			Main        string `json:"main"`
			Description string `json:"description"`
			Icon        string `json:"icon"`
		} `json:"weather"`
		Main struct {
			Temperature float64 `json:"temp"`
			Pressure    float64 `json:"pressure"`
			Humidity    float64 `json:"humidity"`
		} `json:"main"`
		Wind struct {
			Speed   float64 `json:"speed"`
			Degrees float64 `json:"deg"`
		} `json:"wind"`
		Name       string  `json:"name"`
		TimeZone   int     `json:"timezone"`
		Visibility float64 `json:"visibility"`
	}

	var res response
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}

	w := &Weather{}
	w.CityName = res.Name
	w.TimeZone = res.TimeZone
	w.Visibility = res.Visibility
	w.Temperature = res.Main.Temperature
	w.Pressure = res.Main.Pressure
	w.Humidity = res.Main.Humidity
	w.WindSpeed = res.Wind.Speed
	w.WindDegrees = res.Wind.Degrees

	// @NOTE: Maybe take all?
	if len(res.Weather) > 0 {
		w.Conditions = res.Weather[0].Description
		w.Icon = res.Weather[0].Icon
	}

	return w, nil
}

func weatherIconToString(icon string) string {
	// @TODO: Improve these...
	// https://openweathermap.org/weather-conditions
	switch icon {
	case "01d":
		return "☀️" // clear sky day
	case "02d":
		return "⛅" // few clouds day
	case "03d":
		return "☁️" // scattered clouds day
	case "04d":
		return "☁️" // broken clouds day
	case "09d":
		return "🌧️" // shower rain day
	case "10d":
		return "🌦️" // rain day
	case "11d":
		return "⛈️" // thunderstorm day
	case "13d":
		return "❄️" // snow day
	case "50d":
		return "🌫️" // mist day
	case "01n":
		return "🌑" // clear sky night
	case "02n":
		return "⛅" // few clouds night
	case "03n":
		return "☁️" // scattered clouds night
	case "04n":
		return "☁️" // broken clouds night
	case "09n":
		return "🌧️" // shower rain night
	case "10n":
		return "🌦️" // rain night
	case "11n":
		return "⛈️" // thunderstorm night
	case "13n":
		return "❄️" // snow night
	case "50n":
		return "🌫️" // mist night
	default:
		return ""
	}
}

func display(w io.Writer, wt *Weather, opt *options) {
	temperatureSymbol, windSpeedSymbol := "C", "m/s"
	if opt.units == "imperial" {
		temperatureSymbol, windSpeedSymbol = "F", "mi/h"
	}

	weatherSymbol := weatherIconToString(wt.Icon)

	if opt.verbose {
		t := time.Now().UTC().Add(time.Duration(wt.TimeZone) * time.Second)

		fmt.Printf("%s %s\n", wt.CityName, t.Format(time.Stamp))
		fmt.Printf("========================\n")
		fmt.Printf("condition: %s %s\n", weatherSymbol, wt.Conditions)
		fmt.Printf("temperature: %.0f°%s\n", wt.Temperature, temperatureSymbol)
		fmt.Printf("pressure: %.0f hPa\n", wt.Pressure)
		fmt.Printf("humidity: %.1f%%\n", wt.Humidity)
		fmt.Printf("wind: %.0f° %.1f %s\n", wt.WindDegrees, wt.WindSpeed, windSpeedSymbol)
	} else {
		fmt.Printf("%s %0.f°%s %s %s\n", wt.CityName, wt.Temperature, temperatureSymbol, weatherSymbol, wt.Conditions)
	}
}

func main() {
	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "weather displays the current weather of a given city.\n\n")
		fmt.Fprintf(w, "usage:\n")
		fmt.Fprintf(w, "\tweather [options] <city>\n\n")
		fmt.Fprintf(w, "options:\n")
		flag.PrintDefaults()
	}

	opt := options{units: "metric"}

	flag.StringVar(&opt.apiKey, "key", os.Getenv("OPENWEATHER_API_KEY"), "OpenWeather API key")
	flag.BoolVar(&opt.verbose, "v", false, "verbose output")
	flag.Func("units", "units of measurement (metric|imperial)", func(value string) error {
		if value != "metric" && value != "imperial" {
			return errors.New("unit must be 'metric' or 'imperial'")
		}
		opt.units = value
		return nil
	})

	flag.Parse()

	opt.city = strings.Join(flag.Args(), " ")

	if opt.apiKey == "" {
		exitWithError("OpenWeather API key is required")
	}

	if strings.TrimSpace(opt.city) == "" {
		exitWithError("city name is required")
	}

	w, err := fetchWeather(opt.apiKey, opt.city, opt.units)
	if err != nil {
		exitWithError(err.Error())
	}

	display(os.Stdout, w, &opt)
}
