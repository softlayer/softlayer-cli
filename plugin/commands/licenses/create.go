package licenses

import (
	"github.com/spf13/cobra"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type CreateCommand struct {
	*metadata.SoftlayerCommand
	Command         *cobra.Command
	LicensesManager managers.LicensesManager
	Key             string
	Datacenter      string
}

func NewCreateCommand(sl *metadata.SoftlayerCommand) *CreateCommand {
	thisCmd := &CreateCommand{
		SoftlayerCommand: sl,
		LicensesManager:  managers.NewLicensesManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "create",
		Short: T("Order/Create License."),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().StringVar(&thisCmd.Key, "key", "", T("The VMware License Key. To get the required package you can use the command sl licenses create-options Package. E.g VMWARE_VSAN_ENTERPRISE_TIER_III_65_124_TB_6_X_2  [required]"))
	cobraCmd.Flags().StringVar(&thisCmd.Datacenter, "datacenter", "", T("Datacenter shortname  [required]"))
	
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *CreateCommand) Run(args []string) error {
	outputFormat := cmd.GetOutputFlag()

	if cmd.Datacenter == "" || cmd.Key == "" {
		return slErr.NewInvalidUsageError(T("This command requires two arguments."))
	}

	table := cmd.UI.Table([]string{T("Name"), T("Value")})

	orderLicense, err := cmd.LicensesManager.CreateLicense(cmd.Datacenter, cmd.Key)
	if err != nil {
		return slErr.NewAPIError(T("Failed to create the license."), err.Error(), 2)
	}
	table.Add("Id", utils.FormatIntPointer(orderLicense.OrderId))
	table.Add("Created", utils.FormatSLTimePointer(orderLicense.OrderDate))

	utils.PrintTable(cmd.UI, table, outputFormat)
	return nil
}
