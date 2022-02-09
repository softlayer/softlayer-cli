package managers

import (
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/filter"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
)

type LicensesManager interface {
	CreateLicensesOptions() ([]datatypes.Product_Package, error)
}

type licensesManager struct {
	ProductService services.Product_Package
	PackageName                      string
}

func NewLicensesManager(session *session.Session) *licensesManager {
	return &licensesManager{
		services.GetProductPackageService(session),
		"SOFTWARE_LICENSE_PACKAGE",
	}
}

func (l licensesManager) CreateLicensesOptions() ([]datatypes.Product_Package, error) {
	filters := filter.New(filter.Path("keyName").Eq(l.PackageName))
	return l.ProductService.Mask("id,keyName,name,items[prices],regions[location[location[groups]]]").Filter(filters.Build()).GetAllObjects()
}