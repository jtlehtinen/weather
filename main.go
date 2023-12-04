package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

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
}
