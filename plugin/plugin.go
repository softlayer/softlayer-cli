package plugin

import (
	"fmt"

	"os"
	"reflect"
	"strings"

	trace "github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/trace"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/urfave/cli"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/client"
	slError "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/version"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/account"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/autoscale"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/block"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/callapi"
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
	commandMetadata "github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/metadata"
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

var (
	COMMAND_HELP_TEMPLATE = T("NAME:") + `
{{.Name}} - {{.Usage}}{{with .ShortName}}
` + T("ALIAS:") + `
   {{.}}{{end}}

` + T("USAGE:") + `
   {{.Description}}
{{with .Flags}}
` + T("OPTIONS:") + `
{{range .}}   {{.}}
{{end}}{{end}}
`
)

func (sl *SoftlayerPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name:       version.PLUGIN_SOFTLAYER,
		Namespaces: Namespaces(),
		Commands:   GetPluginCommands(getCLITopCommands()),
	}
}

type SoftlayerPlugin struct {
	ui terminal.UI
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
	// initCustomizedHelp(context)
	cli.CommandHelpTemplate = COMMAND_HELP_TEMPLATE

	app := cli.NewApp()
	app.Name = context.CLIName() + "sl "
	app.Usage = T(version.PLUGIN_SOFTLAYER_USAGE)
	app.Version = version.PLUGIN_VERSION

	for _, cmd := range getCLITopCommands() {
		cliCommand := cli.Command{
			Category:    cmd.Category,
			Name:        cmd.Name,
			Description: cmd.Description,
			Usage:       strings.Replace(cmd.Usage, "${COMMAND_NAME}", context.CLIName(), -1),
			Flags:       cmd.Flags,
		}
		if len(cmd.Subcommands) == 0 {
			action := GetCommandAction(context, sl.ui)
			if action != nil {
				cliCommand.Action = action
			}
		} else {
			for _, subCmd := range cmd.Subcommands {
				cliCommand.Subcommands = append(cliCommand.Subcommands,
					cli.Command{
						Category:    subCmd.Category,
						Name:        subCmd.Name,
						Description: subCmd.Description,
						Usage:       strings.Replace(subCmd.Usage, "${COMMAND_NAME}", context.CLIName(), -1),
						Flags:       subCmd.Flags,
						Action:      GetCommandAction(context, sl.ui),
					})
			}
		}
		app.Commands = append(app.Commands, cliCommand)
	}
	err := app.Run(append(strings.Split(context.CommandNamespace(), " "), args...))
	if err != nil {
		sl.ui.Failed(err.Error())
		os.Exit(1)
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

func getCLITopCommands() []cli.Command {
	return []cli.Command{
		autoscale.AutoScaleMetaData(),
		block.BlockMetaData(),
		file.FileMetaData(),
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
		callapi.CallAPIMetadata(),
		tags.TagsMetaData(),
		dedicatedhost.DedicatedhostMetaData(),
		virtual.VSMetaData(),
		account.AccountMetaData(),
		reports.ReportsMetaData(),
	}
}
