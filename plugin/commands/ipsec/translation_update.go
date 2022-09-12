package ipsec

import (
	"strconv"

	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type UpdateTranslationCommand struct {
	*metadata.SoftlayerCommand
	IPSECManager managers.IPSECManager
	Command      *cobra.Command
	StaticIp     string
	RemoteIp     string
	Note         string
}

func NewUpdateTranslationCommand(sl *metadata.SoftlayerCommand) (cmd *UpdateTranslationCommand) {
	thisCmd := &UpdateTranslationCommand{
		SoftlayerCommand: sl,
		IPSECManager:     managers.NewIPSECManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "translation-update " + T("CONTEXT_ID") + " " + T("TRANSLATION_ID"),
		Short: T("Update an address translation for an IPSec"),
		Long: T(`${COMMAND_NAME} sl ipsec translation-update CONTEXT_ID TRANSLATION_ID [OPTIONS]

  Update an address translation for an IPSEC tunnel context.

  A separate configuration request should be made to realize changes on
  network devices.`),
		Args: metadata.TwoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().StringVarP(&thisCmd.StaticIp, "static-ip", "s", "", T("Static IP address[required]"))
	cobraCmd.Flags().StringVarP(&thisCmd.RemoteIp, "remote-ip", "r", "", T("Remote IP address[required]"))
	cobraCmd.Flags().StringVarP(&thisCmd.Note, "note", "n", "", T("Note"))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *UpdateTranslationCommand) Run(args []string) error {
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

	outputFormat := cmd.GetOutputFlag()

	resp, err := cmd.IPSECManager.UpdateTranslation(contextId, translationId, cmd.StaticIp, cmd.RemoteIp, cmd.Note)
	if err != nil {
		return errors.NewAPIError(T("Failed to update translation with ID {{.TransID}} in IPSec {{.ID}}.",
			map[string]interface{}{"TransID": translationId, "ID": contextId}), err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, resp)
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("Updated translation with ID {{.TransID}} in IPSec {{.ID}}.",
		map[string]interface{}{"TransID": translationId, "ID": contextId}))
	return nil
}
