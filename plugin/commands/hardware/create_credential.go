package hardware

import (
	"strconv"
	"strings"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/spf13/cobra"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type CreateCredentialCommand struct {
	*metadata.SoftlayerCommand
	HardwareManager managers.HardwareServerManager
	Command         *cobra.Command
	Username        string
	Password        string
	Notes           string
	Software        string
}

func NewCreateCredentialCommand(sl *metadata.SoftlayerCommand) (cmd *CreateCredentialCommand) {
	thisCmd := &CreateCredentialCommand{
		SoftlayerCommand: sl,
		HardwareManager:  managers.NewHardwareServerManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "create-credential " + T("IDENTIFIER"),
		Short: T("Create a password for a software component."),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().StringVarP(&thisCmd.Username, "username", "U", "", T("The username part of the username/password pair."))
	cobraCmd.Flags().StringVarP(&thisCmd.Password, "password", "P", "", T("The password part of the username/password pair."))
	cobraCmd.Flags().StringVarP(&thisCmd.Notes, "notes", "n", "", T("A note string stored for this username/password pair."))
	cobraCmd.Flags().StringVar(&thisCmd.Software, "software", "", T("The name of this specific piece of software."))

	//#nosec G104 -- This is a false positive
	cobraCmd.MarkFlagRequired("username")
	//#nosec G104 -- This is a false positive
	cobraCmd.MarkFlagRequired("password")
	//#nosec G104 -- This is a false positive
	cobraCmd.MarkFlagRequired("software")

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *CreateCredentialCommand) Run(args []string) error {
	outputFormat := cmd.GetOutputFlag()

	hardwareId, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.NewInvalidSoftlayerIdInputError("Hardware ID")
	}

	hardware, err := cmd.HardwareManager.GetHardware(hardwareId, "mask[softwareComponents[softwareLicense[softwareDescription]]]")
	if err != nil {
		return errors.NewAPIError(T("Failed to get hardware server: {{.ID}}.", map[string]interface{}{"ID": hardwareId}), err.Error(), 2)
	}
	softwareComponents := hardware.SoftwareComponents
	softwareId := 0

	for _, softwareComponent := range softwareComponents {
		if strings.ToLower(*softwareComponent.SoftwareLicense.SoftwareDescription.Name) == strings.ToLower(strings.Trim(cmd.Software, " ")) {
			softwareId = *softwareComponent.Id
		}
	}
	if softwareId == 0 {
		return errors.NewInvalidUsageError(T("Software not found"))
	}

	softwareComponentPasswordTemplate := datatypes.Software_Component_Password{
		Notes:      sl.String(cmd.Notes),
		Password:   sl.String(cmd.Password),
		SoftwareId: sl.Int(softwareId),
		Username:   sl.String(cmd.Username),
	}

	softwareCredential, err := cmd.HardwareManager.CreateSoftwareCredential(softwareComponentPasswordTemplate)
	if err != nil {
		return errors.NewAPIError(T("Failed to create Software Credential."), err.Error(), 2)
	}

	table := cmd.UI.Table([]string{T("Name"), T("Value")})
	table.Add("Software Credential Id", utils.FormatIntPointer(softwareCredential.Id))
	table.Add("Created", utils.FormatSLTimePointer(softwareCredential.CreateDate))
	table.Add("Username", utils.FormatStringPointer(softwareCredential.Username))
	table.Add("Password", utils.FormatStringPointer(softwareCredential.Password))
	notes := "-"
	if softwareCredential.Notes != nil {
		notes = *softwareCredential.Notes
	}
	table.Add("Notes", notes)

	utils.PrintTable(cmd.UI, table, outputFormat)
	return nil
}
