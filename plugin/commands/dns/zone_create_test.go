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

var _ = Describe("Zone create", func() {
	var (
		fakeUI         *terminal.FakeUI
		fakeDNSManager *testhelpers.FakeDNSManager
		cmd            *dns.ZoneCreateCommand
		cliCommand     cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeDNSManager = new(testhelpers.FakeDNSManager)
		cmd = dns.NewZoneCreateCommand(fakeUI, fakeDNSManager)
		cliCommand = cli.Command{
			Name:        dns.DnsZoneCreateMetaData().Name,
			Description: dns.DnsZoneCreateMetaData().Description,
			Usage:       dns.DnsZoneCreateMetaData().Usage,
			Flags:       dns.DnsZoneCreateMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Zone create", func() {
		Context("zone create without zone name", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument.")).To(BeTrue())
			})
		})
		Context("Zone create with server fails", func() {
			BeforeEach(func() {
				fakeDNSManager.CreateZoneReturns(datatypes.Dns_Domain{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc.com")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to create zone: abc.com.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})

		Context("Zone create", func() {
			BeforeEach(func() {
				fakeDNSManager.CreateZoneReturns(datatypes.Dns_Domain{
					Id:   sl.Int(123456),
					Name: sl.String("abc.com"),
				}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc.com")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Zone abc.com was created."}))
			})
		})
	})
})
