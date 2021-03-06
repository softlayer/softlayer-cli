package virtual

import (
	"errors"
	"io/ioutil"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	bmxErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type EditCommand struct {
	UI                   terminal.UI
	VirtualServerManager managers.VirtualServerManager
}

func NewEditCommand(ui terminal.UI, virtualServerManager managers.VirtualServerManager) (cmd *EditCommand) {
	return &EditCommand{
		UI:                   ui,
		VirtualServerManager: virtualServerManager,
	}
}

func (cmd *EditCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return bmxErr.NewInvalidUsageError("This command requires one argument.")
	}

	id, err := utils.ResolveVirtualGuestId(c.Args()[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("Virtual server ID")
	}
	var userData, tagString string
	var publicSpeed, privateSpeed *int

	if c.IsSet("u") && c.IsSet("F") {
		return bmxErr.NewExclusiveFlagsError("[-u|--userdata]", "[-F|--userfile]")
	}

	if c.IsSet("public-speed") {
		publicSpeedInt := c.Int("public-speed")
		publicSpeed = &publicSpeedInt
		if publicSpeedInt != 0 && publicSpeedInt != 10 && publicSpeedInt != 100 && publicSpeedInt != 1000 && publicSpeedInt != 10000 {
			return bmxErr.NewInvalidUsageError("Public network interface speed must be in: 0, 10, 100, 1000, 10000 (Mbps).")
		}
	}
	if c.IsSet("private-speed") {
		privateSpeedInt := c.Int("private-speed")
		privateSpeed = &privateSpeedInt
		if privateSpeedInt != 0 && privateSpeedInt != 10 && privateSpeedInt != 100 && privateSpeedInt != 1000 && privateSpeedInt != 10000 {
			return bmxErr.NewInvalidUsageError("Private network interface speed must be in: 0, 10, 100, 1000, 10000 (Mbps).")
		}
	}

	if tags := c.StringSlice("tag"); c.IsSet("tag") {
		tagString = utils.StringSliceToString(tags)
	}

	if c.IsSet("u") {
		userData = c.String("u")
	}
	if c.IsSet("F") {
		userfile := c.String("F")
		content, err := ioutil.ReadFile(userfile) // #nosec
		if err != nil {
			return cli.NewExitError(T("Failed to read user data file: {{.File}}.\n", map[string]interface{}{"File": userfile})+err.Error(), 1)
		}
		userData = string(content)
	}

	successes, messages := cmd.VirtualServerManager.EditInstance(id, c.String("H"), c.String("D"), userData, tagString, publicSpeed, privateSpeed)
	var multiErrors []error
	for index, success := range successes {
		if success == true {
			cmd.UI.Ok()
			cmd.UI.Print(messages[index])
		} else {
			multiErrors = append(multiErrors, errors.New(messages[index]))
		}
	}
	if len(multiErrors) > 0 {
		return cli.NewExitError(cli.NewMultiError(multiErrors...).Error(), 2)
	}
	return nil
}

func VSEditMetaData() cli.Command {
	return cli.Command{
		Category:    "vs",
		Name:        "edit",
		Description: T("Edit a virtual server instance's details"),
		Usage: T(`${COMMAND_NAME} sl vs edit IDENTIFIER [OPTIONS]
	
EXAMPLE:
   ${COMMAND_NAME} sl vs edit 12345678 -D ibm.com -H myapp --tag testcli --public-speed 1000
   This command updates virtual server instance with ID 12345678 and set its domain to be "ibm.com", hostname to "myapp", tag to "testcli", 
   and public network port speed to 1000 Mbps.`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "D,domain",
				Usage: T("Domain portion of the FQDN"),
			},
			cli.StringFlag{
				Name:  "H,hostname",
				Usage: T("Host portion of the FQDN. example: server"),
			},
			cli.StringSliceFlag{
				Name:  "g,tag",
				Usage: T("Tags to set or empty string to remove all"),
			},
			cli.StringFlag{
				Name:  "u,userdata",
				Usage: T("User defined metadata string"),
			},
			cli.StringFlag{
				Name:  "F,userfile",
				Usage: T("Read userdata from file"),
			},
			cli.IntFlag{
				Name:  "public-speed",
				Usage: T("Public port speed, options are: 0,10,100,1000,10000"),
			},
			cli.IntFlag{
				Name:  "private-speed",
				Usage: T("Private port speed, options are: 0,10,100,1000,10000"),
			},
		},
	}
}
