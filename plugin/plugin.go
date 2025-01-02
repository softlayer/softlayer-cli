package plugin

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"
	"bytes"
	"text/template"

	trace "github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/trace"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/softlayer/softlayer-go/session"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/client"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/account"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/bandwidth"
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

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/meta"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/nas"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/objectstorage"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/order"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/placementgroup"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/reports"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/search"
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
`{{if .HasLocalFlags}} [` + T("OPTIONS") + `] {{RequiredFlags .LocalFlags}} {{end}}

{{.Long}}`

// https://github.ibm.com/ibmcloud-cli/bluemix-cli/blob/master/bluemix/cli/help.go#L68
// Copied/pasted because I don't want to import the whole bluemix/cli lib just for this
var BX_TEMPLATE = `{{"NAME:" | T | HeaderColor}}
  {{.Name}}{{with .Aliases}}{{range .}}, {{.}}{{end}}{{end}} - {{.Short}}

{{"USAGE:" | T | HeaderColor}}
  {{UsageCommandString . }}{{if .HasAvailableSubCommands}}{{$cmds := .Commands}}{{if eq (len .Groups) 0}}
{{"Available Commands:" | HeaderColor}}{{range $cmds}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{else}}{{range $group := .Groups}}
{{.Title}}{{range $cmds}}{{if (and (eq .GroupID $group.ID) (or .IsAvailableCommand (eq .Name "help")))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if not .AllChildCommandsHaveGroup}}
Additional Commands:{{range $cmds}}{{if (and (eq .GroupID "") (or .IsAvailableCommand (eq .Name "help")))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}
{{"OPTIONS:" | T | HeaderColor}}
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{if .HasAvailableInheritedFlags}}
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{end}}
`

// Adds HeaderColor to string in a template
func HeaderColor(text string) string {
	return terminal.HeaderColor(text)
}

// Cobra Default Template is: https://github.com/spf13/cobra/blob/v1.8.1/command.go#L546
// Used to mark flags that are required in the Usage String
func RequiredFlags(flags *pflag.FlagSet) string {
	requiredFlags := ""
	flags.VisitAll(func(pflag *pflag.Flag) {
		flagName := pflag.Name
		if pflag.Shorthand != "" {
			flagName = pflag.Shorthand + "," + pflag.Name
		}

		// Check if this flag is Required.
		// Copied logic from https://github.com/spf13/cobra/blob/v1.8.1/command.go#L1149
		// There is also an annotation for mutually exclusive we might want to look into.
		requiredAnnotation, found := pflag.Annotations[cobra.BashCompOneRequiredFlag]
		if found && requiredAnnotation[0] == "true" {
			requiredFlags = fmt.Sprintf("%s--%s <%s> ", requiredFlags, flagName, strings.ToUpper(pflag.Value.Type()))
		} 
	})
	return requiredFlags
}

// Since we overwrite the cobra Usage template, we need to build it manually.
func UsageCommandString(cmd *cobra.Command) string {
	var buf bytes.Buffer
	var templateFuncs = template.FuncMap{
		"RequiredFlags": RequiredFlags,
	}
	usage := template.New("usage")
	usage.Funcs(templateFuncs)
	template.Must(usage.Parse(USEAGE_TEMPLATE))
	err := usage.Execute(&buf, cmd)
	if err != nil {
		return err.Error()
	}
	return buf.String()
}

func (sl *SoftlayerPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name:       metadata.NS_SL_NAME,
		Namespaces: Namespaces(),
		// TODO change this to convert cobra commands to pluginCommands... maybe see if another plugin does this already???
		Commands:      cobraToCLIMeta(GetTopCobraCommand(sl.ui, sl.session), metadata.NS_SL_NAME),
		Version:       metadata.GetVersion(),
		SDKVersion:    metadata.GetSDKVersion(),
		MinCliVersion: metadata.GetMinCLI(),
	}
}

type SoftlayerPlugin struct {
	ui      terminal.UI
	session *session.Session
}

func (sl *SoftlayerPlugin) Run(context plugin.PluginContext, args []string) {

	trace.Logger = trace.NewLogger(context.Trace())
	terminal.UserAskedForColors = context.ColorEnabled()
	terminal.InitColorSupport()
	sl.ui = terminal.NewStdUI()
	sl.session, _ = client.NewSoftlayerClientSessionFromConfig(context)

	cobraCommand := GetTopCobraCommand(sl.ui, sl.session)

	// When the command comes in from the ibmcloud-cli it has `sl` in the Namespace, which we need to remove
	args = append(strings.Split(context.CommandNamespace(), " "), args...)
	if args[0] == "sl" || args[0] == "" {
		args = args[1:]
	}
	// Gives Cobra the args we were given
	cobraCommand.SetArgs(args)
	cobraErr := cobraCommand.Execute()
	if cobraErr != nil {
		cobraErrorString := fmt.Sprintf("%v", cobraErr)
		// Since we surpress the help message on errors, lets show the help message if the error is 'unknown flag'
		helpTextTriggers := []string{
			"unknown flag",
			"unknown command",
			"unknown shorthand flag",
			"required flag(s)",
			T("Incorrect Usage: "),
			T("Invalid input for")}
		for _, trigger := range helpTextTriggers {
			if strings.Contains(cobraErrorString, trigger) {
				realCommand, _, _ := cobraCommand.Find(args)
				_ = realCommand.Help()
			}
		}
		sl.ui.Failed(terminal.FailureColor(TranslateError(cobraErrorString)))
		os.Exit(1)
	}

}

// This function helps to translate errors coming from Cobra, the common ones in any case.
// If you update this, update the version in testhelpers/fake_command_runner.go as well.
// Or make this a util if we update it a lot
func TranslateError(errorMessage string) string {
	if strings.HasPrefix(errorMessage, "unknown command") {
		// If the 'command' is a number it won't have "" around it, like:
		r, _ := regexp.Compile(`unknown command "?(\w+)"? `)
		matches := r.FindStringSubmatch(errorMessage)
		subs := map[string]interface{}{"CMD": ""}
		if len(matches) >= 2 {
			subs["CMD"] = matches[1]
		} else {
			subs["CMD"], _ = strings.CutPrefix(errorMessage, "unknown command ")
		}
		
		return T("Unknown Command '{{.CMD}}'",subs)
	} else if strings.HasPrefix(errorMessage, "unknown flag") {
		r, _ := regexp.Compile(`unknown flag: (\S+)`)
		matches := r.FindStringSubmatch(errorMessage)
		subs := map[string]interface{}{"CMD": matches[1]}
		return T("Unknown Flag '{{.CMD}}'", subs)
	} else if strings.HasPrefix(errorMessage, "unknown shorthand flag") {
		r, _ := regexp.Compile(`unknown shorthand flag: '(\S+)'`)
		matches := r.FindStringSubmatch(errorMessage)
		subs := map[string]interface{}{"CMD": matches[1]}
		return T("Unknown Flag '{{.CMD}}'", subs)
	} else if strings.HasPrefix(errorMessage, "required flag(s)") {
		r, _ := regexp.Compile(`("[0-9A-Za-z\-]+")`)
		matches := r.FindAllStringSubmatch(errorMessage, -1)
		missingFlags := make([]string, len(matches))
		for i, flag := range matches {
			this_flag := strings.ReplaceAll(flag[0], `"`, "")
			subs := map[string]interface{}{"CMD": fmt.Sprintf("--%s", this_flag)}
			missingFlags[i] = T("Incorrect Usage: '{{.CMD}}' is required", subs)
		}
		return strings.Join(missingFlags, "\n")
	} else {
		return T(errorMessage)
	}
}

func Namespaces() []plugin.Namespace {
	return []plugin.Namespace{
		metadata.SoftlayerNamespace(),
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
		search.SearchNamespace(),
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
		bandwidth.BandwidthNamespace(),
	}
}

func cobraFlagToPlugin(flagSet *pflag.FlagSet) []plugin.Flag {
	var pluginFlags []plugin.Flag
	flagSet.VisitAll(func(pflag *pflag.Flag) {
		flagName := pflag.Name
		if pflag.Shorthand != "" {
			flagName = pflag.Shorthand + "," + pflag.Name
		}
		flagDesc := pflag.Usage
		if !defaultIsZeroValue(pflag) {
			flagDesc = fmt.Sprintf("%s (%s: %s)", pflag.Usage, T("Default"), pflag.DefValue)
		}
		hasValue := true
		if reflect.TypeOf(pflag.Value).String() == "*pflag.boolValue" {
			hasValue = false
		}
		// Check if this flag is Required.
		// Copied logic from https://github.com/spf13/cobra/blob/v1.8.1/command.go#L1149
		// There is also an annotation for mutually exclusive we might want to look into.
		requiredAnnotation, found := pflag.Annotations[cobra.BashCompOneRequiredFlag]
		if found && requiredAnnotation[0] == "true" {
			// Some flags have [required] hard coded in the description, so skip these
			if !strings.Contains(flagDesc, T("required")) {
				flagDesc = fmt.Sprintf("%s [%s]", flagDesc, T("required"))	
			}
			
		} 
		thisFlag := plugin.Flag{
			Name:        flagName,
			Description: flagDesc,
			HasValue:    hasValue,
			Hidden:      pflag.Hidden,
		}
		pluginFlags = append(pluginFlags, thisFlag)
	})
	return pluginFlags
}

// Copied from https://github.com/spf13/pflag/blob/master/flag.go#L538
// Because its a private function for some reason.
func defaultIsZeroValue(f *pflag.Flag) bool {
	switch f.DefValue {
	case "false":
		return true
	case "0", "0s":
		return true
	case "<nil>":
		return true
	case "":
		return true
	case "[]":
		return true
	// Used when 0 is a value users can input
	case "-1":
		return true
	default:
		return false
	}
}

func cobraToCLIMeta(topCommand *cobra.Command, namespace string) []plugin.Command {
	var pluginCommands []plugin.Command
	// Custom Usage to ibmcloud CLI prints out a nice messages for us

	for _, cliCmd := range topCommand.Commands() {
		if len(cliCmd.Commands()) > 0 {
			pluginCommands = append(pluginCommands, cobraToCLIMeta(cliCmd, namespace+" "+cliCmd.Use)...)
		} else {
			thisCmd := plugin.Command{
				Namespace:   namespace,
				Name:        cliCmd.Name(),
				Description: cliCmd.Short,
				Usage:       UsageCommandString(cliCmd),
				Flags: 		 cobraFlagToPlugin(cliCmd.Flags()),
			}
			pluginCommands = append(pluginCommands, thisCmd)
		}
	}

	return pluginCommands
}

func GetTopCobraCommand(ui terminal.UI, session *session.Session) *cobra.Command {

	slCommand := metadata.NewSoftlayerCommand(ui, session)
	helpFlag := false
	cobraCmd := &cobra.Command{
		Use:           "sl",
		Short:         T("Manage Classic infrastructure services"),
		Long:          T("Manage Classic infrastructure services"),
		RunE:          nil,
		SilenceUsage:  true, // Surpresses help text on errors
		SilenceErrors: true,
	}
	// This is to mock the `ibmcloud` usage string. Not perfect, but its close to what you can expect
	cobra.AddTemplateFunc("UsageCommandString", UsageCommandString)
	cobra.AddTemplateFunc("HeaderColor", HeaderColor)
	cobra.AddTemplateFunc("T", T)
	cobraCmd.SetUsageTemplate(BX_TEMPLATE)
	cobraCmd.SetHelpTemplate(`{{.UsageString}}`)
	versionCommand := &cobra.Command{
		Use:   "version",
		Short: T("Print the version of the sl plugin"),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(metadata.PLUGIN_VERSION)
		},
	}
	cobraCmd.AddCommand(versionCommand)

	// Persistent Flags
	cobraCmd.PersistentFlags().Var(slCommand.OutputFlag, "output", T("Specify output format, only JSON is supported now."))
	// This is needed so we can translate the help text
	cobraCmd.PersistentFlags().BoolVarP(&helpFlag, "help", "h", false, T("Usage information."))

	// Commands
	cobraCmd.AddCommand(callapi.NewCallAPICommand(slCommand).Command) // single command
	cobraCmd.AddCommand(account.SetupCobraCommands(slCommand))
	cobraCmd.AddCommand(bandwidth.SetupCobraCommands(slCommand))
	cobraCmd.AddCommand(email.SetupCobraCommands(slCommand))
	cobraCmd.AddCommand(image.SetupCobraCommands(slCommand))
	cobraCmd.AddCommand(hardware.SetupCobraCommands(slCommand))
	cobraCmd.AddCommand(ipsec.SetupCobraCommands(slCommand))
	cobraCmd.AddCommand(reports.SetupCobraCommands(slCommand))
	cobraCmd.AddCommand(eventlog.SetupCobraCommands(slCommand))
	cobraCmd.AddCommand(user.SetupCobraCommands(slCommand))
	cobraCmd.AddCommand(nas.SetupCobraCommands(slCommand))
	cobraCmd.AddCommand(cdn.SetupCobraCommands(slCommand))
	cobraCmd.AddCommand(dns.SetupCobraCommands(slCommand))
	cobraCmd.AddCommand(order.SetupCobraCommands(slCommand))
	cobraCmd.AddCommand(search.SetupCobraCommands(slCommand))
	cobraCmd.AddCommand(security.SetupCobraCommands(slCommand))
	cobraCmd.AddCommand(ticket.SetupCobraCommands(slCommand))
	cobraCmd.AddCommand(placementgroup.SetupCobraCommands(slCommand))
	cobraCmd.AddCommand(securitygroup.SetupCobraCommands(slCommand))
	cobraCmd.AddCommand(tags.SetupCobraCommands(slCommand))
	cobraCmd.AddCommand(block.SetupCobraCommands(slCommand))
	cobraCmd.AddCommand(loadbal.SetupCobraCommands(slCommand))
	cobraCmd.AddCommand(file.SetupCobraCommands(slCommand))
	cobraCmd.AddCommand(licenses.SetupCobraCommands(slCommand))
	cobraCmd.AddCommand(firewall.SetupCobraCommands(slCommand))
	cobraCmd.AddCommand(dedicatedhost.SetupCobraCommands(slCommand))
	cobraCmd.AddCommand(globalip.SetupCobraCommands(slCommand))
	cobraCmd.AddCommand(objectstorage.SetupCobraCommands(slCommand))
	cobraCmd.AddCommand(vlan.SetupCobraCommands(slCommand))
	cobraCmd.AddCommand(virtual.SetupCobraCommands(slCommand))
	cobraCmd.AddCommand(subnet.SetupCobraCommands(slCommand))
	cobraCmd.AddCommand(meta.NewMetaCommand(slCommand).Command) // single use command.

	return cobraCmd
}
