package bandwidth

import (
	"strconv"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/spf13/cobra"

	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

var LOCATION_GROUPS1 = map[string]string{
	"SJC/DAL/WDC/TOR/MON": "US/Canada",
	"AMS/LON/MAD/PAR":     "AMS/LON/MAD/PAR",
	"SNG/HKG/OSA/TOK":     "SNG/HKG/JPN",
	"SYD":                 "AUS",
	"MEX":                 "MEX",
	"SAO":                 "BRA",
	"CHE":                 "IND",
	"MIL":                 "ITA",
	"SEO":                 "KOR",
	"FRA":                 "FRA",
}

type EditCommand struct {
	*metadata.SoftlayerCommand
	BandwidthManager managers.BandwidthManager
	Command          *cobra.Command
	Name             string
}

func NewEditCommand(sl *metadata.SoftlayerCommand) (cmd *EditCommand) {
	thisCmd := &EditCommand{
		SoftlayerCommand: sl,
		BandwidthManager: managers.NewBandwidthManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "pools-edit " + T("IDENTIFIER"),
		Short: T("Edit bandwidth pool."),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().StringVar(&thisCmd.Name, "name", "", T("Pool name."))

	//#nosec G104 -- This is a false positive
	cobraCmd.MarkFlagRequired("name")

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *EditCommand) Run(args []string) error {
	bandwidthPoolId, err := strconv.Atoi(args[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("IDENTIFIER")
	}

	locationGroup, err := cmd.BandwidthManager.GetLocationGroup()
	if err != nil {
		return err
	}
	idLocationGroup := finId_LocationGroup(locationGroup, LOCATION_GROUPS1[cmd.Name])

	_, err = cmd.BandwidthManager.EditPool(bandwidthPoolId, idLocationGroup, cmd.Name)
	if err != nil {
		return err

	}
	cmd.UI.Ok()
	subs := map[string]interface{}{"bandwidthPoolId": bandwidthPoolId}
	cmd.UI.Print(T("Bandwidth pool {{.bandwidthPoolId}} was edited successfully.", subs))
	return nil
}

func finId_LocationGroup(locations []datatypes.Location_Group, regionName string) int {
	for _, location := range locations {
		if utils.FormatStringPointer(location.Name) == regionName {
			return *location.Id
		}
	}
	return 0
}
