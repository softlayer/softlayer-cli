package cdn

import (
	"github.com/spf13/cobra"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type PurgeCommand struct {
	*metadata.SoftlayerCommand
	CdnManager managers.CdnManager
	Command    *cobra.Command
}

func NewPurgeCommand(sl *metadata.SoftlayerCommand) *PurgeCommand {
	thisCmd := &PurgeCommand{
		SoftlayerCommand: sl,
		CdnManager:       managers.NewCdnManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "purge " + T("IDENTIFIER") + " " + T("PATH"),
		Short: T("Creates a purge record and also initiates the purge call."),
		Long: T(`${COMMAND_NAME} sl cdn purge
Example:
${COMMAND_NAME} sl cdn purge 9779455 /article/file.txt"`),
		Args: metadata.TwoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *PurgeCommand) Run(args []string) error {
	cdnId := args[0]
	path := args[1]
	outputFormat := cmd.GetOutputFlag()

	datas, err := cmd.CdnManager.Purge(cdnId, path)
	if err != nil {
		return errors.NewAPIError(T("Failed to purge."), err.Error(), 2)
	}

	table := cmd.UI.Table([]string{
		T("Date"),
		T("Path"),
		T("Saved"),
		T("Status"),
	})
	for _, data := range datas {
		table.Add(
			utils.FormatStringToTime(data.Date),
			utils.FormatStringPointer(data.Path),
			utils.FormatStringPointer(data.Saved),
			utils.FormatStringPointer(data.Status),
		)
	}
	utils.PrintTable(cmd.UI, table, outputFormat)

	return nil
}
