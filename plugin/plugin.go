package plugin

import (
	"os"
	"reflect"
	"strings"

	"github.com/urfave/cli"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/trace"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"

	"github.ibm.com/cgallo/softlayer-cli/version"
)

// plugin name
const PLUGIN_NAME = "slcli"

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

type SLPlugin struct {
	ui terminal.UI
}

func (p *SLPlugin) GetMetadata() plugin.PluginMetadata {

	metadata := plugin.PluginMetadata{

		Name: PLUGIN_NAME,

		Version: plugin.VersionType{
			Major: version.PLUGIN_MAJOR_VERSION,
			Minor: version.PLUGIN_MINOR_VERSION,
			Build: version.PLUGIN_BUILD_VERSION,
		},

		MinCliVersion: plugin.VersionType{
			Major: 0,
			Minor: 5,
			Build: 0,
		},

		PrivateEndpointSupported: false,
	}

	metadata.Commands = []plugin.Command{
		{
			Name:        "hello",
			Alias:       "hi",
			Description: "This is just a SLCLI test",
			Usage:       "ibmcloud slcli",
		},
	}

	// ADD THIS BACK FOR REAL COMMANDS
	// for _, cmd := range getCommands() {
	// 	cmdMeta := cmd.GetMetadata()
	// 	metadata.Commands = append(metadata.Commands, plugin.Command{
	// 		Namespace:   cmdMeta.Namespace,
	// 		Name:        cmdMeta.Name,
	// 		Description: cmdMeta.Description,
	// 		Usage:       commandUsage(cmdMeta, nil),
	// 		Flags:       convertToPluginFlags(cmdMeta.Flags),
	// 	})
	// }
	return metadata
}

func (p *SLPlugin) Run(context plugin.PluginContext, args []string) {

	trace.Logger = trace.NewLogger(context.Trace())

	terminal.UserAskedForColors = context.ColorEnabled()
	terminal.InitColorSupport()
	p.ui = terminal.NewStdUI()

	cli.CommandHelpTemplate = COMMAND_HELP_TEMPLATE

	app := cli.NewApp()
	app.Name = "SLCLI"
	app.Version = version.PLUGIN_VERSION

	fmt.Println("Hi, this is my first plugin for IBM Cloud CLI.")
}

func commandUsage(meta command.CommandMetadata, context plugin.PluginContext) string {
	if context == nil {
		// TODO: tricky! sdk should provide a way to replace cli binary name
		return strings.Replace(meta.Usage, "${BINARY_NAME}", "${COMMAND_NAME} "+meta.Namespace, -1)
	}
	return strings.Replace(meta.Usage, "${BINARY_NAME}", context.CLIName()+" "+meta.Namespace, -1)
}

// func getCommands() []command.Command {
// 	return []command.Command{
// 		// new(commands.Search),
// 		new(commands.Entry),
// 		new(commands.EntryCreate),
// 		new(commands.EntryUpdate),
// 		new(commands.EntryVisibility),
// 		new(commands.EntryVisibilitySet),
// 		new(commands.EntryDelete),
// 		new(commands.Marketplace),
// 		new(commands.Service),
// 		new(commands.GetLocations),
// 		new(commands.GetPricing),
// 		new(commands.EntryCopy),
// 		new(commands.Blacklist),
// 	}
// }

func convertToPluginFlags(flags []cli.Flag) []plugin.Flag {
	var ret []plugin.Flag
	for _, f := range flags {
		ret = append(ret, plugin.Flag{
			Name:        f.GetName(),
			Description: reflect.ValueOf(f).FieldByName("Usage").String(),
		})
	}
	return ret
}
