package securitygroup

import (
	"errors"
	"strconv"
	"strings"

	bmxErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"

	"github.com/spf13/cobra"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type InterfaceAddCommand struct {
	*metadata.SoftlayerCommand
	NetworkManager   managers.NetworkManager
	VSManager        managers.VirtualServerManager
	Command          *cobra.Command
	NetworkComponent int
	Server           int
	Interface        string
}

func NewInterfaceAddCommand(sl *metadata.SoftlayerCommand) (cmd *InterfaceAddCommand) {
	thisCmd := &InterfaceAddCommand{
		SoftlayerCommand: sl,
		NetworkManager:   managers.NewNetworkManager(sl.Session),
		VSManager:		  managers.NewVirtualServerManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "interface-add " + T("SECURITYGROUP_ID"),
		Short: T("Attach an interface to a security group"),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().IntVarP(&thisCmd.NetworkComponent, "network-component", "n", 0, T("The network component ID to associate with the security group"))
	cobraCmd.Flags().IntVarP(&thisCmd.Server, "server", "s", 0, T(" The server ID to associate with the security group"))
	cobraCmd.Flags().StringVarP(&thisCmd.Interface, "interface", "i", "", T("The interface of the server to associate (public/private)"))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *InterfaceAddCommand) Run(args []string) error {
	groupID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("Security group ID")
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
	err = cmd.NetworkManager.AttachSecurityGroupComponent(groupID, componentID)
	if err != nil {
		return slErrors.NewAPIError(T("Failed to add network component {{.ComponentID}} to security group {{.GroupID}}.\n",
			map[string]interface{}{"GroupID": groupID, "ComponentID": componentID}), err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Network component {{.ComponentID}} is added to security group {{.GroupID}}.",
		map[string]interface{}{"GroupID": groupID, "ComponentID": componentID}))
	return nil
}

func ValidateArgs(networkComponent int, serverId int, serverInterface string) error {
	useComponent := networkComponent != 0 && (serverId == 0 && serverInterface == "")
	useServer := networkComponent == 0 && (serverId != 0 && serverInterface != "")
	if (useComponent && useServer) || (!useComponent && !useServer) {
		return bmxErr.NewInvalidUsageError(T("Must set either -n|--network-component or both -s|--server and -i|--interface"))
	}
	if useServer && strings.ToLower(serverInterface) != "public" && strings.ToLower(serverInterface) != "private" {
		return bmxErr.NewInvalidUsageError(T("-i|--interface must be either public or private"))
	}
	return nil
}

func GetComponentId(vsManager managers.VirtualServerManager, networkComponent int, serverId int, serverInterface string) (int, error) {
	useServer := networkComponent == 0 && (serverId != 0 && serverInterface != "")

	if useServer {
		vs, err := vsManager.GetInstance(serverId, "primaryBackendNetworkComponent[id,port], primaryNetworkComponent[id,port]")
		if err != nil {
			return 0, err
		}

		if strings.ToLower(serverInterface) == "public" {
			if vs.PrimaryNetworkComponent != nil && vs.PrimaryNetworkComponent.Id != nil {
				return *vs.PrimaryNetworkComponent.Id, nil
			}
		} else {
			if vs.PrimaryBackendNetworkComponent != nil && vs.PrimaryBackendNetworkComponent.Id != nil {
				return *vs.PrimaryBackendNetworkComponent.Id, nil			
			}
		}

		return 0, errors.New(
			T("Instance {{.ServerID}} has {{.Count}} {{.Interface}} interface.",
			map[string]interface{}{"ServerID": serverId, "Interface": serverInterface, "Count": 0}))

	}
	return networkComponent, nil
}
