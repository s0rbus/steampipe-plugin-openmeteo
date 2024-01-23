# Table: openmeteo_hourlyweather

24 Hourly weather forecast for a given location

## Examples

### Get hourly forecast for 1 day (the default is 7 days) for a given location (lat/long)
The default hourly forecast is 24 hours over 7 days. The API allows changing the number of days between 1 and 16 days, but this plugin currently only allows changing it in the range 1-7.

```sql
select
   to_char(time, 'YYYY-MM-DD HH24:MI:SS') as date_time,
   temperature_2m,
   precipitation,
   wind_speed_10m
from 
   openmeteo_hourly_weather
where
   latitude = 51.5053
and
   longitude = -0.0755;
and
   forecastdays = 1;
order by
   time;
```

```
+---------------------+----------------+---------------+----------------+
| date_time           | temperature_2m | precipitation | wind_speed_10m |
+---------------------+----------------+---------------+----------------+
| 2024-01-23 00:00:00 | 7.6            | 0             | 10.1           |
| 2024-01-23 01:00:00 | 7              | 0             | 8.6            |
| 2024-01-23 02:00:00 | 6.4            | 0             | 8.2            |
| 2024-01-23 03:00:00 | 6.1            | 0             | 8.4            |
| 2024-01-23 04:00:00 | 5.6            | 0             | 7.6            |
| 2024-01-23 05:00:00 | 5.1            | 0             | 7.6            |
| 2024-01-23 06:00:00 | 5.3            | 0             | 6.3            |
| 2024-01-23 07:00:00 | 5.6            | 0             | 6              |
| 2024-01-23 08:00:00 | 6.4            | 0             | 6              |
| 2024-01-23 09:00:00 | 7.1            | 0             | 7.3            |
| 2024-01-23 10:00:00 | 8.1            | 0.9           | 9.1            |
| 2024-01-23 11:00:00 | 9.2            | 0.5           | 10.7           |
| 2024-01-23 12:00:00 | 9.8            | 0.5           | 13.3           |
| 2024-01-23 13:00:00 | 9.9            | 0             | 15.3           |
| 2024-01-23 14:00:00 | 10.5           | 0             | 15.1           |
| 2024-01-23 15:00:00 | 11.5           | 0             | 14.5           |
| 2024-01-23 16:00:00 | 12.2           | 0             | 15.5           |
| 2024-01-23 17:00:00 | 12.6           | 0             | 15.7           |
| 2024-01-23 18:00:00 | 12.6           | 0             | 14.4           |
| 2024-01-23 19:00:00 | 12.6           | 0             | 15.5           |
| 2024-01-23 20:00:00 | 12.9           | 0             | 16.6           |
| 2024-01-23 21:00:00 | 13.3           | 0             | 16             |
| 2024-01-23 22:00:00 | 13.4           | 0             | 16.4           |
| 2024-01-23 23:00:00 | 13.5           | 0             | 16.6           |
+---------------------+----------------+---------------+----------------+
```


### Use pastdays to get previous data. Allowed range is 0-92

```sql
select
   to_char(time, 'YYYY-MM-DD') as date,
   precipitation
from
   openmeteo_hourly_weather
where
   latitude = 55.8612
and
   longitude = -4.2501
and
   pastdays = 2
and
   forecastdays = 1
order by
   1;
```

```

```

