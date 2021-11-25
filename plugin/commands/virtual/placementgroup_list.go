package virtual

import (
	"fmt"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type PlacementGroupListCommand struct {
	UI                   terminal.UI
	VirtualServerManager managers.VirtualServerManager
}

func NewPlacementGroupListCommand(ui terminal.UI, virtualServerManager managers.VirtualServerManager) (cmd *PlacementGroupListCommand) {
	return &PlacementGroupListCommand{
		UI:                   ui,
		VirtualServerManager: virtualServerManager,
	}
}

func (cmd *PlacementGroupListCommand) Run(c *cli.Context) error {
	fmt.Println("hello")
		return nil
}