package managers

import (
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
)

type ReportManager interface {
	GetVirtualGuests(mask string) ([]datatypes.Virtual_Guest, error)
	GetHardwareServers(mask string) ([]datatypes.Hardware, error)
	GetVirtualDedicatedRacks(mask string) ([]datatypes.Network_Bandwidth_Version1_Allotment, error)
	GetMetricTrackingSummaryData(metricTrackingObjectID int, startDateTime datatypes.Time, endDateTime datatypes.Time, validTypes []datatypes.Container_Metric_Data_Type) ([]datatypes.Metric_Tracking_Object_Data, error)
}

type reportManager struct {
	AccountService              services.Account
	MetricTrackingObjectService services.Metric_Tracking_Object
}

func NewReportManager(session *session.Session) *reportManager {
	return &reportManager{
		services.GetAccountService(session),
		services.GetMetricTrackingObjectService(session),
	}
}

//Get virtual guests
//mask: object mask
func (re reportManager) GetVirtualGuests(mask string) ([]datatypes.Virtual_Guest, error) {
	if mask == "" {
		mask = "mask[id,hostname,metricTrackingObjectId,virtualRack[name,id,bandwidthAllotmentTypeId]]"
	}
	return re.AccountService.Mask(mask).GetVirtualGuests()
}

//Get hardwares
//mask: object mask
func (re reportManager) GetHardwareServers(mask string) ([]datatypes.Hardware, error) {
	if mask == "" {
		mask = "mask[id,hostname,metricTrackingObject.id,virtualRack[name,id,bandwidthAllotmentTypeId]]"
	}
	return re.AccountService.Mask(mask).GetHardware()
}

//Get virtual dedicated racks
//mask: object mask
func (re reportManager) GetVirtualDedicatedRacks(mask string) ([]datatypes.Network_Bandwidth_Version1_Allotment, error) {
	if mask == "" {
		mask = "mask[id,name,metricTrackingObjectId]"
	}
	return re.AccountService.Mask(mask).GetVirtualDedicatedRacks()
}

//Get metric tracking object
//id: Metric Tracking Object Id
//mask: object mask
func (re reportManager) GetMetricTrackingSummaryData(metricTrackingObjectID int, startDateTime datatypes.Time, endDateTime datatypes.Time, validTypes []datatypes.Container_Metric_Data_Type) ([]datatypes.Metric_Tracking_Object_Data, error) {
	summaryPeriod := 86400
	return re.MetricTrackingObjectService.Id(metricTrackingObjectID).GetSummaryData(
		&startDateTime,
		&endDateTime,
		validTypes,
		&summaryPeriod,
	)
}
