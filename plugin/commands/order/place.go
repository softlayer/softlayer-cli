package order

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/spf13/cobra"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type PlaceCommand struct {
	*metadata.SoftlayerCommand
	OrderManager managers.OrderManager
	Command      *cobra.Command
	Preset       string
	Verify       bool
	Quantity     int
	Billing      string
	ComplexType  string
	Extras       string
	ForceFlag    bool
}

func NewPlaceCommand(sl *metadata.SoftlayerCommand) (cmd *PlaceCommand) {
	thisCmd := &PlaceCommand{
		SoftlayerCommand: sl,
		OrderManager:     managers.NewOrderManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "place " + T("PACKAGE_KEYNAME") + " " + T("LOCATION") + " " + T("ORDER_ITEM1 ORDER_ITEM2 ORDER_ITEM3 ORDER_ITEM4..."),
		Short: T("Place or verify an order"),
		Long: T(`EXAMPLE: 
	${COMMAND_NAME} sl order place CLOUD_SERVER DALLAS13 GUEST_CORES_4 RAM_16_GB REBOOT_REMOTE_CONSOLE 1_GBPS_PUBLIC_PRIVATE_NETWORK_UPLINKS BANDWIDTH_0_GB_2 1_IP_ADDRESS GUEST_DISK_100_GB_SAN OS_UBUNTU_16_04_LTS_XENIAL_XERUS_MINIMAL_64_BIT_FOR_VSI MONITORING_HOST_PING NOTIFICATION_EMAIL_AND_TICKET AUTOMATED_NOTIFICATION UNLIMITED_SSL_VPN_USERS_1_PPTP_VPN_USER_PER_ACCOUNT NESSUS_VULNERABILITY_ASSESSMENT_REPORTING --billing hourly --extras '{"virtualGuests": [{"hostname": "test", "domain": "softlayer.com"}]}' --complex-type SoftLayer_Container_Product_Order_Virtual_Guest
	This command orders an hourly VSI with 4 CPU, 16 GB RAM, 100 GB SAN disk, Ubuntu 16.04, and 1 Gbps public & private uplink in dal13`),
		Args: metadata.MinimumNArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().StringVar(&thisCmd.Preset, "preset", "", T("The order preset (if required by the package)"))
	cobraCmd.Flags().BoolVar(&thisCmd.Verify, "verify", false, T("Flag denoting whether to verify the order, or not place it"))
	cobraCmd.Flags().IntVar(&thisCmd.Quantity, "quantity", 0, T("The quantity of the item being ordered. This value defaults to 1"))
	cobraCmd.Flags().StringVar(&thisCmd.Billing, "billing", "", T("Billing rate [hourly|monthly], [default: hourly]"))
	cobraCmd.Flags().StringVar(&thisCmd.ComplexType, "complex-type", "", T("The complex type of the order. The type begins with 'SoftLayer_Container_Product_Order_'"))
	cobraCmd.Flags().StringVar(&thisCmd.Extras, "extras", "", T("JSON string that denotes extra data needs to be sent with the order"))
	cobraCmd.Flags().BoolVarP(&thisCmd.ForceFlag, "force", "f", false, T("Force operation without confirmation"))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *PlaceCommand) Run(args []string) error {
	packageKeyname := args[0]
	location := args[1]

	orderItems := []string{}
	if len(args) > 3 {
		orderItems = getOrderItems(args)
	} else {
		orderItems = strings.Split(args[2], ",")
	}

	preset := cmd.Preset

	billingFlag := cmd.Billing
	billing := true
	if billingFlag != "" {
		billingFlag = strings.ToLower(billingFlag)
		if billingFlag != "hourly" && billingFlag != "monthly" {
			return errors.NewInvalidUsageError(T("--billing can only be either hourly or monthly."))
		}
		billing = (billingFlag == "hourly")
	}

	var extrasStruct interface{}
	complexType := cmd.ComplexType
	if _, ok := TYPEMAP[complexType]; ok {
		extrasStruct = TYPEMAP[complexType]
	} else {
		return errors.NewInvalidUsageError(T("Incorrect complex type: {{.Type}}", map[string]interface{}{"Type": complexType}))
	}
	quantity := cmd.Quantity

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

	if cmd.Verify {
		orderPlace, err := cmd.OrderManager.VerifyPlaceOrder(packageKeyname, location, orderItems, complexType, billing, preset, extrasStruct, quantity)
		if err != nil {
			return err
		}
		if outputFormat == "JSON" {
			return utils.PrintPrettyJSON(cmd.UI, orderPlace)
		}
		cmd.PrintOrderVerify(orderPlace, billingFlag)
	} else {
		if !cmd.ForceFlag && outputFormat != "JSON" {
			confirm, err := cmd.UI.Confirm(T("This action will incur charges on your account. Continue?"))
			if err != nil {
				return errors.NewAPIError("", err.Error(), 1)
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

func getOrderItems(args []string) []string {
	orderItems := []string{}
	for i := 2; i < len(args); i++ {
		orderItems = append(orderItems, args[i])
	}
	return orderItems
}
