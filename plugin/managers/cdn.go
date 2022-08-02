package managers

import (
	"strconv"
	"time"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
)

type CdnManager interface {
	GetNetworkCdnMarketplaceConfigurationMapping(mask string) ([]datatypes.Container_Network_CdnMarketplace_Configuration_Mapping, error)
	GetDetailCDN(uniqueId int, mask string) (datatypes.Container_Network_CdnMarketplace_Configuration_Mapping, error)
	GetUsageMetrics(uniqueId int, history int, mask string) (datatypes.Container_Network_CdnMarketplace_Metrics, error)
}

type cdnManager struct {
	CdnService services.Network_CdnMarketplace_Configuration_Mapping
	Session    *session.Session
}

func NewCdnManager(session *session.Session) *cdnManager {
	return &cdnManager{
		CdnService: services.GetNetworkCdnMarketplaceConfigurationMappingService(session),
		Session:    session,
	}
}

/*
SOAP API will return all domains for a particular customer.
https://sldn.softlayer.com/reference/services/SoftLayer_Network_CdnMarketplace_Configuration_Mapping/listDomainMappings/
*/
func (a cdnManager) GetNetworkCdnMarketplaceConfigurationMapping(mask string) ([]datatypes.Container_Network_CdnMarketplace_Configuration_Mapping, error) {
	return a.CdnService.Mask(mask).ListDomainMappings()
}

/*
SOAP API will return the domain mapping based on the uniqueId.
https://sldn.softlayer.com/reference/services/SoftLayer_Network_CdnMarketplace_Configuration_Mapping/listDomainMappingByUniqueId/
*/
func (a cdnManager) GetDetailCDN(uniqueId int, mask string) (datatypes.Container_Network_CdnMarketplace_Configuration_Mapping, error) {
	cdnId := strconv.Itoa(uniqueId)
	cdn, err := a.CdnService.Mask(mask).ListDomainMappingByUniqueId(&cdnId)
	if err != nil {
		return datatypes.Container_Network_CdnMarketplace_Configuration_Mapping{}, err
	}
	return cdn[0], nil
}

/*
SOAP API will return the domain mapping based on the uniqueId.
https://sldn.softlayer.com/reference/services/SoftLayer_Network_CdnMarketplace_Configuration_Mapping/listDomainMappingByUniqueId/
*/
func (a cdnManager) GetUsageMetrics(uniqueId int, history int, mask string) (datatypes.Container_Network_CdnMarketplace_Metrics, error) {
	NetworkCdnMarketplaceMetricsService := services.GetNetworkCdnMarketplaceMetricsService(a.Session)

	startDate := int(time.Now().AddDate(0, 0, (history * -1)).Unix())
	endDate := int(time.Now().Unix())
	frequency := "aggregate"
	cdnId := strconv.Itoa(uniqueId)

	cdn, err := NetworkCdnMarketplaceMetricsService.Mask(mask).GetMappingUsageMetrics(&cdnId, &startDate, &endDate, &frequency)
	if err != nil {
		return datatypes.Container_Network_CdnMarketplace_Metrics{}, err
	}

	return cdn[0], nil
}
