package openmeteo

import (
	"context"
	"time"

	"github.com/hectormalot/omgo"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableOpenmeteoDailyWeather() *plugin.Table {
	return &plugin.Table{
		Name:        "openmeteo_daily_weather",
		Description: "Get Openmeteo daily weather.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "latitude", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "longitude", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "pastdays", Require: plugin.Optional, Operators: []string{"="}},
				{Name: "forecastdays", Require: plugin.Optional, Operators: []string{"="}},
			},
			Hydrate: tableDailyWeatherList,
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
				Name:        "temperature_2m_min",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("Temperature_2m_min"),
				Description: "min temperature at 2m",
			},
			{
				Name:        "temperature_2m_max",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("Temperature_2m_max"),
				Description: "max temperature at 2m",
			},
			{
				Name:        "apparent_temperature_min",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("Apparent_Temperature_min"),
				Description: "Apparent min temperature",
			},
			{
				Name:        "apparent_temperature_max",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("Apparent_Temperature_max"),
				Description: "Apparent max temperature",
			},
			{
				Name:        "precipitation_probability_max",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("PrecipProbMax"),
				Description: "precipitation probability max",
			},
			{
				Name:        "precipitation_sum",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("Precipitation_Sum"),
				Description: "precipitation sum (rain, showers, snow)",
			},
			{
				Name:        "precipitation_hours",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("Precipitation_Hours"),
				Description: "precipitation hours (rain, showers, snow)",
			},
			{
				Name:        "wind_speed_10m_max",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("WindSpeed10mMax"),
				Description: "wind speed 10m max (mph)",
			},
			{
				Name:        "wind_gusts_10m_max",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("WindGusts10mMax"),
				Description: "wind gusts 10m max (mph)",
			},
			{
				Name:        "wind_direction_10m_dominant",
				Type:        proto.ColumnType_DOUBLE,
				Transform:   transform.FromField("WindDirection10mDominant"),
				Description: "wind direction 10m direction (degrees)",
			},
		},
	}
}

func tableDailyWeatherList(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
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

	opts := omgo.Options{WindspeedUnit: "mph", Timezone: "Europe/London", DailyMetrics: []string{
		"temperature_2m_min",
		"temperature_2m_max",
		"apparent_temperature_min",
		"apparent_temperature_max",
		"precipitation_probability_max",
		"precipitation_sum",
		"precipitation_hours",
		"wind_speed_10m_max",
		"wind_gusts_10m_max",
		"wind_direction_10m_dominant",
	},
		PastDays: pd,
	}
	forecast, err := c.Forecast(ctx, loc, &opts)
	if err != nil {
		plugin.Logger(ctx).Error("openmeteo_daily_weather.tableDailyWeatherList", "forecast api error", err)
		return nil, err
	}

	type Row struct {
		Time                     time.Time
		Temperature_2m_min       float64
		Temperature_2m_max       float64
		Apparent_Temperature_min float64
		Apparent_Temperature_max float64
		PrecipProbMax            float64
		Precipitation_Sum        float64
		Precipitation_Hours      float64
		WindSpeed10mMax          float64
		WindGusts10mMax          float64
		WindDirection10mDominant float64
	}

	//At the moment the Open-Meteo client does not implement forecast days option so
	//implementing it manually for now. This means can only do up to default (7), not up to
	//API max of 16 days
	dayLimit := time.Now().AddDate(0, 0, fd-1)
	for i := 0; i < len(forecast.DailyTimes); i++ {
		if fd != 7 && forecast.DailyTimes[i].After(dayLimit) {
			break
		}
		d.StreamListItem(ctx, Row{
			Time:                     forecast.DailyTimes[i],
			Temperature_2m_min:       forecast.DailyMetrics["temperature_2m_min"][i],
			Temperature_2m_max:       forecast.DailyMetrics["temperature_2m_max"][i],
			Apparent_Temperature_min: forecast.DailyMetrics["apparent_temperature_min"][i],
			Apparent_Temperature_max: forecast.DailyMetrics["apparent_temperature_max"][i],
			PrecipProbMax:            forecast.DailyMetrics["precipitation_probability_max"][i],
			Precipitation_Sum:        forecast.DailyMetrics["precipitation_sum"][i],
			Precipitation_Hours:      forecast.DailyMetrics["precipitation_hours"][i],
			WindSpeed10mMax:          forecast.DailyMetrics["wind_speed_10m_max"][i],
			WindGusts10mMax:          forecast.DailyMetrics["wind_gusts_10m_max"][i],
			WindDirection10mDominant: forecast.DailyMetrics["wind_direction_10m_dominant"][i],
		})
		// Check if context has been cancelled or if the limit has been hit (if specified)
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	return nil, nil
}
