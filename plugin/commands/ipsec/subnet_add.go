package ipsec

import (
	"strconv"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type AddSubnetCommand struct {
	UI           terminal.UI
	IPSECManager managers.IPSECManager
}

func NewAddSubnetCommand(ui terminal.UI, ipsecManager managers.IPSECManager) (cmd *AddSubnetCommand) {
	return &AddSubnetCommand{
		UI:           ui,
		IPSECManager: ipsecManager,
	}
}

func (cmd *AddSubnetCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	args0 := c.Args()[0]
	contextId, err := strconv.Atoi(args0)
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Context ID")
	}
	subnetId := c.Int("s")
	subnetType := c.String("t")
	if subnetType != "internal" && subnetType != "remote" && subnetType != "service" {
		return errors.NewInvalidUsageError(T("The subnet type has to be either internal, or remote or service."))
	}
	networkIdentifier := c.String("n")
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
		return cli.NewExitError(T("Failed to get IPSec with ID {{.ID}}.\n", map[string]interface{}{"ID": contextId})+err.Error(), 2)
	}
	if createRemote {
		ids := strings.Split(networkIdentifier, "/")
		id := ids[0]
		cidr, _ := strconv.Atoi(ids[1])
		subnet, err := cmd.IPSECManager.CreateRemoteSubnet(*context.AccountId, id, cidr)
		if err != nil {
			return cli.NewExitError(T("Failed to create subnet with {{.ID}}.\n", map[string]interface{}{"ID": networkIdentifier})+err.Error(), 2)
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
	return cli.NewExitError(T("Failed to add {{.Type}} subnet #{{.ID}} to IPSec {{.ContextID}}.\n",
		map[string]interface{}{"Type": subnetType, "ID": subnetId, "ContextID": contextId})+err.Error(), 2)
}

func IpsecSubnetAddMetaData() cli.Command {
	return cli.Command{
		Category:    "ipsec",
		Name:        "subnet-add",
		Description: T("Add a subnet to an IPSec tunnel context"),
		Usage: T(`${COMMAND_NAME} sl ipsec subnet-add CONTEXT_ID [OPTIONS] 

  Add a subnet to an IPSEC tunnel context.

  A subnet id may be specified to link to the existing tunnel context.

  Otherwise, a network identifier in CIDR notation should be specified,
  indicating that a subnet resource should first be created before
  associating it with the tunnel context. Note that this is only supported
  for remote subnets, which are also deleted upon failure to attach to a
  context.

  A separate configuration request should be made to realize changes on
  network devices.`),
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "s,subnet-id",
				Usage: T("Subnet identifier to add, required"),
			},
			cli.StringFlag{
				Name:  "t,subnet-type",
				Usage: T("Subnet type to add. Options are: internal,remote,service[required]"),
			},
			cli.StringFlag{
				Name:  "n,network",
				Usage: T("Subnet network identifier to create"),
			},
		},
	}
}
