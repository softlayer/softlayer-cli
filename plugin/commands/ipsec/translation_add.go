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

type AddTranslationCommand struct {
	*metadata.SoftlayerCommand
	IPSECManager managers.IPSECManager
	Command      *cobra.Command
	StaticIp     string
	RemoteIp     string
	Note         string
}

func NewAddTranslationCommand(sl *metadata.SoftlayerCommand) (cmd *AddTranslationCommand) {
	thisCmd := &AddTranslationCommand{
		SoftlayerCommand: sl,
		IPSECManager:     managers.NewIPSECManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "translation-add " + T("CONTEXT_ID"),
		Short: T("Add an address translation to an IPSec tunnel"),
		Long: T(`${COMMAND_NAME} sl ipsec translation-add CONTEXT_ID [OPTIONS]

  Add an address translation to an IPSEC tunnel context.

  A separate configuration request should be made to realize changes on
  network devices.`),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().StringVarP(&thisCmd.StaticIp, "static-ip", "s", "", T("Static IP address[required]"))
	cobraCmd.Flags().StringVarP(&thisCmd.RemoteIp, "remote-ip", "r", "", T("Remote IP address[required]"))
	cobraCmd.Flags().StringVarP(&thisCmd.Note, "note", "n", "", T("Note value"))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *AddTranslationCommand) Run(args []string) error {
	args0 := args[0]
	contextId, err := strconv.Atoi(args0)
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Context ID")
	}

	outputFormat := cmd.GetOutputFlag()

	staticIp := cmd.StaticIp
	if staticIp == "" {
		return errors.NewMissingInputError("-s|--static-ip")
	}
	remoteIp := cmd.RemoteIp
	if remoteIp == "" {
		return errors.NewMissingInputError("-r|--remote-ip")
	}
	_, err = cmd.IPSECManager.GetTunnelContext(contextId, "")
	if err != nil {
		return errors.NewAPIError(T("Failed to get IPSec with ID {{.ID}}.\n", map[string]interface{}{"ID": contextId}), err.Error(), 2)
	}
	translation, err := cmd.IPSECManager.CreateTranslation(contextId, staticIp, remoteIp, cmd.Note)
	if err != nil {
		return errors.NewAPIError(T("Failed to create translation for IPSec with ID {{.ID}}.\n", map[string]interface{}{"ID": contextId}), err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, translation)
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("Created translation from {{.StaticIP}} to {{.RemoteIP}} #{{.ID}}.",
		map[string]interface{}{"StaticIP": staticIp, "RemoteIP": remoteIp, "ID": *translation.Id}))
	return nil
}
