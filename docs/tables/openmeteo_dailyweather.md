# Table: openmeteo_dailyweather

Daily weather forecast for a given location

## Examples

### Get daily forecast for 7 days for a given location (lat/long)

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

### The default daily forecast is 7 days. The API allows changing this between 1 and 16 days, but this plugin currently only allows changing it in the range 1-7

```sql
select
   to_char(time, 'YYYY-MM-DD') as date,
   precipitation_hours
from
   openmeteo_daily_weather
where
   latitude = 55.8612
and
   longitude = -4.2501
and
   forecastdays = 2
order by
   1;
```

```
+------------+---------------------+
| date       | precipitation_hours |
+------------+---------------------+
| 2024-01-23 | 15                  |
| 2024-01-24 | 7                   |
+------------+---------------------+
```

### Use pastdays to get previous data. Allowed range is 0-92

```sql
select
   to_char(time, 'YYYY-MM-DD') as date,
   precipitation_hours
from
   openmeteo_daily_weather
where
   latitude = 55.8612
and
   longitude = -4.2501
and
   pastdays = 2
order by
   1;
```

```
+------------+---------------------+
| date       | precipitation_hours |
+------------+---------------------+
| 2024-01-21 | 19                  |
| 2024-01-22 | 15                  |
| 2024-01-23 | 15                  |
| 2024-01-24 | 7                   |
| 2024-01-25 | 14                  |
| 2024-01-26 | 8                   |
| 2024-01-27 | 11                  |
| 2024-01-28 | 13                  |
| 2024-01-29 | 0                   |
+------------+---------------------+
```

