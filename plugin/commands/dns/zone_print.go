package dns

import (
	"strings"

	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type ZonePrintCommand struct {
	*metadata.SoftlayerCommand
	DNSManager managers.DNSManager
	Command    *cobra.Command
}

func NewZonePrintCommand(sl *metadata.SoftlayerCommand) *ZonePrintCommand {
	thisCmd := &ZonePrintCommand{
		SoftlayerCommand: sl,
		DNSManager:       managers.NewDNSManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "zone-print " + T("ZONE"),
		Short: T("zone-print."),
		Long: T(`${COMMAND_NAME} sl dns zone-print ZONE

EXAMPLE:
	${COMMAND_NAME} sl dns zone-print ibm.com
	This command prints zone that is named ibm.com, and in BIND format.`),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *ZonePrintCommand) Run(args []string) error {
	zoneName := args[0]
	zoneID, err := cmd.DNSManager.GetZoneIdFromName(zoneName)
	if err != nil {
		return errors.NewAPIError(T("Failed to get zone ID from zone name: {{.Zone}}.\n", map[string]interface{}{"Zone": zoneName}), err.Error(), 2)
	}
	fileContent, err := cmd.DNSManager.DumpZone(zoneID)
	if err != nil {
		return errors.NewAPIError(T("Failed to dump content for zone: {{.Zone}}.\n", map[string]interface{}{"Zone": zoneName}), err.Error(), 2)
	}
	//TODO need to test on other platforms about line break
	lines := strings.Split(fileContent, "\\n")
	for _, line := range lines {
		cmd.UI.Print(strings.Trim(line, "\""))
	}
	return nil
}
