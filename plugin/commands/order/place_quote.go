package order

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/spf13/cobra"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type PlaceQuoteCommand struct {
	*metadata.SoftlayerCommand
	OrderManager managers.OrderManager
	Command      *cobra.Command
	Preset       string
	Name         string
	SendEmail    bool
	ComplexType  string
	Extras       string
}

func NewPlaceQuoteCommand(sl *metadata.SoftlayerCommand) (cmd *PlaceQuoteCommand) {
	thisCmd := &PlaceQuoteCommand{
		SoftlayerCommand: sl,
		OrderManager:     managers.NewOrderManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "place-quote " + T("PACKAGE_KEYNAME") + " " + T("LOCATION") + " " + T("ORDER_ITEM1 ORDER_ITEM2 ORDER_ITEM3 ORDER_ITEM4..."),
		Short: T("Place a quote"),
		Long: T(`EXAMPLE: 
    ${COMMAND_NAME} sl order place-quote CLOUD_SERVER DALLAS13 GUEST_CORES_4 RAM_16_GB REBOOT_REMOTE_CONSOLE 1_GBPS_PUBLIC_PRIVATE_NETWORK_UPLINKS BANDWIDTH_0_GB_2 1_IP_ADDRESS GUEST_DISK_100_GB_SAN OS_UBUNTU_16_04_LTS_XENIAL_XERUS_MINIMAL_64_BIT_FOR_VSI MONITORING_HOST_PING NOTIFICATION_EMAIL_AND_TICKET AUTOMATED_NOTIFICATION UNLIMITED_SSL_VPN_USERS_1_PPTP_VPN_USER_PER_ACCOUNT NESSUS_VULNERABILITY_ASSESSMENT_REPORTING --extras '{"virtualGuests": [{"hostname": "test", "domain": "softlayer.com"}]}' --complex-type SoftLayer_Container_Product_Order_Virtual_Guest --name "foobar" --send-email
    This command places a quote for a VSI with 4 CPU, 16 GB RAM, 100 GB SAN disk, Ubuntu 16.04, and 1 Gbps public & private uplink in datacenter dal13`),
		Args: metadata.MinimumNArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().StringVar(&thisCmd.Preset, "preset", "", T("The order preset (if required by the package)"))
	cobraCmd.Flags().StringVar(&thisCmd.Name, "name", "", T("A custom name to be assigned to the quote (optional)"))
	cobraCmd.Flags().BoolVar(&thisCmd.SendEmail, "send-email", false, T("The quote will be sent to the associated email address"))
	cobraCmd.Flags().StringVar(&thisCmd.ComplexType, "complex-type", "", T("The complex type of the order. The type begins with 'SoftLayer_Container_Product_Order_'"))
	cobraCmd.Flags().StringVar(&thisCmd.Extras, "extras", "", T("JSON string that denotes extra data needs to be sent with the order"))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *PlaceQuoteCommand) Run(args []string) error {
	packageKeyname := args[0]
	location := args[1]

	orderItems := []string{}
	if len(args) > 3 {
		orderItems = getOrderItems(args)
	} else {
		orderItems = strings.Split(args[2], ",")
	}

	preset := cmd.Preset
	name := cmd.Name

	var extrasStruct interface{}
	complexType := cmd.ComplexType
	if _, ok := TYPEMAP[complexType]; ok {
		extrasStruct = TYPEMAP[complexType]
	} else {
		return errors.NewInvalidUsageError(T("Incorrect complex type: {{.Type}}", map[string]interface{}{"Type": complexType}))
	}

	extras := cmd.Extras
	if strings.HasPrefix(extras, "@") {
		extrasbytes, err := ioutil.ReadFile(strings.TrimPrefix(extras, "@"))
		if err != nil {
			return errors.NewInvalidUsageError(fmt.Sprintf("%s %s: %s", T("failed reading file"), extras, err))
		}
		extras = string(extrasbytes)
	}
	if extras != "" {
		err := json.Unmarshal([]byte(extras), &extrasStruct)
		if err != nil {
			return errors.NewInvalidUsageError(fmt.Sprintf(T("Unable to unmarshal extras json: %s\n"), err.Error()))
		}
	}

	outputFormat := cmd.GetOutputFlag()

	sendEmail := cmd.SendEmail
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
