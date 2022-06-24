package file

import (
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

func FileDuplicateConvertStatusMetaData() cli.Command {
	return cli.Command{
		Category:    "file",
		Name:        "duplicate-convert-status",
		Description: T("Get status for split or move completed percentage of a given file storage duplicate volume."),
		Usage: T(`${COMMAND_NAME} sl file duplicate-convert-status [OPTIONS] VOLUME_ID

EXAMPLE:
   ${COMMAND_NAME} sl file duplicate-convert-status 12345678`),
		Flags: []cli.Flag{
			metadata.OutputFlag(),
		},
	}
}
