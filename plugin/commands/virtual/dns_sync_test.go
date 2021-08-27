package virtual_test

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
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/virtual"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("VS dns sync", func() {
	var (
		fakeUI         *terminal.FakeUI
		fakeVSManager  *testhelpers.FakeVirtualServerManager
		fakeDNSManager *testhelpers.FakeDNSManager
		cmd            *virtual.DnsSyncCommand
		cliCommand     cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeVSManager = new(testhelpers.FakeVirtualServerManager)
		fakeDNSManager = new(testhelpers.FakeDNSManager)
		cmd = virtual.NewDnsSyncCommand(fakeUI, fakeVSManager, fakeDNSManager)
		cliCommand = cli.Command{
			Name:        metadata.VSDNSSyncMetaData().Name,
			Description: metadata.VSDNSSyncMetaData().Description,
			Usage:       metadata.VSDNSSyncMetaData().Usage,
			Flags:       metadata.VSDNSSyncMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("VS dns sync", func() {
		Context("VS dns sync without ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument.")).To(BeTrue())
			})
		})
		Context("VS dns sync with wrong VS ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Invalid input for 'Virtual server ID'. It must be a positive integer.")).To(BeTrue())
			})
		})

		Context("VS dns sync without -f", func() {
			It("return error", func() {
				fakeUI.Inputs("No")
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Attempt to update DNS records for virtual server instance: 1234. Continue?"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Aborted"}))
			})
		})

		Context("VS dns sync with server fails", func() {
			BeforeEach(func() {
				fakeVSManager.GetInstanceReturns(datatypes.Virtual_Guest{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-f")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to get virtual server instance: 1234.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})
		Context("VS dns sync with fails to get zoneid", func() {
			BeforeEach(func() {
				fakeVSManager.GetInstanceReturns(datatypes.Virtual_Guest{
					Domain: sl.String("wilma.com"),
				}, nil)
				fakeDNSManager.GetZoneIdFromNameReturns(0, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-f")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to get zone ID from zone name: wilma.com")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})
		Context("VS dns sync with fails to sync A record", func() {
			BeforeEach(func() {
				fakeVSManager.GetInstanceReturns(datatypes.Virtual_Guest{
					Id:     sl.Int(1234),
					Domain: sl.String("wilma.com"),
				}, nil)
				fakeDNSManager.GetZoneIdFromNameReturns(9, nil)
				fakeDNSManager.SyncARecordReturns(errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-f")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to synchronize A record for virtual server instance: 1234.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})
		Context("VS dns sync with succeed to sync A record", func() {
			BeforeEach(func() {
				fakeVSManager.GetInstanceReturns(datatypes.Virtual_Guest{
					Id:     sl.Int(1234),
					Domain: sl.String("wilma.com"),
				}, nil)
				fakeDNSManager.GetZoneIdFromNameReturns(9, nil)
				fakeDNSManager.SyncARecordReturns(nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Synchronized A record for virtual server instance: 1234."}))
			})
		})
		Context("VS dns sync with fails to sync AAAA record", func() {
			BeforeEach(func() {
				fakeVSManager.GetInstanceReturns(datatypes.Virtual_Guest{
					Id:     sl.Int(1234),
					Domain: sl.String("wilma.com"),
				}, nil)
				fakeDNSManager.GetZoneIdFromNameReturns(9, nil)
				fakeDNSManager.SyncAAAARecordReturns(errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-f", "--aaaa-record")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to synchronize AAAA record for virtual server instance: 1234.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})
		Context("VS dns sync with succeed to sync AAAAA record", func() {
			BeforeEach(func() {
				fakeVSManager.GetInstanceReturns(datatypes.Virtual_Guest{
					Id:     sl.Int(1234),
					Domain: sl.String("wilma.com"),
				}, nil)
				fakeDNSManager.GetZoneIdFromNameReturns(9, nil)
				fakeDNSManager.SyncAAAARecordReturns(nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-f", "--aaaa-record")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Synchronized AAAA record for virtual server instance: 1234."}))
			})
		})
		Context("VS dns sync with fails to sync ptr record", func() {
			BeforeEach(func() {
				fakeVSManager.GetInstanceReturns(datatypes.Virtual_Guest{
					Id:     sl.Int(1234),
					Domain: sl.String("wilma.com"),
				}, nil)
				fakeDNSManager.GetZoneIdFromNameReturns(9, nil)
				fakeDNSManager.SyncPTRRecordReturns(errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-f", "--ptr")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to synchronize PTR record for virtual server instance: 1234.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})
		Context("VS dns sync with succeed to sync AAAAA record", func() {
			BeforeEach(func() {
				fakeVSManager.GetInstanceReturns(datatypes.Virtual_Guest{
					Id:     sl.Int(1234),
					Domain: sl.String("wilma.com"),
				}, nil)
				fakeDNSManager.GetZoneIdFromNameReturns(9, nil)
				fakeDNSManager.SyncPTRRecordReturns(nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-f", "--ptr")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Synchronized PTR record for virtual server instance: 1234."}))
			})
		})
	})
})
