package hardware

import (
	"errors"
	"io/ioutil"
	"strconv"

	bmxErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"

	"github.com/spf13/cobra"
	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type EditCommand struct {
	*metadata.SoftlayerCommand
	HardwareManager managers.HardwareServerManager
	Command         *cobra.Command
	Hostname        string
	Domain          string
	Tag             []string
	Userfile        string
	Userdata        string
	PublicSpeed     int
	PrivateSpeed    int
}

func NewEditCommand(sl *metadata.SoftlayerCommand) (cmd *EditCommand) {
	thisCmd := &EditCommand{
		SoftlayerCommand: sl,
		HardwareManager:  managers.NewHardwareServerManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "edit " + T("IDENTIFIER"),
		Short: T("Edit hardware server details"),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().StringVarP(&thisCmd.Hostname, "hostname", "H", "", T("Host portion of the FQDN"))
	cobraCmd.Flags().StringVarP(&thisCmd.Domain, "domain", "D", "", T("Domain portion of the FQDN"))
	cobraCmd.Flags().StringSliceVarP(&thisCmd.Tag, "tag", "g", []string{}, T("Tags to set or empty string to remove all (multiple occurrence permitted)."))
	cobraCmd.Flags().StringVarP(&thisCmd.Userfile, "userfile", "F", "", T("Read userdata from file"))
	cobraCmd.Flags().StringVarP(&thisCmd.Userdata, "userdata", "u", "", T("User defined metadata string"))
	cobraCmd.Flags().IntVarP(&thisCmd.PublicSpeed, "public-speed", "p", 0, T("Public port speed, options are: 0,10,100,1000,10000"))
	cobraCmd.Flags().IntVarP(&thisCmd.PrivateSpeed, "private-speed", "v", 0, T("Private port speed, options are: 0,10,100,1000,10000"))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *EditCommand) Run(args []string) error {
	hardwareId, err := strconv.Atoi(args[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("Hardware server ID")
	}
	var userData, tagString string
	var publicSpeed, privateSpeed int

	if cmd.Userdata != "" && cmd.Userfile != "" {
		return bmxErr.NewInvalidUsageError(T("[-u|--userdata] is not allowed with [-F|--userfile]."))
	}

	if cmd.PublicSpeed != 0 {
		publicSpeed = cmd.PublicSpeed
		if publicSpeed != 0 && publicSpeed != 10 && publicSpeed != 100 && publicSpeed != 1000 && publicSpeed != 10000 {
			return bmxErr.NewInvalidUsageError(T("Public network interface speed must be in: 0, 10, 100, 1000, 10000 (Mbps)."))
		}
	}
	if cmd.PrivateSpeed != 0 {
		privateSpeed = cmd.PrivateSpeed
		if privateSpeed != 0 && privateSpeed != 10 && privateSpeed != 100 && privateSpeed != 1000 && privateSpeed != 10000 {
			return bmxErr.NewInvalidUsageError(T("Private network interface speed must be in: 0, 10, 100, 1000, 10000 (Mbps)."))
		}
	}

	if tags := cmd.Tag; len(tags) != 0 {
		tagString = utils.StringSliceToString(tags)
	}

	if cmd.Userdata != "" {
		userData = cmd.Userdata
	}
	if cmd.Userfile != "" {
		userfile := cmd.Userfile
		content, err := ioutil.ReadFile(userfile) // #nosec
		if err != nil {
			return slErrors.NewAPIError(T("Failed to read user data file: {{.File}}.\n", map[string]interface{}{"File": userfile}), err.Error(), 1)
		}
		userData = string(content)
	}

	successes, messages := cmd.HardwareManager.Edit(hardwareId, userData, cmd.Hostname, cmd.Domain, "", tagString, publicSpeed, privateSpeed)
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
		for _, errorMessage := range multiErrors {
			cmd.UI.Failed(errorMessage.Error())
		}
	}
	return nil
}
