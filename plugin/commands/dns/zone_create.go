package dns

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type ZoneCreateCommand struct {
	UI         terminal.UI
	DNSManager managers.DNSManager
}

func NewZoneCreateCommand(ui terminal.UI, dnsManager managers.DNSManager) (cmd *ZoneCreateCommand) {
	return &ZoneCreateCommand{
		UI:         ui,
		DNSManager: dnsManager,
	}
}

func (cmd *ZoneCreateCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	zoneName := c.Args()[0]

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	resp, err := cmd.DNSManager.CreateZone(zoneName)
	if err != nil {
		return cli.NewExitError(T("Failed to create zone: {{.Zone}}.\n", map[string]interface{}{"Zone": zoneName})+err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, resp)
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("Zone {{.Zone}} was created.", map[string]interface{}{"Zone": zoneName}))
	return nil
}
