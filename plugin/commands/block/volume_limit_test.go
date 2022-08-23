package block_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/block"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Volume Limits", func() {
	var (
		fakeUI             *terminal.FakeUI
		FakeStorageManager *testhelpers.FakeStorageManager
		cliCommand         *block.VolumeLimitCommand
		fakeSession        *session.Session
		slCommand          *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		FakeStorageManager = new(testhelpers.FakeStorageManager)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = block.NewVolumeLimitCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.StorageManager = FakeStorageManager
	})

	Describe("Volume Limits", func() {
		Context("Testing Volume Limits", func() {
			BeforeEach(func() {
				datacenterName := "Global"
				maxcount := 100
				provisionedCount := 200
				FakeStorageManager.GetVolumeCountLimitsReturns([]datatypes.Container_Network_Storage_DataCenterLimits_VolumeCountLimitContainer{
					datatypes.Container_Network_Storage_DataCenterLimits_VolumeCountLimitContainer{
						DatacenterName:        &datacenterName,
						MaximumAvailableCount: &maxcount,
						ProvisionedCount:      &provisionedCount,
					},
				}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Datacenter"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("MaximumAvailableCount"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("100"))
			})
		})
		Context("Testing Volume Limits Error", func() {
			BeforeEach(func() {
				FakeStorageManager.GetVolumeCountLimitsReturns(nil, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(fakeUI.Outputs()).NotTo(ContainSubstring("OK"))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
		})
	})
})
