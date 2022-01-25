package image_test

import (
	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/image"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
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
			Name:        metadata.ImageDatacenterMetaData().Name,
			Description: metadata.ImageDatacenterMetaData().Description,
			Usage:       metadata.ImageDatacenterMetaData().Usage,
			Flags:       metadata.ImageDatacenterMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("image datacenter", func() {
		Context("without id", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "--add")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("flag needs an argument: -add"))
			})
		})

		Context("with wrong id", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abcd")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Image ID'. It must be a positive integer."))
			})
		})

		Context("add successfully", func() {
			It("return no error using location id", func() {
				err := testhelpers.RunCommand(cliCommand, "123456", "--add", "265592")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"The location was added successfully!"}))
			})
		})

		Context("add successfully", func() {
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "123456", "--add", "dal05")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"The location was added successfully!"}))
			})
		})

		Context("remove successfully", func() {
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "123456", "--remove", "dal05")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"The location was removed successfully!"}))
			})
		})

	})
})
