package order

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/urfave/cli"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type PlaceQuoteCommand struct {
	UI           terminal.UI
	OrderManager managers.OrderManager
	Context      plugin.PluginContext
}

func NewPlaceQuoteCommand(ui terminal.UI, orderManager managers.OrderManager, context plugin.PluginContext) (cmd *PlaceQuoteCommand) {
	return &PlaceQuoteCommand{
		UI:           ui,
		OrderManager: orderManager,
		Context:      context,
	}
}

func (cmd *PlaceQuoteCommand) Run(c *cli.Context) error {
	if c.NArg() != 3 {
		return errors.NewInvalidUsageError(T("This command requires three arguments."))
	}
	packageKeyname := c.Args()[0]
	location := c.Args()[1]
	orderItems := strings.Split(c.Args()[2], ",")

	preset := c.String("preset")
	name := c.String("name")

	var extrasStruct interface{}
	complexType := c.String("complex-type")
	if _, ok := TYPEMAP[complexType]; ok {
		extrasStruct = TYPEMAP[complexType]
	} else {
		return errors.NewInvalidUsageError(T("Incorrect complex type: {{.Type}}", map[string]interface{}{"Type": complexType}))
	}

	extras := c.String("extras")
	if strings.HasPrefix(extras, "@") {
		extrasbytes, err := ioutil.ReadFile(strings.TrimPrefix(extras, "@"))
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("%s %s: %s", T("failed reading file"), extras, err), 1)
		}
		extras = string(extrasbytes)
	}
	if extras != "" {
		err := json.Unmarshal([]byte(extras), &extrasStruct)
		if err != nil {
			return cli.NewExitError(fmt.Sprintf(T("Unable to unmarshal extras json: %s\n"), err.Error()), 1)
		}
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	sendEmail := c.Bool("send-email")
	placeQuote, err := cmd.OrderManager.PlaceQuote(packageKeyname, location, orderItems, complexType, name, preset, extrasStruct, sendEmail)
	if err != nil {
		return err
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, placeQuote)
	}

	cmd.Print(placeQuote)

	return nil
}

func (cmd *PlaceQuoteCommand) Print(placeQuote datatypes.Container_Product_Order_Receipt) {

	table := cmd.UI.Table([]string{"", ""})
	table.Add("ID", utils.FormatIntPointer(placeQuote.Quote.Id))
	table.Add("Name", utils.FormatStringPointer(placeQuote.Quote.Name))
	table.Add("Created", utils.FormatSLTimePointer(placeQuote.OrderDate))
	table.Add("Expires", utils.FormatSLTimePointer(placeQuote.Quote.ExpirationDate))
	table.Add("Status", utils.FormatStringPointer(placeQuote.Quote.Status))
	table.Print()
}
