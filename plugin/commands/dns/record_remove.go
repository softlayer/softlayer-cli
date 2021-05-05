package dns

import (
	"strconv"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type RecordRemoveCommand struct {
	UI         terminal.UI
	DNSManager managers.DNSManager
}

func NewRecordRemoveCommand(ui terminal.UI, dnsManager managers.DNSManager) (cmd *RecordRemoveCommand) {
	return &RecordRemoveCommand{
		UI:         ui,
		DNSManager: dnsManager,
	}
}

func (cmd *RecordRemoveCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}

	recordID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("Record ID")
	}

	err = cmd.DNSManager.DeleteResourceRecord(recordID)
	if err != nil {
		if strings.Contains(err.Error(), slErrors.SL_EXP_OBJ_NOT_FOUND) {
			return cli.NewExitError(T("Unable to find resource record with ID: {{.RecordID}}.\n", map[string]interface{}{"RecordID": recordID})+err.Error(), 0)
		}
		return cli.NewExitError(T("Failed to delete resource record: {{.RecordID}}.\n", map[string]interface{}{"RecordID": recordID})+err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Resource record {{.ID}} was removed.", map[string]interface{}{"ID": recordID}))
	return nil
}
