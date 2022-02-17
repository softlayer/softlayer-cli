package securitygroup

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type InterfaceRemoveCommand struct {
	UI             terminal.UI
	NetworkManager managers.NetworkManager
	VSManager      managers.VirtualServerManager
}

func NewInterfaceRemoveCommand(ui terminal.UI, networkManager managers.NetworkManager, vsManager managers.VirtualServerManager) (cmd *InterfaceRemoveCommand) {
	return &InterfaceRemoveCommand{
		UI:             ui,
		NetworkManager: networkManager,
		VSManager:      vsManager,
	}
}

func (cmd *InterfaceRemoveCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	groupID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Security group ID")
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
	err = cmd.NetworkManager.DetachSecurityGroupComponent(groupID, componentID)
	if err != nil {
		return cli.NewExitError(T("Failed to remove network component {{.ComponentID}} from security group {{.GroupID}}.\n",
			map[string]interface{}{"GroupID": groupID, "ComponentID": componentID})+err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Network component {{.ComponentID}} is removed from security group {{.GroupID}}.",
		map[string]interface{}{"GroupID": groupID, "ComponentID": componentID}))
	return nil
}

func SecurityGroupInterfaceRemoveMetaData() cli.Command {
	return cli.Command{
		Category:    "securitygroup",
		Name:        "interface-remove",
		Description: T("Detach an interface from a security group"),
		Usage:       "${COMMAND_NAME} sl securitygroup interface-remove SECURITYGROUP_ID [OPTIONS]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "n,network-component",
				Usage: T("The network component to remove from the security group"),
			},
			cli.StringFlag{
				Name:  "s,server",
				Usage: T(" The server ID to remove from the security group"),
			},
			cli.StringFlag{
				Name:  "i,interface",
				Usage: T("The interface of the server to remove (public or private)"),
			},
		},
	}
}
