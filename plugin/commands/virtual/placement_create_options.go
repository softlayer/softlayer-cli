package virtual

import (
	"github.com/spf13/cobra"

	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type PlacementGroupCreateOptionsCommand struct {
	*metadata.SoftlayerCommand
	VirtualServerManager managers.VirtualServerManager
	Command              *cobra.Command
}

func NewPlacementGroupCreateOptionsCommand(sl *metadata.SoftlayerCommand) (cmd *PlacementGroupCreateOptionsCommand) {
	thisCmd := &PlacementGroupCreateOptionsCommand{
		SoftlayerCommand:     sl,
		VirtualServerManager: managers.NewVirtualServerManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "placementgroup-create-options",
		Short: T("Get List options for creating a placement group.."),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *PlacementGroupCreateOptionsCommand) Run(args []string) error {
	datacenters, err := cmd.VirtualServerManager.GetDatacenters()
	if err != nil {
		return slErrors.NewAPIError("Internal error.", err.Error(), 2)
	}
	tableRegion := cmd.UI.Table([]string{T("Datacenter"), T("Hostname"), T("BackendRouterId")})
	for _, datacenter := range datacenters {
		routers, err := cmd.VirtualServerManager.GetAvailablePlacementRouters(utils.IntPointertoInt(datacenter.Id))
		if err != nil {
			return slErrors.NewAPIError("Internal error.", err.Error(), 2)
		}
		for _, routerAvalaible := range routers {
			tableRegion.Add(utils.FormatStringPointer(datacenter.LongName), utils.FormatStringPointer(routerAvalaible.Hostname), utils.FormatIntPointer(routerAvalaible.Id))
		}
	}

	rules, err := cmd.VirtualServerManager.GetRules()
	if err != nil {
		return slErrors.NewAPIError("Internal error.", err.Error(), 2)
	}
	tableRules := cmd.UI.Table([]string{T("Id"), T("Rule")})
	for _, rule := range rules {
		tableRules.Add(utils.FormatIntPointer(rule.Id), utils.FormatStringPointer(rule.KeyName))

	}

	outputFormat := cmd.GetOutputFlag()
	utils.PrintTable(cmd.UI, tableRegion, outputFormat)
	utils.PrintTable(cmd.UI, tableRules, outputFormat)

	return nil
}
