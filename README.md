# Steampipe plugin for Open-Meteo

https://open-meteo.com/

Use SQL to query weather forecast data from Open-Meteo.

The plugin utilises the Go Open-Meteo client 'omgo' - https://github.com/HectorMalot/omgo/

This is in early development stages. Currently some hourly and daily forecast data is available.

Open_Meteo does not require registration or an API KEY but non-commercial use is limited to 10,000 requests per day. Rate limiting to address this has not yet been implemented/configured in the plugin.

Once installed, run a query (you must provide latitude and longitude):

```sql
select
   time,
   temperature_2m,
   precipitation_probability,
   precipitation,
   cloud_cover,
   wind_speed_10m
from
   openmeteo_hourly_weather
where
   latitude = 51.5053
and
   longitude = -0.0755;
```

### Using another plugin for location
If you have a location plugin installed, for example s0rbus/locationiq, then you can combine the two using a join so that you can retrieve weather forecast data by placename

```sql
select
   to_char(t1.time, 'YYYY-MM-DD HH24:MI:SS') as date_time,
   t2.address,
   t1.temperature_2m,
   t1.precipitation
from
   openmeteo_hourly_weather as t1,
   locationiq_place2latlong as t2
where
   t1.forecastdays = 1
and
   t2.placequery = 'Cambridge, UK'
and
   t1.latitude = t2.latitude and t1.longitude = t2.longitude
order by
   time;
```

```
+---------------------+-------------------------------------------------------------------------------------+----------------+---------------+
| date_time           | address                                                                             | temperature_2m | precipitation |
+---------------------+-------------------------------------------------------------------------------------+----------------+---------------+
| 2024-01-23 00:00:00 | Cambridge, Cambridgeshire, Cambridgeshire and Peterborough, England, United Kingdom | 6.2            | 0             |
| 2024-01-23 01:00:00 | Cambridge, Cambridgeshire, Cambridgeshire and Peterborough, England, United Kingdom | 5.8            | 0             |
| 2024-01-23 02:00:00 | Cambridge, Cambridgeshire, Cambridgeshire and Peterborough, England, United Kingdom | 5.4            | 0             |
| 2024-01-23 03:00:00 | Cambridge, Cambridgeshire, Cambridgeshire and Peterborough, England, United Kingdom | 4.9            | 0             |
| 2024-01-23 04:00:00 | Cambridge, Cambridgeshire, Cambridgeshire and Peterborough, England, United Kingdom | 4.6            | 0             |
| 2024-01-23 05:00:00 | Cambridge, Cambridgeshire, Cambridgeshire and Peterborough, England, United Kingdom | 4.2            | 0             |
| 2024-01-23 06:00:00 | Cambridge, Cambridgeshire, Cambridgeshire and Peterborough, England, United Kingdom | 4.4            | 0             |
| 2024-01-23 07:00:00 | Cambridge, Cambridgeshire, Cambridgeshire and Peterborough, England, United Kingdom | 4.6            | 0             |
| 2024-01-23 08:00:00 | Cambridge, Cambridgeshire, Cambridgeshire and Peterborough, England, United Kingdom | 4.9            | 0             |
| 2024-01-23 09:00:00 | Cambridge, Cambridgeshire, Cambridgeshire and Peterborough, England, United Kingdom | 5.4            | 0             |
| 2024-01-23 10:00:00 | Cambridge, Cambridgeshire, Cambridgeshire and Peterborough, England, United Kingdom | 6.4            | 0.3           |
| 2024-01-23 11:00:00 | Cambridge, Cambridgeshire, Cambridgeshire and Peterborough, England, United Kingdom | 7.4            | 0.6           |
| 2024-01-23 12:00:00 | Cambridge, Cambridgeshire, Cambridgeshire and Peterborough, England, United Kingdom | 8.2            | 1             |
| 2024-01-23 13:00:00 | Cambridge, Cambridgeshire, Cambridgeshire and Peterborough, England, United Kingdom | 8.9            | 0             |
| 2024-01-23 14:00:00 | Cambridge, Cambridgeshire, Cambridgeshire and Peterborough, England, United Kingdom | 9.4            | 0             |
| 2024-01-23 15:00:00 | Cambridge, Cambridgeshire, Cambridgeshire and Peterborough, England, United Kingdom | 10.1           | 0             |
| 2024-01-23 16:00:00 | Cambridge, Cambridgeshire, Cambridgeshire and Peterborough, England, United Kingdom | 11.3           | 0             |
| 2024-01-23 17:00:00 | Cambridge, Cambridgeshire, Cambridgeshire and Peterborough, England, United Kingdom | 12.6           | 0             |
+---------------------+-------------------------------------------------------------------------------------+----------------+---------------+
```