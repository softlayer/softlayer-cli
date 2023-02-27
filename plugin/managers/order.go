package managers

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/filter"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
)

const (
	CATEGORY_MASK = "id,isRequired,itemCategory[id, name, categoryCode]"
	ITEM_MASK     = "id, keyName, description, itemCategory, categories, prices"
	PACKAGE_MASK  = "id, name, keyName, isActive, type"
	PRESET_MASK   = "id, name, keyName, description, categories, prices, locations"
)

type OrderManager interface {
	ListCategories(packageKeyname string) ([]datatypes.Product_Package_Order_Configuration, error)
	ListItems(packageKeyname string, keyword, category string) ([]datatypes.Product_Item, error)
	ListPackage(keyword, packageType string) ([]datatypes.Product_Package, error)
	PackageLocation(packageKeyname string) ([]datatypes.Location_Region, error)
	VerifyPlaceOrder(packageKeyname, location string, itemKeynames []string, complexType string, hourly bool, presetKeyname string, extras interface{}, quantity int) (datatypes.Container_Product_Order, error)
	PlaceOrder(packageKeyname, location string, itemKeynames []string, complexType string, hourly bool, presetKeyname string, extras interface{}, quantity int) (datatypes.Container_Product_Order_Receipt, error)
	PlaceQuote(packageKeyname, location string, itemKeynames []string, complexType, name, presetKeyname string, extras interface{}, sendEmail bool) (datatypes.Container_Product_Order_Receipt, error)
	GetPresetbyKey(packageKeyname, presetKeyname string) (datatypes.Product_Package_Preset, error)
	ListPreset(packageKeyname, keyword string) ([]datatypes.Product_Package_Preset, error)
	GetPackageByKey(packagenamem, mask string) (datatypes.Product_Package, error)
	GenerateOrder(packageKeyname, location string, itemKeynames []string, complexType string, hourly bool, presetKeyname string, extras interface{}, quantity int) (interface{}, error)
	GetPresetPrices(presetId int) (datatypes.Product_Package_Preset, error)
	GetLocation(location string) (string, error)
	GetPriceIdList(packageKeyname string, itemKeynames []string, presetCore float64) ([]int, error)
	GetActiveQuotes(mask string) ([]datatypes.Billing_Order_Quote, error)
	GetQuote(quoteId int, mask string) (datatypes.Billing_Order_Quote, error)
	SaveQuote(quoteId int) (datatypes.Billing_Order_Quote, error)
	VerifyOrder(quoteId int, extra datatypes.Container_Product_Order) (datatypes.Container_Product_Order, error)
	OrderQuote(quoteId int, extra datatypes.Container_Product_Order) (datatypes.Container_Product_Order_Receipt, error)
	GetRecalculatedOrderContainer(quoteId int) (datatypes.Container_Product_Order, error)
	GetOrderDetail(orderId int, mask string) (datatypes.Billing_Order, error)
}

type orderManager struct {
	PackageService           services.Product_Package
	OrderService             services.Product_Order
	LocationService          services.Location_Datacenter
	PackagePreset            services.Product_Package_Preset
	AccountService           services.Account
	BillingOrderQuoteService services.Billing_Order_Quote
	Session                  *session.Session
}

func NewOrderManager(session *session.Session) *orderManager {
	return &orderManager{
		services.GetProductPackageService(session),
		services.GetProductOrderService(session),
		services.GetLocationDatacenterService(session),
		services.GetProductPackagePresetService(session),
		services.GetAccountService(session),
		services.GetBillingOrderQuoteService(session),
		session,
	}
}

//Get a single package with a given key.
//If no packages are found, returns None
//packageKeyname: string representing the package key name we are interested in.
//mask: Mask to specify the properties we want to retrieve
func (i orderManager) GetPackageByKey(packageKeyname, mask string) (datatypes.Product_Package, error) {
	filters := filter.New(filter.Path("keyName").Eq(packageKeyname))
	packages, err := i.PackageService.Filter(filters.Build()).Mask(mask).GetAllObjects()
	if err != nil {
		return datatypes.Product_Package{}, err
	}
	if len(packages) == 0 {
		return datatypes.Product_Package{}, errors.New(T("Package {{.Package}} does not exist.", map[string]interface{}{"Package": packageKeyname}))
	}
	return packages[len(packages)-1], nil
}

//Get details about an image
//image: The ID of the image.
func (i orderManager) ListCategories(packageKeyname string) ([]datatypes.Product_Package_Order_Configuration, error) {

	packages, err := i.GetPackageByKey(packageKeyname, "id")
	if err != nil {
		return nil, err
	}
	return i.PackageService.Id(*packages.Id).Mask(CATEGORY_MASK).GetConfiguration()
}

func (i orderManager) ListItems(packageKeyname string, keyword, category string) ([]datatypes.Product_Item, error) {

	filters := filter.New()

	if keyword != "" {
		filters = append(filters, filter.Path("items.description").Contains(keyword))
	}
	if category != "" {
		filterContainsValue := filter.Path("items.categories.categoryCode")
		filterContainsValue.Op = "_="
		filterContainsValue.Val = category
		filters = append(filters, filterContainsValue)
	}

	packages, err := i.GetPackageByKey(packageKeyname, "id")
	if err != nil {
		return nil, err
	}
	return i.PackageService.Id(*packages.Id).Filter(filters.Build()).Mask(ITEM_MASK).GetItems()
}

func (i orderManager) ListPackage(keyword, packageType string) ([]datatypes.Product_Package, error) {

	filters := filter.New(filter.Path("type.keyName").NotEq("BLUEMIX_SERVICE"))
	if keyword != "" {
		filters = append(filters, filter.Path("name").Contains(keyword))
	}
	if packageType != "" {
		filters = append(filters, filter.Path("type.keyName").Eq(packageType))
	}

	packages, err := i.PackageService.Filter(filters.Build()).Mask(PACKAGE_MASK).GetAllObjects()
	if err != nil {
		return nil, err
	}

	packageIterms := []datatypes.Product_Package{}
	for _, packageIerm := range packages {
		if packageIerm.IsActive != nil {
			if packageIerm.IsActive != nil && *packageIerm.IsActive != 0 {
				packageIterms = append(packageIterms, packageIerm)
			}
		}
	}
	return packageIterms, nil
}

func (i orderManager) PackageLocation(packageKeyname string) ([]datatypes.Location_Region, error) {
	mask := "mask[description, keyname, locations]"
	packages, err := i.GetPackageByKey(packageKeyname, "id")
	if err != nil {
		return nil, err
	}
	return i.PackageService.Id(*packages.Id).Mask(mask).GetRegions()
}

func (i orderManager) ListPreset(packageKeyname, keyword string) ([]datatypes.Product_Package_Preset, error) {
	filter := filter.New(filter.Path("activePresets.name").Contains(keyword))
	packages, err := i.GetPackageByKey(packageKeyname, "id")
	if err != nil {
		return nil, err
	}
	accPresets, err := i.PackageService.Id(*packages.Id).Filter(filter.Build()).Mask(PRESET_MASK).GetAccountRestrictedActivePresets()
	if err != nil {
		return nil, err
	}
	activePresets, err := i.PackageService.Id(*packages.Id).Filter(filter.Build()).Mask(PRESET_MASK).GetActivePresets()
	if err != nil {
		return nil, err
	}
	activePresets = append(activePresets, accPresets...)
	return activePresets, nil
}
func (i orderManager) VerifyPlaceOrder(packageKeyname, location string, itemKeynames []string, complexType string, hourly bool, presetKeyname string, extras interface{}, quantity int) (datatypes.Container_Product_Order, error) {
	order, err := i.GenerateOrder(packageKeyname, location, itemKeynames, complexType, hourly, presetKeyname, extras, quantity)
	if err != nil {
		return datatypes.Container_Product_Order{}, err
	}
	return i.OrderService.VerifyOrder(order)
}

func (i orderManager) PlaceOrder(packageKeyname, location string, itemKeynames []string, complexType string, hourly bool, presetKeyname string, extras interface{}, quantity int) (datatypes.Container_Product_Order_Receipt, error) {
	order, err := i.GenerateOrder(packageKeyname, location, itemKeynames, complexType, hourly, presetKeyname, extras, quantity)
	if err != nil {
		return datatypes.Container_Product_Order_Receipt{}, err

	}
	return i.OrderService.PlaceOrder(order, sl.Bool(false))

}
func (i orderManager) PlaceQuote(packageKeyname, location string, itemKeynames []string, complexType, name, presetKeyname string, extras interface{}, sendEmail bool) (datatypes.Container_Product_Order_Receipt, error) {
	order, err := i.GenerateOrder(packageKeyname, location, itemKeynames, complexType, false, presetKeyname, extras, 1)

	orderValue := reflect.ValueOf(order).Elem()
	orderValue.FieldByName("QuoteName").Set(reflect.ValueOf(&name))
	orderValue.FieldByName("SendQuoteEmailFlag").Set(reflect.ValueOf(&sendEmail))
	if err != nil {
		return datatypes.Container_Product_Order_Receipt{}, err
	}
	return i.OrderService.PlaceQuote(&order)
}

func (i orderManager) GenerateOrder(packageKeyname, location string, itemKeynames []string, complexType string, hourly bool, presetKeyname string, extras interface{}, quantity int) (interface{}, error) {
	packages, err := i.GetPackageByKey(packageKeyname, "id")
	if err != nil {
		return nil, err
	}
	order := reflect.ValueOf(extras).Elem()
	order.FieldByName("PackageId").Set(reflect.ValueOf(packages.Id))
	locationString, err := i.GetLocation(location)
	if err != nil {
		return nil, err
	}
	if quantity == 0 {
		quantity = 1
	}
	order.FieldByName("Location").Set(reflect.ValueOf(&locationString))
	order.FieldByName("Quantity").Set(reflect.ValueOf(&quantity))
	order.FieldByName("UseHourlyPricing").Set(reflect.ValueOf(&hourly))
	var presetCore float64
	if presetKeyname != "" {
		presetId, err := i.GetPresetbyKey(packageKeyname, presetKeyname)
		if err != nil {
			return nil, err

		}
		presetItems, err := i.GetPresetPrices(*presetId.Id)
		if err != nil {
			return nil, err
		}

		for _, item := range presetItems.Prices {
			if item.Item != nil && item.Item.ItemCategory != nil && item.Item.ItemCategory.CategoryCode != nil && *item.Item.ItemCategory.CategoryCode == "guest_core" {
				if item.Item.Capacity != nil {
					presetCore = float64(*item.Item.Capacity)
				}
			}
		}
		order.FieldByName("PresetId").Set(reflect.ValueOf(presetId.Id))
	}
	if complexType == "" {
		return nil, errors.New("A complex type must be specified with the order")
	}
	order.FieldByName("ComplexType").Set(reflect.ValueOf(&complexType))

	priceIds, err := i.GetPriceIdList(packageKeyname, itemKeynames, presetCore)
	if err != nil {
		return nil, err
	}

	var productItemPrice []datatypes.Product_Item_Price
	for _, priceId := range priceIds {
		priceId := priceId
		productItemPrice = append(productItemPrice, datatypes.Product_Item_Price{Id: &priceId})
	}
	order.FieldByName("Prices").Set(reflect.ValueOf(productItemPrice))

	return extras, nil
}

func (i orderManager) GetLocation(location string) (string, error) {
	_, err := strconv.Atoi(location)
	if err == nil {
		return location, nil
	}
	mask := "mask[id,name,regions[keyname]]"
	match, err := regexp.MatchString(`^[a-zA-Z]{3}[0-9]{2}$`, location)
	if err != nil {
		return "", err
	}
	filters := filter.New()
	if match {
		filters = append(filters, filter.Path("name").Eq(location))
	} else {
		filters = append(filters, filter.Path("regions.keyname").Eq(location))
	}
	dataCenter, err := i.LocationService.Mask(mask).Filter(filters.Build()).GetDatacenters()
	if err != nil {
		return "", err
	}
	if len(dataCenter) != 1 {
		return "", errors.New(fmt.Sprintf("%s: %s", T("Unable to find location"), location))
	}
	return strconv.Itoa(*dataCenter[0].Id), nil
}

func (i orderManager) GetPresetbyKey(packageKeyname, presetKeyname string) (datatypes.Product_Package_Preset, error) {
	filterContainsValue1 := filter.Path("activePresets.keyName")
	filterContainsValue1.Op = "_="
	filterContainsValue1.Val = presetKeyname
	filters := filter.New(filterContainsValue1)

	filterContainsValue2 := filter.Path("accountRestrictedActivePresets.keyName")
	filterContainsValue2.Op = "_="
	filterContainsValue2.Val = presetKeyname
	filters = append(filters, filterContainsValue2)

	packages, err := i.GetPackageByKey(packageKeyname, "id")
	if err != nil {
		return datatypes.Product_Package_Preset{}, err
	}
	accPresets, err := i.PackageService.Id(*packages.Id).Filter(filters.Build()).Mask(PRESET_MASK).GetAccountRestrictedActivePresets()
	if err != nil {
		return datatypes.Product_Package_Preset{}, err
	}
	activePresets, err := i.PackageService.Id(*packages.Id).Filter(filters.Build()).Mask(PRESET_MASK).GetActivePresets()
	if err != nil {
		return datatypes.Product_Package_Preset{}, err
	}
	activePresets = append(activePresets, accPresets...)

	if len(activePresets) == 0 {
		return datatypes.Product_Package_Preset{}, errors.New(T("Preset {{.Preset}} does not exist in package {{.Package}}", map[string]interface{}{"Preset": presetKeyname, "Package": packageKeyname}))
	}
	return activePresets[0], nil
}

func (i orderManager) GetPriceIdList(packageKeyname string, itemKeynames []string, presetCore float64) ([]int, error) {

	mask := "id, itemCategory, keyName, prices[categories]"
	packages, err := i.GetPackageByKey(packageKeyname, "id")
	if err != nil {
		return nil, err
	}
	items, err := i.PackageService.Id(*packages.Id).Mask(mask).GetItems()
	if err != nil {
		return nil, err
	}
	var prices []int
	categoryDict := map[string]int{"gpu0": -1, "pcie_slot0": -1}

	for _, itemKeyname := range itemKeynames {
		var newItems []datatypes.Product_Item
		var matchingItem datatypes.Product_Item
		var priceId int
		for _, i := range items {
			if i.KeyName != nil && *i.KeyName == itemKeyname {
				newItems = append(newItems, i)
			}
		}
		if len(newItems) != 0 {
			matchingItem = newItems[0]
		} else {
			return nil, errors.New(T("Item {{.Item}} does not exist for package {{.Package}}", map[string]interface{}{"Item": itemKeyname, "Package": packageKeyname}))
		}
		var itemCategory string
		if matchingItem.ItemCategory != nil && matchingItem.ItemCategory.CategoryCode != nil {
			itemCategory = *matchingItem.ItemCategory.CategoryCode
		}
		if _, ok := categoryDict[itemCategory]; !ok {
			for _, p := range matchingItem.Prices {
				if ((p.LocationGroupId != nil && *p.LocationGroupId == 0) || p.LocationGroupId == nil) && p.Id != nil {
					var capacityMin, capacityMax int
					if p.CapacityRestrictionMinimum != nil {
						capacityMin, err = strconv.Atoi(*p.CapacityRestrictionMinimum)
						if err != nil {
							capacityMin = -1
						}
					}
					if p.CapacityRestrictionMaximum != nil {
						capacityMax, err = strconv.Atoi(*p.CapacityRestrictionMaximum)
						if err != nil {
							capacityMax = -1
						}
					}
					// Some prices might have a specific TermLengh (in months), select only the 0 month price
					if p.TermLength == nil || *p.TermLength == 0 {
						// No Capacity restrictions to check
						if capacityMin == -1 || presetCore == 0 {
							priceId = *p.Id
							// Get restricted Price
						} else if float64(capacityMin) <= presetCore && presetCore <= float64(capacityMax) {
							priceId = *p.Id
						}
					}

				}
			}
		} else {
			var PriceIdList []int
			categoryDict[itemCategory] += 1
			categoryCode := itemCategory[:len(itemCategory)-1] + strconv.Itoa(categoryDict[itemCategory])
			for _, p := range matchingItem.Prices {
				if p.LocationGroupId != nil && *p.LocationGroupId == 0 && len(p.Categories) > 0 && p.Categories[0].CategoryCode != nil && *p.Categories[0].CategoryCode == categoryCode {
					PriceIdList = append(PriceIdList, *p.Id)
				}
			}
			priceId = PriceIdList[0]
		}
		prices = append(prices, priceId)
	}
	return prices, nil
}

func (i orderManager) GetPresetPrices(presetId int) (datatypes.Product_Package_Preset, error) {
	mask := "mask[prices[item]]"
	prices, err := i.PackagePreset.Id(presetId).Mask(mask).GetObject()
	if err != nil {
		return datatypes.Product_Package_Preset{}, nil
	}
	return prices, nil
}

//Returns active quotes on your account
func (i orderManager) GetActiveQuotes(mask string) ([]datatypes.Billing_Order_Quote, error) {
	if mask == "" {
		mask = "mask[order[id,items[id,package[id,keyName]]]]"
	}
	return i.AccountService.Mask(mask).GetActiveQuotes()
}

//Returns active quote detail on your account
func (i orderManager) GetQuote(quoteId int, mask string) (datatypes.Billing_Order_Quote, error) {
	if mask == "" {
		mask = "mask[order[id,items[package[id,keyName]]]]"
	}
	return i.BillingOrderQuoteService.Id(quoteId).Mask(mask).GetObject()
}

//Save quote
func (i orderManager) SaveQuote(quoteId int) (datatypes.Billing_Order_Quote, error) {
	return i.BillingOrderQuoteService.Id(quoteId).SaveQuote()
}

// Verify quote
// quoteId: ID of quote
// extra: container with extras to verify
func (i orderManager) VerifyOrder(quoteId int, extra datatypes.Container_Product_Order) (datatypes.Container_Product_Order, error) {
	return i.BillingOrderQuoteService.Id(quoteId).VerifyOrder(&extra)
}

// Place order from a quote
// quoteId: ID of quote
// extra: container with extras to order
func (i orderManager) OrderQuote(quoteId int, extra datatypes.Container_Product_Order) (datatypes.Container_Product_Order_Receipt, error) {
	return i.BillingOrderQuoteService.Id(quoteId).PlaceOrder(&extra)
}

// Get Recalculated Order Container
// quoteId: ID of quote
func (i orderManager) GetRecalculatedOrderContainer(quoteId int) (datatypes.Container_Product_Order, error) {
	orderBeingPlacedFlag := false
	userOrderData := datatypes.Container_Product_Order{}
	return i.BillingOrderQuoteService.Id(quoteId).GetRecalculatedOrderContainer(&userOrderData, &orderBeingPlacedFlag)
}

// Return order detail
//int orderId: The order identifier.
//string mask: The object mask.
func (i orderManager) GetOrderDetail(orderId int, mask string) (datatypes.Billing_Order, error) {
	if mask == "" {
		mask = `mask[orderTotalAmount,orderApprovalDate,
		initialInvoice[id,amount,invoiceTotalAmount,
		invoiceTopLevelItems[id, description, hostName, domainName, oneTimeAfterTaxAmount,
		recurringAfterTaxAmount, createDate,
		categoryCode,
		category[name],
		location[name],
		children[id, category[name], description, oneTimeAfterTaxAmount,recurringAfterTaxAmount]]],
		items[description],userRecord[displayName,userStatus]]`
	}
	billingOrderService := services.GetBillingOrderService(i.Session)
	return billingOrderService.Id(orderId).Mask(mask).GetObject()
}
