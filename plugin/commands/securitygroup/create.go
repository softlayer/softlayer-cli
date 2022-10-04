package securitygroup

import (
	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type CreateCommand struct {
	*metadata.SoftlayerCommand
	NetworkManager managers.NetworkManager
	Command        *cobra.Command
	Name           string
	Description    string
}

func NewCreateCommand(sl *metadata.SoftlayerCommand) (cmd *CreateCommand) {
	thisCmd := &CreateCommand{
		SoftlayerCommand: sl,
		NetworkManager:   managers.NewNetworkManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "create",
		Short: T("Create a security group"),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().StringVarP(&thisCmd.Name, "name", "n", "", T("The name of the security group"))
	cobraCmd.Flags().StringVarP(&thisCmd.Description, "description", "d", "", T("The description of the security group"))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *CreateCommand) Run(args []string) error {
	outputFormat := cmd.GetOutputFlag()

	group, err := cmd.NetworkManager.CreateSecurityGroup(cmd.Name, cmd.Description)
	if err != nil {
		return errors.NewAPIError(T("Failed to create security group with name {{.Name}}.\n",
			map[string]interface{}{"Name": cmd.Name}), err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, group)
	}
	table := cmd.UI.Table([]string{T("Name"), T("Value")})
	table.Add(T("ID"), utils.FormatIntPointer(group.Id))
	table.Add(T("Name"), utils.FormatStringPointer(group.Name))
	table.Add(T("Description"), utils.FormatStringPointer(group.Description))
	table.Add(T("Created"), utils.FormatSLTimePointer(group.CreateDate))
	table.Print()
	return nil
}
