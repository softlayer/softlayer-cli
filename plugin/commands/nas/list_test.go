package nas_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/nas"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("nas list", func() {
	var (
		fakeUI                       *terminal.FakeUI
		fakeNasNetworkStorageManager *testhelpers.FakeNasNetworkStorageManager
		cmd                          *nas.ListCommand
		cliCommand                   cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeNasNetworkStorageManager = new(testhelpers.FakeNasNetworkStorageManager)
		cmd = nas.NewListCommand(fakeUI, fakeNasNetworkStorageManager)
		cliCommand = cli.Command{
			Name:        nas.NasListMetaData().Name,
			Description: nas.NasListMetaData().Description,
			Usage:       nas.NasListMetaData().Usage,
			Flags:       nas.NasListMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("nas list", func() {

		Context("Return error", func() {
			It("Set invalid output", func() {
				err := testhelpers.RunCommand(cliCommand, "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeNasNetworkStorageManager.ListNasNetworkStoragesReturns([]datatypes.Network_Storage{}, errors.New("Failed to get NAS Network Storages."))
			})
			It("Failed get NAS Network Storages", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get NAS Network Storages."))
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				fakerNasNetworkStorages := []datatypes.Network_Storage{
					datatypes.Network_Storage{
						Id: sl.Int(111111),
						ServiceResource: &datatypes.Network_Service_Resource{
							Datacenter: &datatypes.Location{
								Name: sl.String("lon06"),
							},
						},
						CapacityGb:                      sl.Int(20),
						ServiceResourceBackendIpAddress: sl.String("abc-lon0601c-tr.azf.network.com"),
					},
				}
				fakeNasNetworkStorageManager.ListNasNetworkStoragesReturns(fakerNasNetworkStorages, nil)
			})
			It("List NAS Network Storages", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("111111"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("lon06"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("20GB"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("abc-lon0601c-tr.azf.network.com"))
			})
		})

	})
})
