package virtual

import (
	"github.com/spf13/cobra"

	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type CreateHostCommand struct {
	*metadata.SoftlayerCommand
	VirtualServerManager managers.VirtualServerManager
	NetworkManager       managers.NetworkManager
	Command              *cobra.Command
	Datacenter           string
	Domain               string
	Hostname             string
	Billing              string
	VlanPrivate          int
	Size                 string
	Force                bool
}

func NewCreateHostCommand(sl *metadata.SoftlayerCommand) (cmd *CreateHostCommand) {
	thisCmd := &CreateHostCommand{
		SoftlayerCommand:     sl,
		VirtualServerManager: managers.NewVirtualServerManager(sl.Session),
		NetworkManager:       managers.NewNetworkManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "host-create",
		Short: T("Create a host for dedicated virtual servers"),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	cobraCmd.Flags().StringVarP(&thisCmd.Datacenter, "datacenter", "d", "", T("Datacenter shortname [required]"))
	cobraCmd.Flags().StringVarP(&thisCmd.Domain, "domain", "D", "", T("Domain portion of the FQDN [required]"))
	cobraCmd.Flags().StringVarP(&thisCmd.Hostname, "hostname", "H", "", T("Host portion of the FQDN [required]"))
	cobraCmd.Flags().StringVarP(&thisCmd.Billing, "billing", "b", "hourly", T("Billing rate. Default is: hourly. Options are: hourly, monthly"))
	cobraCmd.Flags().IntVarP(&thisCmd.VlanPrivate, "vlan-private", "v", 0, T("The ID of the private VLAN on which you want the dedicated host placed. See: '${COMMAND_NAME} sl vlan list' for reference"))
	cobraCmd.Flags().StringVarP(&thisCmd.Size, "size", "s", "", T("Size of the dedicated host, currently only one size is available: 56_CORES_X_242_RAM_X_1_4_TB"))
	cobraCmd.Flags().BoolVarP(&thisCmd.Force, "force", "f", false, T("Force operation without confirmation"))

	return thisCmd
}

func (cmd *CreateHostCommand) Run(args []string) error {
	size := managers.HOST_DEFAULT_SIZE
	if cmd.Size != "" {
		size = cmd.Size
	}

	hostname := cmd.Hostname
	if hostname == "" {
		return slErrors.NewMissingInputError("-H|--hostname")
	}
	domain := cmd.Domain
	if domain == "" {
		return slErrors.NewMissingInputError("-D|--domain")
	}
	datacenter := cmd.Datacenter
	if datacenter == "" {
		return slErrors.NewMissingInputError("-d|--datacenter")
	}
	billing := cmd.Billing

	if billing != "hourly" && billing != "monthly" {
		return slErrors.NewInvalidUsageError(T("[-b|--billing] has to be either hourly or monthly."))
	}

	vlanId := cmd.VlanPrivate
	if vlanId == 0 {
		return slErrors.NewMissingInputError("-v|--vlan-private")
	}

	outputFormat := cmd.GetOutputFlag()

	vlan, err := cmd.NetworkManager.GetVlan(vlanId, "mask[id,primaryRouter[id]]")
	if err != nil {
		return slErrors.NewAPIError(T("Failed to get vlan {{.VlanId}}.\n", map[string]interface{}{"VlanId": vlanId}), err.Error(), 2)
	}
	if !cmd.Force && outputFormat != "JSON" {
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
		return slErrors.NewAPIError(T("Failed to get vlan primary router ID."), "", 2)
	}
	orderReceipt, err := cmd.VirtualServerManager.CreateDedicatedHost(size, hostname, domain, datacenter, billing, *vlan.PrimaryRouter.Id)
	if err != nil {
		return slErrors.NewAPIError(T("Failed to create dedicated host.\n"), err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, orderReceipt)
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("The order {{.OrderID}} was placed.", map[string]interface{}{"OrderID": *orderReceipt.OrderId}))
	cmd.UI.Print(T("You may run '{{.CommandName}} sl vs host-list --order {{.OrderID}}' to find this dedicated host after it is ready.",
		map[string]interface{}{"OrderID": *orderReceipt.OrderId, "CommandName": "ibmcloud"}))
	return nil
}
