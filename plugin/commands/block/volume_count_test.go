package block_test

import (
	"errors"
	"fmt"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/block"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var listVolumeReturns = []datatypes.Network_Storage{
	datatypes.Network_Storage{
		ServiceResource: &datatypes.Network_Service_Resource{
			Datacenter: &datatypes.Location{Name: sl.String("dal10"), RegionCount: sl.Uint(4)},
		},
	},
}
var _ = Describe("Volume cancel", func() {
	var (
		fakeUI             *terminal.FakeUI
		cliCommand         *block.VolumeCountCommand
		fakeSession        *session.Session
		slCommand          *metadata.SoftlayerStorageCommand
		FakeStorageManager *testhelpers.FakeStorageManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		FakeStorageManager = new(testhelpers.FakeStorageManager)
		slCommand = metadata.NewSoftlayerStorageCommand(fakeUI, fakeSession, "block")
		cliCommand = block.NewVolumeCountCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.StorageManager = FakeStorageManager
	})

	Describe("Volume Count Tests", func() {
		Context("Volume Count Too Many Args", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("unknown command 1234 for volume-count"))
			})
		})
		Context("Volume Count Happy Path", func() {
			BeforeEach(func() {
				FakeStorageManager.ListVolumesReturns(listVolumeReturns, nil)
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				fmt.Printf("\nOUTPUT:\n %v\n", fakeUI.Outputs())
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("dal10         1"))
			})
		})
		Context("Volume Count Happy Path with Datacenter", func() {
			BeforeEach(func() {
				FakeStorageManager.ListVolumesReturns(listVolumeReturns, nil)
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--datacenter", "dal10")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("dal10         1"))
			})
		})
		Context("Volume cancel SL API Error", func() {
			BeforeEach(func() {
				FakeStorageManager.ListVolumesReturns(nil, errors.New("SoftLayer_Exception_ObjectNotFound"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to list volumes on your account."))
			})
		})
		Context("sortby count", func() {
			BeforeEach(func() {
				FakeStorageManager.ListVolumesReturns(listVolumeReturns, nil)
			})
			It("Sorting the Count Column", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "Count")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Data center   Count"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("dal10         1"))
			})
		})
		Context("sortby blank value passed Error", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("flag needs an argument: --sortby"))
			})
		})
		Context("sortby wrong value passed Error", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "zc,.//")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: --sortby zc,.// is not supported."))
			})
		})
	})
})
