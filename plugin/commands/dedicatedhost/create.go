package dedicatedhost

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/urfave/cli"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type CreateCommand struct {
	UI                   terminal.UI
	DedicatedHostManager managers.DedicatedHostManager
	NetworkManager       managers.NetworkManager
	Context              plugin.PluginContext
}

func NewCreateCommand(ui terminal.UI, dedicatedHostManager managers.DedicatedHostManager, networkManager managers.NetworkManager, context plugin.PluginContext) (cmd *CreateCommand) {
	return &CreateCommand{
		UI:                   ui,
		DedicatedHostManager: dedicatedHostManager,
		NetworkManager:       networkManager,
		Context:              context,
	}
}

func DedicatedhostCreateMetaData() cli.Command {
	return cli.Command{
		Category:    "dedicatedhost",
		Name:        "create",
		Description: T("Create a dedicatedhost"),
		Usage:       "${COMMAND_NAME} sl dedicatedhost create [OPTIONS]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "H,hostname",
				Usage: T("Host portion of the FQDN [required]"),
			},
			cli.StringFlag{
				Name:  "D,domain",
				Usage: T("Domain portion of the FQDN [required]"),
			},
			cli.StringFlag{
				Name:  "d,datacenter",
				Usage: T("Datacenter shortname [required]"),
			},
			cli.StringFlag{
				Name:  "s,size",
				Usage: T("Size of the dedicated host, currently only one size is available: 56_CORES_X_242_RAM_X_1_4_TB"),
			},
			cli.StringFlag{
				Name:  "b,billing",
				Usage: T("Billing rate. Default is: hourly. Options are: hourly, monthly"),
			},
			cli.StringFlag{
				Name:  "v,vlan-private",
				Usage: T("The ID of the private VLAN on which you want the dedicated host placed. See: '${COMMAND_NAME} sl vlan list' for reference"),
			},
			cli.BoolFlag{
				Name:  "test",
				Usage: T("Do not actually create the dedicatedhost"),
			},
			metadata.ForceFlag(),
			metadata.OutputFlag(),
		},
	}
}

func (cmd *CreateCommand) Run(c *cli.Context) error {
	size := managers.HOST_DEFAULT_SIZE
	if c.IsSet("size") {
		size = c.String("size")
	}
	hostname := c.String("H")
	if hostname == "" {
		return errors.NewMissingInputError("-H|--hostname")
	}
	domain := c.String("D")
	if domain == "" {
		return errors.NewMissingInputError("-D|--domain")
	}
	datacenter := c.String("d")
	if datacenter == "" {
		return errors.NewMissingInputError("-d|--datacenter")
	}
	billing := "hourly"
	if c.IsSet("b") {
		billing = c.String("b")
		if billing != "hourly" && billing != "monthly" {
			return errors.NewInvalidUsageError(T("[-b|--billing] has to be either hourly or monthly."))
		}
	}
	vlanId := c.Int("v")
	if vlanId == 0 {
		return errors.NewMissingInputError("-v|--vlan-private")
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	vlan, err := cmd.NetworkManager.GetVlan(vlanId, "mask[id,primaryRouter[id,datacenter]]")
	if err != nil {
		return cli.NewExitError(T("Failed to get vlan {{.VlanId}}.\n", map[string]interface{}{"VlanId": vlanId})+err.Error(), 2)
	}
	if !c.IsSet("f") && !c.IsSet("force") && outputFormat != "JSON" && !c.IsSet("test") {
		confirm, err := cmd.UI.Confirm(T("This action will incur charges on your account. Continue?"))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}
	if vlan.PrimaryRouter == nil || vlan.PrimaryRouter.Id == nil {
		return cli.NewExitError(T("Failed to get vlan primary router ID."), 2)
	}

	vlanDatacenter := vlan.PrimaryRouter.Datacenter.Name
	if *vlanDatacenter != datacenter {
		return cli.NewExitError(T("The vlan is located at: {{.VLAN}}, Please add a valid private vlan according the datacenter selected.", map[string]interface{}{"VLAN": *vlanDatacenter}), 2)
	}

	orderTemplate, err := cmd.DedicatedHostManager.GenerateOrderTemplate(size, hostname, domain, datacenter, billing, *vlan.PrimaryRouter.Id)
	if err != nil {
		return cli.NewExitError(T("Failed to generate the order template.\n")+err.Error(), 2)
	}

	if c.IsSet("test") {
		orderReceipt, err := cmd.DedicatedHostManager.VerifyInstanceCreation(orderTemplate)
		if err != nil {
			return cli.NewExitError(T("Failed to verify virtual server creation.\n")+err.Error(), 2)
		}
		if outputFormat == "JSON" {
			return utils.PrintPrettyJSON(cmd.UI, orderReceipt)
		}
		cmd.UI.Ok()
		cmd.UI.Print(T("The order is correct."))
	} else {
		orderReceipt, err := cmd.DedicatedHostManager.OrderInstance(orderTemplate)
		if err != nil {
			return cli.NewExitError(T("Failed to Order the dedicatedhost.\n")+err.Error(), 2)
		}
		if outputFormat == "JSON" {
			return utils.PrintPrettyJSON(cmd.UI, orderReceipt)
		}
		cmd.UI.Ok()
		cmd.UI.Print(T("The order {{.OrderID}} was placed.", map[string]interface{}{"OrderID": *orderReceipt.OrderId}))
	}

	return nil
}
