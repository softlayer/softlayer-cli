package file

import (
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

func FileVolumeSetNoteMetaData() cli.Command {
	return cli.Command{
		Category:    "file",
		Name:        "volume-set-note",
		Description: T("Set note for an existing file storage volume."),
		Usage: T(`${COMMAND_NAME} sl file volume-set-note [OPTIONS] VOLUME_ID

EXAMPLE:
   ${COMMAND_NAME} sl file volume-set-note 12345678 --note "this is my note"`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "n,note",
				Usage: T("Public notes related to a Storage volume  [required]"),
			},
			metadata.OutputFlag(),
		},
	}
}
