package tags

import (
	"fmt"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
	"github.com/softlayer/softlayer-go/datatypes"
)


type ListCommand struct {
	UI 			terminal.UI
	TagsManager managers.TagsManager
}

// Structures for easily organizing tag output for either table or JSON
type DetailedTag struct {
	Name 	string
	Tags 	[]TagInformation
}

type TagInformation struct {
	Id 	int
	TagType string
	ResourceName string
}

func NewListCommand(ui terminal.UI, tagsManager managers.TagsManager) (cmd *ListCommand) {
	return &ListCommand{
		UI:             ui,
		TagsManager: 	tagsManager,
	}	
}

func (cmd *ListCommand) Run(c *cli.Context) error {
	outputFormat, _ := metadata.CheckOutputFormat(c, cmd.UI)
	// Print Detailed output, only get tags with associated devices
	if c.IsSet("d") || c.IsSet("details") {

		tags, err := cmd.TagsManager.ListTags()
		if err != nil {
			return cli.NewExitError(err.Error(), 2)
		}
		details := BuildDetailedTagTable(tags, cmd.TagsManager)
		if outputFormat == "JSON" {

			return utils.PrintPrettyJSON(cmd.UI, details)
		}

		if tags != nil {
			for _,tag := range details {
				cmd.UI.Print(tag.Name)
				tagTable := cmd.UI.Table([]string{T("Id"), T("Type"), T("Resource")})
					for _, reference := range tag.Tags{
						tagTable.Add(fmt.Sprintf("%v",reference.Id), reference.TagType, reference.ResourceName)
					}
					tagTable.Print()
				}
		}


	// Prints all tags in a simple table
	} else {
		tags, err1 := cmd.TagsManager.ListTags()
		if err1  != nil {
			return cli.NewExitError(err1.Error(), 2)
		}

		emptyTags, err2 := cmd.TagsManager.ListEmptyTags()
		
		if err2  != nil {
			return cli.NewExitError(err2.Error(), 2)
		}

		tags = append(tags, emptyTags...)
		if outputFormat == "JSON" {
			return utils.PrintPrettyJSON(cmd.UI, tags)
		}

		if tags != nil {

			tagTable := cmd.UI.Table([]string{T("Id"), T("Name"), T("Devices")})

			for _,tag := range tags {
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
	for _,tag := range tags {
		thisTag := DetailedTag{Name: utils.FormatStringPointer(tag.Name), Tags: []TagInformation{}}

		references, err := tagsManager.GetTagReferences(*tag.Id)
		if err != nil {
			thisTagDetail := TagInformation{
				Id: 0,
				TagType: "API ERROR",
				ResourceName: err.Error(),
			}
			thisTag.Tags = append(thisTag.Tags, thisTagDetail)
		} else {
			for _, reference := range references{
				referenceType := "None"
				if reference.TagType != nil {
					referenceType = utils.FormatStringPointer(reference.TagType.KeyName)
				}
				
				referenceId := *reference.ResourceTableId
				resource := tagsManager.ReferenceLookup(referenceType, referenceId)
				thisTagDetail := TagInformation{
					Id: referenceId,
					TagType: referenceType,
					ResourceName: resource,
				}
				thisTag.Tags = append(thisTag.Tags, thisTagDetail)
			}
		}
		detailedTags = append(detailedTags, thisTag)
	}
	return detailedTags
}