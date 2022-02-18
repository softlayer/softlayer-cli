package security

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/softlayer/softlayer-go/session"
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

func GetCommandActionBindings(context plugin.PluginContext, ui terminal.UI, session *session.Session) map[string]func(c *cli.Context) error {
	securityManager := managers.NewSecurityManager(session)

	CommandActionBindings := map[string]func(c *cli.Context) error{
		"security-cert-add": func(c *cli.Context) error {
			return NewCertAddCommand(ui, securityManager).Run(c)
		},
		"security-cert-download": func(c *cli.Context) error {
			return NewCertDownloadCommand(ui, securityManager).Run(c)
		},
		"security-cert-edit": func(c *cli.Context) error {
			return NewCertEditCommand(ui, securityManager).Run(c)
		},
		"security-cert-list": func(c *cli.Context) error {
			return NewCertListCommand(ui, securityManager).Run(c)
		},
		"security-cert-remove": func(c *cli.Context) error {
			return NewCertRemoveCommand(ui, securityManager).Run(c)
		},
		"security-sshkey-add": func(c *cli.Context) error {
			return NewKeyAddCommand(ui, securityManager).Run(c)
		},
		"security-sshkey-edit": func(c *cli.Context) error {
			return NewKeyEditCommand(ui, securityManager).Run(c)
		},
		"security-sshkey-list": func(c *cli.Context) error {
			return NewKeyListCommand(ui, securityManager).Run(c)
		},
		"security-sshkey-print": func(c *cli.Context) error {
			return NewKeyPrintCommand(ui, securityManager).Run(c)
		},
		"security-sshkey-remove": func(c *cli.Context) error {
			return NewKeyRemoveCommand(ui, securityManager).Run(c)
		},
	}
	return CommandActionBindings
}

func SecurityNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  "sl",
		Name:        "security",
		Aliases:     []string{"ssl", "sshkey"},
		Description: T("Classic infrastructure SSH Keys and SSL Certificates"),
	}
}

func SecurityMetaData() cli.Command {
	return cli.Command{
		Category:    "sl",
		Name:        "security",
		Aliases:     []string{"ssl", "sshkey"},
		Description: T("Classic infrastructure SSH Keys and SSL Certificates"),
		Usage:       "${COMMAND_NAME} sl security",
		Subcommands: []cli.Command{
			SecuritySSHKeyAddMetaData(),
			SecuritySSHKeyEditMetaData(),
			SecuritySSHKeyListMetaData(),
			SecuritySSHKeyPrintMetaData(),
			SecuritySSHKeyRemoveMetaData(),
			SecuritySSLCertAddMetaData(),
			SecuritySSLCertDownloadMetaData(),
			SecuritySSLCertEdit(),
			SecuritySSLCertListMetaData(),
			SecuritySSLCertRemove(),
		},
	}
}
