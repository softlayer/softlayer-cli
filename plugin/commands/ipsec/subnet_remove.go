package ipsec

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type RemoveSubnetCommand struct {
	UI           terminal.UI
	IPSECManager managers.IPSECManager
}

func NewRemoveSubnetCommand(ui terminal.UI, ipsecManager managers.IPSECManager) (cmd *RemoveSubnetCommand) {
	return &RemoveSubnetCommand{
		UI:           ui,
		IPSECManager: ipsecManager,
	}
}

func (cmd *RemoveSubnetCommand) Run(c *cli.Context) error {
	if c.NArg() != 3 {
		return errors.NewInvalidUsageError(T("This command requires three arguments."))
	}
	args0 := c.Args()[0]
	contextId, err := strconv.Atoi(args0)
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Context ID")
	}
	args1 := c.Args()[1]
	subnetId, err := strconv.Atoi(args1)
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Subnet ID")
	}
	subnetType := c.Args()[2]
	if subnetType != "internal" && subnetType != "remote" && subnetType != "service" {
		return errors.NewInvalidUsageError(T("The subnet type has to be either internal, or remote or service."))
	}
	_, err = cmd.IPSECManager.GetTunnelContext(contextId, "")
	if err != nil {
		return cli.NewExitError(T("Failed to get IPSec with ID {{.ID}}.\n", map[string]interface{}{"ID": contextId})+err.Error(), 2)
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
	return cli.NewExitError(T("Failed to remove {{.Type}} subnet #{{.ID}} from IPSec {{.ContextID}}.\n",
		map[string]interface{}{"Type": subnetType, "ID": subnetId, "ContextID": contextId})+err.Error(), 2)
}
