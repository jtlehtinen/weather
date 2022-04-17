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
	apiKey   string
	units    string
	verbose  bool
	cityName string
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
}

func makeRequestURL(cityName, units, apiKey string) string {
	cityName = url.QueryEscape(cityName)
	apiKey = url.QueryEscape(apiKey) // Just in case user inputs nonsense...
	return fmt.Sprintf("%s?q=%s&units=%s&appid=%s", BASE_URL, cityName, units, apiKey)
}

func sendRequest(apiKey, cityName, units string) (*Weather, error) {
	u := makeRequestURL(cityName, units, apiKey)

	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("request status %d %s", resp.StatusCode, http.StatusText(resp.StatusCode)))
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
	}

	return w, nil
}

func display(w io.Writer, wt *Weather, opt *options) {
	temperatureSymbol, windSpeedSymbol := "C", "m/s"
	if opt.units == "imperial" {
		temperatureSymbol, windSpeedSymbol = "F", "mi/h"
	}

	if opt.verbose {
		t := time.Now().UTC().Add(time.Duration(wt.TimeZone) * time.Second)

		fmt.Printf("%s %s\n", wt.CityName, t.Format(time.Stamp))
		fmt.Printf("========================\n")
		fmt.Printf("condition: %s\n", wt.Conditions)
		fmt.Printf("temperature: %.0f °%s\n", wt.Temperature, temperatureSymbol)
		fmt.Printf("pressure: %.0f hPa\n", wt.Pressure)
		fmt.Printf("humidity: %.1f%%\n", wt.Humidity)
		fmt.Printf("wind: %.0f° %.1f %s\n", wt.WindDegrees, wt.WindSpeed, windSpeedSymbol)
	} else {
		fmt.Printf("%s  %s  %0.f °%s\n", wt.CityName, wt.Conditions, wt.Temperature, temperatureSymbol)
	}
}

func usageAndExit(errmsg string) {
	if errmsg != "" {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n\n", errmsg)
	}
	flag.Usage()
	os.Exit(2)
}

func main() {
	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "weather displays the current weather of a city.\n\n")
		fmt.Fprintf(w, "USAGE:\n")
		fmt.Fprintf(w, "\tweather [OPTIONS] <CITY-NAME>\n\n")
		fmt.Fprintf(w, "OPTIONS:\n")
		flag.PrintDefaults()
	}

	opt := options{units: "metric"}

	flag.StringVar(&opt.apiKey, "key", os.Getenv("OPENWEATHER_API_KEY"), "openweather api key")
	flag.BoolVar(&opt.verbose, "v", false, "verbose output")

	flag.Func("units", "units of measurement (metric|imperial)", func(value string) error {
		if value != "metric" && value != "imperial" {
			return errors.New("unit must be 'metric' or 'imperial'\n")
		}
		opt.units = value
		return nil
	})

	flag.Parse()

	opt.cityName = strings.Join(flag.Args(), " ")

	if opt.apiKey == "" {
		usageAndExit("openweather api key is required")
	}

	if strings.TrimSpace(opt.cityName) == "" {
		usageAndExit("city name is required")
	}

	w, err := sendRequest(opt.apiKey, opt.cityName, opt.units)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}

	display(os.Stdout, w, &opt)
}
