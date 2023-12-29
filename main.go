package main

import (
    "github.com/turbot/steampipe-plugin-sdk/v5/plugin"
    "github.com/s0rbus/steampipe-plugin-openmeteo/openmeteo"
)

func main() {
    plugin.Serve(&plugin.ServeOpts{PluginFunc: openmeteo.Plugin})
}
