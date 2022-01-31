package account

import (

	"github.com/softlayer/softlayer-go/session"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"

)

type BandwidthPoolsCommand struct {
	UI             terminal.UI
	Session  	   *session.Session
}

func NewBandwidthPoolsCommand(ui terminal.UI, session *session.Session) (cmd *BandwidthPoolsCommand) {
	return &BandwidthPoolsCommand{
		UI:             ui,
		Session: session,
	}
}

func BandwidthPoolsMetaData() cli.Command {
	return cli.Command{
		Category: "account",
		Name:	  "bandwidth-pools",
		Description: T("lists bandwidth pools"),
		Usage: T("INSERT USAGE HERE"),
		Flags: []cli.Flag{
			metadata.OutputFlag(),
		},

	}
}

func (cmd *BandwidthPoolsCommand) Run(c *cli.Context) error {
	table := cmd.UI.Table([]string{"TEST",})
	table.Print()
	return nil
}
