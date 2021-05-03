package managers_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"

	"github.ibm.com/cgallo/softlayer-cli/plugin/managers"
	"github.ibm.com/cgallo/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("User", func() {
	var (
		fakeSLSession *session.Session
		UserManager   managers.UserManager
	)
	BeforeEach(func() {
		fakeSLSession = testhelpers.NewFakeSoftlayerSession(nil)
		UserManager = managers.NewUserManager(fakeSLSession)
	})

	Describe("ListUsers", func() {
		Context("ListUsers succ", func() {
			BeforeEach(func() {
				filenames := []string{
					"SoftLayer_Account_getUsers",
				}
				fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
				UserManager = managers.NewUserManager(fakeSLSession)
			})
			It("Return users", func() {
				users, err := UserManager.ListUsers("")
				Expect(err).ToNot(HaveOccurred())
				Expect(*users[0].Username).To(Equal("IBM2782"))
				Expect(*users[0].HardwareCount).To(Equal(uint(177)))
			})
		})
	})

	Describe("GetUser", func() {
		Context("GetUser succ", func() {
			BeforeEach(func() {
				filenames := []string{
					"SoftLayer_User_Customer_getObject",
				}
				fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
				UserManager = managers.NewUserManager(fakeSLSession)
			})
			It("GetUser users", func() {
				user, err := UserManager.GetUser(22, "mask")
				Expect(err).ToNot(HaveOccurred())
				Expect(*user.Username).To(Equal("test"))
				Expect(*user.LocaleId).To(Equal(1))
			})
		})
	})
	Describe("GetCurrentUser", func() {
		Context("GetCurrentUser succ", func() {
			BeforeEach(func() {
				filenames := []string{
					"SoftLayer_Account_getCurrentUser",
				}
				fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
				UserManager = managers.NewUserManager(fakeSLSession)
			})
			It("GetCurrentUser users", func() {
				user, err := UserManager.GetCurrentUser()
				Expect(err).ToNot(HaveOccurred())
				Expect(*user.Username).To(Equal("test"))
				Expect(*user.LocaleId).To(Equal(1))
			})
		})
	})
	Describe("GetAllPermission", func() {
		Context("GetAllPermission succ", func() {
			BeforeEach(func() {
				filenames := []string{
					"SoftLayer_User_Customer_CustomerPermission_Permission_getAllObjects",
				}
				fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
				UserManager = managers.NewUserManager(fakeSLSession)
			})
			It("GetAllPermission users", func() {
				permission, err := UserManager.GetAllPermission()
				Expect(err).ToNot(HaveOccurred())
				Expect(*permission[0].KeyName).To(Equal("ACCESS_ALL_DEDICATEDHOSTS"))
				Expect(*permission[0].Name).To(Equal("Access Virtual DedicatedHosts"))
			})
		})
	})
	Describe("GetUserPermissions", func() {
		Context("GetUserPermissions succ", func() {
			BeforeEach(func() {
				filenames := []string{
					"SoftLayer_User_Customer_getPermissions",
				}
				fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
				UserManager = managers.NewUserManager(fakeSLSession)
			})
			It("GetUserPermissions users", func() {
				permission, err := UserManager.GetUserPermissions(11)
				Expect(err).ToNot(HaveOccurred())
				Expect(*permission[0].KeyName).To(Equal("ACCESS_ALL_DEDICATEDHOSTS"))
				Expect(*permission[0].Name).To(Equal("Access Virtual DedicatedHosts"))
			})
		})
	})

	Describe("GetLogins", func() {
		Context("GetLogins succ", func() {
			BeforeEach(func() {
				filenames := []string{
					"SoftLayer_User_Customer_getLoginAttempts",
				}
				fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
				UserManager = managers.NewUserManager(fakeSLSession)
			})
			It("GetLogins users", func() {
				var t time.Time
				authentication, err := UserManager.GetLogins(11, t)
				Expect(err).ToNot(HaveOccurred())
				Expect(*authentication[0].SuccessFlag).To(Equal(true))
				Expect(len(authentication)).To(Equal(12))
			})
		})
	})
})
