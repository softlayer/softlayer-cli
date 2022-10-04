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

type CredentialsCommand struct {
	*metadata.SoftlayerCommand
	HardwareManager managers.HardwareServerManager
	Command         *cobra.Command
}

func NewCredentialsCommand(sl *metadata.SoftlayerCommand) (cmd *CredentialsCommand) {
	thisCmd := &CredentialsCommand{
		SoftlayerCommand: sl,
		HardwareManager:  managers.NewHardwareServerManager(sl.Session),
	}


	cobraCmd := &cobra.Command{
		Use:   "credentials " + T("IDENTIFIER"),
		Short: T("List hardware server credentials"),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *CredentialsCommand) Run(args []string) error {
	hardwareID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Hardware server ID")
	}

	outputFormat := cmd.GetOutputFlag()

	hardware, err := cmd.HardwareManager.GetHardware(hardwareID, "")
	if err != nil {
		return errors.NewAPIError(T("Failed to get hardware server {{.ID}}.\n", map[string]interface{}{"ID": hardwareID}), err.Error(), 2)
	}

	if outputFormat == "JSON" {
		if hardware.OperatingSystem != nil {
			return utils.PrintPrettyJSON(cmd.UI, hardware.OperatingSystem.Passwords)
		}
		return utils.PrintPrettyJSON(cmd.UI, hardware.OperatingSystem)
	}

	if hardware.OperatingSystem != nil && hardware.OperatingSystem.Passwords != nil {
		table := cmd.UI.Table([]string{T("Username"), T("Password")})
		for _, item := range hardware.OperatingSystem.Passwords {
			if item.Username != nil && item.Password != nil {
				table.Add(*item.Username, *item.Password)
			}
		}
		table.Print()
		return nil
	}
	return errors.NewInvalidUsageError(T("Failed to find credentials of hardware server {{.ID}}.", map[string]interface{}{"ID": hardwareID}))
}
