package managers_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("DedicatedhostManager", func() {
	var (
		fakeSLSession        *session.Session
		dedicatedhostManager managers.DedicatedHostManager
	)

	BeforeEach(func() {
		fakeSLSession = testhelpers.NewFakeSoftlayerSession(nil)
		dedicatedhostManager = managers.NewDedicatedhostManager(fakeSLSession)
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
				fmt.Println(err)
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
})
