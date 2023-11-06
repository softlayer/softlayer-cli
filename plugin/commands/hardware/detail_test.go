package hardware_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/hardware"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("hardware detail", func() {
	var (
		fakeUI      *terminal.FakeUI
		cliCommand  *hardware.DetailCommand
		fakeSession *session.Session
		fakeHandler *testhelpers.FakeTransportHandler
		slCommand   *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession(nil)
		fakeHandler = testhelpers.GetSessionHandler(fakeSession)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = hardware.NewDetailCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
	})

	Describe("hardware detail", func() {
		Context("hardware detail without ID", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})
		})
		Context("hardware detail with wrong hardware ID", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Hardware server ID'. It must be a positive integer."))
			})
		})

		Context("SoftLayer_Hardware_Server::getObject error", func() {
			BeforeEach(func() {
				fakeHandler.AddApiError("SoftLayer_Hardware_Server", "getObject", 500, "Internal Server Error")
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get hardware server: 1234."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
		})

		Context("SoftLayer_Hardware_Server::getHardDrives", func() {
			BeforeEach(func() {
				fakeHandler.AddApiError("SoftLayer_Hardware_Server", "getHardDrives", 500, "Failed to get the hard drives detail")
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get the hard drives detail"))
			})
		})

		Context("SoftLayer_Hardware_Server::getBandwidthAllotmentDetail", func() {
			BeforeEach(func() {
				fakeHandler.AddApiError("SoftLayer_Hardware_Server", "getBandwidthAllotmentDetail", 500, "Failed to get bandwidth allotment detail")
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get bandwidth allotment detail"))
			})
		})

		Context("SoftLayer_Hardware_Server::getBillingCycleBandwidthUsage", func() {
			BeforeEach(func() {
				fakeHandler.AddApiError("SoftLayer_Hardware_Server", "getBillingCycleBandwidthUsage", 500, "Failed to get billing cycle bandwidth usage")
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get billing cycle bandwidth usage"))
			})
		})

		Context("Happy Path", func() {

			It("No Options", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("GUID               81434794-af69-44d5-bb97-6b6f43454eee"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Seagate Enterprise Capacity 3.5 V5   2000.00 GB   zc21brfd"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("IPMI IP            10.93.138.222"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Private IP         10.93.138.202"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Owner              SL123456"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Note               My golang note"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Public    0.232080   0.101300   20000"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("RAM                Micron / 8GB DDR4 1Rx8 / 2400 ECC NON Reg"))
			})
			It("Price, Passwords, and Components", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "--passwords", "--price", "--components")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Administrator   ThisPasswordISFake"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("root            FakePassword"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("SN06               2023-02-19T06:00:07Z   HARD_DRIVE"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Lenovo / 3943PAJ / Systemx3250-M6 / Intel Xeon SingleProc SATA / 1"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("5.60               2020-09-24T19:46:29Z   REMOTE_MGMT_CARD"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("FQDN               bardcabero.testedit.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Datacenter         dal10"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("ID                 218027"))
			})
		})

		Context("Issue #649", func() {
			BeforeEach(func() {
				fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
				slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
				cliCommand = hardware.NewDetailCommand(slCommand)
			})
			It("return hardware detail", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "218027")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Last transaction   - -"))
			})
		})
	})
})
