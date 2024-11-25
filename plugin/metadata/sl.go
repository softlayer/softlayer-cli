package metadata

import (
	"strconv"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"

	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
)

var (
	LIMIT                  = 50
	NS_SL_NAME             = "sl"
	OutputFlagName         = "output"
	PLUGIN_VERSION         = "1.5.4"
	PLUGIN_SOFTLAYER       = "sl"
	PLUGIN_SOFTLAYER_USAGE = "Classic Infrastructure"
	UsageAgentHeader       = "ibmcloud sl v" + PLUGIN_VERSION
)

const OutputJSON = "JSON"
const OutputCSV = "CSV"

var SupportedOutputFormat = []string{
	OutputJSON,
	OutputCSV,
	//define supported output format here in UPPER case...
}

// SoftLayer Base Command
type SoftlayerCommand struct {
	UI         terminal.UI
	Session    *session.Session
	OutputFlag *CobraOutputFlag
}

func NewSoftlayerCommand(ui terminal.UI, session *session.Session) *SoftlayerCommand {
	return &SoftlayerCommand{
		UI:         ui,
		Session:    session,
		OutputFlag: &CobraOutputFlag{""},
	}
}
func (slcmd *SoftlayerCommand) GetOutputFlag() string {
	return slcmd.OutputFlag.String()
}

// SoftLayer Storage Command
type SoftlayerStorageCommand struct {
	*SoftlayerCommand
	StorageI18n map[string]interface{}
	StorageType string
}

func NewSoftlayerStorageCommand(ui terminal.UI, session *session.Session, storageType string) *SoftlayerStorageCommand {
	return &SoftlayerStorageCommand{
		SoftlayerCommand: NewSoftlayerCommand(ui, session),
		StorageI18n:      map[string]interface{}{"storageType": storageType},
		StorageType:      storageType,
	}
}

func (slcmd *SoftlayerStorageCommand) GetStorageType() string {
	return slcmd.StorageType
}

func SoftlayerNamespace() plugin.Namespace {
	return plugin.Namespace{
		Name:        NS_SL_NAME,
		Description: T("Manage Classic infrastructure services"),
	}
}

// A custom flag type so we can do type checking like expected.
// Basically this just calls strings.ToUpper on --output
type CobraOutputFlag struct {
	Value string
}

func (o *CobraOutputFlag) String() string {
	return o.Value
}

func (o *CobraOutputFlag) Set(p string) error {
	p = strings.ToUpper(p)
	for _, supported := range SupportedOutputFormat {
		if p == supported {
			o.Value = p
			return nil
		}
	}
	return errors.NewInvalidUsageError(T("Invalid output format, only JSON is supported now."))
}

func (o *CobraOutputFlag) Type() string {
	return "string"
}

func GetVersion() plugin.VersionType {
	versionSplit := strings.Split(PLUGIN_VERSION, ".")
	var err error
	major, minor, revision := 0, 0, 0
	// Error checking here seems a bit much, but a mistake in the version string would cause a crash otherwise.
	if len(versionSplit) == 3 {
		if major, err = strconv.Atoi(versionSplit[0]); err != nil {
			major = 99
		}
		if minor, err = strconv.Atoi(versionSplit[1]); err != nil {
			minor = 99
		}
		if revision, err = strconv.Atoi(versionSplit[2]); err != nil {
			revision = 99
		}
	}
	return plugin.VersionType{Major: major, Minor: minor, Build: revision}

}

// Might be a way to read this from go.mod, or something?
func GetSDKVersion() plugin.VersionType {
	return plugin.VersionType{Major: 0, Minor: 9, Build: 0}
}

func GetMinCLI() plugin.VersionType {
	return plugin.VersionType{Major: 2, Minor: 12, Build: 0}
}
