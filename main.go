package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

type options struct {
	apiKey   string
	units    string
	verbose  bool
	cityName string
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
		fmt.Fprintf(w, "USAGE:\n")
		fmt.Fprintf(w, "\tweather [OPTIONS] <CITY-NAME>\n")
		fmt.Fprintf(w, "OPTIONS:\n")
		flag.PrintDefaults()
	}

	opt := options{units: "metric"}

	flag.StringVar(&opt.apiKey, "key", os.Getenv("OPENWEATHER_API_KEY"), "openweather api key (required)")
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
}
