package order

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/urfave/cli"

	"github.ibm.com/cgallo/softlayer-cli/plugin/errors"
	. "github.ibm.com/cgallo/softlayer-cli/plugin/i18n"
	"github.ibm.com/cgallo/softlayer-cli/plugin/managers"
	"github.ibm.com/cgallo/softlayer-cli/plugin/metadata"
	"github.ibm.com/cgallo/softlayer-cli/plugin/utils"
)

type PlaceCommand struct {
	UI           terminal.UI
	OrderManager managers.OrderManager
	Context      plugin.PluginContext
}

func NewPlaceCommand(ui terminal.UI, orderManager managers.OrderManager, context plugin.PluginContext) (cmd *PlaceCommand) {
	return &PlaceCommand{
		UI:           ui,
		OrderManager: orderManager,
		Context:      context,
	}
}

func (cmd *PlaceCommand) Run(c *cli.Context) error {
	if c.NArg() != 3 {
		return errors.NewInvalidUsageError(T("This command requires three arguments."))
	}
	packageKeyname := c.Args()[0]
	location := c.Args()[1]
	orderItems := strings.Split(c.Args()[2], ",")

	preset := c.String("preset")

	billingFlag := c.String("billing")
	billing := true
	if billingFlag != "" {
		billingFlag = strings.ToLower(billingFlag)
		if billingFlag != "hourly" && billingFlag != "monthly" {
			return errors.NewInvalidUsageError(T("--billing can only be either hourly or monthly."))
		}
		billing = (billingFlag == "hourly")
	}

	var extrasStruct interface{}
	complexType := c.String("complex-type")
	if _, ok := TYPEMAP[complexType]; ok {
		extrasStruct = TYPEMAP[complexType]
	} else {
		return errors.NewInvalidUsageError(T("Incorrect complex type: {{.Type}}", map[string]interface{}{"Type": complexType}))
	}
	quantity := c.Int("quantity")

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

	if c.Bool("verify") {
		orderPlace, err := cmd.OrderManager.VerifyPlaceOrder(packageKeyname, location, orderItems, complexType, billing, preset, extrasStruct, quantity)
		if err != nil {
			return err
		}
		if outputFormat == "JSON" {
			return utils.PrintPrettyJSON(cmd.UI, orderPlace)
		}
		cmd.PrintOrderVerify(orderPlace, billingFlag)
	} else {
		if !c.IsSet("f") && outputFormat != "JSON" {
			confirm, err := cmd.UI.Confirm(T("This action will incur charges on your account. Continue?"))
			if err != nil {
				return cli.NewExitError(err.Error(), 1)
			}
			if !confirm {
				cmd.UI.Print(T("Aborted."))
				return nil
			}
		}
		orderPlace, err := cmd.OrderManager.PlaceOrder(packageKeyname, location, orderItems, complexType, billing, preset, extrasStruct, quantity)
		if err != nil {
			return err
		}
		if outputFormat == "JSON" {
			return utils.PrintPrettyJSON(cmd.UI, orderPlace)
		}
		cmd.PrintOrder(orderPlace)
	}
	return nil
}

func (cmd *PlaceCommand) PrintOrderVerify(orderPlace datatypes.Container_Product_Order, billingFlag string) {
	table := cmd.UI.Table([]string{T("keyName"), T("description"), T("cost")})

	for _, Price := range orderPlace.Prices {
		var cost datatypes.Float64
		if billingFlag == "hourly" {
			if Price.HourlyRecurringFee != nil {
				cost = *Price.HourlyRecurringFee
			}
		} else {
			if Price.RecurringFee != nil {
				cost = *Price.RecurringFee
			}
		}
		table.Add(utils.FormatStringPointer(Price.Item.KeyName),
			utils.FormatStringPointer(Price.Item.Description),
			strconv.FormatFloat(float64(cost), 'f', -1, 32))
	}
	table.Print()
}

func (cmd *PlaceCommand) PrintOrder(orderPlace datatypes.Container_Product_Order_Receipt) {

	table := cmd.UI.Table([]string{"", ""})
	table.Add("ID", utils.FormatIntPointer(orderPlace.OrderId))
	table.Add("Created", utils.FormatSLTimePointer(orderPlace.OrderDate))
	table.Add("Status", utils.FormatStringPointer(orderPlace.PlacedOrder.Status))
	table.Print()
}
