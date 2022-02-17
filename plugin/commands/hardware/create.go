package hardware

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type CreateCommand struct {
	UI              terminal.UI
	HardwareManager managers.HardwareServerManager
	Context         plugin.PluginContext
}

func NewCreateCommand(ui terminal.UI, hardwareManager managers.HardwareServerManager, context plugin.PluginContext) (cmd *CreateCommand) {
	return &CreateCommand{
		UI:              ui,
		HardwareManager: hardwareManager,
		Context:         context,
	}
}

func (cmd *CreateCommand) Run(c *cli.Context) error {
	params := make(map[string]interface{})
	if c.IsSet("m") {
		templateFile := c.String("m")
		if _, err := os.Stat(templateFile); os.IsNotExist(err) {
			return errors.NewInvalidUsageError(T("Template file: {{.Location}} does not exist.",
				map[string]interface{}{"Location": templateFile}))
		}
		content, err := ioutil.ReadFile(templateFile) // #nosec
		if err != nil {
			return cli.NewExitError(T("Failed to read template file: {{.File}}.\n", map[string]interface{}{"File": templateFile})+err.Error(), 1)
		}
		err = json.Unmarshal(content, &params)
		if err != nil {
			return cli.NewExitError(T("Failed to unmarshal template file: {{.File}}.\n", map[string]interface{}{"File": templateFile})+err.Error(), 1)
		}
	} else {
		if !c.IsSet("s") {
			return errors.NewMissingInputError("-s|--size")
		}
		params["size"] = c.String("s")
		if !c.IsSet("H") {
			return errors.NewMissingInputError("-H|--hostname")
		}
		params["hostname"] = c.String("H")
		if !c.IsSet("D") {
			return errors.NewMissingInputError("-D|--domain")
		}
		params["domain"] = c.String("D")
		if !c.IsSet("o") {
			return errors.NewMissingInputError("-o|--os")
		}
		params["osName"] = c.String("o")
		if !c.IsSet("d") {
			return errors.NewMissingInputError("-d|--datacenter")
		}
		params["datacenter"] = c.String("d")
		if !c.IsSet("p") {
			return errors.NewMissingInputError("-p|--port-speed")
		}
		params["portSpeed"] = c.Int("p")

		params["billing"] = "hourly"
		if c.IsSet("b") {
			params["billing"] = c.String("b")
			if params["billing"] != "hourly" && params["billing"] != "monthly" {
				return errors.NewInvalidUsageError(T("-b|--billing has to be either hourly or monthly."))
			}
		}
		params["noPublic"] = false
		if c.IsSet("n") {
			params["noPublic"] = true
		}
		params["postInstallURL"] = c.String("i")
		params["ssheKeys"] = c.IntSlice("k")
		params["extras"] = c.StringSlice("e")
	}

	productPackage, err := cmd.HardwareManager.GetPackage()
	if err != nil {
		return cli.NewExitError(T("Failed to get product package for hardware server.\n"+err.Error()), 2)
	}
	orderTemplate, err := cmd.HardwareManager.GenerateCreateTemplate(productPackage, params)
	if err != nil {
		return err
	}

	if c.IsSet("test") {
		result, err := cmd.HardwareManager.VerifyOrder(orderTemplate)
		if err != nil {
			return cli.NewExitError(T("Failed to verify this order.\n")+err.Error(), 2)
		}
		table := cmd.UI.Table([]string{T("item"), T("cost")})
		total := 0.0
		for _, price := range result.Prices {
			if price.RecurringFee != nil && price.Item != nil && price.Item.Description != nil {
				rate := float64(*price.RecurringFee)
				total = total + float64(*price.RecurringFee)
				table.Add(*price.Item.Description, fmt.Sprintf("%.2f", rate))
			}
		}
		table.Add(T("Total monthly cost"), fmt.Sprintf("%.2f", total))
		table.Print()
		cmd.UI.Print(T("Prices reflected here are retail and do not take account level discounts and are not guaranteed."))
		return nil
	} else if c.IsSet("export") {
		content, err := json.Marshal(params)
		if err != nil {
			return cli.NewExitError(T("Failed to marshal hardware server template.\n")+err.Error(), 1)
		}
		export := c.String("export")
		// #nosec G306: write on customer machine
		err = ioutil.WriteFile(export, content, 0644)
		if err != nil {
			return cli.NewExitError(T("Failed to write hardware server template file to: {{.Template}}.\n",
				map[string]interface{}{"Template": export})+err.Error(), 1)
		}
		cmd.UI.Ok()
		cmd.UI.Print(T("Hardware server template is exported to: {{.Template}}.", map[string]interface{}{"Template": export}))
		return nil
	} else {
		if !c.IsSet("f") && !c.IsSet("force") {
			confirm, err := cmd.UI.Confirm(T("This action will incur charges on your account. Continue?"))
			if err != nil {
				return cli.NewExitError(err.Error(), 1)
			}
			if !confirm {
				cmd.UI.Print(T("Aborted."))
				return nil
			}
		}
		orderReceipt, err := cmd.HardwareManager.PlaceOrder(orderTemplate)
		if err != nil {
			return cli.NewExitError(T("Failed to place this order.\n")+err.Error(), 2)
		}
		cmd.UI.Ok()
		cmd.UI.Print(T("Order {{.OrderID}} was placed.", map[string]interface{}{"OrderID": *orderReceipt.OrderId}))
		cmd.UI.Print("")
		table := cmd.UI.Table([]string{T("item"), T("cost")})
		total := 0.0
		if orderReceipt.OrderDetails != nil && orderReceipt.OrderDetails.Prices != nil && len(orderReceipt.OrderDetails.Prices) > 0 {
			for _, price := range orderReceipt.OrderDetails.Prices {
				if price.RecurringFee != nil && price.Item != nil && price.Item.Description != nil {
					rate := float64(*price.RecurringFee)
					total = total + float64(*price.RecurringFee)
					table.Add(*price.Item.Description, fmt.Sprintf("%.2f", rate))
				}
			}
			table.Add(T("Total monthly cost"), fmt.Sprintf("%.2f", total))
		}
		cmd.UI.Print(T("Run '{{.CommandName}} sl hardware list --order {{.OrderID}}' to find this hardware server after it is ready.",
			map[string]interface{}{"OrderID": *orderReceipt.OrderId, "CommandName": cmd.Context.CLIName()}))
		return nil
	}
}

func HardwareCreateMetaData() cli.Command {
	return cli.Command{
		Category:    "hardware",
		Name:        "create",
		Description: T("Order/create a hardware server"),
		Usage: `${COMMAND_NAME} sl hardware create [OPTIONS] 
	See '${COMMAND_NAME} sl hardware create-options' for valid options.`,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "H,hostname",
				Usage: T("Host portion of the FQDN[required]"),
			},
			cli.StringFlag{
				Name:  "D,domain",
				Usage: T("Domain portion of the FQDN[required]"),
			},
			cli.StringFlag{
				Name:  "s,size",
				Usage: T("Hardware size[required]"),
			},
			cli.StringFlag{
				Name:  "o,os",
				Usage: T("OS install code[required]"),
			},
			cli.StringFlag{
				Name:  "d,datacenter",
				Usage: T("Datacenter shortname[required]"),
			},
			cli.IntFlag{
				Name:  "p,port-speed",
				Usage: T("Port speed[required]"),
			},
			cli.StringFlag{
				Name:  "b,billing",
				Usage: T("Billing rate, either hourly or monthly, default is hourly if not specified"),
			},
			cli.StringFlag{
				Name:  "i,post-install",
				Usage: T("Post-install script to download"),
			},
			cli.IntSliceFlag{
				Name:  "k,key",
				Usage: T("SSH keys to add to the root user, multiple occurrence allowed"),
			},
			cli.BoolFlag{
				Name:  "n,no-public",
				Usage: T("Private network only"),
			},
			cli.StringSliceFlag{
				Name:  "e,extra",
				Usage: T("Extra options, multiple occurrence allowed"),
			},
			cli.BoolFlag{
				Name:  "t,test",
				Usage: T("Do not actually create the virtual server"),
			},
			cli.StringFlag{
				Name:  "m,template",
				Usage: T("A template file that defaults the command-line options"),
			},
			cli.StringFlag{
				Name:  "x,export",
				Usage: T("Exports options to a template file"),
			},
			metadata.ForceFlag(),
		},
	}
}
