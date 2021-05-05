package dns

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type ZoneListCommand struct {
	UI         terminal.UI
	DNSManager managers.DNSManager
}

func NewZoneListCommand(ui terminal.UI, dnsManager managers.DNSManager) (cmd *ZoneListCommand) {
	return &ZoneListCommand{
		UI:         ui,
		DNSManager: dnsManager,
	}
}

func (cmd *ZoneListCommand) Run(c *cli.Context) error {
	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	domains, err := cmd.DNSManager.ListZones()
	if err != nil {
		return cli.NewExitError(T("Failed to list zones on your account.\n")+err.Error(), 2)
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
