package metadata

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/urfave/cli"
	. "github.ibm.com/cgallo/softlayer-cli/plugin/i18n"
)

var (
	NS_IMAGE_NAME  = "image"
	CMD_IMAGE_NAME = "image"

	CMD_IMG_DELETE_NAME = "delete"
	CMD_IMG_DETAIL_NAME = "detail"
	CMD_IMG_EDIT_NAME   = "edit"
	CMD_IMG_EXPORT_NAME = "export"
	CMD_IMG_IMPORT_NAME = "import"
	CMD_IMG_LIST_NAME   = "list"
)

func ImageNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  NS_SL_NAME,
		Name:        NS_IMAGE_NAME,
		Description: T("Classic infrastructure Compute images"),
	}
}

func ImageMetaData() cli.Command {
	return cli.Command{
		Category:    NS_SL_NAME,
		Name:        CMD_IMAGE_NAME,
		Description: T("Classic infrastructure Compute images"),
		Usage:       "${COMMAND_NAME} sl image",
		Subcommands: []cli.Command{
			ImageDelMetaData(),
			ImageDetailMetaData(),
			ImageEditMetaData(),
			ImageExportMetaData(),
			ImageImportMetaData(),
			ImageListMetaData(),
		},
	}
}

func ImageDelMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_IMAGE_NAME,
		Name:        CMD_IMG_DELETE_NAME,
		Description: T("Delete an image "),
		Usage: T(`${COMMAND_NAME} sl image delete IDENTIFIER

EXAMPLE: 
   ${COMMAND_NAME} sl image delete 12345678
   This command deletes image with ID 12345678.`),
	}
}
func ImageDetailMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_IMAGE_NAME,
		Name:        CMD_IMG_DETAIL_NAME,
		Description: T("Get details for an image"),
		Usage: T(`${COMMAND_NAME} sl image detail IDENTIFIER [OPTIONS]

EXAMPLE: 
   ${COMMAND_NAME} sl image detail 12345678
   This command gets details for image with ID 12345678.`),
		Flags: []cli.Flag{
			OutputFlag(),
		},
	}
}

func ImageEditMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_IMAGE_NAME,
		Name:        CMD_IMG_EDIT_NAME,
		Description: T("Edit details of an image"),
		Usage: T(`${COMMAND_NAME} sl image edit IDENTIFIER [OPTIONS]

EXAMPLE: 
   ${COMMAND_NAME} sl image edit 12345678 --name ubuntu16 --note testing --tag staging
   This command edits an image with ID 12345678 and set its name to "ubuntu16", note to "testing", and tag to "staging".`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "name",
				Usage: T("Name of the image"),
			},
			cli.StringFlag{
				Name:  "note",
				Usage: T("Add notes for the image"),
			},
			cli.StringFlag{
				Name:  "tag",
				Usage: T("Tags for the image"),
			},
		},
	}
}

func ImageExportMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_IMAGE_NAME,
		Name:        CMD_IMG_EXPORT_NAME,
		Description: T("Export an image to an object storage"),
		Usage:       T("${COMMAND_NAME} sl image export IDENTIFIER URI API_KEY\n  IDENTIFIER: ID of the image\n  URI: The URI for an object storage object (.vhd/.iso file) of the format: cos://<regionName>/<bucketName>/<objectPath>\n  API_KEY: The IBM Cloud API Key with access to IBM Cloud Object Storage instance."),
	}
}

func ImageImportMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_IMAGE_NAME,
		Name:        CMD_IMG_IMPORT_NAME,
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

func ImageListMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_IMAGE_NAME,
		Name:        CMD_IMG_LIST_NAME,
		Description: T("List all images on your account"),
		Usage: T(`${COMMAND_NAME} sl image list [OPTIONS]

EXAMPLE: 
   ${COMMAND_NAME} sl image list --public
   This command list all public images on current account.`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "name",
				Usage: T("Filter on image name"),
			},
			cli.BoolFlag{
				Name:  "public",
				Usage: T("Display only public images"),
			},
			cli.BoolFlag{
				Name:  "private",
				Usage: T("Display only private images"),
			},
			OutputFlag(),
		},
	}
}
