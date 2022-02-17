package dns_test

import (
	"errors"
	"strings"
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/dns"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Zone list", func() {
	var (
		fakeUI         *terminal.FakeUI
		fakeDNSManager *testhelpers.FakeDNSManager
		cmd            *dns.ZoneListCommand
		cliCommand     cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeDNSManager = new(testhelpers.FakeDNSManager)
		cmd = dns.NewZoneListCommand(fakeUI, fakeDNSManager)
		cliCommand = cli.Command{
			Name:        dns.DnsZoneListMetaData().Name,
			Description: dns.DnsZoneListMetaData().Description,
			Usage:       dns.DnsZoneListMetaData().Usage,
			Flags:       dns.DnsZoneListMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Zone list", func() {
		Context("Zone list with server fails", func() {
			BeforeEach(func() {
				fakeDNSManager.ListZonesReturns([]datatypes.Dns_Domain{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to list zones on your account.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
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
				err := testhelpers.RunCommand(cliCommand, "")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				//TODO ID column is colored, unable to verify as a simple string
				Expect(strings.Contains(results[1], "1745153")).To(BeTrue())
				Expect(strings.Contains(results[1], "bcr01.dal06.bluemix.ibmcsf.net   2014111108   2014-11-12T04:45:17Z")).To(BeTrue())
				Expect(strings.Contains(results[2], "1745158")).To(BeTrue())
				Expect(strings.Contains(results[2], "bluemix.ibmcsf.net               2016101304   2016-10-13T09:25:07Z")).To(BeTrue())
				Expect(strings.Contains(results[3], "1745152")).To(BeTrue())
				Expect(strings.Contains(results[3], "dal06.bluemix.ibmcsf.net         2014121600   2014-12-16T14:31:16Z")).To(BeTrue())
			})
		})
	})
})
