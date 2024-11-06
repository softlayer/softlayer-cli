package managers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("DedicatedhostManager", func() {
	var (
		fakeSLSession        *session.Session
		fakeHandler    		*testhelpers.FakeTransportHandler
		dedicatedhostManager managers.DedicatedHostManager
	)

	BeforeEach(func() {
		fakeSLSession = testhelpers.NewFakeSoftlayerSession(nil)
		fakeHandler = testhelpers.GetSessionHandler(fakeSLSession)
		dedicatedhostManager = managers.NewDedicatedhostManager(fakeSLSession)
	})
	AfterEach(func() {
		fakeHandler.ClearApiCallLogs()
		fakeHandler.ClearErrors()
	})

	Describe("Genereate a dedicatedhost order template", func() {
		Context("not found Package", func() {
			BeforeEach(func() {
				filenames := []string{"getAllObjects_dedicatedhostNotFoundPackage"}
				fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
				dedicatedhostManager = managers.NewDedicatedhostManager(fakeSLSession)
			})
			It("it returns dedicatedhost order template", func() {
				_, err := dedicatedhostManager.GenerateOrderTemplate("56_CORES_X_242_RAM_X_1_4_TB", "test", "test.com", "ams01", "hourly", 1234567)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Ordering package is not found"))
			})
		})

		Context("Genereate a dedicatedhost order template", func() {
			BeforeEach(func() {
				filenames := []string{"getAllObjects_dedicatedhost"}
				fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
				dedicatedhostManager = managers.NewDedicatedhostManager(fakeSLSession)
			})
			It("it returns dedicatedhost order template", func() {
				dedicatedhostTemplate, err := dedicatedhostManager.GenerateOrderTemplate("56_CORES_X_242_RAM_X_1_4_TB", "test", "test.com", "ams01", "hourly", 1234567)
				Expect(err).NotTo(HaveOccurred())
				Expect(*dedicatedhostTemplate.Hardware[0].Domain).To(Equal("test.com"))
				Expect(*dedicatedhostTemplate.Hardware[0].Hostname).To(Equal("test"))
			})
		})
	})

	Describe("Verify Dedicatehhost Instance Creation", func() {
		Context("Verify Dedicatehhost Instance Creation", func() {
			BeforeEach(func() {
				filenames := []string{"verifyOrder_dedicatedhost"}
				fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
				dedicatedhostManager = managers.NewDedicatedhostManager(fakeSLSession)
			})
			It("it returns dedicatedhost verify response", func() {
				dedicatedhostTemplate, _ := dedicatedhostManager.GenerateOrderTemplate("56_CORES_X_242_RAM_X_1_4_TB", "test", "test.com", "ams01", "hourly", 1234567)
				verifyOrder, err := dedicatedhostManager.VerifyInstanceCreation(dedicatedhostTemplate)
				Expect(err).NotTo(HaveOccurred())
				Expect(*verifyOrder.Hardware[0].Hostname).To(Equal("test"))
				Expect(*verifyOrder.Hardware[0].Domain).To(Equal("test.com"))
			})
		})
	})

	Describe("Order a Dedicatehhost Instance", func() {
		Context("Order a Dedicatehhost Instance", func() {
			BeforeEach(func() {
				filenames := []string{"placeOrder_dedicatedhost"}
				fakeSLSession = testhelpers.NewFakeSoftlayerSession(filenames)
				dedicatedhostManager = managers.NewDedicatedhostManager(fakeSLSession)
			})
			It("it returns dedicatedhost verify response", func() {
				dedicatedhostTemplate, _ := dedicatedhostManager.GenerateOrderTemplate("56_CORES_X_242_RAM_X_1_4_TB", "test", "test.com", "ams01", "hourly", 1234567)
				placeOrder, err := dedicatedhostManager.OrderInstance(dedicatedhostTemplate)
				Expect(err).NotTo(HaveOccurred())
				Expect(*placeOrder.OrderId).To(Equal(1111111))
			})
		})
	})
	Describe("Dedicated Host Manager Simple functions", func() {
		Context("DeleteHost", func() {
			It("it returns dedicatedhost verify response", func() {
				err := dedicatedhostManager.DeleteHost(12345)
				Expect(err).NotTo(HaveOccurred())
				apiCalls := fakeHandler.ApiCallLogs
				Expect(len(apiCalls)).To(Equal(1))
				Expect(apiCalls[0]).To(MatchFields(IgnoreExtras, Fields{
					"Service": Equal("SoftLayer_Virtual_DedicatedHost"),
					"Method":  Equal("deleteObject"),
					"Options": PointTo(MatchFields(IgnoreExtras, Fields{"Id": PointTo(Equal(12345))})),
				}))
			})
		})
	})
})
