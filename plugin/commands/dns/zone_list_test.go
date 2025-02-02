package dns_test

import (
	"errors"
	"strings"
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/dns"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Zone list", func() {
	var (
		fakeUI         *terminal.FakeUI
		cliCommand     *dns.ZoneListCommand
		fakeSession    *session.Session
		slCommand      *metadata.SoftlayerCommand
		fakeDNSManager *testhelpers.FakeDNSManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = dns.NewZoneListCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		fakeDNSManager = new(testhelpers.FakeDNSManager)
		cliCommand.DNSManager = fakeDNSManager
	})

	Describe("Zone list", func() {
		Context("Zone list with server fails", func() {
			BeforeEach(func() {
				fakeDNSManager.ListZonesReturns([]datatypes.Dns_Domain{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to list zones on your account."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
		})

		Context("Zone list", func() {
			BeforeEach(func() {
				time1, _ := time.Parse(time.RFC3339, "2014-11-12T04:45:17Z")
				time2, _ := time.Parse(time.RFC3339, "2016-10-13T09:25:07Z")
				time3, _ := time.Parse(time.RFC3339, "2014-12-16T14:31:16Z")
				fakeDNSManager.ListZonesReturns([]datatypes.Dns_Domain{
					datatypes.Dns_Domain{
						Id:         sl.Int(1745153),
						Name:       sl.String("bcr01.dal06.bluemix.ibmcsf.net"),
						Serial:     sl.Int(2014111108),
						UpdateDate: sl.Time(time1),
					},
					datatypes.Dns_Domain{
						Id:         sl.Int(1745158),
						Name:       sl.String("bluemix.ibmcsf.net"),
						Serial:     sl.Int(2016101304),
						UpdateDate: sl.Time(time2),
					},
					datatypes.Dns_Domain{
						Id:         sl.Int(1745152),
						Name:       sl.String("dal06.bluemix.ibmcsf.net"),
						Serial:     sl.Int(2014121600),
						UpdateDate: sl.Time(time3),
					},
				}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				//TODO ID column is colored, unable to verify as a simple string
				Expect(results[1]).To(ContainSubstring("1745153"))
				Expect(results[1]).To(ContainSubstring("bcr01.dal06.bluemix.ibmcsf.net   2014111108   2014-11-12T04:45:17Z"))
				Expect(results[2]).To(ContainSubstring("1745158"))
				Expect(results[2]).To(ContainSubstring("bluemix.ibmcsf.net               2016101304   2016-10-13T09:25:07Z"))
				Expect(results[3]).To(ContainSubstring("1745152"))
				Expect(results[3]).To(ContainSubstring("dal06.bluemix.ibmcsf.net         2014121600   2014-12-16T14:31:16Z"))
			})
		})
	})
})
