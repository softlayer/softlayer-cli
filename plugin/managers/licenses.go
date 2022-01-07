package managers

import (
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
)

type LicensesManager interface {
	CreateLicensesOptions() ([]datatypes.Product_Item, error)
}

type licensesManager struct {
	ProductService services.Product_Package
}

func NewLicensesManager(session *session.Session) *licensesManager {
	return &licensesManager{
		services.GetProductPackageService(session),
	}
}

func (l licensesManager) CreateLicensesOptions() ([]datatypes.Product_Item, error) {
	return l.ProductService.Id(301).GetItems()
}