package managers

import (
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/filter"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
)

type LicensesManager interface {
	CreateLicensesOptions() ([]datatypes.Product_Package, error)
	CreateLicense(datacenter string, itemKeyName string) (datatypes.Container_Product_Order_Receipt, error)
	CancelItem(key string, immediate bool) error
}

type licensesManager struct {
	LicensesService services.Product_Package
	Session         *session.Session
}

func NewLicensesManager(session *session.Session) *licensesManager {
	return &licensesManager{
		LicensesService: services.GetProductPackageService(session),
		Session:         session,
	}
}

func (l licensesManager) CreateLicensesOptions() ([]datatypes.Product_Package, error) {
	PackageName := "SOFTWARE_LICENSE_PACKAGE"
	filters := filter.New(filter.Path("keyName").Eq(PackageName))
	return l.LicensesService.Mask("id,keyName,name,items[prices],regions[location[location[groups]]]").Filter(filters.Build()).GetAllObjects()
}

//Add a license to the account using the request placeOrder.
//datacenter: short name of datacenter.
//itemKeyName: name from a specific item price.
//https://sldn.softlayer.com/reference/services/SoftLayer_Product_Order/placeOrder/
func (l licensesManager) CreateLicense(datacenter string, itemKeyName string) (datatypes.Container_Product_Order_Receipt, error) {
	BillingOrderService := services.GetProductOrderService(l.Session)
	orderManager := NewOrderManager(l.Session)

	licensePackageKeyName := "SOFTWARE_LICENSE_PACKAGE"
	packageLicenseItemId, err := orderManager.GetPackageByKey(licensePackageKeyName, "")
	if err != nil {
		return datatypes.Container_Product_Order_Receipt{}, err
	}

	locationId, err := orderManager.GetLocation(datacenter)
	if err != nil {
		return datatypes.Container_Product_Order_Receipt{}, err
	}

	arrayItemKeyName := []string{itemKeyName}
	itemPriceId, err := orderManager.GetPriceIdList(licensePackageKeyName, arrayItemKeyName, 0)
	if err != nil {
		return datatypes.Container_Product_Order_Receipt{}, err
	}

	licenseOrder := datatypes.Container_Product_Order_Software_License{
		Container_Product_Order: datatypes.Container_Product_Order{
			ComplexType: sl.String("SoftLayer_Container_Product_Order_Software_License"),
			Location:    sl.String(locationId),
			Prices: []datatypes.Product_Item_Price{
				datatypes.Product_Item_Price{
					Id: sl.Int(itemPriceId[len(itemPriceId)-1]),
				},
			},
			PackageId:        sl.Int(*packageLicenseItemId.Id),
			Quantity:         sl.Int(1),
			UseHourlyPricing: sl.Bool(false),
		},
	}
	return BillingOrderService.PlaceOrder(&licenseOrder, sl.Bool(false))
}
