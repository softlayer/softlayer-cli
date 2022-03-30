package managers_test

import (
	// "strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("PlacementGroup Tests", func() {
	var (
		fakeSLSession *session.Session
		fakeHandler   *testhelpers.FakeTransportHandler
		pgmanager     managers.PlaceGroupManager
	)
	BeforeEach(func() {
		fakeSLSession = testhelpers.NewFakeSoftlayerSession(nil)
		fakeHandler = testhelpers.GetSessionHandler(fakeSLSession)
		pgmanager = managers.NewPlaceGroupManager(fakeSLSession)
	})
	AfterEach(func() {
		fakeHandler.ClearApiCallLogs()
		fakeHandler.ClearErrors()
	})

	Describe("PlaceGroupManager.List()", func() {
		Context("Test API Results", func() {
			It("Success", func() {
				result, err := pgmanager.List("")
				Expect(err).NotTo(HaveOccurred())
				Expect(len(result)).To(Equal(2))
				apiCalls := fakeHandler.ApiCallLogs
				Expect(len(apiCalls)).To(Equal(1))
				Expect(apiCalls[0].Service).To(Equal("SoftLayer_Account"))
				Expect(apiCalls[0].Method).To(Equal("getPlacementGroups"))
				slOptions := apiCalls[0].Options
				Expect(slOptions.Filter).To(ContainSubstring(`"id":{"operation":"orderBy","options":[{"name":"sort","value":["DESC"]}]}`))
			})
			It("API Error", func() {
				fakeHandler.AddApiError("SoftLayer_Account", "getPlacementGroups", 500, "placement-group error")
				_, err := pgmanager.List("")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("placement-group error: placement-group error (HTTP 500)"))
			})
		})
	})
})
