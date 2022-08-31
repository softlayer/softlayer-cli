package order

import (
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/spf13/cobra"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type PackageLocationCommand struct {
	*metadata.SoftlayerCommand
	OrderManager managers.OrderManager
	Command      *cobra.Command
}

func NewPackageLocationCommand(sl *metadata.SoftlayerCommand) (cmd *PackageLocationCommand) {
	thisCmd := &PackageLocationCommand{
		SoftlayerCommand: sl,
		OrderManager:     managers.NewOrderManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "package-locations " + T("PACKAGE_KEYNAME"),
		Short: T("List datacenters a package can be ordered in"),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *PackageLocationCommand) Run(args []string) error {
	packageKeyname := args[0]

	outputFormat := cmd.GetOutputFlag()

	locations, err := cmd.OrderManager.PackageLocation(packageKeyname)
	if err != nil {
		return errors.NewAPIError(T("Failed to list package locations.\n"), err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, locations)
	}

	cmd.Print(locations)
	return nil
}

func (cmd *PackageLocationCommand) Print(locations []datatypes.Location_Region) {
	table := cmd.UI.Table([]string{T("id"), T("dc"), T("description"), T("keyName")})

	for _, region := range locations {
		for _, datacenter := range region.Locations {
			table.Add(utils.FormatIntPointer(datacenter.Location.Id),
				utils.FormatStringPointer(datacenter.Location.Name),
				utils.FormatStringPointer(region.Description),
				utils.FormatStringPointer(region.Keyname))
		}
	}
	table.Print()
}
