package callapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type CallAPICommand struct {
	UI             terminal.UI
	CallAPIManager managers.CallAPIManager
}

func NewCallAPICommand(ui terminal.UI, callAPIManager managers.CallAPIManager) (cmd *CallAPICommand) {
	return &CallAPICommand{
		UI:             ui,
		CallAPIManager: callAPIManager,
	}
}

func (cmd *CallAPICommand) Run(c *cli.Context) error {
	if c.NArg() != 2 {
		return errors.NewInvalidUsageError(T("This command requires two arguments."))
	}

	args := c.Args()
	var err error
	var output []byte

	parameters := ""

	var options sl.Options

	if c.IsSet("init") {
		initparam := c.Int("init")
		options.Id = &initparam
	}

	if c.IsSet("mask") {
		mask := c.String("mask")
		if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
			mask = fmt.Sprintf("mask[%s]", mask)
		}
		options.Mask = mask
	}

	if c.IsSet("limit") {
		limit := c.Int("limit")
		options.Limit = &limit
	}

	if c.IsSet("offset") {
		offset := c.Int("offset")
		options.Offset = &offset
	}

	if c.IsSet("parameters") {
		parameters = c.String("parameters")
	}

	if c.IsSet("filter") {
		options.Filter = c.String("filter")
	} else {
		// Set Filter by Default to maintein a order in the requests, or replace when send another filter in the command
		options.Filter = DefaultFilter(args[1])
	}

	output, err = cmd.CallAPIManager.CallAPI(args[0], args[1], options, parameters)
	if err != nil {
		return err
	}

	var out bytes.Buffer
	err = json.Indent(&out, output, "", "\t")

	if err != nil {
		_, err := cmd.UI.Writer().Write(output)
		if err != nil {
			return err
		}
	} else {
		_, err := out.WriteTo(cmd.UI.Writer())
		if err != nil {
			return err
		}
	}

	return nil
}

func DefaultFilter(method string) string {
	sort := "ASC"
	getWord := method[0:3]
	methodToFilter := method[3:]
	if methodToFilter == "AllObjects" {
		return fmt.Sprintf(`{"id":{"operation":"orderBy","options":[{"name":"sort","value":["%s"]}]}}`, sort)
	}
	if getWord == "get" {
		r, _ := utf8.DecodeRuneInString(methodToFilter)
		methodToFilter = string(unicode.ToLower(r)) + methodToFilter[1:]
	}
	return fmt.Sprintf(`{"%s":{"id":{"operation":"orderBy","options":[{"name":"sort","value":["%s"]}]}}}`, methodToFilter, sort)
}

func GetCommandActionBindings(context plugin.PluginContext, ui terminal.UI, session *session.Session) map[string]func(c *cli.Context) error {
	callAPIManager := managers.NewCallAPIManager(session)

	CommandActionBindings := map[string]func(c *cli.Context) error{
		"sl-call-api": func(c *cli.Context) error {
			return NewCallAPICommand(ui, callAPIManager).Run(c)
		},
	}

	return CommandActionBindings
}

func CallAPIMetadata() cli.Command {
	return cli.Command{
		Category:    "sl",
		Name:        "call-api",
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
