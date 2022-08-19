package tags

import (
	"fmt"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type ListCommand struct {
	*metadata.SoftlayerCommand
	TagsManager managers.TagsManager
	Command     *cobra.Command
	Detail      bool
}

// Structures for easily organizing tag output for either table or JSON
type DetailedTag struct {
	Name string
	Tags []TagInformation
}

type TagInformation struct {
	Id           int
	TagType      string
	ResourceName string
}

func NewListCommand(sl *metadata.SoftlayerCommand) (cmd *ListCommand) {
	thisCmd := &ListCommand{
		SoftlayerCommand: sl,
		TagsManager:      managers.NewTagsManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "list",
		Short: T("List all tags currently on your account"),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().BoolVarP(&thisCmd.Detail, "detail", "d", false, T("List information about devices using the tag."))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *ListCommand) Run(args []string) error {
	outputFormat := cmd.GetOutputFlag()
	// Print Detailed output, only get tags with associated devices
	if cmd.Detail {

		tags, err := cmd.TagsManager.ListTags()
		if err != nil {
			return errors.NewAPIError("Failed to get tags", err.Error(), 2)
		}
		details := BuildDetailedTagTable(tags, cmd.TagsManager)
		if outputFormat == "JSON" {

			return utils.PrintPrettyJSON(cmd.UI, details)
		}

		if tags != nil {
			for _, tag := range details {
				cmd.UI.Print(tag.Name)
				tagTable := cmd.UI.Table([]string{T("Id"), T("Type"), T("Resource")})
				for _, reference := range tag.Tags {
					tagTable.Add(fmt.Sprintf("%v", reference.Id), reference.TagType, reference.ResourceName)
				}
				tagTable.Print()
			}
		}

		// Prints all tags in a simple table
	} else {
		tags, err1 := cmd.TagsManager.ListTags()
		if err1 != nil {
			return errors.NewAPIError("Failed to get tags", err1.Error(), 2)
		}

		emptyTags, err2 := cmd.TagsManager.ListEmptyTags()

		if err2 != nil {
			return errors.NewAPIError("Failed to get empty tags", err2.Error(), 2)
		}

		tags = append(tags, emptyTags...)
		if outputFormat == "JSON" {
			return utils.PrintPrettyJSON(cmd.UI, tags)
		}

		if tags != nil {

			tagTable := cmd.UI.Table([]string{T("Id"), T("Name"), T("Devices")})

			for _, tag := range tags {
				tagTable.Add(
					utils.FormatIntPointer(tag.Id),
					utils.FormatStringPointer(tag.Name),
					utils.FormatUIntPointer(tag.ReferenceCount),
				)
			}

			tagTable.Print()
		}
	}

	return nil
}

func BuildDetailedTagTable(tags []datatypes.Tag, tagsManager managers.TagsManager) []DetailedTag {

	detailedTags := []DetailedTag{}
	for _, tag := range tags {
		thisTag := DetailedTag{Name: utils.FormatStringPointer(tag.Name), Tags: []TagInformation{}}

		references, err := tagsManager.GetTagReferences(*tag.Id)
		if err != nil {
			thisTagDetail := TagInformation{
				Id:           0,
				TagType:      "API ERROR",
				ResourceName: err.Error(),
			}
			thisTag.Tags = append(thisTag.Tags, thisTagDetail)
		} else {
			for _, reference := range references {
				referenceType := "None"
				if reference.TagType != nil {
					referenceType = utils.FormatStringPointer(reference.TagType.KeyName)
				}

				referenceId := *reference.ResourceTableId
				resource := tagsManager.ReferenceLookup(referenceType, referenceId)
				thisTagDetail := TagInformation{
					Id:           referenceId,
					TagType:      referenceType,
					ResourceName: resource,
				}
				thisTag.Tags = append(thisTag.Tags, thisTagDetail)
			}
		}
		detailedTags = append(detailedTags, thisTag)
	}
	return detailedTags
}
