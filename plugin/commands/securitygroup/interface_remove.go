package securitygroup

import (
	"strconv"

	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type InterfaceRemoveCommand struct {
	*metadata.SoftlayerCommand
	NetworkManager   managers.NetworkManager
	VSManager        managers.VirtualServerManager
	Command          *cobra.Command
	NetworkComponent int
	Server           int
	Interface        string
}

func NewInterfaceRemoveCommand(sl *metadata.SoftlayerCommand) (cmd *InterfaceRemoveCommand) {
	thisCmd := &InterfaceRemoveCommand{
		SoftlayerCommand: sl,
		NetworkManager:   managers.NewNetworkManager(sl.Session),
		VSManager:		  managers.NewVirtualServerManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "interface-remove " + T("SECURITYGROUP_ID"),
		Short: T("Detach an interface from a security group"),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().IntVarP(&thisCmd.NetworkComponent, "network-component", "n", 0, T("The network component to remove from the security group"))
	cobraCmd.Flags().IntVarP(&thisCmd.Server, "server", "s", 0, T(" The server ID to remove from the security group"))
	cobraCmd.Flags().StringVarP(&thisCmd.Interface, "interface", "i", "", T("The interface of the server to remove (public or private)"))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *InterfaceRemoveCommand) Run(args []string) error {
	groupID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Security group ID")
	}

	networkComponent := cmd.NetworkComponent
	serverID := cmd.Server
	serverInterface := cmd.Interface
	err = ValidateArgs(networkComponent, serverID, serverInterface)
	if err != nil {
		return err
	}
	componentID, err := GetComponentId(cmd.VSManager, networkComponent, serverID, serverInterface)
	if err != nil {
		return err
	}
	err = cmd.NetworkManager.DetachSecurityGroupComponent(groupID, componentID)
	if err != nil {
		return errors.NewAPIError(T("Failed to remove network component {{.ComponentID}} from security group {{.GroupID}}.\n",
			map[string]interface{}{"GroupID": groupID, "ComponentID": componentID}), err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Network component {{.ComponentID}} is removed from security group {{.GroupID}}.",
		map[string]interface{}{"GroupID": groupID, "ComponentID": componentID}))
	return nil
}
