package virtual

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
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

func VSCaptureMetaData() cli.Command {
	return cli.Command{
		Category:    "vs",
		Name:        "capture",
		Description: T("Capture virtual server instance into an image"),
		Usage: T(`${COMMAND_NAME} sl vs capture IDENTIFIER [OPTIONS]
	
EXAMPLE:
   ${COMMAND_NAME} sl vs capture 12345678 -n mycloud --all --note testing
   This command captures virtual server instance with ID of 12345678 with all disks into an image named "mycloud" with note "testing".`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "n,name",
				Usage: T("Name of the image [required]"),
			},
			cli.BoolFlag{
				Name:  "all",
				Usage: T("Capture all disks that belong to the virtual server"),
			},
			cli.StringFlag{
				Name:  "note",
				Usage: T("Add a note to be associated with the image"),
			},
			metadata.OutputFlag(),
		},
	}
}
