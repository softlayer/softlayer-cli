package dedicatedhost

import (
	"github.com/spf13/cobra"
	"github.com/urfave/cli"

	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type CreateCommand struct {
	*metadata.SoftlayerCommand
	DedicatedHostManager managers.DedicatedHostManager
	NetworkManager       managers.NetworkManager
	Command              *cobra.Command
	Hostname             string
	Domain               string
	Datacenter           string
	Size                 string
	Billing              string
	VlanPrivate          int
	Test                 bool
	Force                bool
}

func NewCreateCommand(sl *metadata.SoftlayerCommand) *CreateCommand {
	thisCmd := &CreateCommand{
		SoftlayerCommand:     sl,
		DedicatedHostManager: managers.NewDedicatedhostManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "create",
		Short: T("Create a dedicatedhost"),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().StringVar(&thisCmd.Hostname, "hostname", "", T("Host portion of the FQDN [required]"))
	cobraCmd.Flags().StringVar(&thisCmd.Domain, "domain", "", T("Domain portion of the FQDN [required]"))
	cobraCmd.Flags().StringVar(&thisCmd.Datacenter, "datacenter", "", T("Datacenter shortname [required]"))
	cobraCmd.Flags().StringVar(&thisCmd.Size, "size", "", T("Size of the dedicated host, currently only one size is available: 56_CORES_X_242_RAM_X_1_4_TB"))
	cobraCmd.Flags().StringVar(&thisCmd.Billing, "billing", "", T("Billing rate. Default is: hourly. Options are: hourly, monthly"))
	cobraCmd.Flags().IntVar(&thisCmd.VlanPrivate, "vlan-private", 0, T("The ID of the private VLAN on which you want the dedicated host placed. See: '${COMMAND_NAME} sl vlan list' for reference"))
	cobraCmd.Flags().BoolVar(&thisCmd.Test, "test", false, T("Do not actually create the dedicatedhost"))
	cobraCmd.Flags().BoolVar(&thisCmd.Force, "force", false, T("Force operation without confirmation"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *CreateCommand) Run(args []string) error {
	size := managers.HOST_DEFAULT_SIZE
	if cmd.Size != "" {
		size = cmd.Size
	}
	hostname := cmd.Hostname
	if hostname == "" {
		return slErr.NewMissingInputError("--hostname")
	}
	domain := cmd.Domain
	if domain == "" {
		return slErr.NewMissingInputError("--domain")
	}
	datacenter := cmd.Datacenter
	if datacenter == "" {
		return slErr.NewMissingInputError("--datacenter")
	}
	billing := "hourly"
	if cmd.Billing != "" {
		billing = cmd.Billing
		if billing != "hourly" && billing != "monthly" {
			return slErr.NewInvalidUsageError(T("[--billing] has to be either hourly or monthly."))
		}
	}
	vlanId := cmd.VlanPrivate
	if vlanId == 0 {
		return slErr.NewMissingInputError("--vlan-private")
	}

	outputFormat := cmd.GetOutputFlag()

	vlan, err := cmd.NetworkManager.GetVlan(vlanId, "mask[id,primaryRouter[id,datacenter]]")
	if err != nil {
		return slErr.NewAPIError(T("Failed to get vlan {{.VlanId}}.", map[string]interface{}{"VlanId": vlanId}), err.Error(), 2)
	}
	if !cmd.Force {
		confirm, err := cmd.UI.Confirm(T("This action will incur charges on your account. Continue?"))
		if err != nil {
			return err
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
		return slErr.NewAPIError(T("Failed to generate the order template."), err.Error(), 2)
	}

	if cmd.Test {
		orderReceipt, err := cmd.DedicatedHostManager.VerifyInstanceCreation(orderTemplate)
		if err != nil {
			return slErr.NewAPIError(T("Failed to verify virtual server creation."), err.Error(), 2)
		}
		if outputFormat == "JSON" {
			return utils.PrintPrettyJSON(cmd.UI, orderReceipt)
		}
		cmd.UI.Ok()
		cmd.UI.Print(T("The order is correct."))
	} else {
		orderReceipt, err := cmd.DedicatedHostManager.OrderInstance(orderTemplate)
		if err != nil {
			return slErr.NewAPIError(T("Failed to Order the dedicatedhost."), err.Error(), 2)
		}
		if outputFormat == "JSON" {
			return utils.PrintPrettyJSON(cmd.UI, orderReceipt)
		}
		cmd.UI.Ok()
		cmd.UI.Print(T("The order {{.OrderID}} was placed.", map[string]interface{}{"OrderID": *orderReceipt.OrderId}))
	}

	return nil
}
