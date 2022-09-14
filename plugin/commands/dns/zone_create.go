package dns

import (
	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type ZoneCreateCommand struct {
	*metadata.SoftlayerCommand
	DNSManager managers.DNSManager
	Command    *cobra.Command
}

func NewZoneCreateCommand(sl *metadata.SoftlayerCommand) *ZoneCreateCommand {
	thisCmd := &ZoneCreateCommand{
		SoftlayerCommand: sl,
		DNSManager:       managers.NewDNSManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "zone-create " + T("ZONE"),
		Short: T("Create a zone"),
		Long: T(`${COMMAND_NAME} sl dns zone-create ZONE [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl dns zone-create ibm.com 
   This command creates a zone that is named ibm.com.`),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *ZoneCreateCommand) Run(args []string) error {
	zoneName := args[0]

	outputFormat := cmd.GetOutputFlag()

	resp, err := cmd.DNSManager.CreateZone(zoneName)
	if err != nil {
		return errors.NewAPIError(T("Failed to create zone: {{.Zone}}.\n", map[string]interface{}{"Zone": zoneName}), err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, resp)
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("Zone {{.Zone}} was created.", map[string]interface{}{"Zone": zoneName}))
	return nil
}
