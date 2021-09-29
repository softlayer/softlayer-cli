package virtual_test

import (
	"errors"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/virtual"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
	"strings"
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
			Name:        metadata.VSCapacityListMetaData().Name,
			Description: metadata.VSCapacityListMetaData().Description,
			Usage:       metadata.VSCapacityListMetaData().Usage,
			Flags:       metadata.VSCapacityListMetaData().Flags,
			Action:      cmd.Run,
		}
	})
	Describe("VS capacity-list", func() {
		Context("VS capacity-list with wrong parameters", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "--column", "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "flag provided but not defined: -column")).To(BeTrue())
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
				Expect(strings.Contains(err.Error(), "flag provided but not defined: -column")).To(BeTrue())
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
