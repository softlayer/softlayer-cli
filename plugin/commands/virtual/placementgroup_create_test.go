package virtual_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/softlayer/softlayer-go/session"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/virtual"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("VS placementgroup-create", func() {
	var (
		fakeUI        *terminal.FakeUI
		cliCommand    *virtual.PlacementGroupCreateCommand
		fakeSession   *session.Session
		slCommand     *metadata.SoftlayerCommand
		fakeVSManager *testhelpers.FakeVirtualServerManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeVSManager = new(testhelpers.FakeVirtualServerManager)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = virtual.NewPlacementGroupCreateCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.VirtualServerManager = fakeVSManager
	})
	Describe("VS placementgroup-create", func() {
		Context("placementgroup-create option checks", func() {
			It("Missing name options", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '-n, --name' is required"))
			})
			It("Missing backendrouter options", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--name", "testName")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '-b, --backend-router-id' is required"))
			})
			It("Missing rule options", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--name", "testName", "-b", "999")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '-r, --rule-id' is required"))
			})
		})
		Context("placementgroup-create successfull", func() {
			BeforeEach(func() {
				// fakePlacementGroup is from the placementgroup detail test
				fakeVSManager.PlacementCreateReturns(fakePlacementGroup, nil)
			})
			It("return successfully", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--name=TestName", "-b=123", "-r=999")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Successfully created placement group: ID: 123456, Name: test."))
			})
		})
	})

})
