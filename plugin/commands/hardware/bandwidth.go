package hardware

import (
	"fmt"
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/virtual"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type BandwidthCommand struct {
	UI      terminal.UI
	Manager managers.HardwareServerManager
}

func NewBandwidthCommand(ui terminal.UI, manager managers.HardwareServerManager) (cmd *BandwidthCommand) {
	return &BandwidthCommand{
		UI:      ui,
		Manager: manager,
	}
}

func (cmd *BandwidthCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	VsID, err := utils.ResolveVirtualGuestId(c.Args()[0])
	if err != nil {
		return errors.NewInvalidSoftlayerIdInputError("Virtual server ID")
	}

	var start, end string
	var startDate, endDate time.Time

	if c.IsSet("start") {
		start = c.String("start")
		startDate, err = time.Parse(virtual.GetDateFormat(start), start)
		if err != nil {
			return errors.NewInvalidUsageError("Invalid start date: " + err.Error())
		}
	} else {
		startDate = time.Now()
	}
	if c.IsSet("end") {
		end = c.String("end")
		endDate, err = time.Parse(virtual.GetDateFormat(end), end)
		if err != nil {
			return errors.NewInvalidUsageError("Invalid end date: " + err.Error())
		}
	} else {
		endDate = startDate.AddDate(0, -1, 0)
	}

	rollupSeconds := 3600
	if c.IsSet("rollup") {
		rollupSeconds = c.Int("rollup")
	}
	// cmd.UI.Say(fmt.Sprintf("FORMAT: %v, Start: %v (%v), End: %v (%v)\n", GetDateFormat(start), startDate, start, endDate, end))
	bandwidthData, err := cmd.Manager.GetBandwidthData(VsID, startDate, endDate, rollupSeconds)
	if err != nil {
		fmt.Printf("ERR: %v", err)
		return err
	}
	// cmd.UI.Say(fmt.Sprintf("%+v", bandwidthData))

	summaryTable, bandwidthTable := virtual.BuildOutputTable(bandwidthData, cmd.UI)
	summaryTable.Print()
	if !c.IsSet("quite") {
		bandwidthTable.Print()
	}

	return nil
}
