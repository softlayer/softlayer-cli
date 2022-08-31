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

type PackageListCommand struct {
	*metadata.SoftlayerCommand
	OrderManager managers.OrderManager
	Command      *cobra.Command
	Keyword      string
	PackageType  string
}

func NewPackageListCommand(sl *metadata.SoftlayerCommand) (cmd *PackageListCommand) {
	thisCmd := &PackageListCommand{
		SoftlayerCommand: sl,
		OrderManager:     managers.NewOrderManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "package-list",
		Short: T("List packages that can be ordered with the placeOrder API"),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().StringVar(&thisCmd.Keyword, "keyword", "", T("A word (or string) that is used to filter package names"))
	cobraCmd.Flags().StringVar(&thisCmd.PackageType, "package-type", "", T("The keyname for the type of package. For example, BARE_METAL_CPU"))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *PackageListCommand) Run(args []string) error {

	outputFormat := cmd.GetOutputFlag()

	keyword := cmd.Keyword
	packageType := cmd.PackageType

	packages, err := cmd.OrderManager.ListPackage(keyword, packageType)
	if err != nil {
		return errors.NewAPIError(T("Failed to list packages.\n"), err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, packages)
	}

	cmd.Print(packages)
	return nil
}

func (cmd *PackageListCommand) Print(packages []datatypes.Product_Package) {
	table := cmd.UI.Table([]string{T("id"), T("name"), T("keyName"), T("type")})

	for _, pac := range packages {
		table.Add(utils.FormatIntPointer(pac.Id),
			utils.FormatStringPointer(pac.Name),
			utils.FormatStringPointer(pac.KeyName),
			utils.FormatStringPointer(pac.Type.KeyName))
	}
	table.Print()
}
