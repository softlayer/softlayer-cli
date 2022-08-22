package tags

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type DetailCommand struct {
	*metadata.SoftlayerCommand
	TagsManager managers.TagsManager
	Command     *cobra.Command
}

func NewDetailCommand(sl *metadata.SoftlayerCommand) (cmd *DetailCommand) {
	thisCmd := &DetailCommand{
		SoftlayerCommand: sl,
		TagsManager:      managers.NewTagsManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "detail " + T("[TAG NAME]"),
		Short: T("Get information about the resources using the selected tag."),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *DetailCommand) Run(args []string) error {
	tagName := args[0]
	outputFormat := cmd.GetOutputFlag()
	tagDetails, err := cmd.TagsManager.GetTagByTagName(tagName)

	if err != nil {
		return errors.NewAPIError(T("Failed to get tag details"), err.Error(), 2)
	}

	details := BuildDetailedTagTable(tagDetails, cmd.TagsManager)
	if outputFormat == "JSON" {

		return utils.PrintPrettyJSON(cmd.UI, details)
	}

	tagTable := cmd.UI.Table([]string{T("Id"), T("Type"), T("Resource")})
	if tagDetails != nil {
		for _, tag := range details {
			cmd.UI.Print(tag.Name)

			for _, reference := range tag.Tags {
				tagTable.Add(fmt.Sprintf("%v", reference.Id), reference.TagType, reference.ResourceName)
			}

		}
	}
	tagTable.Print()
	return nil

}
