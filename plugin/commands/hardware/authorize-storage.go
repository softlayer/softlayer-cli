package hardware

import (
	"strconv"

	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type AuthorizeStorageCommand struct {
	*metadata.SoftlayerCommand
	HardwareManager managers.HardwareServerManager
	Command         *cobra.Command
	UsernameStorage string
}

func NewAuthorizeStorageCommand(sl *metadata.SoftlayerCommand) (cmd *AuthorizeStorageCommand) {
	thisCmd := &AuthorizeStorageCommand{
		SoftlayerCommand: sl,
		HardwareManager:  managers.NewHardwareServerManager(sl.Session),
	}


	cobraCmd := &cobra.Command{
		Use:   "authorize-storage " + T("IDENTIFIER"),
		Short: T("Authorize File and Block Storage to a Hardware Server"),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().StringVarP(&thisCmd.UsernameStorage, "username-storage", "u", "", T("The storage username to be added to the hardware server."))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *AuthorizeStorageCommand) Run(args []string) error {
	hardwareId, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Hardware server ID")
	}

	outputFormat := cmd.GetOutputFlag()

	storageResult, err := cmd.HardwareManager.AuthorizeStorage(hardwareId, cmd.UsernameStorage)
	if err != nil {
		return errors.NewInvalidUsageError(T("Failed to authorize storage to the hardware server instance: {{.Storage}}.\n{{.Error}}",
			map[string]interface{}{"Storage": cmd.UsernameStorage, "Error": err.Error()}))
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, storageResult)
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("Successfully authorized storage: {{.Storage}} was added.",
		map[string]interface{}{"Storage": cmd.UsernameStorage}))

	return nil
}
