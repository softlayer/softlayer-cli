package autoscale

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type TagCommand struct {
	UI                   terminal.UI
	AutoScaleManager     managers.AutoScaleManager
	VirtualServerManager managers.VirtualServerManager
}

func NewTagCommand(ui terminal.UI, autoScaleManager managers.AutoScaleManager, virtualServerManager managers.VirtualServerManager) (cmd *TagCommand) {
	return &TagCommand{
		UI:                   ui,
		AutoScaleManager:     autoScaleManager,
		VirtualServerManager: virtualServerManager,
	}
}

func (cmd *TagCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}

	autoScaleGroupId, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Autoscale Group ID")
	}

	mask := "mask[id, virtualGuestId, virtualGuest[tagReferences, id, hostname]]"
	autoScaleGroupMembers, err := cmd.AutoScaleManager.GetVirtualGuestMembers(autoScaleGroupId, mask)
	if err != nil {
		return cli.NewExitError(T("Failed to get virtual guest members.\n")+err.Error(), 2)
	}

	tags := c.String("tags")
	if tags != "" {
		cmd.UI.Print(T("New Tags: {{.Tags}}.", map[string]interface{}{"Tags": tags}))
	} else {
		cmd.UI.Print(T("All tags will be removed"))
	}

	for _, virtualGuest := range autoScaleGroupMembers {
		err := cmd.VirtualServerManager.SetTags(*virtualGuest.VirtualGuest.Id, tags)
		if err != nil {
			cmd.UI.Failed(T("Failed set tags for virtual guest {{.VirtualGuestHostname}}.\n", map[string]interface{}{"VirtualGuestHostname": *virtualGuest.VirtualGuest.Hostname})+err.Error(), 2)
		}
		cmd.UI.Print(T("Setting tags for {{.VirtualGuestHostname}}.", map[string]interface{}{"VirtualGuestHostname": *virtualGuest.VirtualGuest.Hostname}))
	}

	cmd.UI.Print(T("Done"))
	return nil
}

func AutoScaleTagMetaData() cli.Command {
	return cli.Command{
		Category:    "autoscale",
		Name:        "tag",
		Description: T("Tags all guests in an autoscale group."),
		Usage: T(`${COMMAND_NAME} sl autoscale tag IDENTIFIER [OPTIONS]

EXAMPLE: 
   ${COMMAND_NAME} sl autoscale tag 123456 --tags 'Use, single or double quotes, if you, want whitespace'
   ${COMMAND_NAME} sl autoscale tag 123456 --tags Otherwise,Just,commas`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "g,tags",
				Usage: T("Tags to set for each guest in this group. Existing tags are overwritten. An empty string will remove all tags."),
			},
		},
	}
}
