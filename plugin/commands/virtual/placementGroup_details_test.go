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
)

var _ = Describe("VS placementgroup-detail", func() {
	var (
		fakeUI        *terminal.FakeUI
		fakeVSManager *testhelpers.FakeVirtualServerManager
		cmd           *virtual.PlacementGroupListCommand
		cliCommand    cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeVSManager = new(testhelpers.FakeVirtualServerManager)
		cmd = virtual.NewPlacementGroupListCommand(fakeUI, fakeVSManager)
		cliCommand = cli.Command{
			Name:        metadata.VSPlacementGroupDetailMetaData().Name,
			Description: metadata.VSPlacementGroupDetailMetaData().Description,
			Usage:       metadata.VSPlacementGroupDetailMetaData().Usage,
			Flags:       metadata.VSPlacementGroupDetailMetaData().Flags,
			Action:      cmd.Run,
		}
		Describe("VS placementgroup-detail", func() {
			Context("placementgroup-Detail without vs ID", func() {
				It("return error", func() {
					err := testhelpers.RunCommand(cliCommand)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument."))
				})
			})
			Context("VS placementgroup-detail with wrong VS ID", func() {
				It("return error", func() {
					err := testhelpers.RunCommand(cliCommand, "abc")
					Expect(err).To(HaveOccurred())
					Expect(strings.Contains(err.Error(), "Placement Group Virtual server ID")).To(BeTrue())
				})
			})
			Context("VS placementgroup-detail successfull", func() {
				BeforeEach(func() {
					fakeVSManager.GetPlacementGroupDetailReturns(datatypes.Virtual_PlacementGroup{
						Id:   sl.Int(123456),
						Name: sl.String("test"),

					}, nil)

				})
				It("return successfully", func() {
					err := testhelpers.RunCommand(cliCommand, "123456")
					Expect(err).NotTo(HaveOccurred())
				})
			})
		})
	})
})
