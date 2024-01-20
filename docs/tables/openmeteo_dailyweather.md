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

### Get the count of content type

```sql
select
  "type",
  count(*)
from
  confluence_content
group by "type";
```


