package tags

import (
	"fmt"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)


type DetailCommand struct {
	UI 			terminal.UI
	TagsManager managers.TagsManager
}


func NewDetailCommand(ui terminal.UI, tagsManager managers.TagsManager) (cmd *DetailCommand) {
	return &DetailCommand{
		UI:             ui,
		TagsManager: 	tagsManager,
	}	
}

func (cmd *DetailCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	tagName := c.Args()[0]
	outputFormat, _ := metadata.CheckOutputFormat(c, cmd.UI)
	tagDetails, err := cmd.TagsManager.GetTagByTagName(tagName)

	if err != nil {
		return cli.NewExitError(err.Error(), 2)
	}

	details :=  BuildDetailedTagTable(tagDetails, cmd.TagsManager)
	if outputFormat == "JSON" {

		return utils.PrintPrettyJSON(cmd.UI, details)
	}

	tagTable := cmd.UI.Table([]string{T("Id"), T("Type"), T("Resource")})
	if tagDetails != nil {
		for _,tag := range details {
			cmd.UI.Print(tag.Name)
			
				for _, reference := range tag.Tags{
					tagTable.Add(fmt.Sprintf("%v",reference.Id), reference.TagType, reference.ResourceName)
				}
				
			}
	}
	tagTable.Print()
	return nil

}

func TagsDetailsMetaData() cli.Command {
	return cli.Command{
		Category:    "tags",
		Name:        "detail",
		Description: T("Get information about the resources using the selected tag."),
		Usage: T(`${COMMAND_NAME} sl tags detail [TAG NAME]

EXAMPLE:
	${COMMAND_NAME} sl tags detail tag1
	Shows all items that are tagged with 'tag1'
`),
		Flags: []cli.Flag{
			metadata.OutputFlag(),
		},
	}
}
