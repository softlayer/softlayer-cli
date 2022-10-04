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

type EditCommand struct {
	*metadata.SoftlayerCommand
	NetworkManager managers.NetworkManager
	Command        *cobra.Command
	Name           string
	Description    string
}

func NewEditCommand(sl *metadata.SoftlayerCommand) (cmd *EditCommand) {
	thisCmd := &EditCommand{
		SoftlayerCommand: sl,
		NetworkManager:   managers.NewNetworkManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "edit " + T("SECURITYGROUP_ID"),
		Short: T("Edit details of a security group"),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().StringVarP(&thisCmd.Name, "name", "n", "", T("The name of the security group"))
	cobraCmd.Flags().StringVarP(&thisCmd.Description, "description", "d", "", T("The description of the security group"))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *EditCommand) Run(args []string) error {
	groupID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Security group ID")
	}

	if cmd.Name == "" && cmd.Description == "" {
		return errors.NewInvalidUsageError(T("Either -n, --name or -d, --description is required to edit security group."))
	}
	err = cmd.NetworkManager.EditSecurityGroup(groupID, cmd.Name, cmd.Description)
	if err != nil {
		return errors.NewAPIError(T("Failed to edit security group {{.ID}}.\n", map[string]interface{}{"ID": groupID}), err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Security group {{.ID}} is updated.", map[string]interface{}{"ID": groupID}))
	return nil
}
