package dns

import (
	"strings"

	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type ZoneDeleteCommand struct {
	*metadata.SoftlayerCommand
	DNSManager managers.DNSManager
	Command    *cobra.Command
}

func NewZoneDeleteCommand(sl *metadata.SoftlayerCommand) *ZoneDeleteCommand {
	thisCmd := &ZoneDeleteCommand{
		SoftlayerCommand: sl,
		DNSManager:       managers.NewDNSManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "zone-delete " + T("ZONE"),
		Short: T("Delete a zone"),
		Long: T(`${COMMAND_NAME} sl dns zone-delete ZONE

EXAMPLE:
	${COMMAND_NAME} sl dns zone-delete ibm.com 
	This command deletes a zone that is named ibm.com.`),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *ZoneDeleteCommand) Run(args []string) error {
	zoneName := args[0]
	zoneID, err := cmd.DNSManager.GetZoneIdFromName(zoneName)
	if err != nil {
		return errors.NewAPIError(T("Failed to get zone ID from zone name: {{.Zone}}.\n", map[string]interface{}{"Zone": zoneName}), err.Error(), 2)
	}
	err = cmd.DNSManager.DeleteZone(zoneID)
	if err != nil {
		if strings.Contains(err.Error(), slErrors.SL_EXP_OBJ_NOT_FOUND) {
			return errors.NewAPIError(T("Unable to find zone with ID: {{.ZoneID}}.\n", map[string]interface{}{"ZoneID": zoneID}), err.Error(), 0)
		}
		return errors.NewAPIError(T("Failed to delete zone: {{.Zone}}.\n", map[string]interface{}{"Zone": zoneName}), err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Zone {{.Zone}} was deleted.", map[string]interface{}{"Zone": zoneName}))
	return nil
}
