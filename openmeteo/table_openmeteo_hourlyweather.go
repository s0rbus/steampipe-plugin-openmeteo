package openmeteo

import (
	"context"
	"time"

	"github.com/hectormalot/omgo"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableOpenmeteoHourlyWeather() *plugin.Table {
	return &plugin.Table{
		Name:        "openmeteo_hourly_weather",
		Description: "Get Openmeteo hourly weather.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "latitude", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "longitude", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "pastdays", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "forecastdays", Require: plugin.Optional, Operators: []string{"="}},
			},
			Hydrate: tableHourlyWeatherList,
		},
		Columns: []*plugin.Column{
			{
				Name:      "latitude",
				Type:      proto.ColumnType_DOUBLE,
				Transform: transform.FromQual("latitude"),
			},
			{
				Name:      "longitude",
				Type:      proto.ColumnType_DOUBLE,
				Transform: transform.FromQual("longitude"),
			},
			{
				Name:        "time",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Time"),
				Description: "observation timestamp",
			},
			{
				Name:        "pastdays",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromQual("pastdays"),
				Description: "number of days before current day to show, default 0",
			},
			{
				Name:        "forecastdays",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromQual("forecastdays"),
				Description: "number of days forecast to show, default 7, max 7",
			},
			{
				Name:        "temperature_2m",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("Temperature_2m"),
				Description: "temperature at 2m",
			},
			{
				Name:        "precipitation_probability",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("PrecipProb"),
				Description: "precipitation probability",
			},
			{
				Name:        "precipitation",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("Precipitation"),
				Description: "precipitation (rain, showers, snow)",
			},
			{
				Name:        "cloud_cover",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("CloudCover"),
				Description: "cloud cover (percentage)",
			},
			{
				Name:        "wind_speed_10m",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("WindSpeed10m"),
				Description: "wind speed 10m (mph)",
			},
		},
	}
}

func tableHourlyWeatherList(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	c, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}
	lat := d.EqualsQuals["latitude"].GetDoubleValue()
	lng := d.EqualsQuals["longitude"].GetDoubleValue()
	loc, err := omgo.NewLocation(lat, lng)
	if err != nil {
		return nil, err
	}
	pd := 0
	if d.Quals["pastdays"] != nil {
		pd = int(d.EqualsQuals["pastdays"].GetInt64Value())
	}
	fd := 7
	if d.Quals["forecastdays"] != nil {
		fd = int(d.EqualsQuals["forecastdays"].GetInt64Value())
		if fd == 0 || fd > 7 {
			fd = 7 //cannot be more than default at the moment, see below
		}
	}
	opts := omgo.Options{WindspeedUnit: "mph", HourlyMetrics: []string{
		"temperature_2m",
		"precipitation_probability",
		"precipitation",
		"cloud_cover",
		"wind_speed_10m",
	},
		PastDays: pd,
	}
	forecast, err := c.Forecast(ctx, loc, &opts)
	if err != nil {
		plugin.Logger(ctx).Error("forecast error", "err", err)
		return nil, err
	}

	type Row struct {
		Time           time.Time
		Temperature_2m float64
		PrecipProb     float64
		Precipitation  float64
		CloudCover     float64
		WindSpeed10m   float64
	}

	//plugin.Logger(ctx).Info("hourlytimes", "len", len(forecast.HourlyTimes))

	//At the moment the Open-Meteo client does not implement forecast days option so
	//implementing it manually for now. This means can only do up to default (7), not up to
	//API max of 16 days
	dayLimit := time.Now().AddDate(0, 0, fd-1)
	for i := 0; i < len(forecast.HourlyTimes); i++ {
		if fd != 7 && forecast.HourlyTimes[i].After(dayLimit) {
			break
		}
		d.StreamListItem(ctx, Row{
			Time:           forecast.HourlyTimes[i],
			Temperature_2m: forecast.HourlyMetrics["temperature_2m"][i],
			PrecipProb:     forecast.HourlyMetrics["precipitation_probability"][i],
			Precipitation:  forecast.HourlyMetrics["precipitation"][i],
			CloudCover:     forecast.HourlyMetrics["cloud_cover"][i],
			WindSpeed10m:   forecast.HourlyMetrics["wind_speed_10m"][i],
		})
	}
	return nil, nil
}
