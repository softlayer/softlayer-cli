package block_test

import (
	"errors"
	"strings"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/block"
	
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Volume options", func() {
	var (
		fakeUI             *terminal.FakeUI
		FakeStorageManager *testhelpers.FakeStorageManager
		cmd                *block.VolumeOptionsCommand
		cliCommand         cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		FakeStorageManager = new(testhelpers.FakeStorageManager)
		cmd = block.NewVolumeOptionsCommand(fakeUI, FakeStorageManager)
		cliCommand = cli.Command{
			Name:        block.BlockVolumeOptionsMetaData().Name,
			Description: block.BlockVolumeOptionsMetaData().Description,
			Usage:       block.BlockVolumeOptionsMetaData().Usage,
			Flags:       block.BlockVolumeOptionsMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Volume options", func() {
		Context("Volume options with server API call fails", func() {
			BeforeEach(func() {
				FakeStorageManager.GetAllDatacentersReturns([]string{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to get all datacenters.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})

		Context("Volume options", func() {
			BeforeEach(func() {
				FakeStorageManager.GetAllDatacentersReturns([]string{
					"ams01", "ams03", "che01", "dal01", "dal02", "dal05", "dal06", "dal07", "dal09", "dal10", "fra02", "hkg02", "hou02", "lon02", "mel01", "mex01", "mil01", "mon01", "osl01", "par01", "sao01", "sea01", "seo01", "sjc01", "sjc03", "sng01", "syd01", "tok02", "tor01", "wdc01", "wdc04"}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Storage Type"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Size (GB)"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OS Type"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"IOPS"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Tier"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Location"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Snapshot Size (GB)"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"performance,endurance"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"20,40,80,100,250,500,1000,2000-3000,4000-7000,8000-9000,10000-12000"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"HYPER_V,LINUX,VMWARE,WINDOWS_2008,WINDOWS_GPT,WINDOWS,XEN"}))
				//Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Size(GB)   20     40     80     100    250    500    1000   2000   4000   8000   12000"}))
				//Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"100    100    100    100    100    100    100    200    300    500    1000"}))
				//Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"1000   2000   4000   6000   6000   6000   6000   6000   6000   6000   6000"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"0.25,2,4,10"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"ams01,ams03,che01,dal01,dal02,dal05,dal06,dal07,dal09,dal10,fra02,hkg02,hou02,lon02,mel01,mex01,mil01,mon01,osl01,par01,sao01,sea01,seo01,sjc01,sjc03,sng01,syd01,tok02,tor01,wdc01,wdc04"}))
				//Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Storage Size(GB)   Available Snapshot Size(GB)"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"0,5,10,20"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"0,5,10,20,40"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"0,5,10,20,40,60,80"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"0,5,10,20,40,60,80,100"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"0,5,10,20,40,60,80,100,150,200,250"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"0,5,10,20,40,60,80,100,150,200,250,300,350,400,450,500"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"0,5,10,20,40,60,80,100,150,200,250,300,350,400,450,500,600,700,1000"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"0,5,10,20,40,60,80,100,150,200,250,300,350,400,450,500,600,700,1000,2000"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"0,5,10,20,40,60,80,100,150,200,250,300,350,400,450,500,600,700,1000,2000,4000"}))

			})
		})
	})
})
