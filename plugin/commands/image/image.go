package image

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/softlayer/softlayer-go/session"
	"github.com/urfave/cli"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

func GetCommandActionBindings(context plugin.PluginContext, ui terminal.UI, session *session.Session) map[string]func(c *cli.Context) error {
	imageManager := managers.NewImageManager(session)

	CommandActionBindings := map[string]func(c *cli.Context) error{
		"image-delete": func(c *cli.Context) error {
			return NewDeleteCommand(ui, imageManager).Run(c)
		},
		"image-detail": func(c *cli.Context) error {
			return NewDetailCommand(ui, imageManager).Run(c)
		},
		"image-edit": func(c *cli.Context) error {
			return NewEditCommand(ui, imageManager).Run(c)
		},
		"image-export": func(c *cli.Context) error {
			return NewExportCommand(ui, imageManager).Run(c)
		},
		"image-import": func(c *cli.Context) error {
			return NewImportCommand(ui, imageManager).Run(c)
		},
		"image-list": func(c *cli.Context) error {
			return NewListCommand(ui, imageManager).Run(c)
		},
		"image-datacenter": func(c *cli.Context) error {
			return NewDatacenterCommand(ui, imageManager).Run(c)
		},
	}

	return CommandActionBindings
}

func ImageNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  "sl",
		Name:        "image",
		Description: T("Classic infrastructure Compute images"),
	}
}

func ImageMetaData() cli.Command {
	return cli.Command{
		Category:    "sl",
		Name:        "image",
		Description: T("Classic infrastructure Compute images"),
		Usage:       "${COMMAND_NAME} sl image",
		Subcommands: []cli.Command{
			ImageDelMetaData(),
			ImageDetailMetaData(),
			ImageEditMetaData(),
			ImageExportMetaData(),
			ImageImportMetaData(),
			ImageListMetaData(),
			ImageDatacenterMetaData(),
		},
	}
}
