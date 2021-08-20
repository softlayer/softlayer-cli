package virtual

import (
	"time"
	"fmt"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type BandwidthCommand struct {
	UI                   terminal.UI
	VirtualServerManager managers.VirtualServerManager
}

func NewBandwidthCommand(ui terminal.UI, virtualServerManager managers.VirtualServerManager) (cmd *BandwidthCommand) {
	return &BandwidthCommand{
		UI:                   ui,
		VirtualServerManager: virtualServerManager,
	}
}

func (cmd *BandwidthCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	VsID, err := utils.ResolveVirtualGuestId(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Virtual server ID")
	}

	startDate := time.Now()
	endDate := startDate.AddDate(0, -1, 0)
	bandwidthData, err := cmd.VirtualServerManager.GetBandwidthData(VsID, startDate, endDate, 3600)
	fmt.Printf("%+v", bandwidthData)
	return nil
}	
