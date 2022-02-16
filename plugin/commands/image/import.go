package image

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type ImportCommand struct {
	UI           terminal.UI
	ImageManager managers.ImageManager
}

func NewImportCommand(ui terminal.UI, imageManager managers.ImageManager) (cmd *ImportCommand) {
	return &ImportCommand{
		UI:           ui,
		ImageManager: imageManager,
	}
}

func (cmd *ImportCommand) Run(c *cli.Context) error {
	if c.NArg() != 3 {
		return errors.NewInvalidUsageError(T("This command requires three arguments."))
	}
	name := c.Args()[0]
	uri := c.Args()[1]
	ibmApiKey := c.Args()[2]

	var note *string
	var osCode *string
	var rootKeyCrn *string
	var wrapperDek *string

	if c.IsSet("note") {
		noteString := c.String("note")
		note = &noteString
	}

	if c.IsSet("os-code") {
		osCodeString := c.String("os-code")
		osCode = &osCodeString
	}

	if c.IsSet("root-key-crn") {
		rootKeyCrnString := c.String("root-key-crn")
		rootKeyCrn = &rootKeyCrnString
	}

	if c.IsSet("wrapped-dek") {
		wrapperDekString := c.String("wrapped-dek")
		wrapperDek = &wrapperDekString
	}

	cloudInit := c.Bool("cloud-init")
	byol := c.Bool("byol")
	isEncrypted := c.Bool("is-encrypted")

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

func ImageImportMetaData() cli.Command {
	return cli.Command{
		Category:    "image",
		Name:        "import",
		Description: T("Import an image from an object storage"),
		Usage:       T("${COMMAND_NAME} sl image import NAME URI API_KEY [--note NOTE] [--os-code OS_CODE] [--root-key-crn ROOT_KEY_CRN] [--wrapper-dek WRAPPER_DEK] [--cloud-init] [--byol] [--is-encrypted]\n  NAME: The image name\n  URI: The URI for an object storage object (.vhd/.iso file) of the format: cos://<regionName>/<bucketName>/<objectPath>\n  API_KEY: The IBM Cloud API Key with access to IBM Cloud Object Storage instance."),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "note",
				Usage: T("The note to be applied to the imported template"),
			},
			cli.StringFlag{
				Name:  "os-code",
				Usage: T("The referenceCode of the operating system software description for the imported VHD, ISO, or RAW image"),
			},
			cli.StringFlag{
				Name:  "root-key-crn",
				Usage: T("CRN of the root key in your KMS instance"),
			},
			cli.StringFlag{
				Name:  "wrapped-dek",
				Usage: T("Wrapped Data Encryption Key provided by IBM KeyProtect. For more info see: https://console.bluemix.net/docs/services/key-protect/wrap-keys.html#wrap-keys"),
			},
			cli.BoolFlag{
				Name:  "cloud-init",
				Usage: T("Specifies if image is cloud-init"),
			},
			cli.BoolFlag{
				Name:  "byol",
				Usage: T("Specifies if image is bring your own license"),
			},
			cli.BoolFlag{
				Name:  "is-encrypted",
				Usage: T("Specifies if image is encrypted"),
			},
		},
	}
}
