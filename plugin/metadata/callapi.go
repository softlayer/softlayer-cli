package metadata

import (
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
)

const (
	CMD_CALLAPI_NAME = "call-api"
)

func CallAPIMetadata() cli.Command {
	return cli.Command{
		Category:    NS_SL_NAME,
		Name:        CMD_CALLAPI_NAME,
		Description: T("Call arbitrary API endpoints"),
		Usage: T(`${COMMAND_NAME} sl call-api SERVICE METHOD [OPTIONS]

EXAMPLE: 
	${COMMAND_NAME} sl call-api SoftLayer_Network_Storage editObject --init 57328245 --parameters '[{"notes":"Testing."}]'
	This command edit a volume notes.
	
	${COMMAND_NAME} sl call-api SoftLayer_User_Customer getObject --init 7051629 --mask "id,firstName,lastName"
	This command show a user detail.
	
	${COMMAND_NAME} sl call-api SoftLayer_Account getVirtualGuests --filter '{"virtualGuests":{"hostname":{"operation":"cli-test"}}}'
	This command list virtual guests.`),
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "init",
				Usage: T("Init parameter"),
			},
			cli.StringFlag{
				Name:  "mask",
				Usage: T("Object mask: use to limit fields returned"),
			},
			cli.StringFlag{
				Name:  "parameters",
				Usage: T("Append parameters to web call"),
			},
			cli.IntFlag{
				Name:  "limit",
				Usage: T("Result limit"),
			},
			cli.IntFlag{
				Name:  "offset",
				Usage: T("Result offset"),
			},
			cli.StringFlag{
				Name:  "filter",
				Usage: T("Object filters"),
			},
		},
	}
}
