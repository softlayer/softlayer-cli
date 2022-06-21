package licenses_test

import (
	"errors"
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"

	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/licenses"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Licenses list Create", func() {
	var (
		fakeUI              *terminal.FakeUI
		fakeLicensesManager *testhelpers.FakeLicensesManager
		cmd                 *licenses.CreateCommand
		cliCommand          cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeLicensesManager = new(testhelpers.FakeLicensesManager)
		cmd = licenses.NewCreateCommand(fakeUI, fakeLicensesManager)
		cliCommand = cli.Command{
			Name:        licenses.CreateMetaData().Name,
			Description: licenses.CreateMetaData().Description,
			Usage:       licenses.CreateMetaData().Usage,
			Flags:       licenses.CreateMetaData().Flags,
			Action:      cmd.Run,
		}
		created, _ := time.Parse(time.RFC3339, "2017-11-08T00:00:00Z")
		testPlaceOrder := datatypes.Container_Product_Order_Receipt{
			OrderId:   sl.Int(123456),
			OrderDate: sl.Time(created),
		}
		fakeLicensesManager.CreateLicenseReturns(testPlaceOrder, nil)
	})

	Describe("Licenses create", func() {
		Context("Licenses create, Invalid Usage", func() {
			It("Set command with an invalid output option", func() {
				err := testhelpers.RunCommand(cliCommand, "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
			It("Set command without any datacenter and keyName", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(`This command requires two arguments.`))
			})
		})

		Context("Licenses create, correct use", func() {
			It("return licenses create", func() {
				err := testhelpers.RunCommand(cliCommand, "--datacenter", "dal05", "--key", "XXX_XXX_XXX")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Name      Value"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Id        123456"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Created   2017-11-08T00:00:00Z"))
			})
			It("return licenses create in format json", func() {
				err := testhelpers.RunCommand(cliCommand, "--datacenter", "dal05", "--key", "XXX_XXX_XXX", "--output", "json")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Name": "Id",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Value": "123456"`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Name": "Created",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Value": "2017-11-08T00:00:00Z"`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`[`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`{`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`}`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`]`))
			})
		})
		Context("Licenses errors", func() {
			It("return license error", func() {
				fakeLicensesManager.CreateLicenseReturns(datatypes.Container_Product_Order_Receipt{}, errors.New("Internal server error"))
				err := testhelpers.RunCommand(cliCommand, "--datacenter", "dal05", "--key", "XXX_XXX_XXX")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to create the license."))
				Expect(err.Error()).To(ContainSubstring("Internal server error"))
			})

		})
	})
})
