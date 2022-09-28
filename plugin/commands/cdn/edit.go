package cdn

import (
	"strconv"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/spf13/cobra"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type EditCommand struct {
	*metadata.SoftlayerCommand
	CdnManager               managers.CdnManager
	Command                  *cobra.Command
	Header                   string
	HttpPort                 int
	HttpsPort                int
	Origin                   string
	RespectHeaders           string
	Cache                    string
	CacheDescription         string
	PerformanceConfiguration string
}

func NewEditCommand(sl *metadata.SoftlayerCommand) *EditCommand {
	thisCmd := &EditCommand{
		SoftlayerCommand: sl,
		CdnManager:       managers.NewCdnManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "edit",
		Short: T("Edit a CDN Account."),
		Long:  T("${COMMAND_NAME} sl cdn edit"),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().StringVar(&thisCmd.Header, "header", "", T("Host header."))
	cobraCmd.Flags().IntVar(&thisCmd.HttpPort, "http-port", 0, T("HTTP port."))
	cobraCmd.Flags().IntVar(&thisCmd.HttpsPort, "https-port", 0, T("HTTPS port."))
	cobraCmd.Flags().StringVar(&thisCmd.Origin, "origin", "", T("Origin server address."))
	cobraCmd.Flags().StringVar(&thisCmd.RespectHeaders, "respect-headers", "", T("Respect headers. The value 1 is On and 0 is Off."))
	cobraCmd.Flags().StringVar(&thisCmd.Cache, "cache", "", T("Cache key optimization. These are the valid options to choose: 'include-all', 'ignore-all', 'include-specified', 'ignore-specified'. If you select 'include-specified' or 'ignore-specified' please add to option cache-description."))
	cobraCmd.Flags().StringVar(&thisCmd.CacheDescription, "cache-description", "", T("In cache option, if you select 'include-specified' or 'ignore-specified', please add a description too using this option e.g --cache include-specified --cache-description description."))
	cobraCmd.Flags().StringVar(&thisCmd.PerformanceConfiguration, "performance-configuration", "", T("Optimize for, 'General web delivery', 'Large file optimization', 'Video on demand optimization', the Dynamic content acceleration option is not added because this has a special configuration."))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *EditCommand) Run(args []string) error {
	cdnId, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("cdn ID")
	}

	if cmd.Header == "" && cmd.HttpPort == 0 && cmd.HttpsPort == 0 && cmd.Origin == "" && cmd.RespectHeaders == "" && cmd.Cache == "" && cmd.PerformanceConfiguration == "" {
		return slErr.NewInvalidUsageError(T("Please pass at least one of the flags."))
	}

	if cmd.RespectHeaders != "" {
		if cmd.RespectHeaders != "0" && cmd.RespectHeaders != "1" {
			return slErr.NewInvalidUsageError(T("Option respect-headers just accept '0' or '1'"))
		}
	}

	if cmd.Cache != "" {
		allowCache := []string{"include-all", "ignore-all", "include-specified", "ignore-specified"}
		if !utils.WordInList(allowCache, cmd.Cache) {
			return slErr.NewInvalidUsageError(T("Option cache just accept: " + utils.ArrayStringToString(allowCache)))
		}
		if cmd.Cache == "include-specified" || cmd.Cache == "ignore-specified" {
			if cmd.CacheDescription == "" {
				return slErr.NewInvalidUsageError(T("cache-description option must be used "))
			}
		}
	}

	if cmd.CacheDescription != "" {
		if cmd.Cache == "" {
			return slErr.NewInvalidUsageError(T("cache-description is only used with the cache option"))
		}
	}

	if cmd.PerformanceConfiguration != "" {
		allowPerformanceConfiguration := []string{"General web delivery", "Large file optimization", "Video on demand optimization"}
		if !utils.WordInList(allowPerformanceConfiguration, cmd.PerformanceConfiguration) {
			return slErr.NewInvalidUsageError(T("Option performance-configuration just accept: " + utils.ArrayStringToString(allowPerformanceConfiguration)))
		}
	}

	outputFormat := cmd.GetOutputFlag()

	cdnEdited, err := cmd.CdnManager.EditCDN(cdnId, cmd.Header, cmd.HttpPort, cmd.HttpsPort, cmd.Origin, cmd.RespectHeaders, cmd.Cache, cmd.Cache, cmd.PerformanceConfiguration)
	if err != nil {
		return errors.NewAPIError(T("Failed to edit CDN. "), err.Error(), 2)
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
