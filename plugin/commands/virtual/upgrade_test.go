package virtual_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/virtual"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("VS upgrade", func() {
	var (
		fakeUI        *terminal.FakeUI
		cliCommand    *virtual.UpgradeCommand
		fakeSession   *session.Session
		slCommand     *metadata.SoftlayerCommand
		fakeVSManager *testhelpers.FakeVirtualServerManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeVSManager = new(testhelpers.FakeVirtualServerManager)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = virtual.NewUpgradeCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.VirtualServerManager = fakeVSManager
	})
	Describe("VS upgrade", func() {
		Context("VS upgrade without ID", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})
		})
		Context("VS upgrade with wrong vs ID", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Virtual server ID'. It must be a positive integer."))
			})
		})
		Context("VS upgrade with wrong parameters", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--private")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Must specify [--cpu] when using [--private]."))
			})
		})
		Context("VS upgrade with wrong parameters", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Must provide [--cpu], [--memory], [--network], [--add-disk], [--resize-disk] or [--flavor] to upgrade."))
			})
		})
		Context("VS upgrade with wrong --resize-disk parameters", func() {
			It("invalid --resize-disk argument", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--resize-disk", "20:2", "-f")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: --resize-disk requires capacity and disk number values separated by one comma."))
			})
			It("invalid --resize-disk capacity", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--resize-disk", "twenty,2", "-f")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for '--resize-disk capacity'. It must be a positive integer."))
			})
			It("invalid --resize-disk disk number", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--resize-disk", "20,two", "-f")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for '--resize-disk disk number'. It must be a positive integer."))
			})
		})
		Context("VS upgrade without -f", func() {
			It("return no error", func() {
				fakeUI.Inputs("No")
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--cpu", "8")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("This action will incur charges on your account. Continue?"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Aborted"))
			})
		})
		Context("VS upgrade with server fails", func() {
			BeforeEach(func() {
				fakeVSManager.UpgradeInstanceReturns(datatypes.Container_Product_Order_Receipt{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--cpu", "8", "-f")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to upgrade virtual server instance: 1234."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
		})
		Context("VS upgrade", func() {
			BeforeEach(func() {
				fakeVSManager.UpgradeInstanceReturns(datatypes.Container_Product_Order_Receipt{
					OrderId: sl.Int(12345678),
				}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--cpu", "8", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Order 12345678 to upgrade virtual server instance: 1234 was placed."))
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--resize-disk", "10,2", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Order 12345678 to upgrade virtual server instance: 1234 was placed."))
			})
		})
	})
})
