package order

import (
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type QuoteCommand struct {
	*metadata.SoftlayerCommand
	OrderManager managers.OrderManager
	ImageManager managers.ImageManager
	Command      *cobra.Command
	Verify       bool
	Quantity     int
	ComplexType  string
	Userdata     string
	Userfile     string
	Postinstall  string
	Key          []int
	Fqdn         []string
	Image        int
}

func NewQuoteCommand(sl *metadata.SoftlayerCommand) (cmd *QuoteCommand) {
	thisCmd := &QuoteCommand{
		SoftlayerCommand: sl,
		OrderManager:     managers.NewOrderManager(sl.Session),
		ImageManager:     managers.NewImageManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "quote " + T("IDENTIFIER"),
		Short: T("View and Order a quote"),
		Long: T(`
EXAMPLE: 
	${COMMAND_NAME} sl order quote 123456 --fqdn testquote.test.com --verify --quantity 1 --postinstall https://mypostinstallscript.com --userdata Myuserdata
	${COMMAND_NAME} sl order quote 123456 --fqdn testquote.test.com --key 111111 --image 222222 --complex-type SoftLayer_Container_Product_Order_Hardware_Server`),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().BoolVar(&thisCmd.Verify, "verify", false, T("If specified, will only show what the quote will order, will NOT place an order [default: False]"))
	cobraCmd.Flags().IntVar(&thisCmd.Quantity, "quantity", 1, T("The quantity of the item being ordered if different from quoted value"))
	cobraCmd.Flags().StringVar(&thisCmd.ComplexType, "complex-type", "", T("The complex type of the order. Starts with 'SoftLayer_Container_Product_Order'.  [default: SoftLayer_Container_Product_Order_Hardware_Server]"))
	cobraCmd.Flags().StringVar(&thisCmd.Userdata, "userdata", "", T("User defined metadata string"))
	cobraCmd.Flags().StringVar(&thisCmd.Userfile, "userfile", "", T("Read userdata from file"))
	cobraCmd.Flags().StringVar(&thisCmd.Postinstall, "postinstall", "", T("Post-install script to download"))
	cobraCmd.Flags().IntSliceVar(&thisCmd.Key, "key", []int{}, T("SSH key Id's to add to the root user. See: 'ibmcloud sl security sshkey-list' for reference (multiple occurrence permitted)"))
	cobraCmd.Flags().StringSliceVar(&thisCmd.Fqdn, "fqdn", []string{}, T("<hostname>.<domain.name.tld> formatted name to use. Specify one fqdn per server (multiple occurrence permitted)  [required]"))
	cobraCmd.Flags().IntVar(&thisCmd.Image, "image", 0, T("Image ID. See: 'ibmcloud sl image list' for reference"))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *QuoteCommand) Run(args []string) error {
	outputFormat := cmd.GetOutputFlag()

	quoteId, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.NewInvalidSoftlayerIdInputError("Quote ID")
	}

	if len(cmd.Fqdn) == 0 {
		return errors.NewMissingInputError("--fqdn")
	}

	if cmd.Userdata != "" && cmd.Userfile != "" {
		return errors.NewExclusiveFlagsError("[--userdata]", "[--userfile]")
	}

	quote, err := cmd.OrderManager.GetQuote(quoteId, "")
	if err != nil {
		return errors.NewAPIError(T("Failed to get Quote."), err.Error(), 2)
	}

	recalculatedOrderContainer, err := cmd.OrderManager.GetRecalculatedOrderContainer(quoteId)
	if err != nil {
		return errors.NewAPIError(T("Failed to get Recalculated Order Container."), err.Error(), 2)
	}

	extra, err := setArguments(cmd, recalculatedOrderContainer)
	if err != nil {
		return err
	}

	packageObject := quote.Order.Items[0].Package
	extra.PackageId = packageObject.Id

	var table terminal.Table
	if cmd.Verify {
		order, err := cmd.OrderManager.VerifyOrder(quoteId, extra)
		if err != nil {
			return errors.NewAPIError(T("Failed to verify Quote.\n"), err.Error(), 2)
		}

		table = cmd.UI.Table([]string{T("KeyName"), T("Description"), T("Cost")})
		for _, price := range order.Prices {
			costKey := "recurringFee"
			if *order.UseHourlyPricing {
				costKey = "hourlyRecurringFee"
			}

			cost := "-"
			if costKey == "hourlyRecurringFee" && price.HourlyRecurringFee != nil {
				cost = utils.FormatSLFloatPointerToFloat(price.HourlyRecurringFee)
			}
			if costKey == "recurringFee" && price.RecurringFee != nil {
				cost = utils.FormatSLFloatPointerToFloat(price.RecurringFee)
			}

			table.Add(
				utils.FormatStringPointer(price.Item.KeyName),
				utils.FormatStringPointer(price.Item.Description),
				cost,
			)
		}
	} else {
		order, err := cmd.OrderManager.OrderQuote(quoteId, extra)
		if err != nil {
			return errors.NewAPIError(T("Failed to order Quote.\n"), err.Error(), 2)
		}

		table = cmd.UI.Table([]string{T("Name"), T("Value")})
		table.Add("Id", utils.FormatIntPointer(order.OrderId))
		table.Add("Created", utils.FormatSLTimePointer(order.OrderDate))
		table.Add("Status", utils.FormatStringPointer(order.PlacedOrder.Status))
	}

	utils.PrintTable(cmd.UI, table, outputFormat)
	return nil
}

func setArguments(cmd *QuoteCommand, recalculatedOrderContainer datatypes.Container_Product_Order) (datatypes.Container_Product_Order, error) {

	quantity := cmd.Quantity
	recalculatedOrderContainer.Quantity = &quantity

	if cmd.Postinstall != "" {
		postinstall := []string{cmd.Postinstall}
		recalculatedOrderContainer.ProvisionScripts = postinstall
	}

	complexType := "SoftLayer_Container_Product_Order_Hardware_Server"
	if cmd.ComplexType != "" {
		complexType = cmd.ComplexType
	}
	recalculatedOrderContainer.ComplexType = &complexType

	servers := []datatypes.Hardware{}
	if len(cmd.Fqdn) != 0 {
		fqdns := cmd.Fqdn
		for _, fqdn := range fqdns {
			fqdnStrings := strings.SplitN(fqdn, ".", 2)
			if len(fqdnStrings) < 2 {
				return datatypes.Container_Product_Order{}, errors.NewInvalidUsageError(fqdn + T(" is not following <hostname>.<domain.name.tld> --fqdn option format"))
			}
			server := datatypes.Hardware{
				Hostname: sl.String(fqdnStrings[0]),
				Domain:   sl.String(fqdnStrings[1]),
			}
			servers = append(servers, server)
		}
		recalculatedOrderContainer.Hardware = servers
	}

	if cmd.Userdata != "" || cmd.Userfile != "" {
		var userData string
		if cmd.Userdata != "" {
			userData = cmd.Userdata
		}
		if cmd.Userfile != "" {
			userfile := cmd.Userfile
			content, err := ioutil.ReadFile(userfile) // #nosec
			if err != nil {
				return datatypes.Container_Product_Order{}, errors.NewInvalidUsageError((T("Failed to read user data from file: {{.File}}.", map[string]interface{}{"File": userfile})))
			}
			userData = string(content)
			for _, hardware := range recalculatedOrderContainer.Hardware {
				hardware.UserData = []datatypes.Hardware_Attribute{datatypes.Hardware_Attribute{Value: &userData}}
			}
		}
	}

	if len(cmd.Key) != 0 {
		keys := cmd.Key
		sshkeysIntegerArray := []int{}
		for _, key := range keys {
			sshkeysIntegerArray = append(sshkeysIntegerArray, key)
		}
		sskeysDatatype := []datatypes.Container_Product_Order_SshKeys{
			datatypes.Container_Product_Order_SshKeys{
				SshKeyIds: sshkeysIntegerArray,
			},
		}
		recalculatedOrderContainer.SshKeys = sskeysDatatype
	}

	if cmd.Image != 0 {
		image, err := cmd.ImageManager.GetImage(cmd.Image)
		if err != nil {
			return datatypes.Container_Product_Order{}, errors.NewAPIError(T("Failed to get Image."), err.Error(), 2)
		}
		recalculatedOrderContainer.ImageTemplateGlobalIdentifier = image.GlobalIdentifier
	}

	return recalculatedOrderContainer, nil
}
