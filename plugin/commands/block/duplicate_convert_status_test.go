package block_test

import (
	"fmt"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/softlayer/softlayer-go/session"

	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/block"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("block duplicate-convert-status", func() {
	var (
		fakeUI      *terminal.FakeUI
		cmd         *block.DuplicateConvertStatusCommand
		cliCommand  cli.Command
		fakeSession *session.Session
		fakeHandler *testhelpers.FakeTransportHandler
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession(nil)
		fakeHandler = testhelpers.GetSessionHandler(fakeSession)

		cmd = block.NewDuplicateConvertStatusCommand(fakeUI, fakeSession)
		cliCommand = cli.Command{
			Name:        block.BlockDuplicateConvertStatusMetaData().Name,
			Description: block.BlockDuplicateConvertStatusMetaData().Description,
			Usage:       block.BlockDuplicateConvertStatusMetaData().Usage,
			Flags:       block.BlockDuplicateConvertStatusMetaData().Flags,
			Action:      cmd.Run,
		}
	})
	AfterEach(func() {
		fakeHandler.ClearApiCallLogs()
		fakeHandler.ClearErrors()
	})
	Describe("block duplicate-convert-status", func() {

		Context("Return error", func() {
			It("Set command without Id", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument."))
			})

			It("Set command with an invalid Id", func() {
				err := testhelpers.RunCommand(cliCommand, "abcde")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Volume ID'. It must be a positive integer."))
			})

			It("Set invalid output", func() {
				err := testhelpers.RunCommand(cliCommand, "123456", "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})

			It("SoftLayer_Network_Storage::getDuplicateConversionStatus() Error", func() {
				fakeHandler.AddApiError("SoftLayer_Network_Storage", "getDuplicateConversionStatus", 500, "ERRRR")
				fmt.Printf("API ERRORS ARE NOW\n%v", fakeHandler.ErrorMap)
				err := testhelpers.RunCommand(cliCommand, "123456")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("ERRRR: ERRRR (HTTP 500)"))
			})

			AfterEach(func() {
				fakeHandler.ClearErrors()
			})
		})

		Context("Happy Path", func() {
			It("Runs without issue", func() {
				err := testhelpers.RunCommand(cliCommand, "123456")
				Expect(err).NotTo(HaveOccurred())
				outputs := fakeUI.Outputs()
				Expect(outputs).To(ContainSubstring("2022-06-13 14:59:17"))
				Expect(outputs).To(ContainSubstring("68"))
				Expect(outputs).To(ContainSubstring("SL02SEVC123456_74"))
			})
		})
	})
})
