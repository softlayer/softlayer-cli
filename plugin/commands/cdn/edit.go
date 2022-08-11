package cdn

import (
	"strconv"

	"github.com/softlayer/softlayer-go/datatypes"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"

	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type EditCommand struct {
	UI         terminal.UI
	CdnManager managers.CdnManager
}

func NewEditCommand(ui terminal.UI, cdnManager managers.CdnManager) (cmd *EditCommand) {
	return &EditCommand{
		UI:         ui,
		CdnManager: cdnManager,
	}
}

func EditMetaData() cli.Command {
	return cli.Command{
		Category:    "cdn",
		Name:        "edit",
		Description: T("Edit a CDN Account."),
		Usage:       T(`${COMMAND_NAME} sl cdn edit`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "header",
				Usage: T("Host header."),
			},
			cli.IntFlag{
				Name:  "http-port",
				Usage: T("HTTP port."),
			},
			cli.IntFlag{
				Name:  "https-port",
				Usage: T("HTTPS port."),
			},
			cli.StringFlag{
				Name:  "origin",
				Usage: T("Origin server address."),
			},
			cli.StringFlag{
				Name:  "respect-headers",
				Usage: T("Respect headers. The value 1 is On and 0 is Off."),
			},
			cli.StringFlag{
				Name:  "cache",
				Usage: T("Cache key optimization. These are the valid options to choose: 'include-all', 'ignore-all', 'include-specified', 'ignore-specified'. If you select 'include-specified' or 'ignore-specified' please add to option cache-description."),
			},
			cli.StringFlag{
				Name:  "cache-description",
				Usage: T("In cache option, if you select 'include-specified' or 'ignore-specified', please add a description too using this option e.g --cache include-specified --cache-description description."),
			},
			cli.StringFlag{
				Name:  "performance-configuration",
				Usage: T("Optimize for, 'General web delivery', 'Large file optimization', 'Video on demand optimization', the Dynamic content acceleration option is not added because this has a special configuration."),
			},
			metadata.OutputFlag(),
		},
	}
}

func (cmd *EditCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return slErr.NewInvalidUsageError(T("This command requires one argument."))
	}

	cdnId, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("cdn ID")
	}

	if !c.IsSet("header") && !c.IsSet("http-port") && !c.IsSet("https-port") && !c.IsSet("origin") && !c.IsSet("respect-headers") && !c.IsSet("cache") && !c.IsSet("performance-configuration") {
		return slErr.NewInvalidUsageError(T("Please pass at least one of the flags."))
	}

	if c.IsSet("respect-headers") {
		if c.String("respect-headers") != "0" && c.String("respect-headers") != "1" {
			return slErr.NewInvalidUsageError(T("Option respect-headers just accept '0' or '1'"))
		}
	}

	if c.IsSet("cache") {
		allowCache := []string{"include-all", "ignore-all", "include-specified", "ignore-specified"}
		if !utils.WordInList(allowCache, c.String("cache")) {
			return slErr.NewInvalidUsageError(T("Option cache just accept: " + utils.ArrayStringToString(allowCache)))
		}
		if c.String("cache") == "include-specified" || c.String("cache") == "ignore-specified" {
			if !c.IsSet("cache-description") {
				return slErr.NewInvalidUsageError(T("cache-description option must be used "))
			}
		}
	}

	if c.IsSet("cache-description") {
		if !c.IsSet("cache") {
			return slErr.NewInvalidUsageError(T("cache-description is only used with the cache option"))
		}
	}

	if c.IsSet("performance-configuration") {
		allowPerformanceConfiguration := []string{"General web delivery", "Large file optimization", "Video on demand optimization"}
		if !utils.WordInList(allowPerformanceConfiguration, c.String("performance-configuration")) {
			return slErr.NewInvalidUsageError(T("Option performance-configuration just accept: " + utils.ArrayStringToString(allowPerformanceConfiguration)))
		}
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	cdnEdited, err := cmd.CdnManager.EditCDN(cdnId, c.String("header"), c.Int("http-port"), c.Int("https-port"), c.String("origin"), c.String("respect-headers"), c.String("cache"), c.String("cache-description"), c.String("performance-configuration"))
	if err != nil {
		return cli.NewExitError(T("Failed to edit CDN. ")+err.Error(), 2)
	}

	PrintEditedCDN(cdnEdited, cmd.UI, outputFormat)
	return nil
}

func PrintEditedCDN(cdnEdited datatypes.Container_Network_CdnMarketplace_Configuration_Mapping, ui terminal.UI, outputFormat string) {

	table := ui.Table([]string{
		T("Name"),
		T("Value"),
	})
	table.Add(T("Create Date"), utils.FormatSLTimePointer(cdnEdited.CreateDate))
	table.Add(T("Header"), utils.FormatStringPointer(cdnEdited.Header))
	if cdnEdited.HttpPort != nil {
		table.Add(T("Http Port"), utils.FormatIntPointer(cdnEdited.HttpPort))
	}
	if cdnEdited.HttpsPort != nil {
		table.Add(T("Https Port"), utils.FormatIntPointer(cdnEdited.HttpsPort))
	}
	table.Add(T("Origin Type"), utils.FormatStringPointer(cdnEdited.OriginType))
	table.Add(T("Performance Configuration"), utils.FormatStringPointer(cdnEdited.PerformanceConfiguration))
	table.Add(T("Protocol"), utils.FormatStringPointer(cdnEdited.Protocol))
	table.Add(T("Respect Headers"), utils.FormatBoolPointer(cdnEdited.RespectHeaders))
	table.Add(T("Unique Id"), utils.FormatStringPointer(cdnEdited.UniqueId))
	table.Add(T("Vendor Name"), utils.FormatStringPointer(cdnEdited.VendorName))
	table.Add(T("Cache key optimization"), utils.FormatStringPointer(cdnEdited.CacheKeyQueryRule))
	table.Add(T("Cname"), utils.FormatStringPointer(cdnEdited.Cname))
	table.Add(T("Origin server address"), utils.FormatStringPointer(cdnEdited.OriginHost))

	utils.PrintTable(ui, table, outputFormat)
}
