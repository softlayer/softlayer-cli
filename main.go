package main

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	sl_plugin "github.ibm.com/SoftLayer/softlayer-cli/plugin"
)

func main() {
	plugin.Start(new(sl_plugin.SoftlayerPlugin))
}
