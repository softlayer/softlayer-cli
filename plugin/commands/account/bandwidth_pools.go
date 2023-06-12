package account

import (
	"fmt"

	"github.com/spf13/cobra"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type BandwidthPoolsCommand struct {
	*metadata.SoftlayerCommand
	AccountManager managers.AccountManager
	Command        *cobra.Command
}

func NewBandwidthPoolsCommand(sl *metadata.SoftlayerCommand) *BandwidthPoolsCommand {
	thisCmd := &BandwidthPoolsCommand{
		SoftlayerCommand: sl,
		AccountManager:   managers.NewAccountManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "bandwidth-pools",
		Short: T("Displays bandwidth pool information."),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *BandwidthPoolsCommand) Run(args []string) error {
	pools, err := cmd.AccountManager.GetBandwidthPools()
	if err != nil {
		return err
	}

	outputFormat := cmd.GetOutputFlag()

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
		serverCount, _ := cmd.AccountManager.GetBandwidthPoolServers(*pool.Id)
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
