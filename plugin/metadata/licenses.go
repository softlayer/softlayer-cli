package metadata

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
)

var (
	NS_LICENSES_NAME  = "licenses"
	CMD_LICENCES_NAME = "licenses"

	CMD_LICENSES_CANCEL_NAME         = "cancel"
	CMD_LICENSES_CREATE_NAME         = "create"
	CMD_LICENSES_CREATE_OPTIONS_NAME = "create-options"
)

func LicensesNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  NS_SL_NAME,
		Name:        NS_LICENSES_NAME,
		Description: T("Classic infrastructure Licenses"),
	}
}

func LicensesMetaData() cli.Command {
	return cli.Command{
		Category:    NS_SL_NAME,
		Name:        NS_LICENSES_NAME,
		Description: T("Classic infrastructure Licenses"),
		Usage:       "${COMMAND_NAME} sl licenses",
		Subcommands: []cli.Command{
			LicensesCreateOptionsMetaData(),
		},
	}
}

func LicensesCreateOptionsMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_LICENCES_NAME,
		Name:        CMD_LICENSES_CREATE_OPTIONS_NAME,
		Description: T("Server order options for a given chassis"),
		Usage:       "${COMMAND_NAME} sl licenses create-options",
		Flags:       []cli.Flag{},
	}
}