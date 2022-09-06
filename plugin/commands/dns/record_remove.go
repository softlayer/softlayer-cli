package dns

import (
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type RecordRemoveCommand struct {
	*metadata.SoftlayerCommand
	DNSManager managers.DNSManager
	Command    *cobra.Command
}

func NewRecordRemoveCommand(sl *metadata.SoftlayerCommand) *RecordRemoveCommand {
	thisCmd := &RecordRemoveCommand{
		SoftlayerCommand: sl,
		DNSManager:       managers.NewDNSManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "record-remove " + T("RECORD_ID"),
		Short: T("Remove resource record from a zone."),
		Long: T(`${COMMAND_NAME} sl dns record-remove RECORD_ID
	
EXAMPLE:
	${COMMAND_NAME} sl dns record-remove 12345678
	This command removes resource record with ID 12345678.`),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *RecordRemoveCommand) Run(args []string) error {
	recordID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("Record ID")
	}

	err = cmd.DNSManager.DeleteResourceRecord(recordID)
	if err != nil {
		if strings.Contains(err.Error(), slErrors.SL_EXP_OBJ_NOT_FOUND) {
			return errors.NewAPIError(T("Unable to find resource record with ID: {{.RecordID}}.\n", map[string]interface{}{"RecordID": recordID}), err.Error(), 0)
		}
		return errors.NewAPIError(T("Failed to delete resource record: {{.RecordID}}.\n", map[string]interface{}{"RecordID": recordID}), err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Resource record {{.ID}} was removed.", map[string]interface{}{"ID": recordID}))
	return nil
}
