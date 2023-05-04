package managers

import (
	"strconv"
	"strings"
	"time"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
)

type CdnManager interface {
	GetNetworkCdnMarketplaceConfigurationMapping() ([]datatypes.Container_Network_CdnMarketplace_Configuration_Mapping, error)
	DeleteCDN(uniqueId string) ([]datatypes.Container_Network_CdnMarketplace_Configuration_Mapping, error)
	GetDetailCDN(uniqueId int, mask string) (datatypes.Container_Network_CdnMarketplace_Configuration_Mapping, error)
	GetUsageMetrics(uniqueId int, history int, mask string) (datatypes.Container_Network_CdnMarketplace_Metrics, error)
	EditCDN(uniqueId int, header string, httpPort int, httpsPort int, origin string, respectHeaders string, cache string, cacheDescription string, performanceConfiguration string) (datatypes.Container_Network_CdnMarketplace_Configuration_Mapping, error)
	CreateCdn(hostname string, originHost string, originType string, http int, https int, bucketName string, cname string, header string, path string, ssl string) ([]datatypes.Container_Network_CdnMarketplace_Configuration_Mapping, error)
	OriginAddCdn(uniqueId string, header string, path string, originHost string, originType string, http int, https int, cacheKey string, optimize string, dynamicPath string, dynamicPrefetch bool, dynamicCompression bool, bucketName string, fileExtension string) ([]datatypes.Container_Network_CdnMarketplace_Configuration_Mapping_Path, error)
}

type cdnManager struct {
	CdnService     services.Network_CdnMarketplace_Configuration_Mapping
	CdnPathService services.Network_CdnMarketplace_Configuration_Mapping_Path
	Session        *session.Session
}

func NewCdnManager(session *session.Session) *cdnManager {
	return &cdnManager{
		CdnService:     services.GetNetworkCdnMarketplaceConfigurationMappingService(session),
		CdnPathService: services.GetNetworkCdnMarketplaceConfigurationMappingPathService(session),
		Session:        session,
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
Delete CDN domain mapping.
https://sldn.softlayer.com/reference/services/SoftLayer_Network_CdnMarketplace_Configuration_Mapping/deleteDomainMapping/
*/
func (a cdnManager) DeleteCDN(uniqueId string) ([]datatypes.Container_Network_CdnMarketplace_Configuration_Mapping, error) {
	return a.CdnService.Mask(mask).DeleteDomainMapping(&uniqueId)
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

/*
https://sldn.softlayer.com/reference/services/SoftLayer_Network_CdnMarketplace_Configuration_Mapping/createDomainMapping/
*/
func (a cdnManager) CreateCdn(hostname string, originHost string, originType string, http int, https int, bucketName string, cname string, header string, path string, ssl string) ([]datatypes.Container_Network_CdnMarketplace_Configuration_Mapping, error) {
	types := map[string]string{
		"server":  "HOST_SERVER",
		"storage": "OBJECT_STORAGE",
	}
	sslCertificate := map[string]string{
		"wilcard": "WILDCARD_CERT",
		"dvSan":   "SHARED_SAN_CERT",
	}

	NewOrigin := datatypes.Container_Network_CdnMarketplace_Configuration_Input{
		Domain:     sl.String(hostname),
		Origin:     sl.String(originHost),
		OriginType: sl.String(types[originType]),
		VendorName: sl.String("akamai"),
	}

	protocol := ""
	if http != 0 {
		protocol = "HTTP"
		NewOrigin.HttpPort = sl.Int(http)
	}
	if https != 0 {
		protocol = "HTTPS"
		NewOrigin.HttpsPort = sl.Int(https)
		NewOrigin.CertificateType = sl.String(sslCertificate[ssl])
	}
	if http != 0 && https != 0 {
		protocol = "HTTP_AND_HTTPS"
	}
	NewOrigin.Protocol = sl.String(protocol)
	if types[originType] == "OBJECT_STORAGE" {
		NewOrigin.BucketName = sl.String(bucketName)
		NewOrigin.Header = sl.String(originHost)
	}
	if cname != "" {
		NewOrigin.Cname = sl.String(cname + ".cdn.appdomain.cloud")
	}
	if header != "" {
		NewOrigin.Header = sl.String(header)
	}
	if path != "" {
		NewOrigin.Path = sl.String("/" + path)
	}

	return a.CdnService.CreateDomainMapping(&NewOrigin)
}

/*
https://sldn.softlayer.com/reference/services/SoftLayer_Network_CdnMarketplace_Configuration_Mapping_Path/createOriginPath/
*/
func (a cdnManager) OriginAddCdn(uniqueId string, header string, path string, originHost string, originType string, http int, https int, cacheKey string, optimize string, dynamicPath string, dynamicPrefetch bool, dynamicCompression bool, bucketName string, fileExtension string) ([]datatypes.Container_Network_CdnMarketplace_Configuration_Mapping_Path, error) {
	types := map[string]string{
		"server":  "HOST_SERVER",
		"storage": "OBJECT_STORAGE",
	}
	performanceConfig := map[string]string{
		"web":     "General web delivery",
		"video":   "Video on demand optimization",
		"file":    "Large file optimization",
		"dynamic": "Dynamic content acceleration",
	}

	NewOrigin := datatypes.Container_Network_CdnMarketplace_Configuration_Input{
		UniqueId:                 sl.String(uniqueId),
		Path:                     sl.String("/" + path),
		Origin:                   sl.String(originHost),
		OriginType:               sl.String(types[originType]),
		CacheKeyQueryRule:        sl.String(cacheKey),
		PerformanceConfiguration: sl.String(performanceConfig[optimize]),
	}

	if optimize == "dynamic" {
		NewOrigin.DynamicContentAcceleration = &datatypes.Container_Network_CdnMarketplace_Configuration_Performance_DynamicContentAcceleration{
			DetectionPath:                 sl.String("/" + dynamicPath),
			PrefetchEnabled:               sl.Bool(dynamicPrefetch),
			MobileImageCompressionEnabled: sl.Bool(dynamicCompression),
		}
	}

	if originType == "storage" {
		NewOrigin.BucketName = sl.String(bucketName)
		NewOrigin.FileExtension = sl.String(fileExtension)
		NewOrigin.Header = sl.String(originHost)
	}

	if header != "" {
		NewOrigin.Header = sl.String(header)
	}
	if http != 0 {
		NewOrigin.HttpPort = sl.Int(http)
	}
	if https != 0 {
		NewOrigin.HttpsPort = sl.Int(https)
	}

	return a.CdnPathService.CreateOriginPath(&NewOrigin)
}
