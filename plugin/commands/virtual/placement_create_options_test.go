package virtual_test

import (
	"errors"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strings"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/virtual"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("VS placementgroup create options", func() {
	var (
		fakeUI        *terminal.FakeUI
		cliCommand    *virtual.PlacementGroupCreateOptionsCommand
		fakeSession   *session.Session
		slCommand     *metadata.SoftlayerCommand
		fakeVSManager *testhelpers.FakeVirtualServerManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeVSManager = new(testhelpers.FakeVirtualServerManager)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = virtual.NewPlacementGroupCreateOptionsCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.VirtualServerManager = fakeVSManager
	})
	Describe("vs placementgroup-create-options", func() {
		Context("VS placementgroup create options with server fails", func() {
			BeforeEach(func() {
				fakeVSManager.GetDatacentersReturns([]datatypes.Location{}, errors.New("Internal Server Error"))
				fakeVSManager.GetAvailablePlacementRoutersReturns([]datatypes.Hardware{}, errors.New("Internal Server Error"))
				fakeVSManager.GetRulesReturns([]datatypes.Virtual_PlacementGroup_Rule{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Internal error."))
			})
		})
		Context("VS placementgroup create options successfull", func() {
			BeforeEach(func() {
				fakeVSManager.GetDatacentersReturns([]datatypes.Location{}, nil)
				fakeVSManager.GetAvailablePlacementRoutersReturns([]datatypes.Hardware{}, nil)
				fakeVSManager.GetRulesReturns([]datatypes.Virtual_PlacementGroup_Rule{}, nil)
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(results[0]).To(ContainSubstring("Datacenter   Hostname   BackendRouterId   "))
				Expect(results[1]).To(ContainSubstring("Id   Rule   "))
			})
		})
	})
})
