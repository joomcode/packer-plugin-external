package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/packer-plugin-sdk/plugin"

	"github.com/joomcode/packer-plugin-external/datasource/json"
	"github.com/joomcode/packer-plugin-external/datasource/raw"
	"github.com/joomcode/packer-plugin-external/version"
)

func main() {
	pps := plugin.NewSet()
	pps.RegisterDatasource("json", new(json.Datasource))
	pps.RegisterDatasource("raw", new(raw.Datasource))
	pps.SetVersion(version.PluginVersion)
	err := pps.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
