package plugin

import (
	"fmt"

	"os"
	"reflect"
	"strings"

	trace "github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/trace"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/urfave/cli"

	"github.com/softlayer/softlayer-go/session"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/client"
	slError "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/version"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/account"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/autoscale"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/block"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/callapi"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/cdn"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/dedicatedhost"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/dns"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/email"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/eventlog"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/file"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/firewall"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/globalip"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/hardware"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/image"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/ipsec"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/licenses"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/loadbal"

	//	commandMetadata "github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/nas"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/objectstorage"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/order"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/placementgroup"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/reports"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/security"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/securitygroup"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/subnet"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/tags"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/ticket"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/user"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/virtual"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/vlan"
)

var USEAGE_TEMPLATE = `${COMMAND_NAME} {{if .HasParent}}{{.Parent.CommandPath}} {{.Use}}{{else}}{{.Use}}{{end}}` +
	`{{if .HasAvailableFlags}} [` + T("OPTIONS") + `] {{end}}
{{.Long}}`

func (sl *SoftlayerPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name:       version.PLUGIN_SOFTLAYER,
		Namespaces: Namespaces(),
		// TODO change this to convert cobra commands to pluginCommands... maybe see if another plugin does this already???
		Commands: cobraToCLIMeta(getTopCobraCommand(sl.ui, sl.session), metadata.NS_SL_NAME),
	}
}

type SoftlayerPlugin struct {
	ui      terminal.UI
	session *session.Session
}

func (sl *SoftlayerPlugin) Run(context plugin.PluginContext, args []string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	trace.Logger = trace.NewLogger(context.Trace())
	terminal.UserAskedForColors = context.ColorEnabled()
	terminal.InitColorSupport()
	sl.ui = terminal.NewStdUI()
	sl.session, _ = client.NewSoftlayerClientSessionFromConfig(context)
	// initCustomizedHelp(context)

	cobraCommand := getTopCobraCommand(sl.ui, sl.session)
	// cobraCommand.SetHelpTemplate(COMMAND_HELP_TEMPLATE)
	// cobraCommand.SetUsageTemplate(USEAGE_TEMPLATE)

	// When the command comes in from the ibmcloud-cli it has `sl` in the Namespace, which we need to remove
	args = append(strings.Split(context.CommandNamespace(), " "), args...)
	if args[0] == "sl" || args[0] == "" {
		args = args[1:]
	}
	// Gives Cobra the args we were given
	cobraCommand.SetArgs(args)
	// fmt.Printf("ARgs: %v\n", args)
	cobraErr := cobraCommand.Execute()
	if cobraErr != nil {
		fmt.Printf("Cobra Error:\n %v", cobraErr)
	} else {
		return
	}

}

func GetCommandAction(pluginContext plugin.PluginContext, ui terminal.UI) func(cliContext *cli.Context) error {
	return func(cliContext *cli.Context) error {
		command := cliContext.Command

		session, err := client.NewSoftlayerClientSessionFromConfig(pluginContext)
		if err != nil {
			return slError.Error_Not_Login(pluginContext)
		}
		actionMaps := GetCommandAcionBindings(pluginContext, ui, session)
		return actionMaps[command.Category+"-"+command.Name](cliContext)
	}
}

func GetPluginCommands(cliCommands []cli.Command) []plugin.Command {
	var pluginCommands []plugin.Command
	for _, cliCmd := range cliCommands {
		if len(cliCmd.Subcommands) > 0 {
			for _, subCmd := range cliCmd.Subcommands {
				subPluginCmd := plugin.Command{
					Namespace:   metadata.SoftlayerNamespace().Name + " " + subCmd.Category,
					Name:        subCmd.Name,
					Description: subCmd.Description,
					Usage:       subCmd.Usage,
					Flags:       convertToPluginFlags(subCmd.Flags),
				}
				pluginCommands = append(pluginCommands, subPluginCmd)
			}
		} else {
			pluginCommand := plugin.Command{
				Namespace:   metadata.SoftlayerNamespace().Name,
				Name:        cliCmd.Name,
				Description: cliCmd.Description,
				Usage:       cliCmd.Usage,
				Flags:       convertToPluginFlags(cliCmd.Flags),
			}
			pluginCommands = append(pluginCommands, pluginCommand)
		}
	}
	return pluginCommands
}

func convertToPluginFlags(flags []cli.Flag) []plugin.Flag {
	var ret []plugin.Flag
	for _, f := range flags {
		ret = append(ret, plugin.Flag{
			Name:        reflect.ValueOf(f).FieldByName("Name").String(),
			Description: reflect.ValueOf(f).FieldByName("Usage").String(),
			HasValue:    reflect.TypeOf(f).String() != "cli.BoolFlag",
			Hidden:      reflect.ValueOf(f).FieldByName("Hidden").Bool(),
		})
	}
	return ret
}

func Namespaces() []plugin.Namespace {
	return []plugin.Namespace{
		metadata.SoftlayerNamespace(),
		autoscale.AutoScaleNamespace(),
		block.BlockNamespace(),
		file.FileNamespace(),
		cdn.CdnNamespace(),
		dns.DnsNamespace(),
		eventlog.EventLogNamespace(),
		firewall.FirewallNamespace(),
		email.EmailNamespace(),
		globalip.GlobalIpNamespace(),
		hardware.HardwareNamespace(),
		image.ImageNamespace(),
		ipsec.IpsecNamespace(),
		licenses.LicensesNamespace(),
		loadbal.LoadbalNamespace(),
		nas.NasNetworkStorageNamespace(),
		security.SecurityNamespace(),
		securitygroup.SecurityGroupNamespace(),
		subnet.SubnetNamespace(),
		ticket.TicketNamespace(),
		placementgroup.PlacementGroupNamespace(),
		objectstorage.ObjectStorageNamespace(),
		order.OrderNamespace(),
		vlan.VlanNamespace(),
		tags.TagsNamespace(),
		user.UserNamespace(),
		dedicatedhost.DedicatedhostNamespace(),
		virtual.VSNamespace(),
		account.AccountNamespace(),
		reports.ReportsNamespace(),
	}
}

/*
func getCLITopCommands() []cli.Command {
	return []cli.Command{
		autoscale.AutoScaleMetaData(),
		block.BlockMetaData(),
		file.FileMetaData(),
		cdn.CdnMetaData(),
		dns.DnsMetaData(),
		eventlog.EventLogMetaData(),
		firewall.FirewallMetaData(),
		email.EmailMetaData(),
		globalip.GlobalIpMetaData(),
		hardware.HardwareMetaData(),
		image.ImageMetaData(),
		ipsec.IpsecMetaData(),
		licenses.LicensesMetaData(),
		loadbal.LoadbalMetaData(),
		commandMetadata.MetadataMetadata(),
		nas.NasNetworkStorageMetaData(),
		security.SecurityMetaData(),
		securitygroup.SecurityGroupMetaData(),
		subnet.SubnetMetaData(),
		ticket.TicketMetaData(),
		vlan.VlanMetaData(),
		placementgroup.PlacementGroupMetaData(),
		objectstorage.ObjectStorageMetaData(),
		order.OrderMetaData(),
		user.UserMetaData(),
		// callapi.CallAPIMetadata(),
		tags.TagsMetaData(),
		dedicatedhost.DedicatedhostMetaData(),
		virtual.VSMetaData(),
		account.AccountMetaData(),
		reports.ReportsMetaData(),
	}
}
*/
func cobraFlagToPlugin(flagSet *pflag.FlagSet) []plugin.Flag {
	var pluginFlags []plugin.Flag
	flagSet.VisitAll(func(pflag *pflag.Flag) {
		thisFlag := plugin.Flag{
			Name:        pflag.Name,
			Description: pflag.Usage,
			HasValue:    false,
			Hidden:      false,
		}
		pluginFlags = append(pluginFlags, thisFlag)
	})
	// TODO, see if its possible to have global values added like VisitAll?
	// outputFlag := plugin.Flag{
	// 	Name: "output",
	// 	Description: "--output=JSON for json output.",
	// 	HasValue: false,
	// 	Hidden: false,
	// }
	// pluginFlags = append(pluginFlags, outputFlag)
	return pluginFlags
}

func cobraToCLIMeta(topCommand *cobra.Command, namespace string) []plugin.Command {
	var pluginCommands []plugin.Command
	// Custom Usage to ibmcloud CLI prints out a nice messages for us
	topCommand.SetUsageTemplate(USEAGE_TEMPLATE)
	for _, cliCmd := range topCommand.Commands() {
		if len(cliCmd.Commands()) > 0 {
			pluginCommands = append(pluginCommands, cobraToCLIMeta(cliCmd, namespace+" "+cliCmd.Use)...)
		} else {
			thisCmd := plugin.Command{
				Namespace:   namespace,
				Name:        cliCmd.Name(),
				Description: cliCmd.Short,
				Usage:       cliCmd.UsageString(),
				Flags:       cobraFlagToPlugin(cliCmd.Flags()),
			}
			pluginCommands = append(pluginCommands, thisCmd)
		}
	}

	// for _, cmd := range pluginCommands {
	// 	fmt.Printf("%v %v\n", cmd.Namespace, cmd.Name)
	// }
	return pluginCommands
}

func getTopCobraCommand(ui terminal.UI, session *session.Session) *cobra.Command {

	slCommand := metadata.NewSoftlayerCommand(ui, session)
	cobraCmd := &cobra.Command{
		Use:   "sl",
		Short: T("Manage Classic infrastructure services"),
		Long:  T("Manage Classic infrastructure services"),
		RunE:  nil,
	}

	// Persistent Flags
	cobraCmd.PersistentFlags().Var(slCommand.OutputFlag, "output", "--output=JSON for json output.")
	// Commands
	cobraCmd.AddCommand(callapi.NewCallAPICommand(slCommand))
	cobraCmd.AddCommand(account.SetupCobraCommands(slCommand))
	cobraCmd.AddCommand(email.SetupCobraCommands(slCommand))
	cobraCmd.AddCommand(reports.SetupCobraCommands(slCommand))
	cobraCmd.AddCommand(eventlog.SetupCobraCommands(slCommand))
	cobraCmd.AddCommand(nas.SetupCobraCommands(slCommand))
	cobraCmd.AddCommand(placementgroup.SetupCobraCommands(slCommand))
	cobraCmd.AddCommand(tags.SetupCobraCommands(slCommand))
	cobraCmd.AddCommand(block.SetupCobraCommands(slCommand))
	cobraCmd.AddCommand(file.SetupCobraCommands(slCommand))
	cobraCmd.AddCommand(licenses.SetupCobraCommands(slCommand))
	cobraCmd.AddCommand(dedicatedhost.SetupCobraCommands(slCommand))

	return cobraCmd
}
