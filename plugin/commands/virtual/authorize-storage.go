package virtual

import (
	"github.com/spf13/cobra"

	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type AuthorizeStorageCommand struct {
	*metadata.SoftlayerCommand
	VirtualServerManager managers.VirtualServerManager
	Command              *cobra.Command
	Username             string
	PortableId           int
}

func NewAuthorizeStorageCommand(sl *metadata.SoftlayerCommand) (cmd *AuthorizeStorageCommand) {
	thisCmd := &AuthorizeStorageCommand{
		SoftlayerCommand:     sl,
		VirtualServerManager: managers.NewVirtualServerManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "authorize-storage " + T("IDENTIFIER"),
		Short: T("Authorize File, Block and Portable Storage to a Virtual Server"),
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
	cobraCmd.Flags().StringVarP(&thisCmd.Username, "username-storage", "u", "", T("The storage username to be added to the virtual server."))
	cobraCmd.Flags().IntVarP(&thisCmd.PortableId, "portable-id", "p", 0, T("The portable storage id to be added to the virtual server"))
	return thisCmd
}

func (cmd *AuthorizeStorageCommand) Run(args []string) error {

	vsID, err := utils.ResolveVirtualGuestId(args[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("Virtual server ID")
	}
	outputFormat := cmd.GetOutputFlag()
	subs := map[string]interface{}{
		"Storage":    cmd.Username,
		"Error":      "",
		"PortableID": cmd.PortableId,
	}

	if cmd.Username != "" {
		storageResult, err := cmd.VirtualServerManager.AuthorizeStorage(vsID, cmd.Username)
		if err != nil {
			subs["Error"] = err.Error()
			// TODO: Remove this error bit from the T() string
			return slErrors.NewAPIError(T("Failed to authorize storage to the virtual server instance: {{.Storage}}.\n{{.Error}}", subs), err.Error(), 2)
		}

		if outputFormat == "JSON" {
			return utils.PrintPrettyJSON(cmd.UI, storageResult)
		}

		cmd.UI.Ok()
		cmd.UI.Print(T("Successfully authorized storage: {{.Storage}} was added.", subs))
	}

	if cmd.PortableId != 0 {
		portableResult, err := cmd.VirtualServerManager.AttachPortableStorage(vsID, cmd.PortableId)
		if err != nil {
			subs["Error"] = err.Error()
			// TODO: Remove this error bit from the T() string
			return slErrors.NewAPIError(T("Failed to authorize portable storage to the virtual server instance: {{.PortableID}}.\n{{.Error}}", subs), err.Error(), 2)
		}

		if outputFormat == "JSON" {
			return utils.PrintPrettyJSON(cmd.UI, portableResult)
		}

		cmd.UI.Ok()
		cmd.UI.Print(T("Successfully authorized storage: {{.PortableID}} was added.", subs))
		table := cmd.UI.Table([]string{T("id"), T("CreateDate")})
		table.Add(utils.FormatIntPointer(portableResult.Id), utils.FormatSLTimePointer(portableResult.CreateDate))
		table.Print()
	}

	return nil
}
