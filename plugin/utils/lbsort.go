package utils

import (
	"github.com/softlayer/softlayer-go/datatypes"
)

type HCTypesByID []datatypes.Network_Application_Delivery_Controller_LoadBalancer_Health_Check_Type

func (a HCTypesByID) Len() int {
	return len(a)
}
func (a HCTypesByID) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a HCTypesByID) Less(i, j int) bool {
	if a[i].Id != nil && a[j].Id != nil {
		return *a[i].Id < *a[j].Id
	}
	return false
}

type RoutingTypesByID []datatypes.Network_Application_Delivery_Controller_LoadBalancer_Routing_Type

func (a RoutingTypesByID) Len() int {
	return len(a)
}
func (a RoutingTypesByID) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a RoutingTypesByID) Less(i, j int) bool {
	if a[i].Id != nil && a[j].Id != nil {
		return *a[i].Id < *a[j].Id
	}
	return false
}

type RoutingMethodsByID []datatypes.Network_Application_Delivery_Controller_LoadBalancer_Routing_Method

func (a RoutingMethodsByID) Len() int {
	return len(a)
}
func (a RoutingMethodsByID) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a RoutingMethodsByID) Less(i, j int) bool {
	if a[i].Id != nil && a[j].Id != nil {
		return *a[i].Id < *a[j].Id
	}
	return false
}

type LBProductItemByPrice []datatypes.Product_Item

func (a LBProductItemByPrice) Len() int {
	return len(a)
}
func (a LBProductItemByPrice) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a LBProductItemByPrice) Less(i, j int) bool {
	if len(a[i].Prices) > 0 && len(a[j].Prices) > 0 && a[i].Prices[0].RecurringFee != nil && a[j].Prices[0].RecurringFee != nil {
		return *a[i].Prices[0].RecurringFee < *a[j].Prices[0].RecurringFee
	}
	return false
}
