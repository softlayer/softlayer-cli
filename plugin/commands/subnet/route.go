package subnet

import (
	"strconv"

	"github.com/spf13/cobra"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type RouteCommand struct {
	*metadata.SoftlayerCommand
	NetworkManager managers.NetworkManager
	Command        *cobra.Command
	Ip					string
	Server			string
	Vsi				string
	Vlan				string
}

func NewRouteCommand(sl *metadata.SoftlayerCommand) *RouteCommand {
	thisCmd := &RouteCommand{
		SoftlayerCommand: sl,
		NetworkManager:   managers.NewNetworkManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "route " + T("IDENTIFIER"),
		Short: T("Change how a secondary subnet is routed."),
		Long: T(`Allows you to change the route of your secondary subnets.
Subnets may be routed as either Static or Portable, and that designation is dictated by the routing destination specified.
Static subnets have an ultimate routing destination of a single IP address but may not be routed to an existing subnet’s IP address whose subnet is routed as a Static.
Portable subnets have an ultimate routing destination of a VLAN.
A subnet can be routed to any resource within the same 'routing region' as the subnet itself, usually limited to a single data center.

See Also: https://sldn.softlayer.com/reference/services/SoftLayer_Network_Subnet/route/

EXAMPLE:
	${COMMAND_NAME} sl subnet route 1234567 --ip 11.22.33.44
	${COMMAND_NAME} sl subnet route 1234567 --server myUniqueHostname<domain.com>
	${COMMAND_NAME} sl subnet route 1234567 --vsi vsiId
	${COMMAND_NAME} sl subnet route 1234567 --vlan vlanId

`),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().StringVarP(&thisCmd.Ip, "ip", "i", "",
		T("A Network_Subnet_IpAddress.id, A dotted-quad IPv4 address, or A full or compressed IPv6 address."))
	cobraCmd.Flags().StringVarP(&thisCmd.Server, "server", "s", "", 
		T(`A Hardware_Server.id or UUID value of the desired server. A value corresponding to a unique
fully-qualified domain name in the format 'hostname<domain>' where < and > are literal, e.g. myhost<mydomain.com>`))
	cobraCmd.Flags().StringVarP(&thisCmd.Vsi, "vsi", "v", "",
		T(`A Virtual_Guest.id or UUID value of the desired server. A value corresponding to a unique
fully-qualified domain name in the format 'hostname<domain>' where < and > are literal, e.g. myhost<mydomain.com>`))
	cobraCmd.Flags().StringVarP(&thisCmd.Vlan, "vlan", "l", "",
		T(`A Network_Vlan.id value of the desired VLAN or A semantic VLAN identifier of the form <data center short name>.<router>.<vlan number>,
eg. dal13.fcr01.1234 - the router name may optionally contain the 'a' or 'b' redundancy qualifier `))
	cobraCmd.MarkFlagsMutuallyExclusive("ip", "server", "vsi", "vlan")
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *RouteCommand) Run(args []string) error {

	var err error
	subnetId, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Subnet ID")
	}
	if cmd.Ip + cmd.Server + cmd.Vsi + cmd.Vlan == "" {
		return slErr.NewMissingInputError("--ip, --server, --vsi or --vlan")
	}


	if cmd.Ip != "" {
		_, err = cmd.NetworkManager.Route(subnetId, "SoftLayer_Network_Subnet_IpAddress", cmd.Ip)
	} else if cmd.Server != "" {
		_, err = cmd.NetworkManager.Route(subnetId, "SoftLayer_Hardware_Server", cmd.Server)
	} else if cmd.Vsi != "" {
		_, err = cmd.NetworkManager.Route(subnetId, "SoftLayer_Virtual_Guest", cmd.Vsi)
	}  else if cmd.Vlan != "" {
		_, err = cmd.NetworkManager.Route(subnetId, "SoftLayer_Network_Vlan", cmd.Vlan)
	}
	if err != nil {
		return err
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("The transaction to route is created, routes will be updated in one or two minutes."))
	return nil
}
