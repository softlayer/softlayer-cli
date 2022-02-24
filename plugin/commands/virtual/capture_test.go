package virtual_test

import (
	"errors"
	"strings"
	"time"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/virtual"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("VS capture", func() {
	var (
		fakeUI        *terminal.FakeUI
		fakeVSManager *testhelpers.FakeVirtualServerManager
		cmd           *virtual.CaptureCommand
		cliCommand    cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeVSManager = new(testhelpers.FakeVirtualServerManager)
		cmd = virtual.NewCaptureCommand(fakeUI, fakeVSManager)
		cliCommand = cli.Command{
			Name:        virtual.VSCaptureMetaData().Name,
			Description: virtual.VSCaptureMetaData().Description,
			Usage:       virtual.VSCaptureMetaData().Usage,
			Flags:       virtual.VSCaptureMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("VS capture", func() {
		Context("VS capture without ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument.")).To(BeTrue())
			})
		})
		Context("VS capture with wrong VS ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Invalid input for 'Virtual server ID'. It must be a positive integer.")).To(BeTrue())
			})
		})

		Context("VS capture without --name", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: '-n|--name' is required")).To(BeTrue())
			})
		})

		Context("VS capture with server fails", func() {
			BeforeEach(func() {
				fakeVSManager.CaptureImageReturns(datatypes.Provisioning_Version1_Transaction{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-n", "myimage")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to capture image for virtual server instance: 1234.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
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
				err := testhelpers.RunCommand(cliCommand, "1234", "-n", "myimage")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"1234"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"2016-12-30T00:00:00Z"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Ongoing"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"12345678"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"false"}))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-n", "myimage", "--all")
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
