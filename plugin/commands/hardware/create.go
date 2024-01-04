package hardware

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type CreateCommand struct {
	*metadata.SoftlayerCommand
	HardwareManager managers.HardwareServerManager
	Command         *cobra.Command
	Hostname        string
	Domain          string
	Size            string
	Os              string
	Datacenter      string
	PortSpeed       int
	Billing         string
	PostInstall     string
	Key             []int
	NoPublic        bool
	Extra           []string
	Test            bool
	Template        string
	Export          string
	ForceFlag       bool
}

func NewCreateCommand(sl *metadata.SoftlayerCommand) (cmd *CreateCommand) {
	thisCmd := &CreateCommand{
		SoftlayerCommand: sl,
		HardwareManager:  managers.NewHardwareServerManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "create",
		Short: T("Order/create a hardware server"),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().StringVarP(&thisCmd.Hostname, "hostname", "H", "", T("Host portion of the FQDN[required]"))
	cobraCmd.Flags().StringVarP(&thisCmd.Domain, "domain", "D", "", T("Domain portion of the FQDN[required]"))
	cobraCmd.Flags().StringVarP(&thisCmd.Size, "size", "s", "", T("Hardware size[required]"))
	cobraCmd.Flags().StringVarP(&thisCmd.Os, "os", "o", "", T("OS install code[required]"))
	cobraCmd.Flags().StringVarP(&thisCmd.Datacenter, "datacenter", "d", "", T("Datacenter shortname[required]"))
	cobraCmd.Flags().IntVarP(&thisCmd.PortSpeed, "port-speed", "p", 0, T("Port speed[required]"))
	cobraCmd.Flags().StringVarP(&thisCmd.Billing, "billing", "b", "", T("Billing rate, either hourly or monthly, default is hourly if not specified"))
	cobraCmd.Flags().StringVarP(&thisCmd.PostInstall, "post-install", "i", "", T("Post-install script to download"))
	cobraCmd.Flags().IntSliceVarP(&thisCmd.Key, "key", "k", []int{}, T("SSH keys to add to the root user, multiple occurrence allowed"))
	cobraCmd.Flags().BoolVarP(&thisCmd.NoPublic, "no-public", "n", false, T("Private network only"))
	cobraCmd.Flags().StringSliceVarP(&thisCmd.Extra, "extra", "e", []string{}, T("Extra options, multiple occurrence allowed"))
	cobraCmd.Flags().BoolVarP(&thisCmd.Test, "test", "t", false, T("Do not actually create the virtual server"))
	cobraCmd.Flags().StringVarP(&thisCmd.Template, "template", "m", "", T("A template file that defaults the command-line options"))
	cobraCmd.Flags().StringVarP(&thisCmd.Export, "export", "x", "", T("Exports options to a template file"))
	cobraCmd.Flags().BoolVarP(&thisCmd.ForceFlag, "force", "f", false, T("Force operation without confirmation"))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *CreateCommand) Run(args []string) error {
	params := make(map[string]interface{})
	if cmd.Template != "" {
		templateFile := cmd.Template
		if _, err := os.Stat(templateFile); os.IsNotExist(err) {
			return errors.NewInvalidUsageError(T("Template file: {{.Location}} does not exist.",
				map[string]interface{}{"Location": templateFile}))
		}
		content, err := ioutil.ReadFile(templateFile) // #nosec
		if err != nil {
			return errors.NewInvalidUsageError(T("Failed to read template file: {{.File}}.\n", map[string]interface{}{"File": templateFile}) + err.Error())
		}
		err = json.Unmarshal(content, &params)
		if err != nil {
			return errors.NewInvalidUsageError(T("Failed to unmarshal template file: {{.File}}.\n", map[string]interface{}{"File": templateFile}) + err.Error())
		}
	} else {
		if cmd.Size == "" {
			return errors.NewMissingInputError("-s|--size")
		}
		params["size"] = cmd.Size
		if cmd.Hostname == "" {
			return errors.NewMissingInputError("-H|--hostname")
		}
		params["hostname"] = cmd.Hostname
		if cmd.Domain == "" {
			return errors.NewMissingInputError("-D|--domain")
		}
		params["domain"] = cmd.Domain
		if cmd.Os == "" {
			return errors.NewMissingInputError("-o|--os")
		}
		params["osName"] = cmd.Os
		if cmd.Datacenter == "" {
			return errors.NewMissingInputError("-d|--datacenter")
		}
		params["datacenter"] = cmd.Datacenter
		if cmd.PortSpeed == 0 {
			return errors.NewMissingInputError("-p|--port-speed")
		}
		params["portSpeed"] = cmd.PortSpeed

		params["billing"] = "hourly"
		if cmd.Billing != "" {
			params["billing"] = cmd.Billing
			if params["billing"] != "hourly" && params["billing"] != "monthly" {
				return errors.NewInvalidUsageError(T("-b|--billing has to be either hourly or monthly."))
			}
		}
		params["noPublic"] = false
		if cmd.NoPublic {
			params["noPublic"] = true
		}
		params["postInstallURL"] = cmd.PostInstall
		params["sshKeys"] = cmd.Key
		params["extras"] = cmd.Extra
	}

	productPackage, err := cmd.HardwareManager.GetPackage()
	if err != nil {
		return errors.NewAPIError(T("Failed to get product package for hardware server.\n"), err.Error(), 2)
	}
	orderTemplate, err := cmd.HardwareManager.GenerateCreateTemplate(productPackage, params)
	if err != nil {
		return err
	}

	if cmd.Test {
		result, err := cmd.HardwareManager.VerifyOrder(orderTemplate)
		if err != nil {
			return errors.NewAPIError(T("Failed to verify this order.\n"), err.Error(), 2)
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
	} else if cmd.Export != "" {
		content, err := json.Marshal(params)
		if err != nil {
			return errors.NewAPIError(T("Failed to marshal hardware server template.\n"), err.Error(), 1)
		}
		export := cmd.Export
		// #nosec G306: write on customer machine
		err = ioutil.WriteFile(export, content, 0644)
		if err != nil {
			return errors.NewAPIError(T("Failed to write hardware server template file to: {{.Template}}.\n",
				map[string]interface{}{"Template": export}), err.Error(), 1)
		}
		cmd.UI.Ok()
		cmd.UI.Print(T("Hardware server template is exported to: {{.Template}}.", map[string]interface{}{"Template": export}))
		return nil
	} else {
		if !cmd.ForceFlag {
			confirm, err := cmd.UI.Confirm(T("This action will incur charges on your account. Continue?"))
			if err != nil {
				return err
			}
			if !confirm {
				cmd.UI.Print(T("Aborted."))
				return nil
			}
		}
		fmt.Printf("ORDER TEMPLATE: %v\n\n", orderTemplate)
		orderReceipt, err := cmd.HardwareManager.PlaceOrder(orderTemplate)
		if err != nil {
			return errors.NewAPIError(T("Failed to place this order.\n"), err.Error(), 2)
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
			map[string]interface{}{"OrderID": *orderReceipt.OrderId, "CommandName": "ibmcloud"}))
		return nil
	}
}
