package subnet

import (
	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type EditCommand struct {
	*metadata.SoftlayerCommand
	NetworkManager managers.NetworkManager
	Command        *cobra.Command
	Note           string
	Tags           string
}

func NewEditCommand(sl *metadata.SoftlayerCommand) *EditCommand {
	thisCmd := &EditCommand{
		SoftlayerCommand: sl,
		NetworkManager:   managers.NewNetworkManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "edit " + T("IDENTIFIER"),
		Short: T("Edit note and tags of a subnet."),
		Long: T(`${COMMAND_NAME} sl subnet edit IDENTIFIER [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl subnet edit 12345678 --note myNote
   ${COMMAND_NAME} sl subnet edit 12345678 --tags tag1`),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().StringVar(&thisCmd.Note, "note", "", T("The note"))
	cobraCmd.Flags().StringVar(&thisCmd.Tags, "tags", "", T("Comma separated list of tags, enclosed in quotes. 'tag1, tag2'"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *EditCommand) Run(args []string) error {
	subnetID, err := utils.ResolveSubnetId(args[0])
	if err != nil {
		return errors.NewInvalidSoftlayerIdInputError("Subnet ID")
	}

	if cmd.Tags == "" && cmd.Note == "" {
		return errors.NewInvalidUsageError(T("Please pass at least one of the flags."))
	}

	if cmd.Tags != "" {
		tags := cmd.Tags
		response, err := cmd.NetworkManager.SetSubnetTags(subnetID, tags)
		if err != nil {
			return errors.NewAPIError(T("Failed to set tags: {{.tags}}.\n", map[string]interface{}{"tags": tags}), err.Error(), 2)
		}
		if response {
			cmd.UI.Ok()
			cmd.UI.Print(T("Set tags successfully"))
		}
	}

	if cmd.Note != "" {
		note := cmd.Note
		response, err := cmd.NetworkManager.SetSubnetNote(subnetID, note)
		if err != nil {
			return errors.NewAPIError(T("Failed to set note: {{.note}}.\n", map[string]interface{}{"note": note}), err.Error(), 2)
		}
		if response {
			cmd.UI.Ok()
			cmd.UI.Print(T("Set note successfully"))
		}
	}

	return nil
}
