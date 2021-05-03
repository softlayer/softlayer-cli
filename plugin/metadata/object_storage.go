package metadata

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	. "github.ibm.com/cgallo/softlayer-cli/plugin/i18n"
)

var (
	NS_OBJECT_STORAGE_NAME = "object-storage"
	//sl-object-storage
	CMD_OBJ_ACC  = "accounts"
	CMD_OBJ_ENPT = "endpoints"
)

func ObjectStorageNamespace() plugin.Namespace {
	return plugin.Namespace{
		Name:        NS_OBJECT_STORAGE_NAME,
		Description: T("Classic infrastructure Object Storage"),
	}
}
