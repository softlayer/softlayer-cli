package bandwidth

import (
	"strconv"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/spf13/cobra"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
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
		Short: T("edit bandwidth pool. "),
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
		return slErrors.NewInvalidSoftlayerIdInputError("Bandwidth Pool ID")
	}
	if cmd.Name == "" {
		return slErrors.NewInvalidUsageError(T("--name must be specified."))
	}
	locationGroup, err := cmd.BandwidthManager.GetLocationGroup()
	if err != nil {
		return errors.NewAPIError(T("Failed to get Location Group."), err.Error(), 2)
	}
	idLocationGroup := finId_LocationGroup(locationGroup, LOCATION_GROUPS1[cmd.Name])

	_, err = cmd.BandwidthManager.EditPool(bandwidthPoolId, idLocationGroup, cmd.Name)
	if err != nil {
		return slErrors.NewAPIError(T("Failed to edit bandwidth with Id: {{.bandwidthPoolId}}.\n", map[string]interface{}{"bandwidthPoolId": bandwidthPoolId}), err.Error(), 2)

	}
	cmd.UI.Ok()
	subs := map[string]interface{}{"bandwidthPoolId": bandwidthPoolId}
	cmd.UI.Print(T("BandwidthPool associated with Id {{.bandwidthPoolId}} was edited successfully.", subs))
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
