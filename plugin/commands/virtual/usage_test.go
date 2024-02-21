package virtual_test

import (
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/virtual"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("VS usage", func() {
	var (
		fakeUI        *terminal.FakeUI
		cliCommand    *virtual.UsageCommand
		fakeSession   *session.Session
		slCommand     *metadata.SoftlayerCommand
		fakeVSManager *testhelpers.FakeVirtualServerManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeVSManager = new(testhelpers.FakeVirtualServerManager)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = virtual.NewUsageCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.VirtualServerManager = fakeVSManager
	})
	Describe("VS usage", func() {
		Context("usage without vs ID", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})
		})
		Context("VS usage with wrong VS ID", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("required flag(s) \"valid-data\" not set"))
			})
		})
		Context("VS usage successfull", func() {
			BeforeEach(func() {
				created, _ := time.Parse(time.RFC3339, "2016-12-25T00:00:00Z")
				fakeVSManager.GetSummaryUsageReturns([]datatypes.Metric_Tracking_Object_Data{
					datatypes.Metric_Tracking_Object_Data{
						Counter:  sl.Float(.053),
						Type:     sl.String("CPU0"),
						DateTime: sl.Time(created),
					},
				}, nil)
			})
			It("return successfully", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "-s", "2015-10-02", "-e", "2016-12-31", "-t", "cpu0")
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})
