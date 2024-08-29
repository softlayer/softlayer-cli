package hardware

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/virtual"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type BandwidthCommand struct {
	*metadata.SoftlayerCommand
	HardwareManager managers.HardwareServerManager
	Command         *cobra.Command
	Start           string
	End             string
	Rollup          int
	quiet           bool
}

func NewBandwidthCommand(sl *metadata.SoftlayerCommand) (cmd *BandwidthCommand) {
	thisCmd := &BandwidthCommand{
		SoftlayerCommand: sl,
		HardwareManager:  managers.NewHardwareServerManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "bandwidth " + T("IDENTIFIER"),
		Short: T("Bandwidth data over date range."),
		Long: T(`${COMMAND_NAME} sl {{.Command}} bandwidth IDENTIFIER [OPTIONS]
Time formats that are either '2006-01-02', '2006-01-02T15:04' or '2006-01-02T15:04-07:00'

Due to some rounding and date alignment details, results here might be slightly different than results in the control portal.
Bandwidth is listed in GB, if no time zone is specified, GMT+0 is assumed.

Example::

   ${COMMAND_NAME} sl {{.Command}} bandwidth 1234 -s 2006-01-02T15:04 -e 2006-01-02T15:04-07:00`, map[string]interface{}{"Command": "hardware"}),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().StringVarP(&thisCmd.Start, "start", "s", "", T("Start date for bandwdith reporting"))
	cobraCmd.Flags().StringVarP(&thisCmd.End, "end", "e", "", T("End date for bandwidth reporting"))
	cobraCmd.Flags().IntVarP(&thisCmd.Rollup, "rollup", "r", 0, T("Number of seconds to report as one data point. 300, 600, 1800, 3600 (default), 43200 or 86400 seconds"))
	cobraCmd.Flags().BoolVarP(&thisCmd.quiet, "quiet", "q", false, T("Only show the summary table."))
	cobraCmd.Flags().SetNormalizeFunc(utils.NormalizeQuietFlag)

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *BandwidthCommand) Run(args []string) error {
	VsID, err := utils.ResolveVirtualGuestId(args[0])
	if err != nil {
		return errors.NewInvalidSoftlayerIdInputError("Virtual server ID")
	}

	var start, end string
	var startDate, endDate time.Time

	if cmd.Start != "" {
		start = cmd.Start
		startDate, err = time.Parse(virtual.GetDateFormat(start), start)
		if err != nil {
			return errors.NewInvalidUsageError("Invalid start date: " + err.Error())
		}
	} else {
		startDate = time.Now()
	}
	if cmd.End != "" {
		end = cmd.End
		endDate, err = time.Parse(virtual.GetDateFormat(end), end)
		if err != nil {
			return errors.NewInvalidUsageError("Invalid end date: " + err.Error())
		}
	} else {
		endDate = startDate.AddDate(0, -1, 0)
	}

	rollupSeconds := 3600
	if cmd.Rollup != 0 {
		rollupSeconds = cmd.Rollup
	}
	// cmd.UI.Say(fmt.Sprintf("FORMAT: %v, Start: %v (%v), End: %v (%v)\n", GetDateFormat(start), startDate, start, endDate, end))
	bandwidthData, err := cmd.HardwareManager.GetBandwidthData(VsID, startDate, endDate, rollupSeconds)
	if err != nil {
		fmt.Printf("ERR: %v", err)
		return err
	}
	// cmd.UI.Say(fmt.Sprintf("%+v", bandwidthData))

	summaryTable, bandwidthTable := virtual.BuildOutputTable(bandwidthData, cmd.UI)
	summaryTable.Print()
	if !cmd.quiet {
		bandwidthTable.Print()
	}

	return nil
}
