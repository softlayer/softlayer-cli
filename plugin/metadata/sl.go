package metadata

import (
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
)

var (
	LIMIT          = 50
	NS_SL_NAME     = "sl"
	OutputFlagName = "output"
)

const OutputJSON = "JSON"

func SoftlayerNamespace() plugin.Namespace {
	return plugin.Namespace{
		Name:        NS_SL_NAME,
		Description: T("Manage Classic infrastructure services"),
	}
}

func ForceFlag() cli.BoolFlag {
	return cli.BoolFlag{
		Name:  "f,force",
		Usage: T("Force operation without confirmation"),
	}
}

func OutputFlag() cli.StringFlag {
	return cli.StringFlag{
		Name:  "output",
		Usage: T("Specify output format, only JSON is supported now."),
	}
}

var SupportedOutputFormat = []string{
	OutputJSON,
	//define supported output format here in UPPER case...
}

func CheckOutputFormat(context *cli.Context, ui terminal.UI) (string, error) {
	if context.IsSet(OutputFlagName) {
		for _, r := range SupportedOutputFormat {
			if r == strings.ToUpper(context.String(OutputFlagName)) {
				return r, nil
			}
		}
		return "", errors.NewInvalidUsageError(T("Invalid output format, only JSON is supported now."))
	}
	return "", nil
}

// QuietFlag is the general `-q, --quiet` flag definition
func QuietFlag() cli.BoolFlag {
	return cli.BoolFlag{
		Name:  "q, quiet",
		Usage: T("Suppress verbose output"),
	}
}