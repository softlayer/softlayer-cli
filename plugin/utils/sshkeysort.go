package utils

import "github.com/softlayer/softlayer-go/datatypes"

type KeyById []datatypes.Security_Ssh_Key

func (a KeyById) Len() int {
	return len(a)
}
func (a KeyById) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a KeyById) Less(i, j int) bool {
	if a[i].Id != nil && a[j].Id != nil {
		return *a[i].Id < *a[j].Id
	}
	return false
}

type KeyByLabel []datatypes.Security_Ssh_Key

func (a KeyByLabel) Len() int {
	return len(a)
}
func (a KeyByLabel) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a KeyByLabel) Less(i, j int) bool {
	if a[i].Label != nil && a[j].Label != nil {
		return *a[i].Label < *a[j].Label
	}
	return false
}

type KeyByFingerprint []datatypes.Security_Ssh_Key

func (a KeyByFingerprint) Len() int {
	return len(a)
}
func (a KeyByFingerprint) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a KeyByFingerprint) Less(i, j int) bool {
	if a[i].Fingerprint != nil && a[j].Fingerprint != nil {
		return *a[i].Fingerprint < *a[j].Fingerprint
	}
	return false
}

type KeyByNotes []datatypes.Security_Ssh_Key

func (a KeyByNotes) Len() int {
	return len(a)
}
func (a KeyByNotes) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a KeyByNotes) Less(i, j int) bool {
	if a[i].Notes != nil && a[j].Notes != nil {
		return *a[i].Notes < *a[j].Notes
	}
	return false
}
