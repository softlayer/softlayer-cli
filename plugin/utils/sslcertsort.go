package utils

import (
	"github.com/softlayer/softlayer-go/datatypes"
)

type CertById []datatypes.Security_Certificate

func (a CertById) Len() int {
	return len(a)
}
func (a CertById) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a CertById) Less(i, j int) bool {
	if a[i].Id != nil && a[j].Id != nil {
		return *a[i].Id < *a[j].Id
	}
	return false
}

type CertByCommonName []datatypes.Security_Certificate

func (a CertByCommonName) Len() int {
	return len(a)
}
func (a CertByCommonName) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a CertByCommonName) Less(i, j int) bool {
	if a[i].CommonName != nil && a[j].CommonName != nil {
		return *a[i].CommonName < *a[j].CommonName
	}
	return false
}

type CertByValidityDays []datatypes.Security_Certificate

func (a CertByValidityDays) Len() int {
	return len(a)
}
func (a CertByValidityDays) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a CertByValidityDays) Less(i, j int) bool {
	if a[i].ValidityDays != nil && a[j].ValidityDays != nil {
		return *a[i].ValidityDays < *a[j].ValidityDays
	}
	return false
}

type CertByNotes []datatypes.Security_Certificate

func (a CertByNotes) Len() int {
	return len(a)
}
func (a CertByNotes) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a CertByNotes) Less(i, j int) bool {
	if a[i].Notes != nil && a[j].Notes != nil {
		return *a[i].Notes < *a[j].Notes
	}
	return false
}
