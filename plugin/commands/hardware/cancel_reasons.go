package hardware

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type CancelReasonsCommand struct {
	UI              terminal.UI
	HardwareManager managers.HardwareServerManager
}

func NewCancelReasonsCommand(ui terminal.UI, hardwareManager managers.HardwareServerManager) (cmd *CancelReasonsCommand) {
	return &CancelReasonsCommand{
		UI:              ui,
		HardwareManager: hardwareManager,
	}
}

func (cmd *CancelReasonsCommand) Run(c *cli.Context) error {
	reasons := cmd.HardwareManager.GetCancellationReasons()
	table := cmd.UI.Table([]string{T("Code"), T("Reason")})
	for key, value := range reasons {
		table.Add(key, value)
	}
	table.Print()
	return nil
}
