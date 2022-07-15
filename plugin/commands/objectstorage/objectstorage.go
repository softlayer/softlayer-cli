package objectstorage

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/softlayer/softlayer-go/session"
	"github.com/urfave/cli"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

func GetCommandActionBindings(context plugin.PluginContext, ui terminal.UI, session *session.Session) map[string]func(c *cli.Context) error {

	objectStorageManager := managers.NewObjectStorageManager(session)

	CommandActionBindings := map[string]func(c *cli.Context) error{
		"object-storage-accounts": func(c *cli.Context) error {
			return NewAccountsCommand(ui, objectStorageManager).Run(c)
		},
		"object-storage-endpoints": func(c *cli.Context) error {
			return NewEndpointsCommand(ui, objectStorageManager).Run(c)
		},
		"object-storage-credential-list": func(c *cli.Context) error {
			return NewCredentialListCommand(ui, objectStorageManager).Run(c)
		},
	}
	return CommandActionBindings
}

func ObjectStorageNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  "sl",
		Name:        "object-storage",
		Description: T("Classic infrastructure Object Storage commands"),
	}
}

func ObjectStorageMetaData() cli.Command {
	return cli.Command{
		Category:    "sl",
		Name:        "object-storage",
		Description: T("Classic infrastructure Object Storage commands"),
		Usage:       "${COMMAND_NAME} sl object-storage",
		Subcommands: []cli.Command{
			AccountsMetaData(),
			EndpointsMetaData(),
			CredentialListMetaData(),
		},
	}
}
