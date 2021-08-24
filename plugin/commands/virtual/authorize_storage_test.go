package virtual_test

import (
	"errors"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/virtual"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Authorize block, portable and file storage to a VS", func() {
	var (
		fakeUI        *terminal.FakeUI
		fakeVSManager *testhelpers.FakeVirtualServerManager
		cmd           *virtual.AuthorizeStorageCommand
		cliCommand    cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeVSManager = new(testhelpers.FakeVirtualServerManager)
		cmd = virtual.NewAuthorizeStorageCommand(fakeUI, fakeVSManager)
		cliCommand = cli.Command{
			Name:        metadata.VSAuthorizeStorageMetaData().Name,
			Description: metadata.VSAuthorizeStorageMetaData().Description,
			Usage:       metadata.VSAuthorizeStorageMetaData().Usage,
			Flags:       metadata.VSAuthorizeStorageMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Authorize Block, File, Portable Storage to a VS", func() {
		Context("Authorize Storage without VS ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument."))
			})
		})
		Context("Authorize Storage with wrong VS ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Virtual server ID'. It must be a positive integer."))
			})
		})

		Context("Authorize storage to a VS", func() {
			BeforeEach(func() {
				fakeVSManager.AuthorizeStorageReturns(true, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "--username-storage", "SL02SL11111111-11")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
			})
		})

		Context("Error Authorize Storage to a VS", func() {
			BeforeEach(func() {
				fakeVSManager.AuthorizeStorageReturns(false, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "--username-storage", "SL02SL111")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to authorize storage to the virtual server instance: SL02SL111.\nInternal Server Error"))
			})
		})

		Context("Authorize Portable Storage to a VS", func() {
			BeforeEach(func() {
				fakeVSManager.AttachPortableStorageReturns(datatypes.Provisioning_Version1_Transaction{
					Id: sl.Int(1234),
				}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "--portable-id", "1234567")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
			})
		})

		Context("Error Authorize Portable Storage to a VS", func() {
			BeforeEach(func() {
				fakeVSManager.AttachPortableStorageReturns(datatypes.Provisioning_Version1_Transaction{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "--portable-id", "1234567")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to authorize portable storage to the virtual server instance: 1234567.\nInternal Server Error"))
			})
		})
	})
})
