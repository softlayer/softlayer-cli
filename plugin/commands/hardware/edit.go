package hardware

import (
	"errors"
	"io/ioutil"
	"strconv"

	bmxErr "github.ibm.com/cgallo/softlayer-cli/plugin/errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	. "github.ibm.com/cgallo/softlayer-cli/plugin/i18n"
	slErrors "github.ibm.com/cgallo/softlayer-cli/plugin/errors"
	"github.ibm.com/cgallo/softlayer-cli/plugin/managers"
	"github.ibm.com/cgallo/softlayer-cli/plugin/utils"
)

type EditCommand struct {
	UI              terminal.UI
	HardwareManager managers.HardwareServerManager
}

func NewEditCommand(ui terminal.UI, hardwareManager managers.HardwareServerManager) (cmd *EditCommand) {
	return &EditCommand{
		UI:              ui,
		HardwareManager: hardwareManager,
	}
}

func (cmd *EditCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return bmxErr.NewInvalidUsageError(T("This command requires one argument."))
	}
	hardwareId, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("Hardware server ID")
	}
	var userData, tagString string
	var publicSpeed, privateSpeed int

	if c.IsSet("u") && c.IsSet("F") {
		return bmxErr.NewInvalidUsageError(T("[-u|--userdata] is not allowed with [-F|--userfile]."))
	}

	if c.IsSet("public-speed") {
		publicSpeed = c.Int("public-speed")
		if publicSpeed != 0 && publicSpeed != 10 && publicSpeed != 100 && publicSpeed != 1000 && publicSpeed != 10000 {
			return bmxErr.NewInvalidUsageError(T("Public network interface speed must be in: 0, 10, 100, 1000, 10000 (Mbps)."))
		}
	}
	if c.IsSet("private-speed") {
		privateSpeed = c.Int("private-speed")
		if privateSpeed != 0 && privateSpeed != 10 && privateSpeed != 100 && privateSpeed != 1000 && privateSpeed != 10000 {
			return bmxErr.NewInvalidUsageError(T("Private network interface speed must be in: 0, 10, 100, 1000, 10000 (Mbps)."))
		}
	}

	if tags := c.StringSlice("tag"); len(tags) != 0 {
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

	successes, messages := cmd.HardwareManager.Edit(hardwareId, userData, c.String("H"), c.String("D"), "", tagString, publicSpeed, privateSpeed)
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
