package dedicatedhost_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/dedicatedhost"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Dedicated host create options", func() {
	var (
		fakeUI                   *terminal.FakeUI
		FakeDedicatedhostManager *testhelpers.FakeDedicatedhostManager
		cmd                      *dedicatedhost.CreateOptionsCommand
		cliCommand               cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		FakeDedicatedhostManager = new(testhelpers.FakeDedicatedhostManager)
		cmd = dedicatedhost.NewCreateOptionsCommand(fakeUI, FakeDedicatedhostManager)
		cliCommand = cli.Command{
			Name:        dedicatedhost.DedicatedhostCreateOptionsMetaData().Name,
			Description: dedicatedhost.DedicatedhostCreateOptionsMetaData().Description,
			Usage:       dedicatedhost.DedicatedhostCreateOptionsMetaData().Usage,
			Flags:       dedicatedhost.DedicatedhostCreateOptionsMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Dedicatedhost create options", func() {
		Context("Dedicatedhost create options with datacenter but not a flavor", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "-d", "ams01")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Both -d|--datacenter and -f|--flavor need to be passed as arguments e.g. ibmcloud sl dedicatedhost create-options -d ams01 -f 56_CORES_X_242_RAM_X_1_4_TB"))
			})
		})

		Context("Dedicatedhost create options with flavor but not a datacenter", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "-f", "56_CORES_X_242_RAM_X_1_4_TB")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Both -d|--datacenter and -f|--flavor need to be passed as arguments e.g. ibmcloud sl dedicatedhost create-options -d ams01 -f 56_CORES_X_242_RAM_X_1_4_TB"))
			})
		})

		Context("Dedicatedhost create options getting vlans available failed", func() {
			BeforeEach(func() {
				FakeDedicatedhostManager.GetVlansOptionsReturns(nil, errors.New("Internal server error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "-d", "ams01", "-f", "56_CORES_X_242_RAM_X_1_4_TB")
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
				err := testhelpers.RunCommand(cliCommand)
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
				err := testhelpers.RunCommand(cliCommand, "-d", "ams01", "-f", "56_CORES_X_242_RAM_X_1_4_TB")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("1234"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("test"))
			})
		})
	})
})
