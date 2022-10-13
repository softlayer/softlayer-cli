package placementgroup

import (
	"sort"

	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type PlacementGroupCreateOptionsCommand struct {
	*metadata.SoftlayerCommand
	PlaceGroupManager managers.PlaceGroupManager
	Command           *cobra.Command
}

func NewPlacementGroupCreateOptionsCommand(sl *metadata.SoftlayerCommand) (cmd *PlacementGroupCreateOptionsCommand) {
	thisCmd := &PlacementGroupCreateOptionsCommand{
		SoftlayerCommand:  sl,
		PlaceGroupManager: managers.NewPlaceGroupManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "create-options",
		Short: T("List options for creating a placement group"),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *PlacementGroupCreateOptionsCommand) Run(args []string) error {

	routers, err := cmd.PlaceGroupManager.GetRouter(nil, "")
	if err != nil {
		return errors.NewAPIError(T("Failed to list routers"), err.Error(), 2)
	}

	cmd.UI.Print("Available Router:")
	if len(routers) == 0 {
		cmd.UI.Print(T("No available router was found."))
	} else {
		sort.Sort(utils.RouterByHostname(routers))
		table := cmd.UI.Table([]string{T("Data Center"), T("Hostname"), T("Backend Router Id")})
		for _, router := range routers {
			center := "-"
			if router.TopLevelLocation != nil {
				center = utils.FormatStringPointer(router.TopLevelLocation.LongName)
			}
			table.Add(center,
				utils.FormatStringPointer(router.Hostname),
				utils.FormatIntPointer(router.Id))
		}
		table.Print()
	}

	rules, err := cmd.PlaceGroupManager.GetRules()
	if err != nil {
		return errors.NewAPIError(T("Failed to list rules"), err.Error(), 2)
	}
	cmd.UI.Print("Rules:")
	if len(rules) == 0 {
		cmd.UI.Print(T("No rules was found."))
	} else {
		ruleTable := cmd.UI.Table([]string{T("ID"), T("Name")})
		for _, rule := range rules {
			ruleTable.Add(utils.FormatIntPointer(rule.Id),
				utils.FormatStringPointer(rule.KeyName),
			)
		}
		ruleTable.Print()
	}

	return nil
}
