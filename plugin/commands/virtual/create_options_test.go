package virtual_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/virtual"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var CreateOptionsReturn = map[string]map[string]string{
	"locations": map[string]string{
		"dal06": "Dalls 6",
		"dal13": "Dallas 13",
	},
	"sizes": map[string]string{
		"size1": "big?",
	},
	"operating_systems": map[string]string{
		"testOS": "This is a test OS",
	},
	"port_speed": map[string]string{
		"100": "100 MBPS",
	},
	"guests": map[string]string{
		"GUEST_01": "A test guest",
	},
	"extras": map[string]string{
		"EXTRA_02": "An Extra Item",
	},
}

var _ = Describe("VS create options", func() {
	var (
		fakeUI        *terminal.FakeUI
		cliCommand    *virtual.CreateOptionsCommand
		fakeSession   *session.Session
		slCommand     *metadata.SoftlayerCommand
		fakeVSManager *testhelpers.FakeVirtualServerManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeVSManager = new(testhelpers.FakeVirtualServerManager)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = virtual.NewCreateOptionsCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.VirtualServerManager = fakeVSManager
	})

	Describe("VS create options", func() {
		Context("VS create options with server fails", func() {
			BeforeEach(func() {
				fakeVSManager.GetCreateOptionsReturns(map[string]map[string]string{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get virtual server creation options."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
		})
		Context("Happy Path", func() {
			BeforeEach(func() {
				fakeVSManager.GetCreateOptionsReturns(CreateOptionsReturn, nil)
			})
			It("Prints a nice table", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("A test guest   GUEST_01"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("An Extra Item   EXTRA_02"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Dallas 13    dal13"))

			})
		})
	})
})
