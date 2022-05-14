# weather

Current weather from the command line using OpenWeather API.

## Usage

The default value for an API key is taken from the OPENWEATHER_API_KEY environment variable.

```
$ weather -key <your-openweather-api-key> paris
#=> Paris  broken clouds  3 °C

$ weather -v -key <your-openweather-api-key> paris
#=>
# Paris Apr  1 20:52:53
# ========================
# condition: broken clouds
# temperature: 3 °C
# pressure: 1011 hPa
# humidity: 81.0%
# wind: 10° 7.7 m/s
```
