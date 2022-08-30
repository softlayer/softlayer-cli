package subnet

import (
	"strconv"

	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type ClearRouteCommand struct {
	*metadata.SoftlayerCommand
	NetworkManager managers.NetworkManager
	Command        *cobra.Command
}

func NewClearRouteCommand(sl *metadata.SoftlayerCommand) *ClearRouteCommand {
	thisCmd := &ClearRouteCommand{
		SoftlayerCommand: sl,
		NetworkManager:   managers.NewNetworkManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "clear-route " + T("IDENTIFIER"),
		Short: T("This interface allows you to remove the route of your Account Owned subnets."),
		Long: T(`${COMMAND_NAME} sl subnet clear-route IDENTIFIER [OPTIONS]

EXAMPLE:
	${COMMAND_NAME} sl subnet clear-route 12345678
	This command allows you to remove the route of your Account Owned subnets.`),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *ClearRouteCommand) Run(args []string) error {
	subnetID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Subnet ID")
	}
	outputFormat := cmd.GetOutputFlag()

	resp, err := cmd.NetworkManager.ClearRoute(subnetID)
	if err != nil {
		return errors.NewAPIError(T("Failed to clear the route for the subnet: {{.ID}}.\n", map[string]interface{}{"ID": subnetID}), err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, resp)
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("The transaction to clear the route is created, routes will be updated in one or two minutes."))
	return nil
}
