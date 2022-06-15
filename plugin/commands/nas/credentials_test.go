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

var _ = Describe("nas credentials", func() {
	var (
		fakeUI                       *terminal.FakeUI
		fakeNasNetworkStorageManager *testhelpers.FakeNasNetworkStorageManager
		cmd                          *nas.CredentialsCommand
		cliCommand                   cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeNasNetworkStorageManager = new(testhelpers.FakeNasNetworkStorageManager)
		cmd = nas.NewCredentialsCommand(fakeUI, fakeNasNetworkStorageManager)
		cliCommand = cli.Command{
			Name:        nas.NasCredentialsMetaData().Name,
			Description: nas.NasCredentialsMetaData().Description,
			Usage:       nas.NasCredentialsMetaData().Usage,
			Flags:       nas.NasCredentialsMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("nas credentials", func() {

		Context("Return error", func() {
			It("Set command without Id", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument."))
			})

			It("Set command with an invalid Id", func() {
				err := testhelpers.RunCommand(cliCommand, "abcde")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Autoscale Group ID'. It must be a positive integer."))
			})

			It("Set invalid output", func() {
				err := testhelpers.RunCommand(cliCommand, "123456", "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeNasNetworkStorageManager.GetNasNetworkStorageReturns(datatypes.Network_Storage{}, errors.New("Failed to get NAS Network Storage."))
			})
			It("Failed get NAS Network Storage", func() {
				err := testhelpers.RunCommand(cliCommand, "123456")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get NAS Network Storage."))
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				fakerNasNetworkStorage := datatypes.Network_Storage{
					Id:       sl.Int(123456),
					Username: sl.String("testuser"),
					Password: sl.String("test1234"),
				}
				fakeNasNetworkStorageManager.GetNasNetworkStorageReturns(fakerNasNetworkStorage, nil)
			})
			It("Display NAS Network Storage credentials", func() {
				err := testhelpers.RunCommand(cliCommand, "123456")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("testuser"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("test1234"))
			})
		})

	})
})
