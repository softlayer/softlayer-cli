package virtual_test

import (
	"errors"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/virtual"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("VS capacity-list", func() {
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
			Name:        virtual.VSPlacementGroupListMetadata().Name,
			Description: virtual.VSPlacementGroupListMetadata().Description,
			Usage:       virtual.VSPlacementGroupListMetadata().Usage,
			Flags:       virtual.VSPlacementGroupListMetadata().Flags,
			Action:      cmd.Run,
		}
	})
	Describe("VS placementgroup-list", func() {
		Context("VS placementgroup-list with wrong parameters", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "--column", "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(  "flag provided but not defined: -column"))
			})
		})
	})
	Describe("VS placementgroup-list", func() {
		Context("Failed to get virtual placementgroup groups on your account.", func() {
			BeforeEach(func() {
				fakeVSManager.PlacementsGroupListReturns(nil, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "--column", "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring( "flag provided but not defined: -column"))
			})
		})
	})
	Describe("VS placementgroup", func() {
		Context("VS capacity-list no error", func() {
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})
