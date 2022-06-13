package block

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type Duplicate_Conversion_Status struct {
	datatypes.Entity
	ActiveConversionStartTime       *string `json:"activeConversionStartTime,omitempty" xmlrpc:"activeConversionStartTime,omitempty"`
	DeDuplicateConversionPercentage *int    `json:"deDuplicateConversionPercentage,omitempty" xmlrpc:"deDuplicateConversionPercentage,omitempty"`
	VolumeUsername                  *string `json:"volumeUsername,omitempty" xmlrpc:"volumeUsername,omitempty"`
}

type DuplicateConvertStatusCommand struct {
	UI      terminal.UI
	Session *session.Session
}

func NewDuplicateConvertStatusCommand(ui terminal.UI, session *session.Session) (cmd *DuplicateConvertStatusCommand) {
	return &DuplicateConvertStatusCommand{
		UI:      ui,
		Session: session,
	}
}

func (cmd *DuplicateConvertStatusCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	volumeID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Volume ID")
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	duplicateConversionStatus, err := getDuplicateConversionStatus(cmd.Session, volumeID)
	if err != nil {
		return cli.NewExitError(T("Failed to get duplicate conversion status of volume {{.VolumeID}}.\n",
			map[string]interface{}{"VolumeID": volumeID})+err.Error(), 2)
	}

	table := cmd.UI.Table([]string{T("Username"), T("Active Conversion Start Timestamp"), T("Completed Percentage")})
	table.Add(
		utils.FormatStringPointer(duplicateConversionStatus.VolumeUsername),
		utils.FormatStringPointer(duplicateConversionStatus.ActiveConversionStartTime),
		utils.FormatIntPointer(duplicateConversionStatus.DeDuplicateConversionPercentage),
	)

	utils.PrintTable(cmd.UI, table, outputFormat)

	return nil
}

func getDuplicateConversionStatus(sess *session.Session, volumeID int) (resp Duplicate_Conversion_Status, err error) {
	mask := "mask[activeConversionStartTime,deDuplicateConversionPercentage,volumeUsername]"
	var options sl.Options
	options.Mask = mask
	options.Id = &volumeID
	err = sess.DoRequest("SoftLayer_Network_Storage", "getDuplicateConversionStatus", nil, &options, &resp)
	return
}

func BlockDuplicateConvertStatusMetaData() cli.Command {
	return cli.Command{
		Category:    "block",
		Name:        "duplicate-convert-status",
		Description: T("Get status for split or move completed percentage of a given block storage duplicate volume."),
		Usage: T(`${COMMAND_NAME} sl block duplicate-convert-status [OPTIONS] VOLUME_ID

EXAMPLE:
   ${COMMAND_NAME} sl block duplicate-convert-status 12345678`),
		Flags: []cli.Flag{
			metadata.OutputFlag(),
		},
	}
}
