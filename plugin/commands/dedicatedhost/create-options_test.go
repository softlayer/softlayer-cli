package dedicatedhost_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/dedicatedhost"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Dedicated host create options", func() {
	var (
		fakeUI                   *terminal.FakeUI
		cliCommand               *dedicatedhost.CreateOptionsCommand
		fakeSession              *session.Session
		slCommand                *metadata.SoftlayerCommand
		FakeDedicatedhostManager *testhelpers.FakeDedicatedHostManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = dedicatedhost.NewCreateOptionsCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		FakeDedicatedhostManager = new(testhelpers.FakeDedicatedHostManager)
		cliCommand.DedicatedHostManager = FakeDedicatedhostManager
	})

	Describe("Dedicatedhost create options", func() {
		Context("Dedicatedhost create options with datacenter but not a flavor", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-d", "ams01")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Both -d|--datacenter and -f|--flavor need to be passed as arguments e.g. ibmcloud sl dedicatedhost create-options -d ams01 -f 56_CORES_X_242_RAM_X_1_4_TB"))
			})
		})

		Context("Dedicatedhost create options with flavor but not a datacenter", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-f", "56_CORES_X_242_RAM_X_1_4_TB")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Both -d|--datacenter and -f|--flavor need to be passed as arguments e.g. ibmcloud sl dedicatedhost create-options -d ams01 -f 56_CORES_X_242_RAM_X_1_4_TB"))
			})
		})

		Context("Dedicatedhost create options getting vlans available failed", func() {
			BeforeEach(func() {
				FakeDedicatedhostManager.GetVlansOptionsReturns(nil, errors.New("Internal server error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-d", "ams01", "-f", "56_CORES_X_242_RAM_X_1_4_TB")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get the vlans available for datacener: ams01 and flavor: 56_CORES_X_242_RAM_X_1_4_TB."))
			})
		})

		Context("Dedicatedhost create options successfully without datacenter and flavor", func() {
			BeforeEach(func() {
				FakeDedicatedhostManager.GetCreateOptionsReturns(map[string]map[string]string{
					managers.LOCATIONS:      map[string]string{"dal10": "Dallas 10"},
					managers.DEDICATED_HOST: map[string]string{"56_CORES_X_242_RAM_X_1_4_TB": "56 Cores X 242 RAM X 1.2 TB"},
				})
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("dal10"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("56_CORES_X_242_RAM_X_1_4_TB"))
			})
		})

		Context("Dedicatedhost create options getting vlans with datacenter and flavor", func() {
			BeforeEach(func() {
				FakeDedicatedhostManager.GetCreateOptionsReturns(map[string]map[string]string{})
				FakeDedicatedhostManager.GetVlansOptionsReturns([]datatypes.Network_Vlan{
					datatypes.Network_Vlan{
						Id:            sl.Int(1234),
						Name:          sl.String("test"),
						PrimaryRouter: &datatypes.Hardware_Router{},
					},
				}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-d", "ams01", "-f", "56_CORES_X_242_RAM_X_1_4_TB")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("1234"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("test"))
			})
		})
	})
})
