package openmeteo

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name:             "steampipe-plugin-openmeteo",
		DefaultTransform: transform.FromGo().NullIfZero(),
		TableMap: map[string]*plugin.Table{
			"openmeteo_hourly_weather": tableOpenmeteoHourlyWeather(),
			"openmeteo_daily_weather":  tableOpenmeteoDailyWeather(),
		},
	}
	return p
}
