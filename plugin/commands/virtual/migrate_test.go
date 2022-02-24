package virtual_test

import (
	"errors"
	"strings"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/virtual"
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
			Name:        virtual.VSMigrateMetaData().Name,
			Description: virtual.VSMigrateMetaData().Description,
			Usage:       virtual.VSMigrateMetaData().Usage,
			Flags:       virtual.VSMigrateMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("VS Migrate", func() {
		Context("VS migrate with correct vs ID but server fails", func() {
			BeforeEach(func() {
				fakeVSManager.MigrateInstanceReturns(datatypes.Provisioning_Version1_Transaction{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "--guest", "1234")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to migrate the virtual server instance.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})

		Context("VS migrate with correct vs ID ", func() {
			BeforeEach(func() {
				fakeVSManager.MigrateInstanceReturns(datatypes.Provisioning_Version1_Transaction{
					Id: sl.Int(1234567),
				}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "--guest", "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"The virtual server is migrating."}))
			})
		})

		Context("Error Migrate a dedicated host", func() {
			BeforeEach(func() {
				fakeVSManager.MigrateDedicatedHostReturns(errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "--guest", "1234567", "--host", "1234")
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
