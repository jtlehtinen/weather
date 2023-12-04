# weather

Fetch and display current weather from the command line using OpenWeather API.

## Usage

The default value for an API key is taken from the OPENWEATHER_API_KEY environment variable. Alternatively the API key can be passed as an argument using the `-key` flag.

```sh
$ weather helsinki
#\=>
# Helsinki -9°C ❄️ snow

$ weather -v helsinki
#\=>
# Helsinki Dec  4 19:14:08
# ========================
# condition: ❄️ snow
# temperature: -9°C
# pressure: 1013 hPa
# humidity: 91.0%
# wind: 354° 4.5 m/s
```
