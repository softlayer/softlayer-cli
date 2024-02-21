package managers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("SecurityManager", func() {
	var (
		fakeSLSession   *session.Session
		securityManager managers.SecurityManager
	)

	BeforeEach(func() {
		fakeSLSession = testhelpers.NewFakeSoftlayerSession(nil)
		securityManager = managers.NewSecurityManager(fakeSLSession)
	})

	Describe("Add ssh key", func() {
		Context("Add ssh key", func() {
			It("returns a ssh key", func() {
				key, err := securityManager.AddSSHKey("key", "lable", "notes")
				Expect(err).ToNot(HaveOccurred())
				Expect(*key.Key).To(Equal("key"))
				Expect(*key.Label).To(Equal("label"))
				Expect(*key.Notes).To(Equal("notes"))
			})
		})
	})

	Describe("delete ssh key", func() {
		Context("delete ssh key", func() {
			It("returns no error", func() {
				err := securityManager.DeleteSSHKey(1234)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe("edit ssh key", func() {
		Context("edit ssh key", func() {
			It("returns no error", func() {
				err := securityManager.EditSSHKey(1234, "changedlabel", "changednotes")
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe("get ssh key", func() {
		Context("get ssh key", func() {
			It("returns a ssh key", func() {
				key, err := securityManager.GetSSHKey(1234)
				Expect(err).ToNot(HaveOccurred())
				Expect(*key.Key).To(Equal("key"))
				Expect(*key.Label).To(Equal("label"))
				Expect(*key.Notes).To(Equal("notes"))
			})
		})
	})

	Describe("list ssh keys", func() {
		Context("list ssh keys", func() {
			It("returns a list of ssh keys", func() {
				keys, err := securityManager.ListSSHKeys("")
				Expect(err).ToNot(HaveOccurred())
				Expect(len(keys) > 0).To(BeTrue())
			})
		})
	})

	Describe("GetSSHKeyIDsFromLabel", func() {
		Context("GetSSHKeyIDsFromLabel", func() {
			It("returns a list of ssh key ids", func() {
				keyids, err := securityManager.GetSSHKeyIDsFromLabel("")
				Expect(err).ToNot(HaveOccurred())
				Expect(len(keyids) > 0).To(BeTrue())
			})
		})
	})

	Describe("list certificates", func() {
		Context("list certificates", func() {
			It("returns a list of certificates", func() {
				certs, err := securityManager.ListCertificates("all")
				Expect(err).ToNot(HaveOccurred())
				Expect(len(certs) > 0).To(BeTrue())
			})
			It("returns a list of certificates", func() {
				certs, err := securityManager.ListCertificates("valid")
				Expect(err).ToNot(HaveOccurred())
				Expect(len(certs) > 0).To(BeTrue())
			})
			It("returns a list of certificates", func() {
				certs, err := securityManager.ListCertificates("expired")
				Expect(err).ToNot(HaveOccurred())
				Expect(len(certs) > 0).To(BeTrue())
			})
			It("returns error", func() {
				certs, err := securityManager.ListCertificates("abc")
				Expect(err).To(HaveOccurred())
				Expect(len(certs) > 0).To(BeFalse())
			})
		})
	})

	Describe("Add certificate", func() {
		Context("Add certificate", func() {
			It("returns a certificate", func() {
				template := datatypes.Security_Certificate{}
				cert, err := securityManager.AddCertificate(template)
				Expect(err).ToNot(HaveOccurred())
				Expect(*cert.CommonName).To(Equal("abc"))
				Expect(*cert.ValidityDays).To(Equal(60))
				Expect(*cert.Notes).To(Equal("notes1"))
			})
		})
	})

	Describe("remove certificate", func() {
		Context("remove certificate", func() {
			It("returns no error", func() {
				err := securityManager.RemoveCertificate(123)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe("edit certificate", func() {
		Context("edit certificate", func() {
			It("returns no error", func() {
				template := datatypes.Security_Certificate{Id: sl.Int(123)}
				err := securityManager.EditCertificate(template)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe("get certificate", func() {
		Context("get certificate", func() {
			It("returns no error", func() {
				cert, err := securityManager.GetCertificate(123)
				Expect(err).ToNot(HaveOccurred())
				Expect(*cert.Id).To(Equal(123))
				Expect(*cert.CommonName).To(Equal("abc"))
				Expect(*cert.ValidityDays).To(Equal(60))
				Expect(*cert.Notes).To(Equal("notes1"))
			})
		})
	})
})
