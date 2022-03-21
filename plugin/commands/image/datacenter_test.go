package image_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/image"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("image datacenter", func() {
	var (
		fakeUI           *terminal.FakeUI
		fakeImageManager *testhelpers.FakeImageManager
		cmd              *image.DatacenterCommand
		cliCommand       cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeImageManager = new(testhelpers.FakeImageManager)
		cmd = image.NewDatacenterCommand(fakeUI, fakeImageManager)
		cliCommand = cli.Command{
			Name:        image.ImageDatacenterMetaData().Name,
			Description: image.ImageDatacenterMetaData().Description,
			Usage:       image.ImageDatacenterMetaData().Usage,
			Flags:       image.ImageDatacenterMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("image datacenter", func() {
		Context("without option argumment", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "123456", "--add")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("This command requires one indentifier, the option and the option arguments."))
			})
		})

		Context("with wrong id", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abcd", "--remove=dal10")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Image ID'. It must be a positive integer."))
			})
		})

		Context("without image id", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("This command requires one indentifier, the option and the option arguments."))
			})
		})

		Context("without options", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "123456")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires --add or --remove option."))
			})
		})

		Context("add successfully", func() {
			It("return no error using location id", func() {
				err := testhelpers.RunCommand(cliCommand, "123456", "--add", "265592")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("The location was added successfully!"))
			})
		})

		Context("add successfully", func() {
			BeforeEach(func() {
				fakerlocations := []datatypes.Location{
					datatypes.Location{
						Id:       sl.Int(138124),
						LongName: sl.String("Dallas 5"),
						Name:     sl.String("dal05"),
					},
				}
				fakeImageManager.GetDatacentersReturns(fakerlocations, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "123456", "--add", "dal05")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("The location was added successfully!"))
			})
		})

		Context("remove successfully", func() {
			BeforeEach(func() {
				fakerlocations := []datatypes.Location{
					datatypes.Location{
						Id:       sl.Int(138124),
						LongName: sl.String("Dallas 5"),
						Name:     sl.String("dal05"),
					},
				}
				fakeImageManager.GetDatacentersReturns(fakerlocations, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "123456", "--remove", "dal05")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("The location was removed successfully!"))
			})
		})

		Context("remove successfully", func() {
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "123456", "--remove", "265592")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("The location was removed successfully!"))
			})
		})

	})
})
