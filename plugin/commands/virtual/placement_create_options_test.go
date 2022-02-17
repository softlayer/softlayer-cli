package virtual_test

import (
	"errors"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/virtual"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
	"strings"

	. "github.com/onsi/gomega"
)

var _ = Describe("VS placementgroup create options", func() {
	var (
		fakeUI        *terminal.FakeUI
		fakeVSManager *testhelpers.FakeVirtualServerManager
		cmd           *virtual.PlacementGroupCreateOptionsCommand
		cliCommand    cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeVSManager = new(testhelpers.FakeVirtualServerManager)
		cmd = virtual.NewPlacementGruopCreateOptionsCommand(fakeUI, fakeVSManager)
		cliCommand = cli.Command{
			Name:        virtual.VSPlacementGroupCreateOptionsMetaData().Name,
			Description: virtual.VSPlacementGroupCreateOptionsMetaData().Description,
			Usage:       virtual.VSPlacementGroupCreateOptionsMetaData().Usage,
			Flags:       virtual.VSPlacementGroupCreateOptionsMetaData().Flags,
			Action:      cmd.Run,
		}
	})
	Describe("vs placementgroup-create-options", func() {
		Context("VS placementgroup create options with server fails", func() {
			BeforeEach(func() {
				fakeVSManager.GetDatacentersReturns([]datatypes.Location{}, errors.New("Internal Server Error"))
				fakeVSManager.GetAvailablePlacementRoutersReturns([]datatypes.Hardware{}, errors.New("Internal Server Error"))
				fakeVSManager.GetRulesReturns([]datatypes.Virtual_PlacementGroup_Rule{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: Internal error.")).To(BeTrue())
			})
		})
		Context("VS placementgroup create options successfull", func() {
			BeforeEach(func() {
				fakeVSManager.GetDatacentersReturns([]datatypes.Location{}, nil)
				fakeVSManager.GetAvailablePlacementRoutersReturns([]datatypes.Hardware{}, nil)
				fakeVSManager.GetRulesReturns([]datatypes.Virtual_PlacementGroup_Rule{}, nil)
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[0], "Datacenter   Hostname   BackendRouterId   ")).To(BeTrue())
				Expect(strings.Contains(results[1], "Id   Rule   ")).To(BeTrue())
			})
		})
	})
})