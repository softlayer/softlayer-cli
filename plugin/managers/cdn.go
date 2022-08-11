package managers

import (
	"strconv"
	"strings"
	"time"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
)

type CdnManager interface {
	GetNetworkCdnMarketplaceConfigurationMapping() ([]datatypes.Container_Network_CdnMarketplace_Configuration_Mapping, error)
	GetDetailCDN(uniqueId int, mask string) (datatypes.Container_Network_CdnMarketplace_Configuration_Mapping, error)
	GetUsageMetrics(uniqueId int, history int, mask string) (datatypes.Container_Network_CdnMarketplace_Metrics, error)
	EditCDN(uniqueId int, header string, httpPort int, httpsPort int, origin string, respectHeaders string, cache string, cacheDescription string, performanceConfiguration string) (datatypes.Container_Network_CdnMarketplace_Configuration_Mapping, error)
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
func (a cdnManager) GetNetworkCdnMarketplaceConfigurationMapping() ([]datatypes.Container_Network_CdnMarketplace_Configuration_Mapping, error) {
	return a.CdnService.ListDomainMappings()
}

/*
Gets the domain mapping based on the uniqueId.
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
Gets the total number of predetermined statistics for direct display for the given mapping.
https://sldn.softlayer.com/reference/services/SoftLayer_Network_CdnMarketplace_Metrics/getMappingUsageMetrics/
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

/*
Updates the Domain Mapping identified by the Unique Id. Following fields are allowed to be changed: originHost, HttpPort/HttpsPort, RespectHeaders, ServeStale

Additionally, bucketName and fileExtension if OriginType is Object Store
https://sldn.softlayer.com/reference/services/SoftLayer_Network_CdnMarketplace_Configuration_Mapping/updateDomainMapping/
*/
func (a cdnManager) EditCDN(uniqueId int, header string, httpPort int, httpsPort int, origin string, respectHeaders string, cache string, cacheDescription string, performanceConfiguration string) (datatypes.Container_Network_CdnMarketplace_Configuration_Mapping, error) {

	cdnId := strconv.Itoa(uniqueId)
	cdn, err := a.CdnService.Mask(mask).ListDomainMappingByUniqueId(&cdnId)
	if err != nil {
		return datatypes.Container_Network_CdnMarketplace_Configuration_Mapping{}, err
	}
	config := datatypes.Container_Network_CdnMarketplace_Configuration_Input{
		UniqueId:   cdn[0].UniqueId,
		OriginType: cdn[0].OriginType,
		Protocol:   cdn[0].Protocol,
		Path:       cdn[0].Path,
		VendorName: cdn[0].VendorName,
		Cname:      cdn[0].Cname,
		Domain:     cdn[0].Domain,
		Origin:     cdn[0].OriginHost,
		Header:     cdn[0].Header,
	}
	if cdn[0].HttpPort != nil {
		config.HttpPort = cdn[0].HttpPort
	}
	if cdn[0].HttpsPort != nil {
		config.HttpsPort = cdn[0].HttpsPort
	}

	if header != "" {
		config.Header = &header
	}

	if httpPort != 0 {
		config.HttpPort = &httpPort
	}

	if httpsPort != 0 {
		config.HttpsPort = &httpsPort
	}
	if origin != "" {
		config.Origin = &origin
	}
	if respectHeaders != "" {
		config.RespectHeaders = &respectHeaders
	}
	if cache != "" {
		if cacheDescription != "" {
			cache := strings.Split(cache, "-")
			value := cache[0] + ": " + cacheDescription
			config.CacheKeyQueryRule = &value
		} else {
			config.CacheKeyQueryRule = &cache
		}
	}
	if performanceConfiguration != "" {
		config.PerformanceConfiguration = &performanceConfiguration
	}

	cdn, err = a.CdnService.UpdateDomainMapping(&config)
	if err != nil {
		return datatypes.Container_Network_CdnMarketplace_Configuration_Mapping{}, err
	}
	return cdn[0], nil
}
