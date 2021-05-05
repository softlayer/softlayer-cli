package block_test

import (
	"errors"
	"strings"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/block"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Volume Limits", func() {
	var (
		fakeUI             *terminal.FakeUI
		FakeStorageManager *testhelpers.FakeStorageManager
		cmd                *block.VolumeLimitCommand
		cliCommand         cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		FakeStorageManager = new(testhelpers.FakeStorageManager)
		cmd = block.NewVolumeLimitCommand(fakeUI, FakeStorageManager)
		cliCommand = cli.Command{
			Name:        metadata.BlockVolumeLimitsMetaData().Name,
			Description: metadata.BlockVolumeLimitsMetaData().Description,
			Usage:       metadata.BlockVolumeLimitsMetaData().Usage,
			Action:      cmd.Run,
		}
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
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Datacenter"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"MaximumAvailableCount"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"ProvisionedCount"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Global"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"100"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"200"}))
			})
		})
		Context("Testing Volume Limits Error", func() {
			BeforeEach(func() {
				FakeStorageManager.GetVolumeCountLimitsReturns(nil, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(fakeUI.Outputs()).NotTo(ContainSubstring("OK"))
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})
	})
})
