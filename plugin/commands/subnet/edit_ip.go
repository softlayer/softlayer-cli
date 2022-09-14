package subnet

import (
	"strconv"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/spf13/cobra"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type EditIpCommand struct {
	*metadata.SoftlayerCommand
	NetworkManager managers.NetworkManager
	Command        *cobra.Command
	Note           string
}

func NewEditIpCommand(sl *metadata.SoftlayerCommand) *EditIpCommand {
	thisCmd := &EditIpCommand{
		SoftlayerCommand: sl,
		NetworkManager:   managers.NewNetworkManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "edit-ip " + T("IDENTIFIER"),
		Short: T("Set the note of the ipAddress."),
		Long: T(`${COMMAND_NAME} sl subnet edit-ip IDENTIFIER [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl subnet edit-ip 11.22.33.44 --note myNote
   ${COMMAND_NAME} sl subnet edit-ip 12345678 --note myNote`),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().StringVar(&thisCmd.Note, "note", "", T("The note"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *EditIpCommand) Run(args []string) error {

	if cmd.Note == "" {
		return errors.NewMissingInputError(T("--note"))
	}

	subnetIpAddressID, err := strconv.Atoi(args[0])
	if err != nil {
		ipAddress := args[0]
		subnetIpAddress, err := cmd.NetworkManager.GetIpByAddress(ipAddress)
		if err != nil {
			return errors.NewAPIError(T("Failed to get Subnet IP by address")+"\n", err.Error(), 2)
		}
		if subnetIpAddress.Id == nil {
			address := map[string]interface{}{"address": ipAddress}
			return cli.NewExitError(T("Unable to find object with IP address: {{.address}}", address), 2)
		}
		subnetIpAddressID = *subnetIpAddress.Id
	}

	note := cmd.Note
	subnetIpAddressTemplate := datatypes.Network_Subnet_IpAddress{
		Note: sl.String(note),
	}
	response, err := cmd.NetworkManager.EditSubnetIpAddress(subnetIpAddressID, subnetIpAddressTemplate)
	if err != nil {
		note := map[string]interface{}{"note": note}
		return errors.NewAPIError(T("Failed to set note: {{.note}}.", note)+"\n", err.Error(), 2)
	}
	if response {
		cmd.UI.Ok()
		cmd.UI.Print(T("Set note successfully"))
	}
	return nil
}
