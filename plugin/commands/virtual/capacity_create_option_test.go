package virtual_test

import (
	"errors"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/virtual"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
	"strings"

	. "github.com/onsi/gomega"
)

var _ = Describe("VS capacity create options", func() {
	var (
		fakeUI        *terminal.FakeUI
		fakeVSManager *testhelpers.FakeVirtualServerManager
		cmd           *virtual.CapacityCreateOptiosCommand
		cliCommand    cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeVSManager = new(testhelpers.FakeVirtualServerManager)
		cmd = virtual.NewCapacityCreateOptiosCommand(fakeUI, fakeVSManager)
		cliCommand = cli.Command{
			Name:        metadata.VSCapacityCreateOptionsMetadata().Name,
			Description: metadata.VSCapacityCreateOptionsMetadata().Description,
			Usage:       metadata.VSCapacityCreateOptionsMetadata().Usage,
			Flags:       metadata.VSCapacityCreateOptionsMetadata().Flags,
			Action:      cmd.Run,
		}
	})
	Describe("vs capacity-create-options", func() {
		Context("VS capacity create options with server fails", func() {
			BeforeEach(func() {
				fakeVSManager.GetCapacityCreateOptionsReturns([]datatypes.Product_Item{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: Internal error.")).To(BeTrue())
			})
		})
		Context("VS capacity create options successfull", func() {
			BeforeEach(func() {
				fakeVSManager.GetCapacityCreateOptionsReturns([]datatypes.Product_Item{}, nil)
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[0], "KeyName   Description   term   Default Hourly Price Per Instance   ")).To(BeTrue())
				Expect(strings.Contains(results[1], "Location   POD   BackendRouterId   ")).To(BeTrue())
			})
		})
	})
})
