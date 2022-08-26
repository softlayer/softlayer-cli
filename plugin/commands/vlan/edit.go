package vlan

import (
	"strconv"

	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type EditCommand struct {
	*metadata.SoftlayerCommand
	NetworkManager managers.NetworkManager
	Command        *cobra.Command
	Name           string
}

func NewEditCommand(sl *metadata.SoftlayerCommand) *EditCommand {
	thisCmd := &EditCommand{
		SoftlayerCommand: sl,
		NetworkManager:   managers.NewNetworkManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "edit " + T("IDENTIFIER"),
		Short: T("Edit the details about a VLAN."),
		Long: T(`${COMMAND_NAME} sl vlan edit IDENTIFIER [OPTIONS]
	
EXAMPLE:
	${COMMAND_NAME} sl vlan edit 12345678 -n myvlan-rename
	This command updates vlan with ID 12345678 and gives it a new name "myvlan-rename".`),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().StringVarP(&thisCmd.Name, "name", "n", "", T("The name of the VLAN"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *EditCommand) Run(args []string) error {
	vlanID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("VLAN ID")
	}

	if cmd.Name == "" {
		return errors.NewMissingInputError("-n|--name")
	}

	err = cmd.NetworkManager.EditVlan(vlanID, cmd.Name)
	if err != nil {
		return errors.NewAPIError(T("Failed to edit VLAN: {{.VlanID}}.\n", map[string]interface{}{"VlanID": vlanID}), err.Error(), 2)
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("VLAN {{.VlanID}} was updated.", map[string]interface{}{"VlanID": vlanID}))
	return nil
}
