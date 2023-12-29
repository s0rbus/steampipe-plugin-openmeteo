package openmeteo

import (
	"context"

	"github.com/hectormalot/omgo"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func connect(ctx context.Context, d *plugin.QueryData) (omgo.Client, error) {
   return omgo.NewClient()
}
