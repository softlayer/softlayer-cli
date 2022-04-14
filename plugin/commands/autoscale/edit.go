package autoscale

import (
	"io/ioutil"
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	bmxErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type EditCommand struct {
	UI               terminal.UI
	AutoScaleManager managers.AutoScaleManager
}

func NewEditCommand(ui terminal.UI, autoScaleManager managers.AutoScaleManager) (cmd *EditCommand) {
	return &EditCommand{
		UI:               ui,
		AutoScaleManager: autoScaleManager,
	}
}

func (cmd *EditCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}

	autoScaleGroupId, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Autoscale Group ID")
	}

	if !c.IsSet("name") && !c.IsSet("min") && !c.IsSet("max") && !c.IsSet("userdata") &&
		!c.IsSet("userfile") && !c.IsSet("cpu") && !c.IsSet("memory") {
		return errors.NewInvalidUsageError(T("Please pass at least one of the flags."))
	}

	if c.IsSet("userdata") && c.IsSet("userfile") {
		return bmxErr.NewExclusiveFlagsError("[--userdata]", "[--userfile]")
	}

	autoSacaleGroupTemplate := datatypes.Scale_Group{}

	if c.IsSet("name") {
		autoSacaleGroupTemplate.Name = sl.String(c.String("name"))
	}

	if c.IsSet("min") {
		autoSacaleGroupTemplate.MinimumMemberCount = sl.Int(c.Int("min"))
	}

	if c.IsSet("max") {
		autoSacaleGroupTemplate.MaximumMemberCount = sl.Int(c.Int("max"))
	}

	if c.IsSet("cpu") || c.IsSet("memory") || c.IsSet("userdata") || c.IsSet("userfile") {
		mask := "mask[virtualGuestMemberTemplate]"
		autoScale, err := cmd.AutoScaleManager.GetScaleGroup(autoScaleGroupId, mask)
		if err != nil {
			return cli.NewExitError(T("Failed to get AutoScale group.\n")+err.Error(), 2)
		}
		autoSacaleGroupTemplate.VirtualGuestMemberTemplate = autoScale.VirtualGuestMemberTemplate

		if c.IsSet("cpu") {
			autoSacaleGroupTemplate.VirtualGuestMemberTemplate.StartCpus = sl.Int(c.Int("cpu"))
		}

		if c.IsSet("memory") {
			autoSacaleGroupTemplate.VirtualGuestMemberTemplate.MaxMemory = sl.Int(c.Int("memory"))
		}

		if c.IsSet("userdata") || c.IsSet("userfile") {
			var userData string
			if c.IsSet("userdata") {
				userData = c.String("userdata")
			}
			if c.IsSet("userfile") {
				userfile := c.String("userfile")
				content, err := ioutil.ReadFile(userfile) // #nosec
				if err != nil {
					return cli.NewExitError((T("Failed to read user data from file: {{.File}}.", map[string]interface{}{"File": userfile})), 2)
				}
				userData = string(content)
			}
			autoSacaleGroupTemplate.VirtualGuestMemberTemplate.UserData = []datatypes.Virtual_Guest_Attribute{datatypes.Virtual_Guest_Attribute{Value: &userData}}
		}

	}

	response, err := cmd.AutoScaleManager.EditScaleGroup(autoScaleGroupId, &autoSacaleGroupTemplate)
	if err != nil {
		return cli.NewExitError(T("Failed to update Auto Scale Group.\n")+err.Error(), 2)
	}

	if response {
		cmd.UI.Ok()
	}
	return nil
}

func AutoScaleEditMetaData() cli.Command {
	return cli.Command{
		Category:    "autoscale",
		Name:        "edit",
		Description: T("Edits an Autoscale group."),
		Usage: T(`${COMMAND_NAME} sl autoscale edit IDENTIFIER [OPTIONS]

EXAMPLE: 
   ${COMMAND_NAME} sl autoscale edit 123456 --name newscalegroupname`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "name",
				Usage: T("TEXT Scale group's name."),
			},
			cli.IntFlag{
				Name:  "min",
				Usage: T("INTEGER Set the minimum number of guests"),
			},
			cli.IntFlag{
				Name:  "max",
				Usage: T("INTEGER Set the maximum number of guests"),
			},
			cli.StringFlag{
				Name:  "userdata",
				Usage: T("TEXT User defined metadata string"),
			},
			cli.StringFlag{
				Name:  "userfile",
				Usage: T("PATH Read userdata from a file"),
			},
			cli.IntFlag{
				Name:  "cpu",
				Usage: T("INTEGER Number of CPUs for new guests (existing not effected)"),
			},
			cli.IntFlag{
				Name:  "memory",
				Usage: T("INTEGER RAM in MB or GB for new guests (existing not effected)"),
			},
		},
	}
}
