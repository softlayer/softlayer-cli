package dns_test

import (
	"errors"
	"strings"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/dns"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Record edit", func() {
	var (
		fakeUI         *terminal.FakeUI
		fakeDNSManager *testhelpers.FakeDNSManager
		cmd            *dns.RecordEditCommand
		cliCommand     cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeDNSManager = new(testhelpers.FakeDNSManager)
		cmd = dns.NewRecordEditCommand(fakeUI, fakeDNSManager)
		cliCommand = cli.Command{
			Name:        dns.DnsRecordEditMetaData().Name,
			Description: dns.DnsRecordEditMetaData().Description,
			Usage:       dns.DnsRecordEditMetaData().Usage,
			Flags:       dns.DnsRecordEditMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Record edit", func() {
		Context("Record edit without zone name", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument.")).To(BeTrue())
			})
		})
		Context("Record edit with wrong zone name", func() {
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

		Context("Record edit with listing records fails", func() {
			BeforeEach(func() {
				fakeDNSManager.ListResourceRecordsReturns([]datatypes.Dns_Domain_ResourceRecord{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc.com")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to list resource records under zone: abc.com.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})

		Context("Record edit with server fails", func() {
			BeforeEach(func() {
				fakeDNSManager.ListResourceRecordsReturns([]datatypes.Dns_Domain_ResourceRecord{
					datatypes.Dns_Domain_ResourceRecord{
						Id:   sl.Int(1234),
						Type: sl.String("a"),
						Host: sl.String("ftp"),
						Data: sl.String("127.0.0.1"),
						Ttl:  sl.Int(900),
					},
				}, nil)
				fakeDNSManager.EditResourceRecordReturns(errors.New("Internal Server Error"))
			})
			It("return  error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc.com", "--data", "127.0.0.2")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to update resource record 1234 under zone abc.com.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})

		Context("Record edit with multiple records and server fails", func() {
			BeforeEach(func() {
				fakeDNSManager.ListResourceRecordsReturns([]datatypes.Dns_Domain_ResourceRecord{
					datatypes.Dns_Domain_ResourceRecord{
						Id:   sl.Int(1234),
						Type: sl.String("a"),
						Host: sl.String("ftp"),
						Data: sl.String("127.0.0.1"),
						Ttl:  sl.Int(900),
					},
					datatypes.Dns_Domain_ResourceRecord{
						Id:   sl.Int(5678),
						Type: sl.String("a"),
						Host: sl.String("mail"),
						Data: sl.String("127.0.0.8"),
						Ttl:  sl.Int(900),
					},
				}, nil)
				fakeDNSManager.EditResourceRecordReturns(errors.New("Internal Server Error"))
			})
			It("return  error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc.com", "--data", "127.0.0.2")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to update resource record 1234 under zone abc.com.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Failed to update resource record 5678 under zone abc.com.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})

		Context("Record edit with different parameters", func() {
			BeforeEach(func() {
				fakeDNSManager.ListResourceRecordsReturns([]datatypes.Dns_Domain_ResourceRecord{
					datatypes.Dns_Domain_ResourceRecord{
						Id:   sl.Int(1234),
						Type: sl.String("a"),
						Host: sl.String("ftp"),
						Data: sl.String("127.0.0.1"),
						Ttl:  sl.Int(900),
					},
				}, nil)
				fakeDNSManager.EditResourceRecordReturns(nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc.com", "--data", "127.0.0.2")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Updated resource record under zone abc.com: ID=1234, type=a, record=ftp, data=127.0.0.2, ttl=900."}))
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc.com", "--ttl", "3600")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Updated resource record under zone abc.com: ID=1234, type=a, record=ftp, data=127.0.0.1, ttl=3600."}))
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc.com", "--by-record", "ftp", "--data", "127.0.0.2", "--ttl", "3600")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Updated resource record under zone abc.com: ID=1234, type=a, record=ftp, data=127.0.0.2, ttl=3600."}))
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc.com", "--by-id", "1234", "--data", "127.0.0.2", "--ttl", "3600")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Updated resource record under zone abc.com: ID=1234, type=a, record=ftp, data=127.0.0.2, ttl=3600."}))
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc.com", "--by-id", "5678", "--data", "127.0.0.2", "--ttl", "3600")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).NotTo(ContainSubstrings([]string{"Updated resource record under zone abc.com: ID=1234, type=a, record=ftp, data=127.0.0.2, ttl=3600."}))
			})
		})
	})
})
