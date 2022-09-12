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

type UpdateCommand struct {
	*metadata.SoftlayerCommand
	IPSECManager         managers.IPSECManager
	Command              *cobra.Command
	Name                 string
	RemotePeer           string
	PresharedKey         string
	Phase1Auth           string
	Phase1Crypto         string
	Phase1Dh             int
	Phase1KeyTtl         int
	Phase2Auth           string
	Phase2Crypto         string
	Phase2Dh             int
	Phase2ForwardSecrecy int
	Phase2KeyTtl         int
}

func NewUpdateCommand(sl *metadata.SoftlayerCommand) (cmd *UpdateCommand) {
	thisCmd := &UpdateCommand{
		SoftlayerCommand: sl,
		IPSECManager:     managers.NewIPSECManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "update " + T("CONTEXT_ID"),
		Short: T("Update tunnel context properties"),
		Long: T(`${COMMAND_NAME} sl ipsec update CONTEXT_ID [OPTIONS]

  Update tunnel context properties.

  Updates are made atomically, so either all are accepted or none are.

  Key life values must be in the range 120-172800.

  Phase 2 perfect forward secrecy must be in the range 0-1.

  A separate configuration request should be made to realize changes on
  network devices.`),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().StringVarP(&thisCmd.Name, "name", "n", "", T("Friendly name"))
	cobraCmd.Flags().StringVarP(&thisCmd.RemotePeer, "remote-peer", "r", "", T("Remote peer IP address"))
	cobraCmd.Flags().StringVarP(&thisCmd.PresharedKey, "preshared-key", "k", "", T("Preshared key"))
	cobraCmd.Flags().StringVarP(&thisCmd.Phase1Auth, "phase1-auth", "a", "", T("Phase 1 authentication. Options are: MD5,SHA1,SHA256"))
	cobraCmd.Flags().StringVarP(&thisCmd.Phase1Crypto, "phase1-crypto", "c", "", T("Phase 1 encryption. Options are: DES,3DES,AES128,AES192,AES256"))
	cobraCmd.Flags().IntVarP(&thisCmd.Phase1Dh, "phase1-dh", "d", 0, T("Phase 1 Diffie-Hellman group. Options are: 0,1,2,5"))
	cobraCmd.Flags().IntVarP(&thisCmd.Phase1KeyTtl, "phase1-key-ttl", "t", 0, T("Phase 1 key life. Range is 120-172800"))
	cobraCmd.Flags().StringVarP(&thisCmd.Phase2Auth, "phase2-auth", "u", "", T("Phase 2 authentication. Options are: MD5,SHA1,SHA256"))
	cobraCmd.Flags().StringVarP(&thisCmd.Phase2Crypto, "phase2-crypto", "y", "", T("Phase 2 encryption. Options are: DES,3DES,AES128,AES192,AES256"))
	cobraCmd.Flags().IntVarP(&thisCmd.Phase2Dh, "phase2-dh", "e", 0, T("Phase 2 Diffie-Hellman group. Options are: 0,1,2,5"))
	cobraCmd.Flags().IntVarP(&thisCmd.Phase2ForwardSecrecy, "phase2-forward-secrecy", "f", 0, T("Phase 2 perfect forward secrecy. Range is 0-1"))
	cobraCmd.Flags().IntVarP(&thisCmd.Phase2KeyTtl, "phase2-key-ttl", "l", 0, T("Phase 2 key life. Range is 120-172800"))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *UpdateCommand) Run(args []string) error {
	args0 := args[0]
	contextId, err := strconv.Atoi(args0)
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Context ID")
	}
	if cmd.Phase1Auth != "" {
		phase1Auth := cmd.Phase1Auth
		if phase1Auth != "MD5" && phase1Auth != "SHA1" && phase1Auth != "SHA256" {
			return errors.NewInvalidUsageError(T("-a|--phase1-auth must be either MD5 or SHA1 or SHA256."))
		}
	}
	if cmd.Phase1Crypto != "" {
		phase1Crypto := cmd.Phase1Crypto
		if phase1Crypto != "DES" && phase1Crypto != "3DES" && phase1Crypto != "AES128" && phase1Crypto != "AES192" && phase1Crypto != "AES256" {
			return errors.NewInvalidUsageError(T("-c|--phase1-crypto must be either DES or 3DES or AES128 or AES192 or AES256."))
		}
	}
	if cmd.Phase1Dh != 0 {
		phase1Dh := cmd.Phase1Dh
		if phase1Dh != 0 && phase1Dh != 1 && phase1Dh != 2 && phase1Dh != 5 {
			return errors.NewInvalidUsageError(T("-d|--phase1-dh must be either 0 or 1 or 2 or 5."))
		}
	}
	if cmd.Phase1KeyTtl != 0 {
		phase1KeyLife := cmd.Phase1KeyTtl
		if phase1KeyLife < 120 || phase1KeyLife > 172800 {
			return errors.NewInvalidUsageError(T("-t|--phase1-key-ttl must be in range 120-172800."))
		}
	}
	if cmd.Phase2Auth != "" {
		phase2Auth := cmd.Phase2Auth
		if phase2Auth != "MD5" && phase2Auth != "SHA1" && phase2Auth != "SHA256" {
			return errors.NewInvalidUsageError(T("-u|--phase2-auth must be either MD5 or SHA1 or SHA256."))
		}
	}
	if cmd.Phase2Crypto != "" {
		phase2Crypto := cmd.Phase2Crypto
		if phase2Crypto != "DES" && phase2Crypto != "3DES" && phase2Crypto != "AES128" && phase2Crypto != "AES192" && phase2Crypto != "AES256" {
			return errors.NewInvalidUsageError(T("-y|--phase2-crypto must be either DES or 3DES or AES128 or AES192 or AES256."))
		}
	}
	if cmd.Phase2Dh != 0 {
		phase2Dh := cmd.Phase2Dh
		if phase2Dh != 0 && phase2Dh != 1 && phase2Dh != 2 && phase2Dh != 5 {
			return errors.NewInvalidUsageError(T("-e|--phase2-dh must be either 0 or 1 or 2 or 5."))
		}
	}
	if cmd.Phase2ForwardSecrecy != 0 {
		phase2ForwardSecrecy := cmd.Phase2ForwardSecrecy
		if phase2ForwardSecrecy != 0 && phase2ForwardSecrecy != 1 {
			return errors.NewInvalidUsageError(T("-f|--phase2-forward-secrecy must be either 0 or 1."))
		}
	}
	if cmd.Phase2KeyTtl != 0 {
		phase2KeyLife := cmd.Phase2KeyTtl
		if phase2KeyLife < 120 || phase2KeyLife > 172800 {
			return errors.NewInvalidUsageError(T("-l|--phase2-key-ttl must be in range 120-172800."))
		}
	}

	outputFormat := cmd.GetOutputFlag()

	resp, err := cmd.IPSECManager.UpdateTunnelContext(contextId,
		cmd.Name,
		cmd.RemotePeer,
		cmd.PresharedKey,
		cmd.Phase1Auth,
		cmd.Phase1Crypto,
		cmd.Phase1Dh,
		cmd.Phase1KeyTtl,
		cmd.Phase2Auth,
		cmd.Phase2Crypto,
		cmd.Phase2Dh,
		cmd.Phase2ForwardSecrecy,
		cmd.Phase2KeyTtl)
	if err != nil {
		return errors.NewAPIError(T("Failed to update IPSec {{.ContextID}}.\n", map[string]interface{}{"ContextID": contextId}), err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, resp)
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("Updated IPSec {{.ContextID}}.", map[string]interface{}{"ContextID": contextId}))
	return nil
}
