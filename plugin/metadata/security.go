package metadata

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
)

var (
	NS_SECURITY_NAME  = "security"
	CMD_SECURITY_NAME = "security"

	//sl security
	CMD_SECURITY_SSHKEY_ADD_NAME       = "sshkey-add"
	CMD_SECURITY_SSHKEY_EDIT_NAME      = "sshkey-edit"
	CMD_SECURITY_SSHKEY_LIST_NAME      = "sshkey-list"
	CMD_SECURITY_SSHKEY_PRINT_NAME     = "sshkey-print"
	CMD_SECURITY_SSHKEY_REMOVE_NAME    = "sshkey-remove"
	CMD_SECURITY_SSLCERT_ADD_NAME      = "cert-add"
	CMD_SECURITY_SSLCERT_DOWNLOAD_NAME = "cert-download"
	CMD_SECURITY_SSLCERT_EDIT_NAME     = "cert-edit"
	CMD_SECURITY_SSLCERT_LIST_NAME     = "cert-list"
	CMD_SECURITY_SSLCERT_REMOVE_NAME   = "cert-remove"
)

func SecurityNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  NS_SL_NAME,
		Name:        NS_SECURITY_NAME,
		Description: T("Classic infrastructure SSH Keys and SSL Certificates"),
	}
}

func SecurityMetaData() cli.Command {
	return cli.Command{
		Category:    NS_SL_NAME,
		Name:        CMD_SECURITY_NAME,
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

func SecuritySSHKeyAddMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_SECURITY_NAME,
		Name:        CMD_SECURITY_SSHKEY_ADD_NAME,
		Description: T("Add a new SSH key"),
		Usage: T(`${COMMAND_NAME} sl security sshkey-add LABEL [OPTIONS]
	
EXAMPLE:
   ${COMMAND_NAME} sl security sshkey-add my_sshkey -f ~/.ssh/id_rsa.pub --note mykey
   This command adds an SSH key from file ~/.ssh/id_rsa.pub with a note "mykey".`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "f,in-file",
				Usage: T("The id_rsa.pub file to import for this key"),
			},
			cli.StringFlag{
				Name:  "k,key",
				Usage: T("The actual SSH key"),
			},
			cli.StringFlag{
				Name:  "note",
				Usage: T("Extra note to be associated with the key"),
			},
			OutputFlag(),
		},
	}
}

func SecuritySSHKeyEditMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_SECURITY_NAME,
		Name:        CMD_SECURITY_SSHKEY_EDIT_NAME,
		Description: T("Edit an SSH key"),
		Usage: T(`${COMMAND_NAME} sl security sshkey-edit IDENTIFIER [OPTIONS]
	
EXAMPLE:
   ${COMMAND_NAME} sl security sshkey-edit 12345678 --label IBMCloud --note testing
   This command updates the SSH key with ID 12345678 and sets label to "IBMCloud" and note to "testing".`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "label",
				Usage: T("The new label for the key"),
			},
			cli.StringFlag{
				Name:  "note",
				Usage: T("New notes for the key"),
			},
		},
	}
}

func SecuritySSHKeyListMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_SECURITY_NAME,
		Name:        CMD_SECURITY_SSHKEY_LIST_NAME,
		Description: T("List SSH keys on your account"),
		Usage: T(`${COMMAND_NAME} sl security sshkey-list [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl security sshkey-list --sortby label
   This command lists all SSH keys on current account and sorts them by label.`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "sortby",
				Usage: T("Column to sort by. Options are: id,label,fingerprint,note"),
			},
			OutputFlag(),
		},
	}
}

func SecuritySSHKeyPrintMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_SECURITY_NAME,
		Name:        CMD_SECURITY_SSHKEY_PRINT_NAME,
		Description: T("Prints out an SSH key to the screen"),
		Usage: T(`${COMMAND_NAME} sl security sshkey-print IDENTIFIER [OPTIONS]
	
EXAMPLE:
   ${COMMAND_NAME} sl security sshkey-print 12345678 -f ~/mykey.pub
   This command shows the ID, label and notes of SSH key with ID 12345678 and write the public key to file: ~/mykey.pub.`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "f,out-file",
				Usage: T("The public SSH key will be written to this file"),
			},
			OutputFlag(),
		},
	}
}

func SecuritySSHKeyRemoveMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_SECURITY_NAME,
		Name:        CMD_SECURITY_SSHKEY_REMOVE_NAME,
		Description: T("Permanently removes an SSH key"),
		Usage: T(`${COMMAND_NAME} sl security sshkey-remove IDENTIFIER [OPTIONS]
	
EXAMPLE:
   ${COMMAND_NAME} sl security sshkey-remove 12345678 -f 
   This command removes the SSH key with ID 12345678 without asking for confirmation.`),
		Flags: []cli.Flag{
			ForceFlag(),
		},
	}
}

func SecuritySSLCertAddMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_SECURITY_NAME,
		Name:        CMD_SECURITY_SSLCERT_ADD_NAME,
		Description: T("Add and upload SSL certificate details"),
		Usage: T(`${COMMAND_NAME} sl security cert-add [OPTIONS]
	
EXAMPLE:
   ${COMMAND_NAME} sl security cert-add --crt ~/ibm.com.cert --key ~/ibm.com.key 
   This command adds certificate file: ~/ibm.com.cert and private key file ~/ibm.com.key for domain ibm.com.`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "crt",
				Usage: T("Certificate file"),
			},
			cli.StringFlag{
				Name:  "csr",
				Usage: T("Certificate Signing Request file"),
			},
			cli.StringFlag{
				Name:  "icc",
				Usage: T("Intermediate Certificate file"),
			},
			cli.StringFlag{
				Name:  "key",
				Usage: T("Private Key file"),
			},
			cli.StringFlag{
				Name:  "notes",
				Usage: T("Additional notes"),
			},
			OutputFlag(),
		},
	}
}

func SecuritySSLCertDownloadMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_SECURITY_NAME,
		Name:        CMD_SECURITY_SSLCERT_DOWNLOAD_NAME,
		Description: T("Download SSL certificate and key files"),
		Usage: T(`${COMMAND_NAME} sl security cert-download IDENTIFIER [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl security cert-download 12345678
   This command downloads four files to current directory for certificate with ID 12345678. The four files are: certificate file, certificate signing request file, intermediate certificate file and private key file.`),
		Flags: []cli.Flag{
			OutputFlag(),
		},
	}
}

func SecuritySSLCertEdit() cli.Command {
	return cli.Command{
		Category:    CMD_SECURITY_NAME,
		Name:        CMD_SECURITY_SSLCERT_EDIT_NAME,
		Description: T("Edit SSL certificate"),
		Usage: T(`${COMMAND_NAME} sl security cert-edit IDENTIFIER [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl security cert-edit 12345678 --key ~/ibm.com.key 
   This command edits certificate with ID 12345678 and updates its private key with file: ~/ibm.com.key.`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "crt",
				Usage: T("Certificate file"),
			},
			cli.StringFlag{
				Name:  "csr",
				Usage: T("Certificate Signing Request file"),
			},
			cli.StringFlag{
				Name:  "icc",
				Usage: T("Intermediate Certificate file"),
			},
			cli.StringFlag{
				Name:  "key",
				Usage: T("Private Key file"),
			},
			cli.StringFlag{
				Name:  "notes",
				Usage: T("Additional notes"),
			},
		},
	}
}

func SecuritySSLCertListMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_SECURITY_NAME,
		Name:        CMD_SECURITY_SSLCERT_LIST_NAME,
		Description: T("List SSL certificates on your account"),
		Usage: T(`${COMMAND_NAME} sl security cert-list [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl security cert-list --status valid --sortby days_until_expire
   This command lists all valid certificates on current account and sort them by validity days.`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "status",
				Usage: T("Show certificates with this status, default is: all, options are: all,valid,expired"),
			},
			cli.StringFlag{
				Name:  "sortby",
				Usage: T("Column to sort by. Options are: id,common_name,days_until_expire,note"),
			},
			OutputFlag(),
		},
	}
}

func SecuritySSLCertRemove() cli.Command {
	return cli.Command{
		Category:    CMD_SECURITY_NAME,
		Name:        CMD_SECURITY_SSLCERT_REMOVE_NAME,
		Description: T("Remove SSL certificate"),
		Usage: T(`${COMMAND_NAME} sl security cert-remove IDENTIFIER [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl security cert-remove 12345678 
   This command removes certificate with ID 12345678.`),
		Flags: []cli.Flag{
			ForceFlag(),
		},
	}
}
