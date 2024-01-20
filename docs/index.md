---
organization: s0rbus
category: ["software development"]
icon_url: "/images/plugins/s0rbus/openmeteo.svg"
brand_color: "#FF8800"
display_name: "Open-Meteo"
short_name: "openmeteo"
description: "Steampipe plugin for querying weather from Open-Meteo."
og_description: "Query Open-Meteo weather forecast with SQL! Open source CLI. No DB required."
og_image: "/images/plugins/s0rbus/open-meteo-social-graphic.png"
---

# Open-Meteo + Steampipe

[Open-Meteo](https://open-meteo.com/ is an open-source weather API offerinf free access for non-commercial use.

[Steampipe](https://steampipe.io) is an open source CLI to instantly query cloud APIs using SQL.

For example:

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

```
+----------------------+----------------+---------------------------+---------------+-------------+----------------+
| time                 | temperature_2m | precipitation_probability | precipitation | cloud_cover | wind_speed_10m |
+----------------------+----------------+---------------------------+---------------+-------------+----------------+
| 2024-01-18T02:00:00Z | -2.9           | 0                         | 0             | 0           | 6.5            |
| 2024-01-19T04:00:00Z | -3             | 0                         | 0             | 0           | 3.9            |
etc
+----------------------+----------------+---------------------------+---------------+-------------+----------------+
```

## Documentation

- **[Table definitions & examples →](https://hub.steampipe.io/plugins/s0rbus/openmeteo/tables)**

## Get started

### Install

Download and install the latest Opem-Meteo plugin:

```bash
steampipe plugin install s0rbus/openmeteo
```

### Credentials

No API key required

### Configuration

Installing the latest Open-Meteo plugin will create a config file (`~/.steampipe/config/openmeteo.spc`) with a single connection named `openmeteo`:

```hcl
connection "openmeteo" {
    plugin    = "s0rbus/openmeteo"
}
```

## Get involved

- Open source: https://github.com/s0rbus/steampipe-plugin-openmeteo
- Community: [Slack Channel](https://join.slack.com/t/steampipe/shared_invite/zt-oij778tv-lYyRTWOTMQYBVAbtPSWs3g)
