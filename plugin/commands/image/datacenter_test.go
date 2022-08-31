package image_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/image"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("image datacenter", func() {
	var (
		fakeUI           *terminal.FakeUI
		cliCommand       *image.DatacenterCommand
		fakeSession      *session.Session
		slCommand        *metadata.SoftlayerCommand
		fakeImageManager *testhelpers.FakeImageManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeImageManager = new(testhelpers.FakeImageManager)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = image.NewDatacenterCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.ImageManager = fakeImageManager
	})

	Describe("image datacenter", func() {
		Context("with wrong id", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abcd", "--remove=dal10")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Image ID'. It must be a positive integer."))
			})
		})

		Context("without options", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires --add or --remove option."))
			})
		})

		Context("add successfully", func() {
			It("return no error using location id", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--add", "265592")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
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
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--add", "dal05")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
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
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--remove", "dal05")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
			})
		})

		Context("remove successfully", func() {
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--remove", "265592")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
			})
		})

	})
})
