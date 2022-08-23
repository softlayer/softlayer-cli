package block_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/softlayer/softlayer-go/session"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/block"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Volume options", func() {
	var (
		fakeUI             *terminal.FakeUI
		FakeStorageManager *testhelpers.FakeStorageManager
		cliCommand         *block.VolumeOptionsCommand
		fakeSession        *session.Session
		slCommand          *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		FakeStorageManager = new(testhelpers.FakeStorageManager)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = block.NewVolumeOptionsCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.StorageManager = FakeStorageManager
	})

	Describe("Volume options", func() {
		Context("Volume options with server API call fails", func() {
			BeforeEach(func() {
				FakeStorageManager.GetAllDatacentersReturns([]string{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get all datacenters."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
		})

		Context("Volume options", func() {
			BeforeEach(func() {
				FakeStorageManager.GetAllDatacentersReturns([]string{
					"ams01", "ams03", "che01", "dal01", "dal02", "dal05", "dal06", "dal07", "dal09", "dal10", "fra02", "hkg02", "hou02", "lon02", "mel01", "mex01", "mil01", "mon01", "osl01", "par01", "sao01", "sea01", "seo01", "sjc01", "sjc03", "sng01", "syd01", "tok02", "tor01", "wdc01", "wdc04"}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Storage Type"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("HYPER_V,LINUX,VMWARE,WINDOWS_2008,WINDOWS_GPT,WINDOWS,XEN"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0,5,10,20,40,60,80,100,150,200,250,300,350,400,450,500,600,700,1000,2000,4000"))

			})
		})
	})
})
