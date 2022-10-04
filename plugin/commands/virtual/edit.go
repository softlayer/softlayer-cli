package virtual

import (
	"errors"
	"io/ioutil"

	"github.com/spf13/cobra"

	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type EditCommand struct {
	*metadata.SoftlayerCommand
	VirtualServerManager managers.VirtualServerManager
	Command              *cobra.Command
	Domain               string
	Hostname             string
	Userdata             string
	Userfile             string
	Tag                  []string
	PublicSpeed          int
	PrivateSpeed         int
}

func NewEditCommand(sl *metadata.SoftlayerCommand) (cmd *EditCommand) {
	thisCmd := &EditCommand{
		SoftlayerCommand:     sl,
		VirtualServerManager: managers.NewVirtualServerManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "edit " + T("IDENTIFIER"),
		Short: T("Edit a virtual server instance's details"),
		Long: T(`${COMMAND_NAME} sl vs edit IDENTIFIER [OPTIONS]
	
EXAMPLE:
   ${COMMAND_NAME} sl vs edit 12345678 -D ibm.com -H myapp --tag testcli --public-speed 1000
   This command updates virtual server instance with ID 12345678 and set its domain to be "ibm.com", hostname to "myapp", tag to "testcli", 
   and public network port speed to 1000 Mbps.`),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	cobraCmd.Flags().StringVarP(&thisCmd.Domain, "domain", "D", "", T("Domain portion of the FQDN"))
	cobraCmd.Flags().StringVarP(&thisCmd.Hostname, "hostname", "H", "", T("Host portion of the FQDN. example: server"))
	cobraCmd.Flags().StringVarP(&thisCmd.Userdata, "userdata", "u", "", T("User defined metadata string"))
	cobraCmd.Flags().StringVarP(&thisCmd.Userfile, "userfile", "F", "", T("Read userdata from file"))
	cobraCmd.Flags().StringSliceVarP(&thisCmd.Tag, "tag", "g", []string{}, T("Tags to set or empty string to remove all"))
	cobraCmd.Flags().IntVar(&thisCmd.PublicSpeed, "public-speed", -1, T("Public port speed, options are: 0,10,100,1000,10000"))
	cobraCmd.Flags().IntVar(&thisCmd.PrivateSpeed, "private-speed", -1, T("Private port speed, options are: 0,10,100,1000,10000"))
	return thisCmd
}

func (cmd *EditCommand) Run(args []string) error {

	id, err := utils.ResolveVirtualGuestId(args[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("Virtual server ID")
	}
	var userData, tagString string
	var publicSpeed, privateSpeed *int

	if cmd.Userdata != "" && cmd.Userfile != "" {
		return slErrors.NewExclusiveFlagsError("[-u|--userdata]", "[-F|--userfile]")
	}

	if cmd.PublicSpeed >= 0 {
		publicSpeedInt := cmd.PublicSpeed
		publicSpeed = &publicSpeedInt
		if publicSpeedInt != 0 && publicSpeedInt != 10 && publicSpeedInt != 100 && publicSpeedInt != 1000 && publicSpeedInt != 10000 {
			return slErrors.NewInvalidUsageError("Public network interface speed must be in: 0, 10, 100, 1000, 10000 (Mbps).")
		}
	}
	if cmd.PrivateSpeed >= 0 {
		privateSpeedInt := cmd.PrivateSpeed
		privateSpeed = &privateSpeedInt
		if privateSpeedInt != 0 && privateSpeedInt != 10 && privateSpeedInt != 100 && privateSpeedInt != 1000 && privateSpeedInt != 10000 {
			return slErrors.NewInvalidUsageError("Private network interface speed must be in: 0, 10, 100, 1000, 10000 (Mbps).")
		}
	}

	if len(cmd.Tag) > 0 {
		tagString = utils.StringSliceToString(cmd.Tag)
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

	successes, messages := cmd.VirtualServerManager.EditInstance(id, cmd.Hostname, cmd.Domain, userData, tagString, publicSpeed, privateSpeed)
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
		return slErrors.CollapseErrors(multiErrors)
	}
	return nil
}
