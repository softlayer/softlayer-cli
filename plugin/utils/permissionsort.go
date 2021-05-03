package utils

import (
	"github.com/softlayer/softlayer-go/datatypes"
)

type PermissionsBykeyName []datatypes.User_Customer_CustomerPermission_Permission

func (a PermissionsBykeyName) Len() int {
	return len(a)
}
func (a PermissionsBykeyName) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a PermissionsBykeyName) Less(i, j int) bool {
	if a[i].KeyName != nil && a[j].KeyName != nil {
		return *a[i].KeyName < *a[j].KeyName
	}
	return true
}
