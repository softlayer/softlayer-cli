package bandwidth

import (
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/spf13/cobra"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

var LOCATION_GROUPS = map[string]string{
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

type PoolsCreateCommand struct {
	*metadata.SoftlayerCommand
	BandwidthManager managers.BandwidthManager
	Command          *cobra.Command
	Name             string
	Region           string
}

func NewPoolsCreateCommand(sl *metadata.SoftlayerCommand) *PoolsCreateCommand {

	thisCmd := &PoolsCreateCommand{
		SoftlayerCommand: sl,
		BandwidthManager: managers.NewBandwidthManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "pools-create",
		Short: T("Create a bandwidth pool."),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().StringVar(&thisCmd.Name, "name", "", T("Pool name."))
	cobraCmd.Flags().StringVar(&thisCmd.Region, "region", "", T("Region selected. Permit:[SJC/DAL/WDC/TOR/MON, AMS/LON/MAD/PAR, SNG/HKG/OSA/TOK, SYD, MEX, SAO, CHE, MIL, SEO, FRA]"))

	//#nosec G104 -- This is a false positive
	cobraCmd.MarkFlagRequired("name")
	//#nosec G104 -- This is a false positive
	cobraCmd.MarkFlagRequired("region")

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *PoolsCreateCommand) Run(args []string) error {

	outputFormat := cmd.GetOutputFlag()

	locationGroup, err := cmd.BandwidthManager.GetLocationGroup()
	if err != nil {
		return errors.NewAPIError(T("Failed to get Location Group."), err.Error(), 2)
	}
	idLocationGroup := finIdLocationGroup(locationGroup, LOCATION_GROUPS[cmd.Region])
	poolCreated, err := cmd.BandwidthManager.CreatePool(cmd.Name, idLocationGroup)
	if err != nil {
		return errors.NewAPIError(T("Failed to get Create Bandwidth Pool."), err.Error(), 2)
	}

	table := cmd.UI.Table([]string{T("Name"), T("Value")})
	table.Add("Id", utils.FormatIntPointer(poolCreated.Id))
	table.Add("Name Pool", utils.FormatStringPointer(poolCreated.Name))
	table.Add("Region", cmd.Region)
	table.Add("Created Date", utils.FormatSLTimePointer(poolCreated.CreateDate))

	utils.PrintTable(cmd.UI, table, outputFormat)
	return nil
}

func finIdLocationGroup(locations []datatypes.Location_Group, regionName string) int {
	for _, location := range locations {
		if utils.FormatStringPointer(location.Name) == regionName {
			return *location.Id
		}
	}
	return 0
}
