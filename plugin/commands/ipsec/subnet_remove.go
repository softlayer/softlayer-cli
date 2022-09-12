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

type RemoveSubnetCommand struct {
	*metadata.SoftlayerCommand
	IPSECManager managers.IPSECManager
	Command      *cobra.Command
}

func NewRemoveSubnetCommand(sl *metadata.SoftlayerCommand) (cmd *RemoveSubnetCommand) {
	thisCmd := &RemoveSubnetCommand{
		SoftlayerCommand: sl,
		IPSECManager:     managers.NewIPSECManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "subnet-remove " + T("CONTEXT_ID") + " " + T("SUBNET_ID") + " " + T("SUBNET_TYPE"),
		Short: T("Remove a subnet from an IPSEC tunnel context"),
		Long: T(`${COMMAND_NAME} sl ipsec subnet-remove CONTEXT_ID SUBNET_ID SUBNET_TYPE 

  Remove a subnet from an IPSEC tunnel context.

  The subnet id to remove must be specified.

  Remote subnets are deleted upon removal from a tunnel context.

  A separate configuration request should be made to realize changes on
  network devices.`),
		Args: metadata.ThreeArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *RemoveSubnetCommand) Run(args []string) error {
	args0 := args[0]
	contextId, err := strconv.Atoi(args0)
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Context ID")
	}
	args1 := args[1]
	subnetId, err := strconv.Atoi(args1)
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Subnet ID")
	}
	subnetType := args[2]
	if subnetType != "internal" && subnetType != "remote" && subnetType != "service" {
		return errors.NewInvalidUsageError(T("The subnet type has to be either internal, or remote or service."))
	}
	_, err = cmd.IPSECManager.GetTunnelContext(contextId, "")
	if err != nil {
		return errors.NewAPIError(T("Failed to get IPSec with ID {{.ID}}.\n", map[string]interface{}{"ID": contextId}), err.Error(), 2)
	}
	succeeded := false
	if subnetType == "internal" {
		err = cmd.IPSECManager.RemoveInternalSubnet(contextId, subnetId)
		if err == nil {
			succeeded = true
		}
	} else if subnetType == "remote" {
		err = cmd.IPSECManager.RemoveRemoteSubnet(contextId, subnetId)
		if err == nil {
			succeeded = true
		}
	} else if subnetType == "service" {
		err = cmd.IPSECManager.RemoveServiceSubnet(contextId, subnetId)
		if err == nil {
			succeeded = true
		}
	}
	if succeeded {
		cmd.UI.Ok()
		cmd.UI.Print(T("Removed {{.Type}} subnet #{{.ID}} from IPSec {{.ContextID}}.",
			map[string]interface{}{"Type": subnetType, "ID": subnetId, "ContextID": contextId}))
		return nil
	}
	return errors.NewAPIError(T("Failed to remove {{.Type}} subnet #{{.ID}} from IPSec {{.ContextID}}.\n",
		map[string]interface{}{"Type": subnetType, "ID": subnetId, "ContextID": contextId}), err.Error(), 2)
}
