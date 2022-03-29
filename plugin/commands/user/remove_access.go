package user

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type RemoveAccessCommand struct {
	UI          terminal.UI
	UserManager managers.UserManager
}

func NewRemoveAccessCommand(ui terminal.UI, userManager managers.UserManager) (cmd *RemoveAccessCommand) {
	return &RemoveAccessCommand{
		UI:          ui,
		UserManager: userManager,
	}
}

func (cmd *RemoveAccessCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one identifier."))
	}

	if !c.IsSet("hardware") && !c.IsSet("virtual") && !c.IsSet("dedicated") {
		return errors.NewInvalidUsageError(T("This command requires one option."))
	}

	userId, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return errors.NewInvalidUsageError(T("User ID should be a number."))
	}

	if c.IsSet("hardware") {
		hardwareId, err := strconv.Atoi(c.String("hardware"))
		if err != nil {
			return errors.NewInvalidUsageError(T("Hardware ID should be a number."))
		} else {
			i18nsubs := map[string]interface{}{"userId": userId, "hardwareId": hardwareId}
			response, err := cmd.UserManager.RemoveHardwareAccess(userId, hardwareId)
			if err != nil {
				return err
			}
			if response {
				cmd.UI.Ok()
				cmd.UI.Print(T("Access was removed from user {{.userId}} to hardware {{.hardwareId}}", i18nsubs))
			} else {
				cli.NewExitError(T("Failed to remove access user{{.userId}} to hardware {{.hardwareId}}", i18nsubs), 2)
			}
		}
	}

	if c.IsSet("dedicated") {
		dedicatedHostId, err := strconv.Atoi(c.String("dedicated"))
		if err != nil {
			return errors.NewInvalidUsageError(T("Dedicated host ID should be a number."))
		} else {
			i18nsubs := map[string]interface{}{"userId": userId, "dedicatedHostId": dedicatedHostId}
			response, err := cmd.UserManager.RemoveDedicatedHostAccess(userId, dedicatedHostId)
			if err != nil {
				return err
			}
			if response {
				cmd.UI.Ok()
				cmd.UI.Print(T("Access was removed from user {{.userId}} to dedicated host {{.dedicatedHostId}}", i18nsubs))
			} else {
				cli.NewExitError(T("Failed to remove access user{{.userId}} to dedicated host {{.dedicatedHostId}}", i18nsubs), 2)
			}
		}
	}

	if c.IsSet("virtual") {
		virtualId, err := strconv.Atoi(c.String("virtual"))
		if err != nil {
			return errors.NewInvalidUsageError(T("Virtual server ID should be a number."))
		} else {
			i18nsubs := map[string]interface{}{"userId": userId, "virtualId": virtualId}
			response, err := cmd.UserManager.RemoveVirtualGuestAccess(userId, virtualId)
			if err != nil {
				return err
			}
			if response {
				cmd.UI.Ok()
				cmd.UI.Print(T("Access was removed from user {{.userId}} to virtual server {{.virtualId}}", i18nsubs))
			} else {
				cli.NewExitError(T("Failed to remove access user{{.userId}} to virtual server {{.virtualId}}", i18nsubs), 2)
			}
		}
	}

	return nil

}

func UserRemoveAccessMataData() cli.Command {
	return cli.Command{
		Category:    "user",
		Name:        "remove-access",
		Description: T("Remove access from a user to an specific device"),
		Usage: T(`${COMMAND_NAME} sl user remove-access IDENTIFIER [OPTION]
	
EXAMPLE: 
   ${COMMAND_NAME} sl user remove-access userId --hardware hardwareId
   ${COMMAND_NAME} sl user remove-access 123456 --hardware 987654`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "hardware",
				Usage: T("ID of hardware to remove access"),
			},
			cli.StringFlag{
				Name:  "virtual",
				Usage: T("ID of virtual server to remove access"),
			},
			cli.StringFlag{
				Name:  "dedicated",
				Usage: T("ID of dedicated host to remove access"),
			},
		},
	}
}
