package ipsec

import (
	"strconv"

	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type RemoveTranslationCommand struct {
	*metadata.SoftlayerCommand
	IPSECManager managers.IPSECManager
	Command      *cobra.Command
}

func NewRemoveTranslationCommand(sl *metadata.SoftlayerCommand) (cmd *RemoveTranslationCommand) {
	thisCmd := &RemoveTranslationCommand{
		SoftlayerCommand: sl,
		IPSECManager:     managers.NewIPSECManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "translation-remove " + T("CONTEXT_ID") + " " + T("TRANSLATION_ID"),
		Short: T("Remove a translation entry from an IPSec"),
		Long: T(`${COMMAND_NAME} sl ipsec translation-remove CONTEXT_ID TRANSLATION_ID 

  Remove a translation entry from an IPSEC tunnel context.

  A separate configuration request should be made to realize changes on
  network devices.`),
		Args: metadata.TwoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *RemoveTranslationCommand) Run(args []string) error {
	args0 := args[0]
	contextId, err := strconv.Atoi(args0)
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Context ID")
	}
	args1 := args[1]
	translationId, err := strconv.Atoi(args1)
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Translation ID")
	}
	_, err = cmd.IPSECManager.GetTranslation(contextId, translationId)
	if err != nil {
		return errors.NewAPIError(T("Failed to get translation with ID {{.TransID}} from IPSec {{.ID}}.",
			map[string]interface{}{"TransID": translationId, "ID": contextId}), err.Error(), 2)
	}
	err = cmd.IPSECManager.RemoveTranslation(contextId, translationId)
	if err != nil {
		return errors.NewAPIError(T("Failed to remove translation with ID {{.TransID}} from IPSec {{.ID}}.",
			map[string]interface{}{"TransID": translationId, "ID": contextId}), err.Error(), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Print(T("Removed translation with ID {{.TransID}} from IPSec {{.ID}}.",
		map[string]interface{}{"TransID": translationId, "ID": contextId}))
	return nil
}
