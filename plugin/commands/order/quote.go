package order

import (
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type QuoteCommand struct {
	UI           terminal.UI
	OrderManager managers.OrderManager
	ImageManager managers.ImageManager
}

func NewQuoteCommand(ui terminal.UI, orderManager managers.OrderManager, imageManager managers.ImageManager) (cmd *QuoteCommand) {
	return &QuoteCommand{
		UI:           ui,
		OrderManager: orderManager,
		ImageManager: imageManager,
	}
}

func (cmd *QuoteCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	quoteId, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return errors.NewInvalidSoftlayerIdInputError("Quote ID")
	}

	if c.IsSet("userdata") && c.IsSet("userfile") {
		return errors.NewExclusiveFlagsError("[--userdata]", "[--userfile]")
	}

	quote, err := cmd.OrderManager.GetQuote(quoteId, "")
	if err != nil {
		return cli.NewExitError(T("Failed to get Quote.\n"+err.Error()), 2)
	}

	recalculatedOrderContainer, err := cmd.OrderManager.GetRecalculatedOrderContainer(quoteId)
	if err != nil {
		return cli.NewExitError(T("Failed to get Recalculated Order Container.\n"+err.Error()), 2)
	}

	extra, err := setArguments(c, cmd, recalculatedOrderContainer)
	if err != nil {
		return err
	}

	packageObject := quote.Order.Items[0].Package
	extra.PackageId = packageObject.Id

	var table terminal.Table
	if c.IsSet("verify") {
		order, err := cmd.OrderManager.VerifyOrder(quoteId, extra)
		if err != nil {
			return cli.NewExitError(T("Failed to verify Quote.\n"+err.Error()), 2)
		}

		table = cmd.UI.Table([]string{T("keyName"), T("description"), T("cost")})
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
			return cli.NewExitError(T("Failed to order Quote.\n"+err.Error()), 2)
		}

		table = cmd.UI.Table([]string{T("Name"), T("Value")})
		table.Add("Id", utils.FormatIntPointer(order.OrderId))
		table.Add("Created", utils.FormatSLTimePointer(order.OrderDate))
		table.Add("Status", utils.FormatStringPointer(order.PlacedOrder.Status))
	}

	utils.PrintTable(cmd.UI, table, outputFormat)
	return nil
}

func setArguments(c *cli.Context, cmd *QuoteCommand, recalculatedOrderContainer datatypes.Container_Product_Order) (datatypes.Container_Product_Order, error) {

	quantity := 1
	if c.IsSet("quantity") {
		quantity = c.Int("quantity")
	}
	recalculatedOrderContainer.Quantity = &quantity

	postinstall := []string{}
	if c.IsSet("postinstall") {
		postinstall = []string{c.String("postinstall")}
		recalculatedOrderContainer.ProvisionScripts = postinstall
	}

	complexType := "SoftLayer_Container_Product_Order_Hardware_Server"
	if c.IsSet("complex-type") {
		complexType = c.String("complex-type")
	}
	recalculatedOrderContainer.ComplexType = &complexType

	servers := []datatypes.Hardware{}
	if c.IsSet("fqdn") {
		fqdns := c.StringSlice("fqdn")
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

	if c.IsSet("userdata") || c.IsSet("userfile") {
		var userData string
		if c.IsSet("userdata") {
			userData = c.String("userdata")
		}
		if c.IsSet("userfile") {
			userfile := c.String("userfile")
			content, err := ioutil.ReadFile(userfile) // #nosec
			if err != nil {
				return datatypes.Container_Product_Order{}, cli.NewExitError((T("Failed to read user data from file: {{.File}}.", map[string]interface{}{"File": userfile})), 2)
			}
			userData = string(content)
			for _, hardware := range recalculatedOrderContainer.Hardware {
				hardware.UserData = []datatypes.Hardware_Attribute{datatypes.Hardware_Attribute{Value: &userData}}
			}
		}
	}

	if c.IsSet("key") {
		keys := c.IntSlice("key")
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

	if c.IsSet("image") {
		image, err := cmd.ImageManager.GetImage(c.Int("image"))
		if err != nil {
			return datatypes.Container_Product_Order{}, cli.NewExitError(T("Failed to get Image.\n"+err.Error()), 2)
		}
		recalculatedOrderContainer.ImageTemplateGlobalIdentifier = image.GlobalIdentifier
	}

	return recalculatedOrderContainer, nil
}

func OrderQuoteMetaData() cli.Command {
	return cli.Command{
		Category:    "order",
		Name:        "quote",
		Description: T("View and Order a quote"),
		Usage: T(`${COMMAND_NAME} sl order quote IDENTIFIER [OPTIONS]

EXAMPLE: 
	${COMMAND_NAME} sl order quote 123456 --fqdn testquote.test.com --verify --quantity 1 --postinstall https://mypostinstallscript.com --userdata Myuserdata
	${COMMAND_NAME} sl order quote 123456 --fqdn testquote.test.com --key 111111 --image 222222 --complex-type SoftLayer_Container_Product_Order_Hardware_Server`),
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "verify",
				Usage: T("If specified, will only show what the quote will order, will NOT place an order [default: False]"),
			},
			cli.IntFlag{
				Name:  "quantity",
				Usage: T("The quantity of the item being ordered if different from quoted value"),
				Value: 1,
			},
			cli.StringFlag{
				Name:  "complex-type",
				Usage: T("The complex type of the order. Starts with 'SoftLayer_Container_Product_Order'.  [default: SoftLayer_Container_Product_Order_Hardware_Server]"),
			},
			cli.StringFlag{
				Name:  "userdata",
				Usage: T("User defined metadata string"),
			},
			cli.StringFlag{
				Name:  "userfile",
				Usage: T("Read userdata from file"),
			},
			cli.StringFlag{
				Name:  "postinstall",
				Usage: T("Post-install script to download"),
			},
			cli.IntSliceFlag{
				Name:  "key",
				Usage: T("SSH key Id's to add to the root user. See: 'ibmcloud sl security sshkey-list' for reference (multiple occurrence permitted)"),
			},
			cli.StringSliceFlag{
				Name:     "fqdn",
				Usage:    T("<hostname>.<domain.name.tld> formatted name to use. Specify one fqdn per server (multiple occurrence permitted)  [required]"),
				Required: true,
			},
			cli.IntFlag{
				Name:  "image",
				Usage: T("Image ID. See: 'ibmcloud sl image list' for reference"),
			},
			metadata.OutputFlag(),
		},
	}
}
