package securitygroup

import (
	"errors"
	"strconv"
	"strings"

	"github.com/softlayer/softlayer-go/datatypes"
	bmxErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type InterfaceAddCommand struct {
	UI             terminal.UI
	NetworkManager managers.NetworkManager
	VSManager      managers.VirtualServerManager
}

func NewInterfaceAddCommand(ui terminal.UI, networkManager managers.NetworkManager, vsManager managers.VirtualServerManager) (cmd *InterfaceAddCommand) {
	return &InterfaceAddCommand{
		UI:             ui,
		NetworkManager: networkManager,
		VSManager:      vsManager,
	}
}

func (cmd *InterfaceAddCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return bmxErr.NewInvalidUsageError(T("This command requires one argument."))
	}
	groupID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("Security group ID")
	}

	networkComponent := c.Int("n")
	serverID := c.Int("s")
	serverInterface := c.String("i")
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
		return cli.NewExitError(T("Failed to add network component {{.ComponentID}} to security group {{.GroupID}}.\n",
			map[string]interface{}{"GroupID": groupID, "ComponentID": componentID})+err.Error(), 2)
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
		vs, err := vsManager.GetInstance(serverId, "networkComponents[id,port]")
		if err != nil {
			return 0, err
		}
		port := 0
		if strings.ToLower(serverInterface) == "public" {
			port = 1
		}
		var component []datatypes.Virtual_Guest_Network_Component
		for _, c := range vs.NetworkComponents {
			if c.Port != nil && *c.Port == port {
				component = append(component, c)
			}
		}
		if len(component) != 1 {
			return 0, errors.New(T("Instance {{.ServerID}} has {{.Count}} {{.Interface}} interface.",
				map[string]interface{}{"ServerID": serverId, "Interface": serverInterface, "Count": len(component)}))
		}
		return *component[0].Id, nil
	}
	return networkComponent, nil
}
