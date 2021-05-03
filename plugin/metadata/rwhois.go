package metadata

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	. "github.ibm.com/cgallo/softlayer-cli/plugin/i18n"
)

var (
	NS_RWHOIS_NAME = "rwhois"
	//sl-rwhois
	CMD_WHO_EDIT = "edit"
	CMD_WHO_SHOW = "show"
)

func RwhoisNamespace() plugin.Namespace {
	return plugin.Namespace{
		Name:        NS_RWHOIS_NAME,
		Description: T("Classic infrastructure Referral Whois"),
	}
}
