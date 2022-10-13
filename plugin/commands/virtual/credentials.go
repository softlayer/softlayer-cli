package virtual

import (
	"github.com/spf13/cobra"

	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type CredentialsCommand struct {
	*metadata.SoftlayerCommand
	VirtualServerManager managers.VirtualServerManager
	Command              *cobra.Command
}

func NewCredentialsCommand(sl *metadata.SoftlayerCommand) (cmd *CredentialsCommand) {
	thisCmd := &CredentialsCommand{
		SoftlayerCommand:     sl,
		VirtualServerManager: managers.NewVirtualServerManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "credentials " + T("IDENTIFIER"),
		Short: T("List virtual server instance credentials"),
		Long: T(`${COMMAND_NAME} sl vs authorize-storage [OPTIONS] IDENTIFIER

EXAMPLE:
   ${COMMAND_NAME} sl vs authorize-storage --username-storage SL01SL30-37 1234567
   Authorize File, Block and Portable Storage to a Virtual Server.`),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *CredentialsCommand) Run(args []string) error {

	vsID, err := utils.ResolveVirtualGuestId(args[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("Virtual server ID")
	}

	outputFormat := cmd.GetOutputFlag()

	virtualGuest, err := cmd.VirtualServerManager.GetInstance(vsID, "operatingSystem[passwords[username,password]]")
	if err != nil {
		return slErrors.NewAPIError(T("Failed to get virtual server instance: {{.VsID}}.\n", map[string]interface{}{"VsID": vsID}), err.Error(), 2)
	}

	if outputFormat == "JSON" {
		if virtualGuest.OperatingSystem != nil {
			return utils.PrintPrettyJSON(cmd.UI, virtualGuest.OperatingSystem.Passwords)
		}
		return utils.PrintPrettyJSON(cmd.UI, virtualGuest.OperatingSystem)
	}

	table := cmd.UI.Table([]string{T("Username"), T("Password"), T("Software"), T("Version")})
	if virtualGuest.OperatingSystem != nil && len(virtualGuest.OperatingSystem.Passwords) > 0 {
		for _, pwd := range virtualGuest.OperatingSystem.Passwords {
			table.Add(
				utils.FormatStringPointer(pwd.Username),
				utils.FormatStringPointer(pwd.Password),
				utils.FormatStringPointer(virtualGuest.OperatingSystem.SoftwareLicense.SoftwareDescription.ReferenceCode),
				utils.FormatStringPointer(virtualGuest.OperatingSystem.SoftwareLicense.SoftwareDescription.Version),
			)
		}
	}
	table.Print()
	return nil
}
