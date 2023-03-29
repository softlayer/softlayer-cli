package cdn

import (
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
	CdnManager managers.CdnManager
	Command    *cobra.Command
	HostName   string
	OriginHost string
	OriginType string
	Http       int
	Https      int
	BucketName string
	CName      string
	Header     string
	Path       string
	Ssl        string
}

func NewOriginAddCommand(sl *metadata.SoftlayerCommand) *OriginAddCommand {
	thisCmd := &OriginAddCommand{
		SoftlayerCommand: sl,
		CdnManager:       managers.NewCdnManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "origin-add",
		Short: T("Create an origin path for an existing CDN mapping."),
		Long: T(`${COMMAND_NAME} sl cdn origin-add
Example:
${COMMAND_NAME} sl cdn origin-add --hostname www.example.com --origin 123.45.67.8 --http 80`),
		Args: metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().StringVar(&thisCmd.HostName, "hostname", "", T("To route requests to your website, enter the hostname for your website, for example, www.example.com or app.example.com. [required]"))
	cobraCmd.Flags().StringVar(&thisCmd.OriginHost, "origin", "", T("Your server IP address or hostname. [required]"))
	cobraCmd.Flags().StringVar(&thisCmd.OriginType, "origin-type", "server", T("The origin type. [Permit: server, storage] Note: If OriginType is storage then OriginHost is take as Endpoint"))
	cobraCmd.Flags().IntVar(&thisCmd.Http, "http", 0, T("Http port"))
	cobraCmd.Flags().IntVar(&thisCmd.Https, "https", 0, T("Https port"))
	cobraCmd.Flags().StringVar(&thisCmd.BucketName, "bucket-name", "", T("Bucket name"))
	cobraCmd.Flags().StringVar(&thisCmd.CName, "cname", "", T("Enter a globally unique subdomain. The full URL becomes the CNAME we use to configure your DNS. If no value is entered, we will generate a CNAME for you."))
	cobraCmd.Flags().StringVar(&thisCmd.Header, "header", "", T("The edge server uses the host header in the HTTP header to communicate with the Origin host. It defaults to Hostname."))
	cobraCmd.Flags().StringVar(&thisCmd.Path, "path", "", T("Give a path relative to the domain provided, which can be used to reach this Origin. For example, 'articles/video' => 'www.example.com/articles/video"))
	cobraCmd.Flags().StringVar(&thisCmd.Ssl, "ssl", "dvSan", T("A DV SAN Certificate allows HTTPS traffic over your personal domain, but it requires a domain validation to prove ownership. A wildcard certificate allows HTTPS traffic only when using the CNAME given. [Permit: dvSan, wilcard]"))

	cobraCmd.MarkFlagRequired("hostname")
	cobraCmd.MarkFlagRequired("origin")
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *OriginAddCommand) Run(args []string) error {
	if cmd.Http == 0 && cmd.Https == 0 {
		return errors.NewMissingInputError("http or https")
	}
	if cmd.OriginType != "server" && cmd.OriginType != "storage" {
		return errors.NewInvalidUsageError("--origintype")
	}
	if cmd.Ssl != "dvSan" && cmd.Ssl != "wilcard" && cmd.Ssl != "" {
		return errors.NewInvalidUsageError("--ssl")
	}

	outputFormat := cmd.GetOutputFlag()

	newCdn, err := cmd.CdnManager.OriginAdd(cmd.HostName, cmd.OriginHost, cmd.OriginType, cmd.Http, cmd.Https, cmd.BucketName, cmd.CName, cmd.Header, cmd.Path, cmd.Ssl)
	if err != nil {
		return errors.NewAPIError(T("Failed to create a CDN."), err.Error(), 2)
	}

	PrintCndCreated(cmd.UI, newCdn, outputFormat)
	return nil
}

func PrintCndCreated(ui terminal.UI, cdn []datatypes.Container_Network_CdnMarketplace_Configuration_Mapping, outputFormat string) {
	table := ui.Table([]string{
		T("Name"),
		T("Value"),
	})
	if len(cdn) > 0 {
		table.Add(T("CDN Unique ID"), utils.FormatStringPointer(cdn[0].UniqueId))
		if cdn[0].BucketName != nil {
			table.Add(T("Bucket Name"), utils.FormatStringPointer(cdn[0].BucketName))
		}
		table.Add(T("Hostname"), utils.FormatStringPointer(cdn[0].Domain))
		table.Add(T("Header"), utils.FormatStringPointer(cdn[0].Header))
		table.Add(T("IBM CNAME"), utils.FormatStringPointer(cdn[0].Cname))
		table.Add(T("Akamai CNAME"), utils.FormatStringPointer(cdn[0].AkamaiCname))
		table.Add(T("Origin Host"), utils.FormatStringPointer(cdn[0].OriginHost))
		table.Add(T("Origin Type"), utils.FormatStringPointer(cdn[0].OriginType))
		table.Add(T("Protocol"), utils.FormatStringPointer(cdn[0].Protocol))
		table.Add(T("Http Port"), utils.FormatIntPointer(cdn[0].HttpPort))
		table.Add(T("Https Port"), utils.FormatIntPointer(cdn[0].HttpsPort))
		table.Add(T("Certificate Type"), utils.FormatStringPointer(cdn[0].CertificateType))
		table.Add(T("Provider"), utils.FormatStringPointer(cdn[0].VendorName))
		table.Add(T("Path"), utils.FormatStringPointer(cdn[0].Path))
	}
	utils.PrintTable(ui, table, outputFormat)
}
