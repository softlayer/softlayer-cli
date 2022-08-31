package image

import (
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/spf13/cobra"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type ImportCommand struct {
	*metadata.SoftlayerCommand
	ImageManager managers.ImageManager
	Command      *cobra.Command
	Note         string
	OsCode       string
	RootKeyCrn   string
	WrappedDek   string
	CloudInit    bool
	Byol         bool
	IsEncrypted  bool
}

func NewImportCommand(sl *metadata.SoftlayerCommand) (cmd *ImportCommand) {
	thisCmd := &ImportCommand{
		SoftlayerCommand: sl,
		ImageManager:     managers.NewImageManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "import " + T("NAME") + " " + T("URI") + " " + T("API_KEY"),
		Short: T("Import an image from an object storage"),
		Long: T(`
EXAMPLE:
	${COMMAND_NAME} sl image import NAME URI API_KEY [--note NOTE] [--os-code OS_CODE] [--root-key-crn ROOT_KEY_CRN] [--wrapper-dek WRAPPER_DEK] [--cloud-init] [--byol] [--is-encrypted]
	NAME: The image name
	URI: The URI for an object storage object (.vhd/.iso file) of the format: cos://<regionName>/<bucketName>/<objectPath>
	API_KEY: The IBM Cloud API Key with access to IBM Cloud Object Storage instance.`),
		Args: metadata.ThreeArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().StringVar(&thisCmd.Note, "note", "", T("The note to be applied to the imported template"))
	cobraCmd.Flags().StringVar(&thisCmd.OsCode, "os-code", "", T("The referenceCode of the operating system software description for the imported VHD, ISO, or RAW image"))
	cobraCmd.Flags().StringVar(&thisCmd.RootKeyCrn, "root-key-crn", "", T("CRN of the root key in your KMS instance"))
	cobraCmd.Flags().StringVar(&thisCmd.WrappedDek, "wrapped-dek", "", T("Wrapped Data Encryption Key provided by IBM KeyProtect. For more info see: https://console.bluemix.net/docs/services/key-protect/wrap-keys.html#wrap-keys"))
	cobraCmd.Flags().BoolVar(&thisCmd.CloudInit, "cloud-init", false, T("Specifies if image is cloud-init"))
	cobraCmd.Flags().BoolVar(&thisCmd.Byol, "byol", false, T("Specifies if image is bring your own license"))
	cobraCmd.Flags().BoolVar(&thisCmd.IsEncrypted, "is-encrypted", false, T("Specifies if image is encrypted"))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *ImportCommand) Run(args []string) error {
	name := args[0]
	uri := args[1]
	ibmApiKey := args[2]

	var note *string
	var osCode *string
	var rootKeyCrn *string
	var wrapperDek *string

	if cmd.Note != "" {
		noteString := cmd.Note
		note = &noteString
	}

	if cmd.OsCode != "" {
		osCodeString := cmd.OsCode
		osCode = &osCodeString
	}

	if cmd.RootKeyCrn != "" {
		rootKeyCrnString := cmd.RootKeyCrn
		rootKeyCrn = &rootKeyCrnString
	}

	if cmd.WrappedDek != "" {
		wrapperDekString := cmd.WrappedDek
		wrapperDek = &wrapperDekString
	}

	cloudInit := cmd.CloudInit
	byol := cmd.Byol
	isEncrypted := cmd.IsEncrypted

	config := datatypes.Container_Virtual_Guest_Block_Device_Template_Configuration{
		Name:                         &name,
		Uri:                          &uri,
		IbmApiKey:                    &ibmApiKey,
		Note:                         note,
		OperatingSystemReferenceCode: osCode,
		CrkCrn:                       rootKeyCrn,
		WrappedDek:                   wrapperDek,
		CloudInit:                    &cloudInit,
		Byol:                         &byol,
		IsEncrypted:                  &isEncrypted,
	}

	resp, err := cmd.ImageManager.ImportImage(config)
	if err != nil {
		return err
	}

	cmd.UI.Ok()
	table := cmd.UI.Table([]string{T("Name"), T("Value")})
	table.Add(T("Name"), utils.FormatStringPointer(resp.Name))
	table.Add(T("ID"), utils.FormatIntPointer(resp.Id))
	table.Add(T("Created Date"), utils.FormatSLTimePointer(resp.CreateDate))
	table.Add(T("GUID"), utils.FormatStringPointer(resp.GlobalIdentifier))
	table.Print()
	return nil
}
