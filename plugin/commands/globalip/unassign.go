package globalip

import (
	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type UnassignCommand struct {
	*metadata.SoftlayerCommand
	NetworkManager managers.NetworkManager
	Command        *cobra.Command
	Details        bool
}

func NewUnassignCommand(sl *metadata.SoftlayerCommand) *UnassignCommand {
	thisCmd := &UnassignCommand{
		SoftlayerCommand: sl,
		NetworkManager:   managers.NewNetworkManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "unassign " + T("IDENTIFIER"),
		Short: T("Unassign a global IP from a target router or device."),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().BoolVar(&thisCmd.Details, "details", false, T("Shows a very detailed list of charges"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *UnassignCommand) Run(args []string) error {
	outputFormat := cmd.GetOutputFlag()

	globalIPID, err := utils.ResolveGloablIPId(args[0])
	if err != nil {
		return errors.NewInvalidSoftlayerIdInputError("Globalip ID")
	}

	resp, err := cmd.NetworkManager.UnassignGlobalIP(globalIPID)
	if err != nil {
		return errors.NewAPIError(T("Failed to unassign global IP {{.ID}}.", map[string]interface{}{"ID": globalIPID}), err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, resp)
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("The transaction to unroute a global IP address is created, routes will be updated in one or two minutes."))
	return nil
}
