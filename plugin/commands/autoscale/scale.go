package autoscale

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type ScaleCommand struct {
	UI               terminal.UI
	AutoScaleManager managers.AutoScaleManager
}

func NewScaleCommand(ui terminal.UI, autoScaleManager managers.AutoScaleManager) (cmd *ScaleCommand) {
	return &ScaleCommand{
		UI:               ui,
		AutoScaleManager: autoScaleManager,
	}
}

func (cmd *ScaleCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}

	autoScaleGroupId, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return errors.NewInvalidSoftlayerIdInputError("Autoscale Group ID")
	}

	if !c.IsSet("amount") {
		return errors.NewMissingInputError("--amount")
	}

	if c.IsSet("up") && c.IsSet("down") {
		return errors.NewExclusiveFlagsError("[--up]", "[--down]")
	}

	if c.IsSet("to") && c.IsSet("by") {
		return errors.NewExclusiveFlagsError("[--to]", "[--by]")
	}

	if !c.IsSet("to") && !c.IsSet("by") {
		return errors.NewInvalidUsageError(T("--to or --by is required"))
	}

	amount := c.Int("amount")

	if c.IsSet("by") {
		delta := amount
		if c.IsSet("down") {
			delta = delta * -1
		}
		scale_Members, err := cmd.AutoScaleManager.Scale(autoScaleGroupId, delta)
		if err != nil {
			return cli.NewExitError(T("Failed to scale Auto Scale Group.")+err.Error(), 2)
		}
		if len(scale_Members) == amount {
			cmd.UI.Ok()
			cmd.UI.Print(T("Auto Scale Group was scaled successfully"))
		}
	}

	if c.IsSet("to") {
		members, err := cmd.AutoScaleManager.GetVirtualGuestMembers(autoScaleGroupId, "")
		if err != nil {
			return cli.NewExitError(T("Failed to get virtual guest members.\n")+err.Error(), 2)
		}
		newMembers, err := cmd.AutoScaleManager.ScaleTo(autoScaleGroupId, amount)
		if err != nil {
			return cli.NewExitError(T("Failed to scale Auto Scale Group.")+err.Error(), 2)
		}
		newMembersNumber := len(newMembers)
		membersNumber := len(members)
		if membersNumber > amount {
			newMembersNumber = newMembersNumber * -1
		}
		if membersNumber+newMembersNumber == amount {
			cmd.UI.Ok()
			cmd.UI.Print(T("Auto Scale Group was scaled successfully"))
		}
	}
	return nil
}

func AutoScaleScaleMetaData() cli.Command {
	return cli.Command{
		Category:    "autoscale",
		Name:        "scale",
		Description: T("Scales an Autoscale group. Bypasses a scale group's cooldown period."),
		Usage: T(`${COMMAND_NAME} sl autoscale scale IDENTIFIER [OPTIONS]

EXAMPLE: 
   ${COMMAND_NAME} sl autoscale scale 123456 --amount 1 --by --down
   ${COMMAND_NAME} sl autoscale scale 123456 --amount 2 --by --up
   ${COMMAND_NAME} sl autoscale scale 123456 --amount 2 --by
   ${COMMAND_NAME} sl autoscale scale 123456 --amount 3  --to`),
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "up",
				Usage: T("Adds guests. (Only set --up or --down) (default)"),
			},
			cli.BoolFlag{
				Name:  "down",
				Usage: T("Removes guests. (Only set --up or --down)"),
			},
			cli.BoolFlag{
				Name:  "by",
				Usage: T("Will add/remove the specified number of guests. (Only set --by or --to)"),
			},
			cli.BoolFlag{
				Name:  "to",
				Usage: T("Will add/remove a number of guests to get the group's guest count to the specified number. (Only set --by or --to)"),
			},
			cli.IntFlag{
				Name:  "amount",
				Usage: T("Number of guests for the scale action."),
			},
		},
	}
}
