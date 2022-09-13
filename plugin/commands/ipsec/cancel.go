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

type CancelCommand struct {
	*metadata.SoftlayerCommand
	IPSECManager managers.IPSECManager
	Command      *cobra.Command
	Immediate    bool
	Reason       string
	ForceFlag    bool
}

func NewCancelCommand(sl *metadata.SoftlayerCommand) (cmd *CancelCommand) {
	thisCmd := &CancelCommand{
		SoftlayerCommand: sl,
		IPSECManager:     managers.NewIPSECManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "cancel " + T("CONTEXT_ID"),
		Short: T("Cancel a IPSec VPN tunnel context"),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().BoolVar(&thisCmd.Immediate, "immediate", false, T("Cancel the IPSec immediately instead of on the billing anniversary"))
	cobraCmd.Flags().StringVar(&thisCmd.Reason, "reason", "", T("An optional reason for cancellation"))
	cobraCmd.Flags().BoolVarP(&thisCmd.ForceFlag, "force", "f", false, T("Force operation without confirmation"))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *CancelCommand) Run(args []string) error {
	args0 := args[0]
	contextId, err := strconv.Atoi(args0)
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Context ID")
	}
	if !cmd.ForceFlag {
		confirm, err := cmd.UI.Confirm(T("This will cancel the IPSec: {{.ContextID}} and cannot be undone. Continue?", map[string]interface{}{"ContextID": contextId}))
		if err != nil {
			return err
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}
	err = cmd.IPSECManager.CancelTunnelContext(contextId, cmd.Immediate, cmd.Reason)
	if err != nil {
		return errors.NewAPIError(T("Failed to cancel IPSec {{.ContextID}}.\n", map[string]interface{}{"ContextID": contextId}), err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("IPSec {{.ContextID}} is cancelled.", map[string]interface{}{"ContextID": contextId}))
	return nil
}
