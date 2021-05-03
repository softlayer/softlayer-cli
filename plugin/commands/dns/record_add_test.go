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
	"github.ibm.com/cgallo/softlayer-cli/plugin/commands/dns"
	"github.ibm.com/cgallo/softlayer-cli/plugin/metadata"
	"github.ibm.com/cgallo/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Record add", func() {
	var (
		fakeUI         *terminal.FakeUI
		fakeDNSManager *testhelpers.FakeDNSManager
		cmd            *dns.RecordAddCommand
		cliCommand     cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeDNSManager = new(testhelpers.FakeDNSManager)
		cmd = dns.NewRecordAddCommand(fakeUI, fakeDNSManager)
		cliCommand = cli.Command{
			Name:        metadata.DnsRecordAddMetaData().Name,
			Description: metadata.DnsRecordAddMetaData().Description,
			Usage:       metadata.DnsRecordAddMetaData().Usage,
			Flags:       metadata.DnsRecordAddMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Record add", func() {
		Context("Record add with not enough parameters", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires four arguments.")).To(BeTrue())
			})
		})
		Context("Record add with not enough parameters", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc.com")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires four arguments.")).To(BeTrue())
			})
		})

		Context("Record add with wrong zone name", func() {
			BeforeEach(func() {
				fakeDNSManager.GetZoneIdFromNameReturns(0, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc.com", "ftp", "a", "127.0.0.1")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to get zone ID from zone name: abc.com.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})

		Context("Record add with server fails", func() {
			BeforeEach(func() {
				fakeDNSManager.GetZoneIdFromNameReturns(123, nil)
				fakeDNSManager.CreateResourceRecordReturns(datatypes.Dns_Domain_ResourceRecord{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc.com", "ftp", "a", "127.0.0.1")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to create resource record under zone abc.com: type=a, record=ftp, data=127.0.0.1, ttl=7200.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc.com", "ftp", "a", "127.0.0.1", "--ttl", "3600")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to create resource record under zone abc.com: type=a, record=ftp, data=127.0.0.1, ttl=3600.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})

		Context("Record add", func() {
			BeforeEach(func() {
				fakeDNSManager.GetZoneIdFromNameReturns(123, nil)
				fakeDNSManager.CreateResourceRecordReturns(datatypes.Dns_Domain_ResourceRecord{
					Id:   sl.Int(1234),
					Type: sl.String("a"),
					Host: sl.String("ftp"),
					Data: sl.String("127.0.0.1"),
					Ttl:  sl.Int(7200),
				}, nil)
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc.com", "ftp", "a", "127.0.0.1")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Created resource record under zone abc.com: ID=1234, type=a, record=ftp, data=127.0.0.1, ttl=7200."}))
			})
		})
		Context("Record add", func() {
			BeforeEach(func() {
				fakeDNSManager.GetZoneIdFromNameReturns(123, nil)
				fakeDNSManager.CreateResourceRecordReturns(datatypes.Dns_Domain_ResourceRecord{
					Id:   sl.Int(1234),
					Type: sl.String("a"),
					Host: sl.String("ftp"),
					Data: sl.String("127.0.0.1"),
					Ttl:  sl.Int(3600),
				}, nil)
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc.com", "ftp", "a", "127.0.0.1", "--ttl", "3600")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Created resource record under zone abc.com: ID=1234, type=a, record=ftp, data=127.0.0.1, ttl=3600."}))
			})
		})
	})
})
