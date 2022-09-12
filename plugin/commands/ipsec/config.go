package ipsec

import (
	"strconv"

	"github.com/spf13/cobra"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type ConfigCommand struct {
	*metadata.SoftlayerCommand
	IPSECManager managers.IPSECManager
	Command      *cobra.Command
}

func NewConfigCommand(sl *metadata.SoftlayerCommand) (cmd *ConfigCommand) {
	thisCmd := &ConfigCommand{
		SoftlayerCommand: sl,
		IPSECManager:     managers.NewIPSECManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "config " + T("CONTEXT_ID"),
		Short: T("Request configuration of a tunnel context"),
		Long: T(`${COMMAND_NAME} sl ipsec config CONTEXT_ID [OPTIONS]

  Request configuration of a tunnel context.

  This action will update the advancedConfigurationFlag on the context
  instance and further modifications against the context will be prevented
  until all changes can be propagated to network devices.`),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *ConfigCommand) Run(args []string) error {
	args0 := args[0]
	contextId, err := strconv.Atoi(args0)
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Context ID")
	}
	err = cmd.IPSECManager.ApplyConfiguration(contextId)
	if err != nil {
		return errors.NewAPIError(T("Failed to enqueue configuration request for IPSec {{.ContextID}}.\n", map[string]interface{}{"ContextID": contextId}), err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Configuration request received for IPSec {{.ContextID}}.", map[string]interface{}{"ContextID": contextId}))
	return nil
}
