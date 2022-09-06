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

type RouteCommand struct {
	*metadata.SoftlayerCommand
	NetworkManager managers.NetworkManager
	Command        *cobra.Command
	TypeId         string
	Type           string
}

func NewRouteCommand(sl *metadata.SoftlayerCommand) *RouteCommand {
	thisCmd := &RouteCommand{
		SoftlayerCommand: sl,
		NetworkManager:   managers.NewNetworkManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "route " + T("IDENTIFIER"),
		Short: T("This interface allows you to change the route of your Account Owned subnets."),
		Long: T(`${COMMAND_NAME} sl subnet route IDENTIFIER [OPTIONS]

EXAMPLE:
	${COMMAND_NAME} sl subnet route --type-id 1234567 --type SoftLayer_Network_Subnet_IpAddress 12345678
	This command allows you to change the route of your Account Owned subnets.`),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().StringVarP(&thisCmd.TypeId, "type-id", "i", "", T("An appropriate identifier for the specified $type, e.g. the identifier of a SoftLayer_Network_Subnet_IpAddress [required]"))
	cobraCmd.Flags().StringVarP(&thisCmd.Type, "type", "t", "", T("Type value in static routing e.g.: SoftLayer_Network_Subnet_IpAddress, SoftLayer_Hardware_Server [required]."))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *RouteCommand) Run(args []string) error {

	subnetID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Subnet ID")
	}

	outputFormat := cmd.GetOutputFlag()

	if cmd.Type == "" {
		return errors.NewInvalidUsageError(T("[-t/--type] is required."))
	}

	if cmd.TypeId == "" {
		return errors.NewInvalidUsageError(T("[-i/--type-id] is required."))
	}

	resp, err := cmd.NetworkManager.Route(subnetID, cmd.Type, cmd.TypeId)
	if err != nil {
		return errors.NewAPIError(T("Failed to route using the type: {{.TYPE}} and identifier: {{.IDENTIFIER}}.\n",
			map[string]interface{}{"TYPE": cmd.Type, "IDENTIFIER": cmd.TypeId}), err.Error(), 2)

	}
	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, resp)
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("The transaction to route is created, routes will be updated in one or two minutes."))
	return nil
}
