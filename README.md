# Steampipe plugin for Open-Meteo

https://open-meteo.com/

Use SQL to query weather forecast data from Open-Meteo.

This is in early development stages. Currently some hourly and daily forecast data is available.

Open_Meteo does not require registration or an API KEY but non-commercial use is limited to 10,000 requests per day. Rate limiting to address this has not yet been implemented/configured in the plugin.

Once installed, run a query (you must provide latitude and longitude):
```
select time,temperature_2m,precipitation_probability,precipitation,cloud_cover,wind_speed_10m from openmeteo_hourly_weather where latitude = 51.5102 and longitude = -0.1181
```
