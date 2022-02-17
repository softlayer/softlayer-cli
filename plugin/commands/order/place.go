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

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
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

func OrderPlaceMetaData() cli.Command {
	return cli.Command{
		Category:    "order",
		Name:        "place",
		Description: T("Place or verify an order"),
		Usage: T(`${COMMAND_NAME} sl order place PACKAGE_KEYNAME LOCATION ORDER_ITEM1,ORDER_ITEM2,ORDER_ITEM3,ORDER_ITEM4... [OPTIONS]
	
	EXAMPLE: 
	${COMMAND_NAME} sl order place CLOUD_SERVER DALLAS13 GUEST_CORES_4,RAM_16_GB,REBOOT_REMOTE_CONSOLE,1_GBPS_PUBLIC_PRIVATE_NETWORK_UPLINKS,BANDWIDTH_0_GB_2,1_IP_ADDRESS,GUEST_DISK_100_GB_SAN,OS_UBUNTU_16_04_LTS_XENIAL_XERUS_MINIMAL_64_BIT_FOR_VSI,MONITORING_HOST_PING,NOTIFICATION_EMAIL_AND_TICKET,AUTOMATED_NOTIFICATION,UNLIMITED_SSL_VPN_USERS_1_PPTP_VPN_USER_PER_ACCOUNT,NESSUS_VULNERABILITY_ASSESSMENT_REPORTING --billing hourly --extras '{"virtualGuests": [{"hostname": "test", "domain": "softlayer.com"}]}' --complex-type SoftLayer_Container_Product_Order_Virtual_Guest
	This command orders an hourly VSI with 4 CPU, 16 GB RAM, 100 GB SAN disk, Ubuntu 16.04, and 1 Gbps public & private uplink in dal13`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "preset",
				Usage: T("The order preset (if required by the package)"),
			},
			cli.BoolFlag{
				Name:  "verify",
				Usage: T("Flag denoting whether to verify the order, or not place it"),
			},
			cli.IntFlag{
				Name:  "quantity",
				Usage: T("The quantity of the item being ordered. This value defaults to 1"),
			},
			cli.StringFlag{
				Name:  "billing",
				Usage: T("Billing rate [hourly|monthly], [default: hourly]"),
			},
			cli.StringFlag{
				Name:  "complex-type",
				Usage: T("The complex type of the order. The type begins with 'SoftLayer_Container_Product_Order_'"),
			},
			cli.StringFlag{
				Name:  "extras",
				Usage: T("JSON string that denotes extra data needs to be sent with the order"),
			},
			metadata.ForceFlag(),
			metadata.OutputFlag(),
		},
	}
}
