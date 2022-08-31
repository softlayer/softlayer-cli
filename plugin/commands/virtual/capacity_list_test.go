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
		cmd           *virtual.CapacityListCommand
		cliCommand    cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeVSManager = new(testhelpers.FakeVirtualServerManager)
		cmd = virtual.NewCapacityListCommand(fakeUI, fakeVSManager)
		cliCommand = cli.Command{
			Name:        virtual.VSCapacityListMetaData().Name,
			Description: virtual.VSCapacityListMetaData().Description,
			Usage:       virtual.VSCapacityListMetaData().Usage,
			Flags:       virtual.VSCapacityListMetaData().Flags,
			Action:      cmd.Run,
		}
	})
	Describe("VS capacity-list", func() {
		Context("VS capacity-list with wrong parameters", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "--column", "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("flag provided but not defined: -column"))
			})
		})
	})
	Describe("VS capacity-list", func() {
		Context("Failed to get virtual Reserved capacity groups on your account.", func() {
			BeforeEach(func() {
				fakeVSManager.CapacityListReturns(nil, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "--column", "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("flag provided but not defined: -column"))
			})
		})
	})
	Describe("VS capacity-list", func() {
		Context("VS capacity-list no error", func() {
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})
