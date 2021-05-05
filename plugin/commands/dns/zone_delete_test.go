package dns_test

import (
	"errors"
	"strings"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/dns"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Zone delete", func() {
	var (
		fakeUI         *terminal.FakeUI
		fakeDNSManager *testhelpers.FakeDNSManager
		cmd            *dns.ZoneDeleteCommand
		cliCommand     cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeDNSManager = new(testhelpers.FakeDNSManager)
		cmd = dns.NewZoneDeleteCommand(fakeUI, fakeDNSManager)
		cliCommand = cli.Command{
			Name:        metadata.DnsZoneDeleteMetaData().Name,
			Description: metadata.DnsZoneDeleteMetaData().Description,
			Usage:       metadata.DnsZoneDeleteMetaData().Usage,
			Flags:       metadata.DnsZoneDeleteMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Zone delete", func() {
		Context("zone delete without zone name", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument.")).To(BeTrue())
			})
		})
		Context("Zone delete with wrong zone name", func() {
			BeforeEach(func() {
				fakeDNSManager.GetZoneIdFromNameReturns(0, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to get zone ID from zone name: abc.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})

		Context("Zone delete with server fails", func() {
			BeforeEach(func() {
				fakeDNSManager.DeleteZoneReturns(errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to delete zone: abc.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())

			})
		})

		Context("Zone delete with zone not found", func() {
			BeforeEach(func() {
				fakeDNSManager.GetZoneIdFromNameReturns(1234, nil)
				fakeDNSManager.DeleteZoneReturns(errors.New("SoftLayer_Exception_ObjectNotFound"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Unable to find zone with ID: 1234.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "SoftLayer_Exception_ObjectNotFound")).To(BeTrue())

			})
		})

		Context("Zone delete", func() {
			BeforeEach(func() {
				fakeDNSManager.GetZoneIdFromNameReturns(1234, nil)
				fakeDNSManager.DeleteZoneReturns(nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Zone abc was deleted."}))
			})
		})
	})
})
