package virtual_test

import (

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/virtual"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
	"strings"
	"time"
)

var _ = Describe("VS usage", func() {
	var (
		fakeUI        *terminal.FakeUI
		fakeVSManager *testhelpers.FakeVirtualServerManager
		cmd           *virtual.UsageCommand
		cliCommand    cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeVSManager = new(testhelpers.FakeVirtualServerManager)
		cmd = virtual.NewUsageCommand(fakeUI, fakeVSManager)
		cliCommand = cli.Command{
			Name:        metadata.VSUsageMetaData().Name,
			Description: metadata.VSUsageMetaData().Description,
			Usage:       metadata.VSUsageMetaData().Usage,
			Flags:       metadata.VSUsageMetaData().Flags,
			Action:      cmd.Run,
		}
	})
	Describe("VS usage", func() {
		Context("usage without vs ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument."))
			})
		})
		Context("VS usage with wrong VS ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Virtual server ID")).To(BeTrue())
			})
		})
		Context("VS usage successfull", func() {
			BeforeEach(func() {
				created, _ := time.Parse(time.RFC3339, "2016-12-25T00:00:00Z")
				fakeVSManager.GetSummaryUsageReturns([]datatypes.Metric_Tracking_Object_Data{
					datatypes.Metric_Tracking_Object_Data{
						Counter: sl.Float(.053),
						Type: sl.String("CPU0"),
						DateTime: sl.Time(created),
						},

				}, nil)
			})
			It("return successfully", func() {
				err := testhelpers.RunCommand(cliCommand, "123456", "-s", "2015-10-02", "-e", "2016-12-31", "-t" ,"cpu0")
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})
