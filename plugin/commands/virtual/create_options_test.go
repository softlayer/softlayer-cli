package virtual_test

import (
	"errors"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/virtual"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("VS create options", func() {
	var (
		fakeUI        *terminal.FakeUI
		fakeVSManager *testhelpers.FakeVirtualServerManager
		cmd           *virtual.CreateOptionsCommand
		cliCommand    cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeVSManager = new(testhelpers.FakeVirtualServerManager)
		cmd = virtual.NewCreateOptionsCommand(fakeUI, fakeVSManager)
		cliCommand = cli.Command{
			Name:        virtual.VSCreateOptionsMetaData().Name,
			Description: virtual.VSCreateOptionsMetaData().Description,
			Usage:       virtual.VSCreateOptionsMetaData().Usage,
			Flags:       virtual.VSCreateOptionsMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("VS create options", func() {
		Context("VS create options with server fails", func() {
			BeforeEach(func() {
				fakeVSManager.GetCreateOptionsReturns(map[string]map[string]string{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to get virtual server creation options.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})
	})
})
