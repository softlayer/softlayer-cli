package ipsec

import (
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type AddSubnetCommand struct {
	*metadata.SoftlayerCommand
	IPSECManager managers.IPSECManager
	Command      *cobra.Command
	SubnetId     int
	SubnetType   string
	Network      string
}

func NewAddSubnetCommand(sl *metadata.SoftlayerCommand) (cmd *AddSubnetCommand) {
	thisCmd := &AddSubnetCommand{
		SoftlayerCommand: sl,
		IPSECManager:     managers.NewIPSECManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "subnet-add " + T("CONTEXT_ID"),
		Short: T("Add a subnet to an IPSec tunnel context"),
		Long: T(`${COMMAND_NAME} sl ipsec subnet-add CONTEXT_ID [OPTIONS] 

  Add a subnet to an IPSEC tunnel context.

  A subnet id may be specified to link to the existing tunnel context.

  Otherwise, a network identifier in CIDR notation should be specified,
  indicating that a subnet resource should first be created before
  associating it with the tunnel context. Note that this is only supported
  for remote subnets, which are also deleted upon failure to attach to a
  context.

  A separate configuration request should be made to realize changes on
  network devices.`),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().IntVarP(&thisCmd.SubnetId, "subnet-id", "s", 0, T("Subnet identifier to add, required"))
	cobraCmd.Flags().StringVarP(&thisCmd.SubnetType, "subnet-type", "t", "", T("Subnet type to add. Options are: internal,remote,service[required]"))
	cobraCmd.Flags().StringVarP(&thisCmd.Network, "network", "n", "", T("Subnet network identifier to create"))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *AddSubnetCommand) Run(args []string) error {
	args0 := args[0]
	contextId, err := strconv.Atoi(args0)
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Context ID")
	}
	subnetId := cmd.SubnetId
	subnetType := cmd.SubnetType
	if subnetType != "internal" && subnetType != "remote" && subnetType != "service" {
		return errors.NewInvalidUsageError(T("The subnet type has to be either internal, or remote or service."))
	}
	networkIdentifier := cmd.Network
	createRemote := false
	if subnetId == 0 {
		if networkIdentifier == "" {
			return errors.NewInvalidUsageError(T("Either -s|--subnet-id or -n|--network must be provided."))
		}
		if subnetType != "remote" {
			return errors.NewInvalidUsageError(T("Unable to create {{.Type}} subnet.", map[string]interface{}{"Type": subnetType}))
		}
		createRemote = true
	}
	context, err := cmd.IPSECManager.GetTunnelContext(contextId, "id,accountId")
	if err != nil {
		return errors.NewAPIError(T("Failed to get IPSec with ID {{.ID}}.\n", map[string]interface{}{"ID": contextId}), err.Error(), 2)
	}
	if createRemote {
		ids := strings.Split(networkIdentifier, "/")
		id := ids[0]
		cidr, _ := strconv.Atoi(ids[1])
		subnet, err := cmd.IPSECManager.CreateRemoteSubnet(*context.AccountId, id, cidr)
		if err != nil {
			return errors.NewAPIError(T("Failed to create subnet with {{.ID}}.\n", map[string]interface{}{"ID": networkIdentifier}), err.Error(), 2)
		}
		subnetId = *subnet.Id
		cmd.UI.Print(T("Created subnet {{.ID}}/{{.CIDR}} #{{.Identifier}}.",
			map[string]interface{}{"ID": id, "CIDR": cidr, "Identifier": *subnet.Id}))
	}
	succeeded := false
	if subnetType == "internal" {
		err = cmd.IPSECManager.AddInternalSubnet(contextId, subnetId)
		if err == nil {
			succeeded = true
		}
	} else if subnetType == "remote" {
		err = cmd.IPSECManager.AddRemoteSubnet(contextId, subnetId)
		if err == nil {
			succeeded = true
		}
	} else if subnetType == "service" {
		err = cmd.IPSECManager.AddServiceSubnet(contextId, subnetId)
		if err == nil {
			succeeded = true
		}
	}
	if succeeded {
		cmd.UI.Ok()
		cmd.UI.Print(T("Added {{.Type}} subnet #{{.ID}} to IPSec {{.ContextID}}.",
			map[string]interface{}{"Type": subnetType, "ID": subnetId, "ContextID": contextId}))
		return nil
	}
	return errors.NewAPIError(T("Failed to add {{.Type}} subnet #{{.ID}} to IPSec {{.ContextID}}.\n",
		map[string]interface{}{"Type": subnetType, "ID": subnetId, "ContextID": contextId}), err.Error(), 2)
}
