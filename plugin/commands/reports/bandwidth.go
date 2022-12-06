package reports

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/spf13/cobra"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type BandwidthCommand struct {
	*metadata.SoftlayerCommand
	ReportManager managers.ReportManager
	Command       *cobra.Command
	Start         string
	End           string
	SortBy        string
	Virtual       bool
	Server        bool
	Pool          bool
}

type metricObject struct {
	id         int
	typeDevice string
	name       string
	pool       string
	data       []datatypes.Metric_Tracking_Object_Data
}

type tableRow struct {
	typeDevice string
	hostname   string
	publicIn   int
	publicOut  int
	privateIn  int
	privateOut int
	pool       string
}

func NewBandwidthCommand(sl *metadata.SoftlayerCommand) *BandwidthCommand {
	thisCmd := &BandwidthCommand{
		SoftlayerCommand: sl,
		ReportManager:    managers.NewReportManager(sl.Session),
	}
	cobraCmd := &cobra.Command{

		Use:   "bandwidth",
		Short: T("Bandwidth report for every pool/server."),
		Long: `EXAMPLE: 
   ${COMMAND_NAME} sl report bandwidth
   ${COMMAND_NAME} sl report bandwidth --server --start 2022-06-07 --end 2022-06-08
   ${COMMAND_NAME} sl report bandwidth --start 2022-06-07 --end 2022-06-08 --sortby privateOut`,
		Args: metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().StringVar(&thisCmd.Start, "start", "", T("datetime in the format 'YYYY-MM-DD' or 'YYYY-MM-DD HH:MM:SS'"))
	cobraCmd.Flags().StringVar(&thisCmd.End, "end", "", T("datetime in the format 'YYYY-MM-DD' or 'YYYY-MM-DD HH:MM:SS'"))
	cobraCmd.Flags().StringVar(&thisCmd.SortBy, "sortby", "", T("Column to sort by (type, hostname, publicIn, publicOut, privateIn, privateOut, pool)[default: publicOut]"))
	cobraCmd.Flags().BoolVar(&thisCmd.Virtual, "virtual", false, T("Show only the bandwidth summary for each virtual server"))
	cobraCmd.Flags().BoolVar(&thisCmd.Server, "server", false, T("Show only the bandwidth summary for each hardware server "))
	cobraCmd.Flags().BoolVar(&thisCmd.Pool, "pool", false, ("Show only the bandwidth pool summary."))
	thisCmd.Command = cobraCmd
	return thisCmd

}

func (cmd *BandwidthCommand) Run(args []string) error {
	outputFormat := cmd.GetOutputFlag()
	var err error
	sortBy := cmd.SortBy
	if sortBy == "" {
		sortBy = "publicout"
	} else {
		sortBy = strings.ToLower(sortBy)
		sortByOptions := []string{"type", "hostname", "publicin", "publicout", "privatein", "privateout", "pool"}
		if !utils.WordInList(sortByOptions, sortBy) {
			return errors.NewInvalidUsageError(T("Invalid --sortBy option."))
		}
	}

	var endDate time.Time
	if cmd.End != "" {
		date := cmd.End
		endDate, err = time.Parse("2006-01-02 15:04:05", date)
		if err != nil {
			endDate, err = time.Parse("2006-01-02", date)
			if err != nil {
				return errors.NewInvalidUsageError(T("Invalid format date to --end."))
			}
		}
	} else {
		endDate = time.Now()
	}

	var startDate time.Time
	if cmd.Start != "" {
		date := cmd.Start
		startDate, err = time.Parse("2006-01-02 15:04:05", date)
		if err != nil {
			startDate, err = time.Parse("2006-01-02", date)
			if err != nil {
				return errors.NewInvalidUsageError(T("Invalid format date to --start."))
			}
		}
	} else {
		startDate = endDate.AddDate(0, -1, 0)
	}

	if !endDate.After(startDate) {
		return errors.NewInvalidUsageError(T("End Date must be greater than Start Date."))
	}

	cmd.UI.Print(T("Generating bandwidth report for {{.startDate}} to {{.endDate}}",
		map[string]interface{}{"startDate": startDate.Format("2006-01-02 15:04:05"), "endDate": endDate.Format("2006-01-02 15:04:05")}))

	metricObjects := []metricObject{}
	if !cmd.Virtual && !cmd.Server && !cmd.Pool {
		metricObjects, err = getVirtualBandwidth(cmd, metricObjects, datatypes.Time{Time: startDate}, datatypes.Time{Time: endDate})
		if err != nil {
			return err
		}
		metricObjects, err = getHardwareBandwidth(cmd, metricObjects, datatypes.Time{Time: startDate}, datatypes.Time{Time: endDate})
		if err != nil {
			return err
		}
		metricObjects, err = getPoolBandwidth(cmd, metricObjects, datatypes.Time{Time: startDate}, datatypes.Time{Time: endDate})
		if err != nil {
			return err
		}
	} else {
		if cmd.Virtual {
			metricObjects, err = getVirtualBandwidth(cmd, metricObjects, datatypes.Time{Time: startDate}, datatypes.Time{Time: endDate})
			if err != nil {
				return err
			}
		}
		if cmd.Server {
			metricObjects, err = getHardwareBandwidth(cmd, metricObjects, datatypes.Time{Time: startDate}, datatypes.Time{Time: endDate})
			if err != nil {
				return err
			}
		}
		if cmd.Pool {
			metricObjects, err = getPoolBandwidth(cmd, metricObjects, datatypes.Time{Time: startDate}, datatypes.Time{Time: endDate})
			if err != nil {
				return err
			}
		}
	}

	tableRows := getTableRows(metricObjects)

	//sort metricObjects array
	switch sortBy {
	case "type":
		sort.Sort(ByType(tableRows))
	case "hostname":
		sort.Sort(ByHostname(tableRows))
	case "publicin":
		sort.Sort(ByPublicIn(tableRows))
	case "publicout":
		sort.Sort(ByPublicOut(tableRows))
	case "privatein":
		sort.Sort(ByPrivateIn(tableRows))
	case "privateout":
		sort.Sort(ByPrivateOut(tableRows))
	case "pool":
		sort.Sort(ByPool(tableRows))
	default:
		sort.Sort(ByPublicOut(tableRows))
	}

	table := cmd.UI.Table([]string{T("type"), T("Hostname"), T("publicIn"), T("publicOut"), T("privateIn"), T("privateOut"), T("pool")})
	for _, row := range tableRows {
		table.Add(
			row.typeDevice,
			row.hostname,
			fmt.Sprintf("%.2f GB", float64(row.publicIn)/1000000000),
			fmt.Sprintf("%.2f GB", float64(row.publicOut)/1000000000),
			fmt.Sprintf("%.2f GB", float64(row.privateIn)/1000000000),
			fmt.Sprintf("%.2f GB", float64(row.privateOut)/1000000000),
			row.pool,
		)
	}

	utils.PrintTable(cmd.UI, table, outputFormat)
	return nil
}

// interface to sort by type
type ByType []tableRow

func (a ByType) Len() int           { return len(a) }
func (a ByType) Less(i, j int) bool { return a[i].typeDevice < a[j].typeDevice }
func (a ByType) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// interface to sort by hostname
type ByHostname []tableRow

func (a ByHostname) Len() int           { return len(a) }
func (a ByHostname) Less(i, j int) bool { return a[i].hostname < a[j].hostname }
func (a ByHostname) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// interface to sort by publicIn
type ByPublicIn []tableRow

func (a ByPublicIn) Len() int           { return len(a) }
func (a ByPublicIn) Less(i, j int) bool { return a[i].publicIn < a[j].publicIn }
func (a ByPublicIn) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// interface to sort by PublicOut
type ByPublicOut []tableRow

func (a ByPublicOut) Len() int           { return len(a) }
func (a ByPublicOut) Less(i, j int) bool { return a[i].publicOut < a[j].publicOut }
func (a ByPublicOut) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// interface to sort by privateIn
type ByPrivateIn []tableRow

func (a ByPrivateIn) Len() int           { return len(a) }
func (a ByPrivateIn) Less(i, j int) bool { return a[i].privateIn < a[j].privateIn }
func (a ByPrivateIn) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// interface to sort by privateOut
type ByPrivateOut []tableRow

func (a ByPrivateOut) Len() int           { return len(a) }
func (a ByPrivateOut) Less(i, j int) bool { return a[i].privateOut < a[j].privateOut }
func (a ByPrivateOut) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// interface to sort by pool
type ByPool []tableRow

func (a ByPool) Len() int           { return len(a) }
func (a ByPool) Less(i, j int) bool { return a[i].pool < a[j].pool }
func (a ByPool) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// Return Bandwidth Summary of each virtual guest on account
func getVirtualBandwidth(cmd *BandwidthCommand, metricObjects []metricObject, start datatypes.Time, end datatypes.Time) ([]metricObject, error) {

	virtualGuests, err := cmd.ReportManager.GetVirtualGuests("")
	if err != nil {
		return metricObjects, errors.NewAPIError(T("Failed to get virtual guests on your account.\n"), err.Error(), 2)
	}

	validTypes := []datatypes.Container_Metric_Data_Type{
		datatypes.Container_Metric_Data_Type{
			KeyName:     sl.String("PUBLICIN_NET_OCTET"),
			Name:        sl.String("publicIn_net_octet"),
			SummaryType: sl.String("sum"),
		},
		datatypes.Container_Metric_Data_Type{
			KeyName:     sl.String("PUBLICOUT_NET_OCTET"),
			Name:        sl.String("publicOut_net_octet"),
			SummaryType: sl.String("sum"),
		},
		datatypes.Container_Metric_Data_Type{
			KeyName:     sl.String("PRIVATEIN_NET_OCTET"),
			Name:        sl.String("privateIn_net_octet"),
			SummaryType: sl.String("sum"),
		},
		datatypes.Container_Metric_Data_Type{
			KeyName:     sl.String("PRIVATEOUT_NET_OCTET"),
			Name:        sl.String("privateOut_net_octet"),
			SummaryType: sl.String("sum"),
		},
	}

	progressBar := utils.ProgressBar(T("Calculating for virtual"), len(virtualGuests))

	for _, virtualGuest := range virtualGuests {
		if virtualGuest.MetricTrackingObjectId != nil {
			metricTrackingSummary, err := cmd.ReportManager.GetMetricTrackingSummaryData(*virtualGuest.MetricTrackingObjectId, start, end, validTypes)
			if err != nil {
				log.Println(T("Failed to get metric tracking summary of Object with Id {{.MetricTrackingObjectId}}.",
					map[string]interface{}{"MetricTrackingObjectId": *virtualGuest.MetricTrackingObjectId}) + err.Error())
				continue
			}
			pool := "-"
			if virtualGuest.VirtualRack != nil {
				if *virtualGuest.VirtualRack.BandwidthAllotmentTypeId == 2 {
					pool = *virtualGuest.VirtualRack.Name
				}
			}
			virtualGuestMetricObject := metricObject{
				id:         *virtualGuest.Id,
				typeDevice: "virtual",
				name:       *virtualGuest.Hostname,
				pool:       pool,
				data:       metricTrackingSummary,
			}
			metricObjects = append(metricObjects, virtualGuestMetricObject)
		}
		progressBar.Add()
	}
	return metricObjects, nil
}

// Return Bandwidth Summary of each hardware server on account
func getHardwareBandwidth(cmd *BandwidthCommand, metricObjects []metricObject, start datatypes.Time, end datatypes.Time) ([]metricObject, error) {

	hardwareServers, err := cmd.ReportManager.GetHardwareServers("")
	if err != nil {
		return metricObjects, errors.NewAPIError(T("Failed to get hardware servers on your account.\n"), err.Error(), 2)
	}

	validTypes := []datatypes.Container_Metric_Data_Type{
		datatypes.Container_Metric_Data_Type{
			KeyName:     sl.String("PUBLICIN"),
			Name:        sl.String("publicIn"),
			SummaryType: sl.String("counter"),
		},
		datatypes.Container_Metric_Data_Type{
			KeyName:     sl.String("PUBLICOUT"),
			Name:        sl.String("publicOut"),
			SummaryType: sl.String("counter"),
		},
		datatypes.Container_Metric_Data_Type{
			KeyName:     sl.String("PRIVATEIN"),
			Name:        sl.String("privateIn"),
			SummaryType: sl.String("counter"),
		},
		datatypes.Container_Metric_Data_Type{
			KeyName:     sl.String("PRIVATEOUT"),
			Name:        sl.String("privateOut"),
			SummaryType: sl.String("counter"),
		},
	}

	progressBar := utils.ProgressBar(T("Calculating for hardware"), len(hardwareServers))
	for _, hardware := range hardwareServers {
		if hardware.MetricTrackingObject != nil {
			if hardware.MetricTrackingObject.Id != nil {
				metricTrackingSummary, err := cmd.ReportManager.GetMetricTrackingSummaryData(*hardware.MetricTrackingObject.Id, start, end, validTypes)
				if err != nil {
					log.Println(T("Failed to get metric tracking summary of Object with Id {{.MetricTrackingObjectId}}.",
						map[string]interface{}{"MetricTrackingObjectId": *hardware.MetricTrackingObject.Id}) + err.Error())
					continue
				}
				pool := "-"
				if hardware.VirtualRack != nil {
					if *hardware.VirtualRack.BandwidthAllotmentTypeId == 2 {
						pool = *hardware.VirtualRack.Name
					}
				}
				virtualGuestMetricObject := metricObject{
					id:         *hardware.Id,
					typeDevice: "hardware",
					name:       *hardware.Hostname,
					pool:       pool,
					data:       metricTrackingSummary,
				}
				metricObjects = append(metricObjects, virtualGuestMetricObject)
			}
		}
		progressBar.Add()
	}
	return metricObjects, nil
}

// Return Bandwidth Summary of each pool on account
func getPoolBandwidth(cmd *BandwidthCommand, metricObjects []metricObject, start datatypes.Time, end datatypes.Time) ([]metricObject, error) {

	pools, err := cmd.ReportManager.GetVirtualDedicatedRacks("")
	if err != nil {
		return metricObjects, errors.NewAPIError(T("Failed to get virtual dedicated racks on your account.\n"), err.Error(), 2)
	}

	validTypes := []datatypes.Container_Metric_Data_Type{
		datatypes.Container_Metric_Data_Type{
			KeyName:     sl.String("PUBLICIN"),
			Name:        sl.String("publicIn"),
			SummaryType: sl.String("sum"),
		},
		datatypes.Container_Metric_Data_Type{
			KeyName:     sl.String("PUBLICOUT"),
			Name:        sl.String("publicOut"),
			SummaryType: sl.String("sum"),
		},
		datatypes.Container_Metric_Data_Type{
			KeyName:     sl.String("PRIVATEIN"),
			Name:        sl.String("privateIn"),
			SummaryType: sl.String("sum"),
		},
		datatypes.Container_Metric_Data_Type{
			KeyName:     sl.String("PRIVATEOUT"),
			Name:        sl.String("privateOut"),
			SummaryType: sl.String("sum"),
		},
	}

	progressBar := utils.ProgressBar(T("Calculating for bandwidth pools"), len(pools))
	for _, pool := range pools {
		if pool.MetricTrackingObjectId != nil {
			metricTrackingSummary, err := cmd.ReportManager.GetMetricTrackingSummaryData(*pool.MetricTrackingObjectId, start, end, validTypes)
			if err != nil {
				log.Println(T("Failed to get metric tracking summary of Object with Id {{.MetricTrackingObjectId}}.",
					map[string]interface{}{"MetricTrackingObjectId": *pool.MetricTrackingObjectId}) + err.Error())
				continue
			}
			virtualGuestMetricObject := metricObject{
				id:         *pool.Id,
				typeDevice: "pool",
				name:       *pool.Name,
				pool:       "-",
				data:       metricTrackingSummary,
			}
			metricObjects = append(metricObjects, virtualGuestMetricObject)
		}
		progressBar.Add()
	}
	return metricObjects, nil
}

func getTableRows(metricObjects []metricObject) []tableRow {
	tableRows := []tableRow{}
	row := tableRow{}
	for _, metricObject := range metricObjects {
		pubIn := getTotalByTypeKey("publicIn_net_octet", metricObject.data)
		pubOut := getTotalByTypeKey("publicOut_net_octet", metricObject.data)
		privateIn := getTotalByTypeKey("privateIn_net_octet", metricObject.data)
		privateOut := getTotalByTypeKey("privateOut_net_octet", metricObject.data)
		row = tableRow{
			typeDevice: metricObject.typeDevice,
			hostname:   metricObject.name,
			publicIn:   pubIn,
			publicOut:  pubOut,
			privateIn:  privateIn,
			privateOut: privateOut,
			pool:       metricObject.pool,
		}
		tableRows = append(tableRows, row)
	}
	return tableRows
}

// Return the total of an specific traffic type
func getTotalByTypeKey(typeKey string, metricObjectData []datatypes.Metric_Tracking_Object_Data) int {
	total := 0
	for _, data := range metricObjectData {
		if *data.Type == typeKey {
			total = total + int(*data.Counter)
		}
	}
	return total
}
