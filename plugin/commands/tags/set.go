package tags

import (
	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type SetCommand struct {
	*metadata.SoftlayerCommand
	TagsManager managers.TagsManager
	Command     *cobra.Command
	Tags        string
	KeyName     string
	ResourceId  int
}

func NewSetCommand(sl *metadata.SoftlayerCommand) (cmd *SetCommand) {
	thisCmd := &SetCommand{
		SoftlayerCommand: sl,
		TagsManager:      managers.NewTagsManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "set",
		Short: T("Set Tags."),
		Long: T(`${COMMAND_NAME} sl tags set [OPTIONS]

EXAMPLE:
	${COMMAND_NAME} sl tags set --tags 'tag1,tag2' --key-name HARDWARE --resource-id 123456
`),
		Args: metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().StringVar(&thisCmd.Tags, "tags", "", T("Comma seperated list of tags, enclosed in quotes. 'tag1,tag2'  [required]"))
	cobraCmd.Flags().StringVar(&thisCmd.KeyName, "key-name", "", T("Key name of a tag type e.g. GUEST, HARDWARE. See slcli tags taggable output.  [required]"))
	cobraCmd.Flags().IntVar(&thisCmd.ResourceId, "resource-id", 0, T("ID of the object being tagged  [required]"))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *SetCommand) Run(args []string) error {
	tags := cmd.Tags
	keyName := cmd.KeyName
	resourceId := cmd.ResourceId

	if tags == "" {
		return errors.NewMissingInputError("--tags")
	}
	if keyName == "" {
		return errors.NewMissingInputError("--key-name")
	}
	if resourceId == 0 {
		return errors.NewMissingInputError("--resource-id")
	}

	response, err := cmd.TagsManager.SetTags(tags, keyName, resourceId)
	if err != nil {
		return errors.NewAPIError(T("Failed to set tags."), err.Error(), 2)
	}
	if response {
		cmd.UI.Ok()
		cmd.UI.Print("Set tags successfully")
	}

	return nil
}
