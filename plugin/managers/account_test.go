package managers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("AccountManager", func() {
	var (
		fakeSLSession 	 *session.Session
		accountManager    managers.AccountManager
	)

	BeforeEach(func() {
		fakeSLSession = testhelpers.NewFakeSoftlayerSession(nil)
		accountManager = managers.NewAccountManager(fakeSLSession)
	})

	Describe("SummaryByDatacenter", func() {
		Context("SummaryByDatacenter", func() {
			It("Returns no errors", func() {
				summary, err := accountManager.SummaryByDatacenter()
				Expect(err).ToNot(HaveOccurred())
				Expect(summary["dal05"]["vlan_count"]).To(Equal(18))
			})
		})
	})
	Describe("GetBandwidthPools()", func() {
		Context("API call is made", func() {
			It("Returns no errors", func() {
				summary, err := accountManager.GetBandwidthPools()
				Expect(err).ToNot(HaveOccurred())
				Expect(*summary[0].Id).To(Equal(309961))
			})
		})
	})
	Describe("GetBandwidthPoolServers()", func() {
		Context("API call is made", func() {
			It("Returns no errors", func() {
				totals, err := accountManager.GetBandwidthPoolServers(12345)
				Expect(err).ToNot(HaveOccurred())
				Expect(totals).To(Equal(3))
			})
		})
	})
})
