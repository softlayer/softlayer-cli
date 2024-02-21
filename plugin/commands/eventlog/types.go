package eventlog

import (
	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type TypesCommand struct {
	*metadata.SoftlayerCommand
	EventLogManager managers.EventLogManager
	Command         *cobra.Command
}

func NewTypesCommand(sl *metadata.SoftlayerCommand) (cmd *TypesCommand) {
	thisCmd := &TypesCommand{
		SoftlayerCommand: sl,
		EventLogManager:  managers.NewEventLogManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "types",
		Short: T("Get Event Log types"),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *TypesCommand) Run(args []string) error {

	outputFormat := cmd.GetOutputFlag()

	types, err := cmd.EventLogManager.GetEventLogTypes()
	if err != nil {
		return errors.NewAPIError(T("Failed to get Event Log types.\n"), err.Error(), 2)
	}

	table := cmd.UI.Table([]string{T("Types")})

	for _, typeEvent := range types {
		table.Add(typeEvent)
	}

	utils.PrintTable(cmd.UI, table, outputFormat)

	return nil
}
