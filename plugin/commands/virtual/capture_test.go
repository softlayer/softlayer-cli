package virtual_test

import (
	"errors"
	"time"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/virtual"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("VS capture", func() {
	var (
		fakeUI        *terminal.FakeUI
		cliCommand    *virtual.CaptureCommand
		fakeSession   *session.Session
		slCommand     *metadata.SoftlayerCommand
		fakeVSManager *testhelpers.FakeVirtualServerManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeVSManager = new(testhelpers.FakeVirtualServerManager)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = virtual.NewCaptureCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.VirtualServerManager = fakeVSManager
	})

	Describe("VS capture", func() {
		Context("VS capture without ID", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument."))
			})
		})
		Context("VS capture with wrong VS ID", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Virtual server ID'. It must be a positive integer."))
			})
		})

		Context("VS capture without --name", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '-n|--name' is required"))
			})
		})

		Context("VS capture with server fails", func() {
			BeforeEach(func() {
				fakeVSManager.CaptureImageReturns(datatypes.Provisioning_Version1_Transaction{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "-n", "myimage")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to capture image for virtual server instance: 1234."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
		})

		Context("VS capture ", func() {
			BeforeEach(func() {
				created, _ := time.Parse(time.RFC3339, "2016-12-30T00:00:00Z")
				fakeVSManager.CaptureImageReturns(datatypes.Provisioning_Version1_Transaction{
					GuestId:    sl.Int(1234),
					Id:         sl.Int(12345678),
					CreateDate: sl.Time(created),
					TransactionStatus: &datatypes.Provisioning_Version1_Transaction_Status{
						Name: sl.String("Ongoing"),
					},
				}, nil)
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "-n", "myimage")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"1234"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"2016-12-30T00:00:00Z"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Ongoing"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"12345678"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"false"}))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234", "-n", "myimage", "--all")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"1234"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"2016-12-30T00:00:00Z"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Ongoing"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"12345678"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"true"}))
			})
		})
	})
})
