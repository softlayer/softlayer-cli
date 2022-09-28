package virtual

import (
	"strconv"

	"github.com/spf13/cobra"

	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type CaptureCommand struct {
	*metadata.SoftlayerCommand
	VirtualServerManager managers.VirtualServerManager
	Command              *cobra.Command
	Name                 string
	All                  bool
	Note                 string
}

func NewCaptureCommand(sl *metadata.SoftlayerCommand) (cmd *CaptureCommand) {
	thisCmd := &CaptureCommand{
		SoftlayerCommand:     sl,
		VirtualServerManager: managers.NewVirtualServerManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "capture " + T("IDENTIFIER"),
		Short: T("Capture virtual server instance into an image"),
		Long: T(`${COMMAND_NAME} sl vs capture IDENTIFIER [OPTIONS]
	
EXAMPLE:
   ${COMMAND_NAME} sl vs capture 12345678 -n mycloud --all --note testing
   This command captures virtual server instance with ID of 12345678 with all disks into an image named "mycloud" with note "testing".`),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	cobraCmd.Flags().StringVarP(&thisCmd.Name, "name", "n", "", T("Name of the image [required]"))
	cobraCmd.Flags().BoolVar(&thisCmd.All, "all", false, T("Capture all disks that belong to the virtual server"))
	cobraCmd.Flags().StringVar(&thisCmd.Note, "note", "", T("Add a note to be associated with the image"))
	return thisCmd
}
func (cmd *CaptureCommand) Run(args []string) error {

	vsID, err := utils.ResolveVirtualGuestId(args[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("Virtual server ID")
	}
	if cmd.Name == "" {
		return slErrors.NewMissingInputError("-n|--name")
	}

	outputFormat := cmd.GetOutputFlag()

	txn, err := cmd.VirtualServerManager.CaptureImage(vsID, cmd.Name, cmd.Note, cmd.All)
	if err != nil {
		return slErrors.NewAPIError(T("Failed to capture image for virtual server instance: {{.VsID}}.\n",
			map[string]interface{}{"VsID": vsID}), err.Error(), 2)
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
	table.Add(T("all_disks"), strconv.FormatBool(cmd.All))
	table.Print()
	return nil
}
