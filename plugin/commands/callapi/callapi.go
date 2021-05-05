package callapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

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
