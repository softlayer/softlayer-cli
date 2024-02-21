package subnet

import (
	"strconv"

	"github.com/spf13/cobra"

	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
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
		Short: T("Removes the routing for a specified subnet, turning it into an unrouted portable subnet."),
		Args:  metadata.OneArgs,
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

	_, err = cmd.NetworkManager.ClearRoute(subnetID)
	if err != nil {
		return err
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("The transaction to clear the route is created, routes will be updated in one or two minutes."))
	return nil
}
