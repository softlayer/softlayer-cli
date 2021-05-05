package virtual

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type CaptureCommand struct {
	UI                   terminal.UI
	VirtualServerManager managers.VirtualServerManager
}

func NewCaptureCommand(ui terminal.UI, virtualServerManager managers.VirtualServerManager) (cmd *CaptureCommand) {
	return &CaptureCommand{
		UI:                   ui,
		VirtualServerManager: virtualServerManager,
	}
}

func (cmd *CaptureCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	vsID, err := utils.ResolveVirtualGuestId(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Virtual server ID")
	}
	if !c.IsSet("name") {
		return errors.NewMissingInputError("-n|--name")
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	txn, err := cmd.VirtualServerManager.CaptureImage(vsID, c.String("name"), c.String("note"), c.Bool("all"))
	if err != nil {
		return cli.NewExitError(T("Failed to capture image for virtual server instance: {{.VsID}}.\n",
			map[string]interface{}{"VsID": vsID})+err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, txn)
	}

	table := cmd.UI.Table([]string{T("name"), T("value")})
	table.Add(T("vs_id"), utils.FormatIntPointer(txn.GuestId))
	table.Add(T("date_time"), utils.FormatSLTimePointer(txn.CreateDate))
	var transaction string
	if txn.TransactionStatus != nil {
		transaction = utils.FormatStringPointer(txn.TransactionStatus.Name)
	}
	table.Add(T("transaction"), transaction)
	table.Add(T("transaction_id"), utils.FormatIntPointer(txn.Id))
	table.Add(T("all_disks"), strconv.FormatBool(c.Bool("all")))
	table.Print()
	return nil
}
