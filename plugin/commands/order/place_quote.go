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

func OrderPlaceQuoteMetaData() cli.Command {
	return cli.Command{
		Category:    "order",
		Name:        "place-quote",
		Description: T("Place a quote"),
		Usage: T(`${COMMAND_NAME} sl order place-quote PACKAGE_KEYNAME LOCATION ORDER_ITEM1,ORDER_ITEM2,ORDER_ITEM3,ORDER_ITEM4... [OPTIONS]

    EXAMPLE: 
    ${COMMAND_NAME} sl order place-quote CLOUD_SERVER DALLAS13 GUEST_CORES_4,RAM_16_GB,REBOOT_REMOTE_CONSOLE,1_GBPS_PUBLIC_PRIVATE_NETWORK_UPLINKS,BANDWIDTH_0_GB_2,1_IP_ADDRESS,GUEST_DISK_100_GB_SAN,OS_UBUNTU_16_04_LTS_XENIAL_XERUS_MINIMAL_64_BIT_FOR_VSI,MONITORING_HOST_PING,NOTIFICATION_EMAIL_AND_TICKET,AUTOMATED_NOTIFICATION,UNLIMITED_SSL_VPN_USERS_1_PPTP_VPN_USER_PER_ACCOUNT,NESSUS_VULNERABILITY_ASSESSMENT_REPORTING --extras '{"virtualGuests": [{"hostname": "test", "domain": "softlayer.com"}]}' --complex-type SoftLayer_Container_Product_Order_Virtual_Guest --name "foobar" --send-email
    This command places a quote for a VSI with 4 CPU, 16 GB RAM, 100 GB SAN disk, Ubuntu 16.04, and 1 Gbps public & private uplink in datacenter dal13`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "preset",
				Usage: T("The order preset (if required by the package)"),
			},
			cli.StringFlag{
				Name:  "name",
				Usage: T("A custom name to be assigned to the quote (optional)"),
			},
			cli.BoolFlag{
				Name:  "send-email",
				Usage: T("The quote will be sent to the associated email address"),
			},
			cli.StringFlag{
				Name:  "complex-type",
				Usage: T("The complex type of the order. The type begins with 'SoftLayer_Container_Product_Order_'"),
			},
			cli.StringFlag{
				Name:  "extras",
				Usage: T("JSON string that denotes extra data needs to be sent with the order"),
			},
			metadata.OutputFlag(),
		},
	}
}
