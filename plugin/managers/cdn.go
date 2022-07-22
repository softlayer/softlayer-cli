package managers

import (
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
)

type CdnManager interface {
	GetNetworkCdnMarketplaceConfigurationMapping(mask string) ([]datatypes.Container_Network_CdnMarketplace_Configuration_Mapping, error)
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
