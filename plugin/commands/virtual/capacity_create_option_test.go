package virtual_test

import (
	"errors"
	"strings"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/virtual"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"

)

var _ = Describe("VS capacity create options", func() {
	var (
		fakeUI        *terminal.FakeUI
		cliCommand    *virtual.CapacityCreateOptionsCommand
		fakeSession   *session.Session
		slCommand     *metadata.SoftlayerCommand
		fakeVSManager *testhelpers.FakeVirtualServerManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeVSManager = new(testhelpers.FakeVirtualServerManager)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = virtual.NewCapacityCreateOptionsCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.VirtualServerManager = fakeVSManager
	})
	Describe("vs capacity-create-options", func() {
		Context("VS capacity create options with server fails", func() {
			BeforeEach(func() {
				fakeVSManager.GetCapacityCreateOptionsReturns([]datatypes.Product_Item{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Internal error."))
			})
		})
		Context("VS capacity create options successfull", func() {
			BeforeEach(func() {
				fakeVSManager.GetCapacityCreateOptionsReturns([]datatypes.Product_Item{}, nil)
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(results[0]).To(ContainSubstring("KeyName   Description   term   Default Hourly Price Per Instance"))
				Expect(results[1]).To(ContainSubstring("Location   POD   BackendRouterId"))
			})
		})
	})
})
