package virtual

import (
	"time"
	"fmt"
	"sort"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"

	"github.com/softlayer/softlayer-go/datatypes"
)

type BandwidthCommand struct {
	UI                   terminal.UI
	VirtualServerManager managers.VirtualServerManager
}

type SummaryDataType struct {
	Name 		string
	Sum 		float64
	Maximum 	float64
	MaxDate 	string

}

func NewBandwidthCommand(ui terminal.UI, virtualServerManager managers.VirtualServerManager) (cmd *BandwidthCommand) {
	return &BandwidthCommand{
		UI:                   ui,
		VirtualServerManager: virtualServerManager,
	}
}

func (cmd *BandwidthCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	VsID, err := utils.ResolveVirtualGuestId(c.Args()[0])
	if err != nil {
		return errors.NewInvalidSoftlayerIdInputError("Virtual server ID")
	}
	
	var start, end string
	var startDate, endDate time.Time

	if c.IsSet("start") {
		start = c.String("start")
		startDate, err = time.Parse(GetDateFormat(start), start)
		if err != nil {
			return errors.NewInvalidUsageError("Invalid start date: " + err.Error())
		}
	} else {
		startDate = time.Now()
	}
	if c.IsSet("end") {
		end = c.String("end")
		endDate, err = time.Parse(GetDateFormat(end), end)
		if err != nil {
			return errors.NewInvalidUsageError("Invalid end date: " + err.Error())
		}
	} else {
		endDate = startDate.AddDate(0, -1, 0)
	}
	
	rollupSeconds := 3600
	if c.IsSet("rollup") {
		rollupSeconds = c.Int("rollup")
	}
	// cmd.UI.Say(fmt.Sprintf("FORMAT: %v, Start: %v (%v), End: %v (%v)\n", GetDateFormat(start), startDate, start, endDate, end))
	bandwidthData, err := cmd.VirtualServerManager.GetBandwidthData(VsID, startDate, endDate, rollupSeconds)
	if err != nil {
		fmt.Printf("ERR: %v", err)
		return err 
	}
	// cmd.UI.Say(fmt.Sprintf("%+v", bandwidthData))
	
	summaryTable, bandwidthTable := BuildOutputTable(bandwidthData, cmd.UI)
	summaryTable.Print()
	if !c.IsSet("quite") {
		bandwidthTable.Print()	
	}
	

	
	return nil
}	

// Borrowed from https://stackoverflow.com/questions/56374333/flexible-date-time-parsing-in-go
func GetDateFormat(inputDate string) string {
	dateFormat := "2006-01-02T15:04-07:00"
	var layout string
    if len(inputDate) < len(dateFormat) {
        layout = dateFormat[:len(inputDate)]
    } else {
        layout = dateFormat
    }
    return layout
}

func BuildOutputTable(trackingData []datatypes.Metric_Tracking_Object_Data, ui terminal.UI) (terminal.Table, terminal.Table) {

	formattedData := make(map[string]map[string]float64)
	summaryData := map[string]SummaryDataType{
		"publicIn_net_octet": SummaryDataType{Name: "Pub In", Maximum: 0.0, Sum: 0.0},
		"publicOut_net_octet": SummaryDataType{Name: "Pub Out", Maximum: 0.0, Sum: 0.0},
		"privateIn_net_octet": SummaryDataType{Name: "Pri In", Maximum: 0.0, Sum: 0.0},
		"privateOut_net_octet": SummaryDataType{Name: "Pri Out", Maximum: 0.0, Sum: 0.0},
	}
	// var sumPubIn, sumPubOut, sumPriIn, sumPriOut float64
	summaryTable := ui.Table([]string{"Type", "Sum GB", "Average MBps", "MAX GB", "Max Date"})
	bandwidthTable := ui.Table([]string{"Date", "Pub In", "Pub Out", "Pri In", "Pri Out"})

	if trackingData == nil || len(trackingData) < 1 {
		summaryTable.Add(T("No data"), "-", "-", "-", "-")
		bandwidthTable.Add(T("No data"), "-", "-", "-", "-")
		return summaryTable, bandwidthTable
	}

	// Groups the data by date, instead of individual datapoints
	for _, point := range trackingData {
		
		theTime := point.DateTime.Format("2006-01-02 15:04")
		if formattedData[theTime] == nil {
			formattedData[theTime] = make(map[string]float64)
		}
		theType := *point.Type
		// value = round(float(point['counter']) / 2 ** 20, 4)
		// Conversion from byte to MB
		formattedData[theTime][theType] = float64(*point.Counter) / 1048576
	
	}

	// This sorts the dates because even though the API returns them sorted, go seems to put them
	// in non-sorted orders when building the formattedData map.
	dateKeys := make([]string, 0, len(formattedData))
	for date, _ := range formattedData {
		dateKeys = append(dateKeys, date)
	}
	sort.Strings(dateKeys)

	for _, date := range dateKeys {
		values := formattedData[date]
		bandwidthTable.Add(
			date,
			fmt.Sprintf("%.4f", values["publicIn_net_octet"]   / 1024),
			fmt.Sprintf("%.4f", values["publicOut_net_octet"]  / 1024),
			fmt.Sprintf("%.4f", values["privateIn_net_octet"]  / 1024),
			fmt.Sprintf("%.4f", values["privateOut_net_octet"] / 1024),
		)
		// Updates for the Summary Table here
		for keyName, summary := range summaryData {
			summary.Sum += values[keyName]
			if summary.Maximum < values[keyName] {
				summary.Maximum = values[keyName]
				summary.MaxDate = date
			}
			summaryData[keyName] = summary
		}
	}

	// Builds summary table
	for _, summary := range summaryData {
		summaryTable.Add(
			summary.Name,
			fmt.Sprintf("%.4f", summary.Sum / 1024),
			fmt.Sprintf("%.4f", summary.Sum / float64(len(trackingData))),
			fmt.Sprintf("%.4f", summary.Maximum / 1024),
			summary.MaxDate,
		)
	}
	
	return summaryTable, bandwidthTable

}