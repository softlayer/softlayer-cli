package virtual_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/virtual"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("VS dns sync", func() {
	var (
		fakeUI         *terminal.FakeUI
		cliCommand     *virtual.DnsSyncCommand
		fakeSession    *session.Session
		slCommand      *metadata.SoftlayerCommand
		fakeVSManager  *testhelpers.FakeVirtualServerManager
		fakeDNSManager *testhelpers.FakeDNSManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeVSManager = new(testhelpers.FakeVirtualServerManager)
		fakeDNSManager = new(testhelpers.FakeDNSManager)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = virtual.NewDnsSyncCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.VirtualServerManager = fakeVSManager
		cliCommand.DNSManager = fakeDNSManager
	})

	Describe("VS dns sync", func() {
		Context("VS dns sync without ID", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})
		})
		Context("VS dns sync with wrong VS ID", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Virtual server ID'. It must be a positive integer."))
			})
		})

		Context("VS dns sync without -f", func() {
			It("return error", func() {
				fakeUI.Inputs("No")
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Attempt to update DNS records for virtual server instance: 1234. Continue?"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Aborted"))
			})
		})

		Context("VS dns sync with server fails", func() {
			BeforeEach(func() {
				fakeVSManager.GetInstanceReturns(datatypes.Virtual_Guest{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "-f")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get virtual server instance: 1234."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
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
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "-f")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get zone ID from zone name: wilma.com"))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
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
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "-f")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to synchronize A record for virtual server instance: 1234."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
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
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Synchronized A record for virtual server instance: 1234."))
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
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "-f", "--aaaa-record")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to synchronize AAAA record for virtual server instance: 1234."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
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
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "-f", "--aaaa-record")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Synchronized AAAA record for virtual server instance: 1234."))
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
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "-f", "--ptr")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to synchronize PTR record for virtual server instance: 1234."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
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
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "-f", "--ptr")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Synchronized PTR record for virtual server instance: 1234."))
			})
		})
	})
})
