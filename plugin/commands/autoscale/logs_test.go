package autoscale_test

import (
	"errors"
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/autoscale"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("autoscale logs", func() {
	var (
		fakeUI               *terminal.FakeUI
		fakeAutoScaleManager *testhelpers.FakeAutoScaleManager
		fakeSecurityManager  *testhelpers.FakeSecurityManager
		cmd                  *autoscale.LogsCommand
		cliCommand           cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeAutoScaleManager = new(testhelpers.FakeAutoScaleManager)
		fakeSecurityManager = new(testhelpers.FakeSecurityManager)
		cmd = autoscale.NewLogsCommand(fakeUI, fakeAutoScaleManager, fakeSecurityManager)
		cliCommand = cli.Command{
			Name:        autoscale.AutoScaleLogsMetaData().Name,
			Description: autoscale.AutoScaleLogsMetaData().Description,
			Usage:       autoscale.AutoScaleLogsMetaData().Usage,
			Flags:       autoscale.AutoScaleLogsMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("autoscale logs", func() {

		Context("Return error", func() {
			It("Set command without Id", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one identifier."))
			})

			It("Set command with an invalid Id", func() {
				err := testhelpers.RunCommand(cliCommand, "abcde")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Autoscale group ID should be a number."))
			})

			It("Set invalid output", func() {
				err := testhelpers.RunCommand(cliCommand, "123456", "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeAutoScaleManager.GetLogsScaleGroupReturns([]datatypes.Scale_Group_Log{}, errors.New("Failed to get AutoScale group logs"))
			})
			It("Failed get scale group logs", func() {
				err := testhelpers.RunCommand(cliCommand, "123456")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get AutoScale group logs"))
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				created, _ := time.Parse(time.RFC3339, "2017-11-08T00:00:00Z")
				fakerLogs := []datatypes.Scale_Group_Log{
					datatypes.Scale_Group_Log{
						CreateDate:  sl.Time(created),
						Description: sl.String("Minimum of 1 is over current count of 0, scaling up to minimum"),
					},
				}
				fakeAutoScaleManager.GetLogsScaleGroupReturns(fakerLogs, nil)
			})

			It("Set command with date-min option", func() {
				err := testhelpers.RunCommand(cliCommand, "123456", "--date-min=2017-01-01")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("2017-11-08T00:00:00Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Minimum of 1 is over current count of 0, scaling up to minimum"))
			})
		})
	})
})
