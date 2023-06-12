package cdn

import (
	"regexp"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/spf13/cobra"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type OriginAddCommand struct {
	*metadata.SoftlayerCommand
	CdnManager         managers.CdnManager
	Command            *cobra.Command
	Header             string
	Path               string
	OriginHost         string
	OriginType         string
	Http               int
	Https              int
	CacheKey           string
	Optimize           string
	DynamicPath        string
	DynamicPrefetch    bool
	DynamicCompression bool
	BuckeName          string
	fileExtension      string
}

func NewOriginAddCommand(sl *metadata.SoftlayerCommand) *OriginAddCommand {
	thisCmd := &OriginAddCommand{
		SoftlayerCommand: sl,
		CdnManager:       managers.NewCdnManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "origin-add " + T("IDENTIFIER"),
		Short: T("Create an origin path for an existing CDN mapping."),
		Long: T(`${COMMAND_NAME} sl cdn origin-add
Example:
${COMMAND_NAME} sl cdn origin-add --origin 123.123.123.123 --path /example/videos --http 80`),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().StringVar(&thisCmd.Header, "header", "", T("The edge server uses the host header in the HTTP header to communicate with the Origin host. It defaults to Hostname."))
	cobraCmd.Flags().StringVar(&thisCmd.Path, "path", "", T("Give a path relative to the domain provided, which can be used to reach this Origin. For example, 'articles/video' => 'www.example.com/articles/video [required]"))
	cobraCmd.Flags().StringVar(&thisCmd.OriginHost, "origin", "", T("Your server IP address or hostname. [required]"))
	cobraCmd.Flags().StringVar(&thisCmd.OriginType, "origin-type", "server", T("The origin type. [Permit: server, storage] Note: If OriginType is storage then OriginHost is take as Endpoint."))
	cobraCmd.Flags().IntVar(&thisCmd.Http, "http", 0, T("Http port. [http or https is required]"))
	cobraCmd.Flags().IntVar(&thisCmd.Https, "https", 0, T("Https port. [http or https is required]"))
	cobraCmd.Flags().StringVar(&thisCmd.CacheKey, "cache-key", "include-all", T("Cache query rules with the following formats: 'include-all', 'ignore-all', 'include: <query-names>', 'ignore: <query-names>'. example <query-names> = 'uuid=1234567 issue=important'."))
	cobraCmd.Flags().StringVar(&thisCmd.Optimize, "optimize", "web", T("Performance configuration. [Permit: web, video, file, dynamic]"))
	cobraCmd.Flags().StringVar(&thisCmd.DynamicPath, "dynamic-path", "", T("The path that Akamai edge servers periodically fetch the test object from. example = /detection-test-object.html"))
	cobraCmd.Flags().BoolVar(&thisCmd.DynamicPrefetch, "prefetching", true, T("Enable or disable the embedded object prefetching feature."))
	cobraCmd.Flags().BoolVar(&thisCmd.DynamicCompression, "compression", true, T("Enable or disable compression of JPEG images for requests over certain network conditions."))
	cobraCmd.Flags().StringVar(&thisCmd.BuckeName, "bucket-name", "", T("Bucket name."))
	cobraCmd.Flags().StringVar(&thisCmd.fileExtension, "file-extensions", "", T("Specify the file extensions that can be stored on the CDN service, separated by commas. For example, 'jpg, pdf, jpeg, png' is a valid list. Leave the flag empty to allow all extensions."))

	//#nosec G104 -- This is a false positive
	cobraCmd.MarkFlagRequired("path")
	//#nosec G104 -- This is a false positive
	cobraCmd.MarkFlagRequired("origin")
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *OriginAddCommand) Run(args []string) error {
	uniqueId := args[0]

	if cmd.Http == 0 && cmd.Https == 0 {
		return errors.NewMissingInputError(T("http or https"))
	}

	if cmd.OriginType != "server" && cmd.OriginType != "storage" {
		return errors.NewInvalidUsageError(T("--origintype"))
	}

	if cmd.Optimize != "web" && cmd.Optimize != "video" && cmd.Optimize != "file" && cmd.Optimize != "dynamic" {
		return errors.NewInvalidUsageError(T("--optimize"))
	}

	permitCacheKey := regexp.MustCompile(`^(ignore|include): \w+`)
	if cmd.CacheKey != "ignore-all" && cmd.CacheKey != "include-all" && !permitCacheKey.MatchString(cmd.CacheKey) {
		return errors.NewInvalidUsageError(T("--cache-key"))
	}

	if cmd.OriginType == "storage" && cmd.BuckeName == "" {
		return errors.NewInvalidUsageError(T("--bucket-name can not be empty"))
	}

	outputFormat := cmd.GetOutputFlag()

	newOrigin, err := cmd.CdnManager.OriginAddCdn(uniqueId, cmd.Header, cmd.Path, cmd.OriginHost, cmd.OriginType, cmd.Http, cmd.Https, cmd.CacheKey, cmd.Optimize, cmd.DynamicPath, cmd.DynamicPrefetch, cmd.DynamicCompression, cmd.BuckeName, cmd.fileExtension)
	if err != nil {
		return errors.NewAPIError(T("Failed to create a Origin."), err.Error(), 2)
	}

	PrintNewOrigin(cmd.UI, newOrigin, outputFormat)
	return nil
}

func PrintNewOrigin(ui terminal.UI, cdn []datatypes.Container_Network_CdnMarketplace_Configuration_Mapping_Path, outputFormat string) {
	table := ui.Table([]string{
		T("Name"),
		T("Value"),
	})
	if len(cdn) > 0 {
		table.Add(T("CDN Unique ID"), utils.FormatStringPointer(cdn[0].MappingUniqueId))
		if cdn[0].BucketName != nil {
			table.Add(T("Bucket Name"), utils.FormatStringPointer(cdn[0].BucketName))
		}
		if cdn[0].FileExtension != nil {
			table.Add(T("File Extension"), utils.FormatStringPointer(cdn[0].FileExtension))
		}
		table.Add(T("Header"), utils.FormatStringPointer(cdn[0].Header))
		table.Add(T("Path"), utils.FormatStringPointer(cdn[0].Path))
		table.Add(T("Origin"), utils.FormatStringPointer(cdn[0].Origin))
		table.Add(T("Origin Type"), utils.FormatStringPointer(cdn[0].OriginType))
		table.Add(T("Http Port"), utils.FormatIntPointer(cdn[0].HttpPort))
		table.Add(T("Https Port"), utils.FormatIntPointer(cdn[0].HttpsPort))
		table.Add(T("Cache Key Rule"), utils.FormatStringPointer(cdn[0].CacheKeyQueryRule))
		table.Add(T("Performance Configuration"), utils.FormatStringPointer(cdn[0].PerformanceConfiguration))
		table.Add(T("Status"), utils.FormatStringPointer(cdn[0].Status))
	}
	utils.PrintTable(ui, table, outputFormat)
}
