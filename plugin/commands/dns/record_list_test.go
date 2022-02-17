package dns_test

import (
	"errors"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/dns"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Record list", func() {
	var (
		fakeUI         *terminal.FakeUI
		fakeDNSManager *testhelpers.FakeDNSManager
		cmd            *dns.RecordListCommand
		cliCommand     cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeDNSManager = new(testhelpers.FakeDNSManager)
		cmd = dns.NewRecordListCommand(fakeUI, fakeDNSManager)
		cliCommand = cli.Command{
			Name:        dns.DnsRecordListMetaData().Name,
			Description: dns.DnsRecordListMetaData().Description,
			Usage:       dns.DnsRecordListMetaData().Usage,
			Flags:       dns.DnsRecordListMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Record list", func() {
		Context("Record list without zone name", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument.")).To(BeTrue())
			})
		})
		Context("Record list with wrong zone name", func() {
			BeforeEach(func() {
				fakeDNSManager.GetZoneIdFromNameReturns(0, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc.com")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to get zone ID from zone name: abc.com.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})

		Context("Record list with server fails", func() {
			BeforeEach(func() {
				fakeDNSManager.GetZoneIdFromNameReturns(1234, nil)
				fakeDNSManager.ListResourceRecordsReturns([]datatypes.Dns_Domain_ResourceRecord{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc.com")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to list resource records under zone: abc.com")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())

			})
		})

		Context("Record list", func() {
			BeforeEach(func() {
				fakeDNSManager.GetZoneIdFromNameReturns(1234, nil)
				fakeDNSManager.ListResourceRecordsReturns([]datatypes.Dns_Domain_ResourceRecord{
					datatypes.Dns_Domain_ResourceRecord{
						Id:   sl.Int(50585314),
						Type: sl.String("txt"),
						Host: sl.String("txt"),
						Data: sl.String("bcr01.dal06.bluemix.ibmcsf.net"),
						Ttl:  sl.Int(900),
					},
					datatypes.Dns_Domain_ResourceRecord{
						Id:   sl.Int(50585306),
						Type: sl.String("soa"),
						Host: sl.String("@"),
						Data: sl.String("ns1.softlayer.com."),
						Ttl:  sl.Int(900),
					},
					datatypes.Dns_Domain_ResourceRecord{
						Id:   sl.Int(50585307),
						Type: sl.String("ns"),
						Host: sl.String("@"),
						Data: sl.String("ns1.softlayer.com."),
						Ttl:  sl.Int(900),
					},
					datatypes.Dns_Domain_ResourceRecord{
						Id:   sl.Int(50585308),
						Type: sl.String("ns"),
						Host: sl.String("@"),
						Data: sl.String("ns2.softlayer.com."),
						Ttl:  sl.Int(900),
					},
					datatypes.Dns_Domain_ResourceRecord{
						Id:   sl.Int(50585313),
						Type: sl.String("mx"),
						Host: sl.String("@"),
						Data: sl.String("mail.dal06.bluemix.ibmcsf.net."),
						Ttl:  sl.Int(900),
					},
					datatypes.Dns_Domain_ResourceRecord{
						Id:   sl.Int(50585305),
						Type: sl.String("a"),
						Host: sl.String("@"),
						Data: sl.String("127.0.0.1 "),
						Ttl:  sl.Int(900),
					},
					datatypes.Dns_Domain_ResourceRecord{
						Id:   sl.Int(50585312),
						Type: sl.String("a"),
						Host: sl.String("ftp"),
						Data: sl.String("127.0.0.1 "),
						Ttl:  sl.Int(86400),
					},
					datatypes.Dns_Domain_ResourceRecord{
						Id:   sl.Int(50585309),
						Type: sl.String("a"),
						Host: sl.String("mail"),
						Data: sl.String("127.0.0.1 "),
						Ttl:  sl.Int(86400),
					},
					datatypes.Dns_Domain_ResourceRecord{
						Id:   sl.Int(50585310),
						Type: sl.String("a"),
						Host: sl.String("webmail"),
						Data: sl.String("127.0.0.1 "),
						Ttl:  sl.Int(86400),
					},
					datatypes.Dns_Domain_ResourceRecord{
						Id:   sl.Int(50585311),
						Type: sl.String("a"),
						Host: sl.String("www"),
						Data: sl.String("127.0.0.1 "),
						Ttl:  sl.Int(86400),
					},
				}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc.com")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "50585314")).To(BeTrue())
				Expect(strings.Contains(results[1], "txt       TXT    900     bcr01.dal06.bluemix.ibmcsf.net")).To(BeTrue())
				Expect(strings.Contains(results[2], "50585306")).To(BeTrue())
				Expect(strings.Contains(results[2], "@         SOA    900     ns1.softlayer.com.")).To(BeTrue())
				Expect(strings.Contains(results[3], "50585307")).To(BeTrue())
				Expect(strings.Contains(results[3], "@         NS     900     ns1.softlayer.com.")).To(BeTrue())
				Expect(strings.Contains(results[4], "50585308")).To(BeTrue())
				Expect(strings.Contains(results[4], "@         NS     900     ns2.softlayer.com.")).To(BeTrue())
				Expect(strings.Contains(results[5], "50585313")).To(BeTrue())
				Expect(strings.Contains(results[5], "@         MX     900     mail.dal06.bluemix.ibmcsf.net.")).To(BeTrue())
				Expect(strings.Contains(results[6], "50585305")).To(BeTrue())
				Expect(strings.Contains(results[6], "@         A      900     127.0.0.1")).To(BeTrue())
				Expect(strings.Contains(results[7], "50585312")).To(BeTrue())
				Expect(strings.Contains(results[7], "ftp       A      86400   127.0.0.1")).To(BeTrue())
				Expect(strings.Contains(results[8], "50585309")).To(BeTrue())
				Expect(strings.Contains(results[8], "mail      A      86400   127.0.0.1")).To(BeTrue())
				Expect(strings.Contains(results[9], "50585310")).To(BeTrue())
				Expect(strings.Contains(results[9], "webmail   A      86400   127.0.0.1")).To(BeTrue())
				Expect(strings.Contains(results[10], "50585311")).To(BeTrue())
				Expect(strings.Contains(results[10], "www       A      86400   127.0.0.1")).To(BeTrue())
			})
		})
	})
})
