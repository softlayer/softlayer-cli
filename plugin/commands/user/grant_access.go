package user

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type GrantAccessCommand struct {
	UI          terminal.UI
	UserManager managers.UserManager
}

func NewGrantAccessCommand(ui terminal.UI, userManager managers.UserManager) (cmd *GrantAccessCommand) {
	return &GrantAccessCommand{
		UI:          ui,
		UserManager: userManager,
	}
}

func (cmd *GrantAccessCommand) Run(c *cli.Context) error {
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
			response, err := cmd.UserManager.AddHardwareAccess(userId, hardwareId)
			if err != nil {
				return err
			}
			if response {
				cmd.UI.Ok()
				cmd.UI.Print(T("Access was granted from user {{.userId}} to hardware {{.hardwareId}}", i18nsubs))
			} else {
				cmd.UI.Print(T("Failed to grant access user{{.userId}} to hardware {{.hardwareId}} device", i18nsubs))
			}
		}
	}

	if c.IsSet("dedicated") {
		dedicatedHostId, err := strconv.Atoi(c.String("dedicated"))
		if err != nil {
			return errors.NewInvalidUsageError(T("Dedicated host ID should be a number."))
		} else {
			i18nsubs := map[string]interface{}{"userId": userId, "dedicatedHostId": dedicatedHostId}
			response, err := cmd.UserManager.AddDedicatedHostAccess(userId, dedicatedHostId)
			if err != nil {
				return err
			}
			if response {
				cmd.UI.Ok()
				cmd.UI.Print(T("Access was granted from user {{.userId}} to dedicated host {{.dedicatedHostId}}", i18nsubs))
			} else {
				cmd.UI.Print(T("Failed to grant access user{{.userId}} to dedicated host {{.dedicatedHostId}} device", i18nsubs))
			}
		}
	}

	if c.IsSet("virtual") {
		virtualId, err := strconv.Atoi(c.String("virtual"))
		if err != nil {
			return errors.NewInvalidUsageError(T("Virtual server ID should be a number."))
		} else {
			i18nsubs := map[string]interface{}{"userId": userId, "virtualId": virtualId}
			response, err := cmd.UserManager.AddVirtualGuestAccess(userId, virtualId)
			if err != nil {
				return err
			}
			if response {
				cmd.UI.Ok()
				cmd.UI.Print(T("Access was granted from user {{.userId}} to virtual server {{.virtualId}}", i18nsubs))
			} else {
				cmd.UI.Print(T("Failed to grant access user{{.userId}} to virtual server {{.virtualId}} device", i18nsubs))
			}
		}
	}

	return nil

}

func UserGrantAccessMataData() cli.Command {
	return cli.Command{
		Category:    "user",
		Name:        "grant-access",
		Description: T("Grant access from a user to an specific device"),
		Usage: T(`${COMMAND_NAME} sl user grant-access IDENTIFIER [OPTION]
	
EXAMPLE: 
   ${COMMAND_NAME} sl user grant-access userId --hardware hardwareId
   ${COMMAND_NAME} sl user grant-access 123456 --hardware 987654`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "hardware",
				Usage: T("ID of hardware to grant access"),
			},
			cli.StringFlag{
				Name:  "virtual",
				Usage: T("ID of virtual server to grant access"),
			},
			cli.StringFlag{
				Name:  "dedicated",
				Usage: T("ID of dedicated host to grant access"),
			},
		},
	}
}
