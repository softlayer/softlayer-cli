package virtual_test

import (
	"errors"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/virtual"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
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
			Name:        metadata.VSCreateOptionsMetaData().Name,
			Description: metadata.VSCreateOptionsMetaData().Description,
			Usage:       metadata.VSCreateOptionsMetaData().Usage,
			Flags:       metadata.VSCreateOptionsMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("VS create options", func() {
		Context("VS create options with server fails", func() {
			BeforeEach(func() {
				fakeVSManager.GetCreateOptionsReturns(datatypes.Container_Virtual_Guest_Configuration{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to get virtual server creation options.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})

		Context("VS create options", func() {
			BeforeEach(func() {
				fakeVSManager.GetCreateOptionsReturns(datatypes.Container_Virtual_Guest_Configuration{}, nil)
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[0], "Name")).To(BeTrue())
				Expect(strings.Contains(results[0], "Value")).To(BeTrue())
			})
		})
	})
})
