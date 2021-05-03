package managers

import (
	"errors"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/filter"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	. "github.ibm.com/cgallo/softlayer-cli/plugin/i18n"
)

//Manages account SSH keys and SSL certificates in SoftLayer.
//See product information here:
//https://knowledgelayer.softlayer.com/procedure/ssh-keys and
//http://www.softlayer.com/ssl-certificates
type SecurityManager interface {
	AddSSHKey(key string, label string, notes string) (datatypes.Security_Ssh_Key, error)
	DeleteSSHKey(keyID int) error
	EditSSHKey(keyID int, label string, notes string) error
	GetSSHKey(keyID int) (datatypes.Security_Ssh_Key, error)
	ListSSHKeys(label string) ([]datatypes.Security_Ssh_Key, error)
	GetSSHKeyIDsFromLabel(label string) ([]int, error)
	ListCertificates(method string) ([]datatypes.Security_Certificate, error)
	AddCertificate(cert datatypes.Security_Certificate) (datatypes.Security_Certificate, error)
	RemoveCertificate(certID int) error
	EditCertificate(cert datatypes.Security_Certificate) error
	GetCertificate(certID int) (datatypes.Security_Certificate, error)
}

type securityManager struct {
	SSHKeyService      services.Security_Ssh_Key
	CertificateService services.Security_Certificate
	AccountService     services.Account
}

func NewSecurityManager(session *session.Session) *securityManager {
	return &securityManager{
		services.GetSecuritySshKeyService(session),
		services.GetSecurityCertificateService(session),
		services.GetAccountService(session),
	}
}

//Adds a new SSH key to the account.
//key: The SSH key to add
//label: The label for the key
//notes: The Notes for the key
func (sm securityManager) AddSSHKey(key string, label string, notes string) (datatypes.Security_Ssh_Key, error) {
	template := datatypes.Security_Ssh_Key{
		Key:   sl.String(key),
		Label: sl.String(label),
		Notes: sl.String(notes),
	}
	return sm.SSHKeyService.CreateObject(&template)
}

//Permanently deletes an SSH key from the account.
//keyID: The ID of the key to delete
func (sm securityManager) DeleteSSHKey(keyID int) error {
	_, err := sm.SSHKeyService.Id(keyID).DeleteObject()
	if err != nil {
		return err
	}
	return nil
}

//Edits information about an SSH key.
//keyID: The ID of the key to delete
//label: The label for the key
//notes: The Notes for the key
func (sm securityManager) EditSSHKey(keyID int, label string, notes string) error {
	template := datatypes.Security_Ssh_Key{}
	if label != "" {
		template.Label = sl.String(label)
	}
	if notes != "" {
		template.Notes = sl.String(notes)
	}
	_, err := sm.SSHKeyService.Id(keyID).EditObject(&template)
	if err != nil {
		return err
	}
	return nil
}

//Returns full information about a single SSH key.
//keyID: The ID of the key to delete
func (sm securityManager) GetSSHKey(keyID int) (datatypes.Security_Ssh_Key, error) {
	return sm.SSHKeyService.Id(keyID).GetObject()
}

//Lists all SSH keys on the account.
//label: The label for the key to be filtered
func (sm securityManager) ListSSHKeys(label string) ([]datatypes.Security_Ssh_Key, error) {
	if label != "" {
		filters := filter.New(filter.Path("sshKeys.label").Eq(label))
		return sm.AccountService.Filter(filters.Build()).GetSshKeys()
	}
	return sm.AccountService.GetSshKeys()
}

//Return sshkey IDs which match the given label.
//label: The label for the key
func (sm securityManager) GetSSHKeyIDsFromLabel(label string) ([]int, error) {
	keys, err := sm.ListSSHKeys(label)
	if err != nil {
		return []int{}, err
	}
	result := []int{}
	for _, key := range keys {
		if key.Id != nil {
			result = append(result, *key.Id)
		}
	}
	return result, nil
}

//List all certificates.
//method: The type of certificates to list. Options are: 'all', 'expired', and 'valid'.
func (sm securityManager) ListCertificates(method string) ([]datatypes.Security_Certificate, error) {
	mask := "mask[id,commonName,validityDays,notes]"
	if method == "" || method == "all" {
		return sm.AccountService.Mask(mask).GetSecurityCertificates()
	} else if method == "expired" {
		return sm.AccountService.Mask(mask).GetExpiredSecurityCertificates()
	} else if method == "valid" {
		return sm.AccountService.Mask(mask).GetValidSecurityCertificates()
	}
	return []datatypes.Security_Certificate{}, errors.New(T("Invalid method."))
}

//Creates a new certificate.
//cert: a template certificate object to be created
func (sm securityManager) AddCertificate(cert datatypes.Security_Certificate) (datatypes.Security_Certificate, error) {
	return sm.CertificateService.CreateObject(&cert)
}

//Removes a certificate.
//certID: a certificate ID to remove
func (sm securityManager) RemoveCertificate(certID int) error {
	_, err := sm.CertificateService.Id(certID).DeleteObject()
	if err != nil {
		return err
	}
	return nil
}

//Updates a certificate with the included options.
//cert: the certificate to be updated
func (sm securityManager) EditCertificate(cert datatypes.Security_Certificate) error {
	if cert.Id != nil {
		_, err := sm.CertificateService.Id(*cert.Id).EditObject(&cert)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New(T("certificate ID not found"))
}

//Gets a certificate with the ID specified.
//certID: a certificate ID to retrieve
func (sm securityManager) GetCertificate(certID int) (datatypes.Security_Certificate, error) {
	return sm.CertificateService.Id(certID).GetObject()
}
