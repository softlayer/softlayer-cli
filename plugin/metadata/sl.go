package metadata

import (
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"

	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
)

var (
	LIMIT          = 50
	NS_SL_NAME     = "sl"
	OutputFlagName = "output"
)

const OutputJSON = "JSON"
const OutputCSV = "CSV"

var SupportedOutputFormat = []string{
	OutputJSON,
	OutputCSV,
	//define supported output format here in UPPER case...
}

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

type SoftlayerStorageCommand struct {
	*SoftlayerCommand
	StorageI18n map[string]interface{}
}

func NewSoftlayerStorageCommand(ui terminal.UI, session *session.Session, storageType string) *SoftlayerStorageCommand {
	return &SoftlayerStorageCommand{
		SoftlayerCommand: NewSoftlayerCommand(ui, session),
		StorageI18n:      map[string]interface{}{"storageType": storageType},
	}
}

func (slcmd *SoftlayerCommand) GetOutputFlag() string {
	return slcmd.OutputFlag.String()
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
