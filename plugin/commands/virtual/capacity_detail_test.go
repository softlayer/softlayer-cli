package virtual_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/virtual"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
	"strings"
)

var _ = Describe("VS capacity-detail", func() {
	var (
		fakeUI        *terminal.FakeUI
		fakeVSManager *testhelpers.FakeVirtualServerManager
		cmd           *virtual.CapacityDetailCommand
		cliCommand    cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeVSManager = new(testhelpers.FakeVirtualServerManager)
		cmd = virtual.NewCapacityDetailCommand(fakeUI, fakeVSManager)
		cliCommand = cli.Command{
			Name:        metadata.VSCapacityDetailMetaData().Name,
			Description: metadata.VSCapacityDetailMetaData().Description,
			Usage:       metadata.VSCapacityDetailMetaData().Usage,
			Flags:       metadata.VSCapacityDetailMetaData().Flags,
			Action:      cmd.Run,
		}
	})
	Describe("VS capacity-detail", func() {
		Context("VS capacity-detail without ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Reserved Capacity Group Virtual server ID")).To(BeTrue())
			})
		})
		Context("VS capacity-detail with wrong VS ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Reserved Capacity Group Virtual server ID")).To(BeTrue())
			})
		})
	})
})
