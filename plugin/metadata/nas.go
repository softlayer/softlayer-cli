package metadata

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
)

var (
	NS_NAS_NAME = "nas"
	//sl-nas
	CMD_NAS_CRED = "credentials"
	CMD_NAS_LIST = "list"
)

func NasNamespace() plugin.Namespace {
	return plugin.Namespace{
		Name:        NS_NAS_NAME,
		Description: T("Classic infrastructure Network Attached Storage"),
	}
}
