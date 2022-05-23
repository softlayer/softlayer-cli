package managers

import (
	"errors"
	"strconv"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/filter"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
)

type LicensesManager interface {
	CreateLicensesOptions() ([]datatypes.Product_Package, error)
	CreateLicense(key string, datacenter string) (datatypes.Container_Product_Order_Receipt, error)
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

	licensePackageKeyName := "SOFTWARE_LICENSE_PACKAGE"
	packageLicenseItemId, err := l.GetPackageId(licensePackageKeyName)
	if err != nil {
		return datatypes.Container_Product_Order_Receipt{}, err
	}

	locationId, err := l.GetLocationId(datacenter)
	if err != nil {
		return datatypes.Container_Product_Order_Receipt{}, err
	}

	itemPriceId, err := l.GetItemPriceId(packageLicenseItemId, itemKeyName)
	if err != nil {
		return datatypes.Container_Product_Order_Receipt{}, err
	}

	licenseOrder := datatypes.Container_Product_Order_Software_License{
		Container_Product_Order: datatypes.Container_Product_Order{
			ComplexType: sl.String("SoftLayer_Container_Product_Order_Software_License"),
			Location:    sl.String(strconv.Itoa(locationId)),
			Prices: []datatypes.Product_Item_Price{
				datatypes.Product_Item_Price{
					Id: sl.Int(itemPriceId),
				},
			},
			PackageId:        sl.Int(packageLicenseItemId),
			Quantity:         sl.Int(1),
			UseHourlyPricing: sl.Bool(false),
		},
	}
	return BillingOrderService.PlaceOrder(&licenseOrder, sl.Bool(false))
}

//Returns location id of datacenter for ProductOrder::placeOrder().
//location: shortname of datacenter
//https://sldn.softlayer.com/reference/services/SoftLayer_Location/getDatacenters/
func (l licensesManager) GetLocationId(location string) (int, error) {
	LocationDatacenterService := services.GetLocationDatacenterService(l.Session)
	filters := filter.New(filter.Path("name").Eq(location))
	datacenters, err := LocationDatacenterService.Mask("longName,id,name").Filter(filters.Build()).GetDatacenters()
	if err != nil {
		return 0, err
	}
	for _, datacenter := range datacenters {
		if datacenter.Name != nil && *datacenter.Name == location {
			return *datacenter.Id, nil
		}
	}
	return 0, errors.New("Invalid datacenter name specified.")
}

//Returns a itemPriceId valid for a package.
//itemPricesId: id from a specific package.
//keyName: name from a specific item price.
//https://sldn.softlayer.com/reference/services/SoftLayer_Product_Package/getItemPrices/
func (l licensesManager) GetItemPriceId(itemPricesId int, keyName string) (int, error) {
	ProductPackageService := services.GetProductPackageService(l.Session)
	itemPrices, err := ProductPackageService.Id(itemPricesId).GetItems()
	if err != nil {
		return 0, err
	}
	for _, itemPrice := range itemPrices {
		if *itemPrice.KeyName == keyName {
			for _, price := range itemPrice.Prices {
				return *price.Id, nil
			}
		}
	}
	return 0, errors.New("Invalid keyName.")
}

//Returns all the active packages. This will give you a basic description of the packages that are currently active
//packageKeyName: name from a specific package.
//https://sldn.softlayer.com/reference/services/SoftLayer_Product_Package/getAllObjects/
func (l licensesManager) GetPackageId(packageKeyName string) (int, error) {
	ProductPackageService := services.GetProductPackageService(l.Session)
	packageItems, err := ProductPackageService.GetAllObjects()
	if err != nil {
		return 0, err
	}
	for _, packageItem := range packageItems {
		if *packageItem.KeyName == packageKeyName {
			return *packageItem.Id, nil
		}
	}
	return 0, errors.New("Invalid package keyName.")
}
