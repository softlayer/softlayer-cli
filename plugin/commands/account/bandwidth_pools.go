package account

import (
	"fmt"
	"github.com/softlayer/softlayer-go/session"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type BandwidthPoolsCommand struct {
	UI      terminal.UI
	Session *session.Session
}

func NewBandwidthPoolsCommand(ui terminal.UI, session *session.Session) (cmd *BandwidthPoolsCommand) {
	return &BandwidthPoolsCommand{
		UI:      ui,
		Session: session,
	}
}

func BandwidthPoolsMetaData() cli.Command {
	return cli.Command{
		Category:    "account",
		Name:        "bandwidth-pools",
		Description: T("lists bandwidth pools"),
		Usage:       T(`${COMMAND_NAME} sl account bandwidth-pools`),
		Flags: []cli.Flag{
			metadata.OutputFlag(),
		},
	}
}

func (cmd *BandwidthPoolsCommand) Run(c *cli.Context) error {
	accountManager := managers.NewAccountManager(cmd.Session)
	pools, err := accountManager.GetBandwidthPools()
	if err != nil {
		return err
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}
	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, pools)
	}

	table := cmd.UI.Table([]string{
		T("ID"),
		T("Pool Name"),
		T("Region"),
		T("Servers"),
		T("Allocation"),
		T("Current Usage"),
		T("Projected Usage"),
	})
	for _, pool := range pools {
		curr_usage, proj_usage, allocation := "-", "-", "-"
		if pool.BillingCyclePublicBandwidthUsage != nil {
			curr_usage = fmt.Sprintf("%.2f GB", float64(*pool.BillingCyclePublicBandwidthUsage.AmountOut))
		}
		if pool.ProjectedPublicBandwidthUsage != nil {
			proj_usage = fmt.Sprintf("%.2f GB", float64(*pool.ProjectedPublicBandwidthUsage))
		}
		if pool.TotalBandwidthAllocated != nil {
			allocation = fmt.Sprintf("%d GB", uint(*pool.TotalBandwidthAllocated))
		}
		serverCount, _ := accountManager.GetBandwidthPoolServers(*pool.Id)
		table.Add(
			utils.FormatIntPointer(pool.Id),
			utils.FormatStringPointer(pool.Name),
			utils.FormatStringPointer(pool.LocationGroup.Name),
			fmt.Sprintf("%d", serverCount),
			allocation,
			curr_usage,
			proj_usage,
		)
	}

	table.Print()

	return nil
}
