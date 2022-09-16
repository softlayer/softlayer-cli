package globalip

import (
	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type AssignCommand struct {
	*metadata.SoftlayerCommand
	NetworkManager managers.NetworkManager
	Command        *cobra.Command
}

func NewAssignCommand(sl *metadata.SoftlayerCommand) *AssignCommand {
	thisCmd := &AssignCommand{
		SoftlayerCommand: sl,
		NetworkManager:   managers.NewNetworkManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "assign " + T("IDENTIFIER TARGET"),
		Short: T("Assign a global IP to a target router or device"),
		Long: T(`${COMMAND_NAME} sl globalip assign IDENTIFIER TARGET [OPTIONS]

EXAMPLE:
	${COMMAND_NAME} sl globalip assign 12345678 9.111.123.456
	This command assigns IP address with ID 12345678 to a target device whose IP address is 9.111.123.456.`),
		Args: metadata.TwoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *AssignCommand) Run(args []string) error {

	globalIPID, err := utils.ResolveGloablIPId(args[0])
	if err != nil {
		return errors.NewInvalidSoftlayerIdInputError("Globalip ID")
	}
	outputFormat := cmd.GetOutputFlag()

	targetIPAddress := args[1]
	resp, err := cmd.NetworkManager.AssignGlobalIP(globalIPID, targetIPAddress)
	if err != nil {
		return errors.NewAPIError(T("Failed to assign global IP {{.IpID}} to target {{.Target}}.\n",
			map[string]interface{}{"IpID": globalIPID, "Target": targetIPAddress}), err.Error(), 2)

	}
	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, resp)
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("The transaction to modify a global IP route is created, routes will be updated in one or two minutes."))
	return nil
}
