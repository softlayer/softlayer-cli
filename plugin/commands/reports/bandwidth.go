package reports

import (
	"fmt"
	"sort"
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type BandwidthCommand struct {
	UI            terminal.UI
	ReportManager managers.ReportManager
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

func NewBandwidthCommand(ui terminal.UI, reportManager managers.ReportManager) (cmd *BandwidthCommand) {
	return &BandwidthCommand{
		UI:            ui,
		ReportManager: reportManager,
	}
}

func (cmd *BandwidthCommand) Run(c *cli.Context) error {
	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	sortBy := ""
	if !c.IsSet("sortby") {
		sortBy = "hostname"
	} else {
		sortBy = c.String("sortby")
		if sortBy != "type" && sortBy != "hostname" && sortBy != "publicIn" && sortBy != "publicOut" && sortBy != "privateIn" &&
			sortBy != "privateOut" && sortBy != "pool" {
			return errors.NewInvalidUsageError(T("Invalid --sortBy option."))
		}
	}

	var endDate time.Time
	if c.IsSet("end") {
		date := c.String("end")
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
	if c.IsSet("start") {
		date := c.String("start")
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

	cmd.UI.Print(T("Generating bandwidth report for {{.startDate}} to {{.endDate}}",
		map[string]interface{}{"startDate": startDate.Format("2006-01-02 15:04:05"), "endDate": endDate.Format("2006-01-02 15:04:05")}))

	metricObjects := []metricObject{}
	if !c.IsSet("virtual") && !c.IsSet("server") && !c.IsSet("pool") {
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
		if c.IsSet("virtual") && c.Bool("virtual") {
			metricObjects, err = getVirtualBandwidth(cmd, metricObjects, datatypes.Time{Time: startDate}, datatypes.Time{Time: endDate})
			if err != nil {
				return err
			}
		}
		if c.IsSet("server") && c.Bool("server") {
			metricObjects, err = getHardwareBandwidth(cmd, metricObjects, datatypes.Time{Time: startDate}, datatypes.Time{Time: endDate})
			if err != nil {
				return err
			}
		}
		if c.IsSet("pool") && c.Bool("pool") {
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
	case "publicIn":
		sort.Sort(ByPublicIn(tableRows))
	case "publicOut":
		sort.Sort(ByPublicOut(tableRows))
	case "privateIn":
		sort.Sort(ByPrivateIn(tableRows))
	case "privateOut":
		sort.Sort(ByPrivateOut(tableRows))
	case "pool":
		sort.Sort(ByPool(tableRows))
	}

	table := cmd.UI.Table([]string{T("type"), T("hostname"), T("publicIn"), T("publicOut"), T("privateIn"), T("privateOut"), T("pool")})
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

//Return Bandwidth Summary of each virtual guest on account
func getVirtualBandwidth(cmd *BandwidthCommand, metricObjects []metricObject, start datatypes.Time, end datatypes.Time) ([]metricObject, error) {

	virtualGuests, err := cmd.ReportManager.GetVirtualGuests("")
	if err != nil {
		return metricObjects, cli.NewExitError(T("Failed to get virtual guests on your account.\n")+err.Error(), 2)
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

	cmd.UI.Print(T("Calculating for virtual"))

	for _, virtualGuest := range virtualGuests {
		if virtualGuest.MetricTrackingObjectId != nil {
			metricTrackingSummary, err := cmd.ReportManager.GetMetricTrackingSummaryData(*virtualGuest.MetricTrackingObjectId, start, end, validTypes)
			if err != nil {
				return metricObjects, cli.NewExitError(T("Failed to get metric tracking summary of Object with Id {{.MetricTrackingObjectId}}.",
					map[string]interface{}{"MetricTrackingObjectId": *virtualGuest.MetricTrackingObjectId})+err.Error(), 2)
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
	}
	return metricObjects, nil
}

//Return Bandwidth Summary of each hardware server on account
func getHardwareBandwidth(cmd *BandwidthCommand, metricObjects []metricObject, start datatypes.Time, end datatypes.Time) ([]metricObject, error) {

	hardwareServers, err := cmd.ReportManager.GetHardwareServers("")
	if err != nil {
		return metricObjects, cli.NewExitError(T("Failed to get hardware servers on your account.\n")+err.Error(), 2)
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

	cmd.UI.Print(T("Calculating for hardware"))

	for _, hardware := range hardwareServers {
		if hardware.MetricTrackingObject != nil {
			if hardware.MetricTrackingObject.Id != nil {
				metricTrackingSummary, err := cmd.ReportManager.GetMetricTrackingSummaryData(*hardware.MetricTrackingObject.Id, start, end, validTypes)
				if err != nil {
					return metricObjects, cli.NewExitError(T("Failed to get metric tracking summary of Object with Id {{.MetricTrackingObjectId}}.",
						map[string]interface{}{"MetricTrackingObjectId": *hardware.MetricTrackingObject.Id})+err.Error(), 2)
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
	}
	return metricObjects, nil
}

//Return Bandwidth Summary of each pool on account
func getPoolBandwidth(cmd *BandwidthCommand, metricObjects []metricObject, start datatypes.Time, end datatypes.Time) ([]metricObject, error) {

	pools, err := cmd.ReportManager.GetVirtualDedicatedRacks("")
	if err != nil {
		return metricObjects, cli.NewExitError(T("Failed to get virtual dedicated racks on your account.\n")+err.Error(), 2)
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

	cmd.UI.Print(T("Calculating for bandwidth pools"))

	for _, pool := range pools {
		if pool.MetricTrackingObjectId != nil {
			metricTrackingSummary, err := cmd.ReportManager.GetMetricTrackingSummaryData(*pool.MetricTrackingObjectId, start, end, validTypes)
			if err != nil {
				return metricObjects, cli.NewExitError(T("Failed to get metric tracking summary of Object with Id {{.MetricTrackingObjectId}}.",
					map[string]interface{}{"MetricTrackingObjectId": *pool.MetricTrackingObjectId})+err.Error(), 2)
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

//Return the total of an specific traffic type
func getTotalByTypeKey(typeKey string, metricObjectData []datatypes.Metric_Tracking_Object_Data) int {
	total := 0
	for _, data := range metricObjectData {
		if *data.Type == typeKey {
			total = total + int(*data.Counter)
		}
	}
	return total
}

func ReportBandwidthMetaData() cli.Command {
	return cli.Command{
		Category:    "report",
		Name:        "bandwidth",
		Description: T("Bandwidth report for every pool/server."),
		Usage: T(`${COMMAND_NAME} sl report bandwidth

EXAMPLE: 
   ${COMMAND_NAME} sl report bandwidth
   ${COMMAND_NAME} sl report bandwidth --server --start 2022-06-07 --end 2022-06-08
   ${COMMAND_NAME} sl report bandwidth --start 2022-06-07 --end 2022-06-08 --sortby privateOut`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "start",
				Usage: T("datetime in the format 'YYYY-MM-DD' or 'YYYY-MM-DD HH:MM:SS'"),
			},
			cli.StringFlag{
				Name:  "end",
				Usage: T("datetime in the format 'YYYY-MM-DD' or 'YYYY-MM-DD HH:MM:SS'"),
			},
			cli.StringFlag{
				Name:  "sortby",
				Usage: T("Column to sort by [default: hostname]"),
			},
			cli.BoolFlag{
				Name:  "virtual",
				Usage: T("Show only the bandwidth summary for each virtual server"),
			},
			cli.BoolFlag{
				Name:  "server",
				Usage: T("Show only the bandwidth summary for each hardware server "),
			},
			cli.BoolFlag{
				Name:  "pool",
				Usage: T("Show only the bandwidth pool summary."),
			},
			metadata.OutputFlag(),
		},
	}
}
