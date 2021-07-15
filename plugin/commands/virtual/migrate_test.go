package virtual_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/virtual"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("VS migrate", func() {
	var (
		fakeUI        *terminal.FakeUI
		fakeVSManager *testhelpers.FakeVirtualServerManager
		cmd           *virtual.MigrateCommand
		cliCommand    cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeVSManager = new(testhelpers.FakeVirtualServerManager)
		cmd = virtual.NewMigrageCommand(fakeUI, fakeVSManager)
		cliCommand = cli.Command{
			Name:        metadata.VSMigrateMataData().Name,
			Description: metadata.VSMigrateMataData().Description,
			Usage:       metadata.VSMigrateMataData().Usage,
			Flags:       metadata.VSMigrateMataData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("VS Migrate", func() {
		Context("Shows VS and Dedicated pendingMigrate data", func() {
			It("List VS and DedicatedHost pendingMigration", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				call1 := fakeVSManager.MigrateInstanceArgsForCall(0)
				Expect(call1).To(Equal([]int{1234567}))
			})
		})
	})
})
