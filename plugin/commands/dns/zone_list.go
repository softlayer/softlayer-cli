package dns

import (
	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type ZoneListCommand struct {
	*metadata.SoftlayerCommand
	DNSManager managers.DNSManager
	Command    *cobra.Command
}

func NewZoneListCommand(sl *metadata.SoftlayerCommand) *ZoneListCommand {
	thisCmd := &ZoneListCommand{
		SoftlayerCommand: sl,
		DNSManager:       managers.NewDNSManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "zone-list",
		Short: T("List all zones on your account."),
		Long: T(`${COMMAND_NAME} sl dns zone-list [OPTIONS]

EXAMPLE:
	${COMMAND_NAME} sl dns zone-list
	This command lists all zones under current account.`),
		Args: metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *ZoneListCommand) Run(args []string) error {
	outputFormat := cmd.GetOutputFlag()

	domains, err := cmd.DNSManager.ListZones()
	if err != nil {
		return errors.NewAPIError(T("Failed to list zones on your account.\n"), err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, domains)
	}
	table := cmd.UI.Table([]string{T("ID"), T("name"), T("serial"), T("updated")})
	for _, domain := range domains {
		table.Add(utils.FormatIntPointer(domain.Id),
			utils.FormatStringPointer(domain.Name),
			utils.FormatIntPointer(domain.Serial),
			utils.FormatSLTimePointer(domain.UpdateDate))
	}
	table.Print()
	return nil
}
