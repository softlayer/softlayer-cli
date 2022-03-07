package reports

import (
	"fmt"
	"github.com/softlayer/softlayer-go/session"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type DCClosuresCommand struct {
	UI      terminal.UI
	Session *session.Session
}

func NewDCClosuresCommand(ui terminal.UI, session *session.Session) (cmd *DCClosuresCommand) {
	return &DCClosuresCommand{
		UI:      ui,
		Session: session,
	}
}

func DCClosuresMetaData() cli.Command {
	return cli.Command{
		Category:    "reports",
		Name:        "datacenter-closures",
		Description: T("Reports which resources are still active in Datacenters that are scheduled to be closed."),
		Usage:       T(`${COMMAND_NAME} sl reports datacenter-closures`),
		Flags: []cli.Flag{
			metadata.OutputFlag(),
		},
	}
}

func (cmd *DCClosuresCommand) Run(c *cli.Context) error {
	fmt.Printf("HELLO WORLD")

	return nil
}
