package security

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/session"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

// Make sure to add new commands to security_test.go as well.
func GetCommandActionBindings(ui terminal.UI, session *session.Session) map[string]func(c *cli.Context) error {
	securityManager := managers.NewSecurityManager(session)

	CommandActionBindings := map[string]func(c *cli.Context) error{

		//security -  ssh keys
		NS_SECURITY_NAME + "-" + CMD_SECURITY_SSHKEY_ADD_NAME: func(c *cli.Context) error {
			return NewKeyAddCommand(ui, securityManager).Run(c)
		},
		NS_SECURITY_NAME + "-" + CMD_SECURITY_SSHKEY_EDIT_NAME: func(c *cli.Context) error {
			return NewKeyEditCommand(ui, securityManager).Run(c)
		},
		NS_SECURITY_NAME + "-" + CMD_SECURITY_SSHKEY_LIST_NAME: func(c *cli.Context) error {
			return NewKeyListCommand(ui, securityManager).Run(c)
		},
		NS_SECURITY_NAME + "-" + CMD_SECURITY_SSHKEY_PRINT_NAME: func(c *cli.Context) error {
			return NewKeyPrintCommand(ui, securityManager).Run(c)
		},
		NS_SECURITY_NAME + "-" + CMD_SECURITY_SSHKEY_REMOVE_NAME: func(c *cli.Context) error {
			return NewKeyRemoveCommand(ui, securityManager).Run(c)
		},
		// security - ssl certs
		NS_SECURITY_NAME + "-" + CMD_SECURITY_SSLCERT_ADD_NAME: func(c *cli.Context) error {
			return NewCertAddCommand(ui, securityManager).Run(c)
		},
		NS_SECURITY_NAME + "-" + CMD_SECURITY_SSLCERT_EDIT_NAME: func(c *cli.Context) error {
			return NewCertEditCommand(ui, securityManager).Run(c)
		},
		NS_SECURITY_NAME + "-" + CMD_SECURITY_SSLCERT_DOWNLOAD_NAME: func(c *cli.Context) error {
			return NewCertDownloadCommand(ui, securityManager).Run(c)
		},
		NS_SECURITY_NAME + "-" + CMD_SECURITY_SSLCERT_REMOVE_NAME: func(c *cli.Context) error {
			return NewCertRemoveCommand(ui, securityManager).Run(c)
		},
		NS_SECURITY_NAME + "-" + CMD_SECURITY_SSLCERT_LIST_NAME: func(c *cli.Context) error {
			return NewCertListCommand(ui, securityManager).Run(c)
		},
	}

	return CommandActionBindings
}
