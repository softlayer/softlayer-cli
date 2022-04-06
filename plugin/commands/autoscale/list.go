package autoscale

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/urfave/cli"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type ListCommand struct {
	UI               terminal.UI
	AutoScaleManager managers.AutoScaleManager
}

func NewListCommand(ui terminal.UI, autoScaleManager managers.AutoScaleManager) (cmd *ListCommand) {
	return &ListCommand{
		UI:               ui,
		AutoScaleManager: autoScaleManager,
	}
}

type Autoscale struct {
	Visibility string
	datatypes.Scale_Group
}

func (cmd *ListCommand) Run(c *cli.Context) error {
	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	mask := "mask[id,name,status,minimumMemberCount,maximumMemberCount,virtualGuestMemberCount]"
	scaleGroups, err := cmd.AutoScaleManager.ListScaleGroups(mask)
	if err != nil {
		return cli.NewExitError(T("Failed to get scale groups.")+err.Error(), 2)
	}

	table := cmd.UI.Table([]string{T("Id"), T("Name"), T("Status"), T("Min/Max"), T("Running")})
	for _, scale := range scaleGroups {
		membercount := strconv.Itoa(*scale.MinimumMemberCount) + "/" + strconv.Itoa(*scale.MaximumMemberCount)
		table.Add(
			utils.FormatIntPointer(scale.Id),
			utils.FormatStringPointer(scale.Name),
			utils.FormatStringPointer(scale.Status.Name),
			membercount,
			utils.FormatUIntPointer(scale.VirtualGuestMemberCount),
		)
	}

	if outputFormat == "JSON" {
		table.PrintJson()
	} else {
		table.Print()
	}
	return nil
}

func AutoScaleListMetaData() cli.Command {
	return cli.Command{
		Category:    "autoscale",
		Name:        "list",
		Description: T("List all Autoscale Groups on your account"),
		Usage: T(`${COMMAND_NAME} sl autoscale list

EXAMPLE: 
   ${COMMAND_NAME} sl autoscale list
   This command list all Autoscale Groups on current account.`),
		Flags: []cli.Flag{
			metadata.OutputFlag(),
		},
	}
}
