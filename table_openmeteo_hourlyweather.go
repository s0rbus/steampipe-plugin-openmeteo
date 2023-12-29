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
	opts := omgo.Options{WindspeedUnit: "mph", HourlyMetrics: []string{
		"temperature_2m",
		"precipitation_probability",
		"precipitation",
		"cloud_cover",
		"wind_speed_10m",
	}}
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

	for i := 0; i < len(forecast.HourlyTimes); i++ {
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
