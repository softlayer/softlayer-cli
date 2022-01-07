package licenses_test

import (
	"errors"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/licenses"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
	"strings"
)

var _ = Describe("licenses create-options", func() {
	var (
		fakeUI              *terminal.FakeUI
		fakeLicensesManager *testhelpers.FakeLicensesManager
		cmd                 *licenses.LicensesOptionsCommand
		cliCommand          cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeLicensesManager = new(testhelpers.FakeLicensesManager)
		cmd = licenses.NewLicensesOptionsCommand(fakeUI, fakeLicensesManager)
		cliCommand = cli.Command{
			Name:        metadata.LicensesCreateOptionsMetaData().Name,
			Description: metadata.LicensesCreateOptionsMetaData().Description,
			Usage:       metadata.LicensesCreateOptionsMetaData().Usage,
			Flags:       metadata.LicensesCreateOptionsMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Context("licenses create options returns error", func() {
		BeforeEach(func() {
			fakeLicensesManager.CreateLicensesOptionsReturns(nil, errors.New("Internal server error"))
		})
		It("return error", func() {
			err := testhelpers.RunCommand(cliCommand)
			Expect(err).To(HaveOccurred())
			Expect(strings.Contains(err.Error(), "Failed to licenses create options.\nInternal server error")).To(BeTrue())
			Expect(strings.Contains(err.Error(), "Internal server error")).To(BeTrue())
		})
	})

})
