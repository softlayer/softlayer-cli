package virtual_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	"time"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/virtual"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var placementDate, _ = time.Parse(time.RFC3339, "2017-01-03T00:00:00Z")
var fakePlacementGroup = datatypes.Virtual_PlacementGroup{
	Id:   sl.Int(123456),
	Name: sl.String("test"),
	BackendRouter: &datatypes.Hardware_Router_Backend{
		Hardware_Router: datatypes.Hardware_Router{
			Hardware_Switch: datatypes.Hardware_Switch{
				Hardware: datatypes.Hardware{
					Id:       sl.Int(1115295),
					Hostname: sl.String("bcr01a.wdc07"),
				},
			},
		},
	},
	Rule:       &datatypes.Virtual_PlacementGroup_Rule{Name: sl.String("Rule_Name")},
	GuestCount: sl.Uint(0),
	CreateDate: sl.Time(placementDate),
}
var _ = Describe("VS placementgroup-detail", func() {
	var (
		fakeUI        *terminal.FakeUI
		cliCommand    *virtual.PlacementGroupDetailsCommand
		fakeSession   *session.Session
		slCommand     *metadata.SoftlayerCommand
		fakeVSManager *testhelpers.FakeVirtualServerManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeVSManager = new(testhelpers.FakeVirtualServerManager)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = virtual.NewPlacementGroupDetailsCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.VirtualServerManager = fakeVSManager
	})
	Describe("VS placementgroup-detail", func() {
		Context("placementgroup-Detail without vs ID", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument."))
			})
		})
		Context("VS placementgroup-detail with wrong VS ID", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Placement Group Virtual server ID"))
			})
		})
		Context("VS placementgroup-detail successfull", func() {
			BeforeEach(func() {
				fakeVSManager.GetPlacementGroupDetailReturns(fakePlacementGroup, nil)
			})
			It("return successfully", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456")
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

})
