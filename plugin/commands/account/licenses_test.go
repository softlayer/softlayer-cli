package account_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/account"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Account list Licenses", func() {
	var (
		fakeUI      *terminal.FakeUI
		cliCommand  *account.LicensesCommand
		fakeSession *session.Session
		slCommand   *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = account.NewLicensesCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
	})

	Describe("Account licenses", func() {
		Context("Account licenses, Invalid Usage", func() {
			It("Set command with an invalid output option", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
		})

		Context("Account licenses, correct use", func() {
			It("return account licenses", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Control Panel Licenses"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Id      Ip_address       Manufacturer   Software                                                       Key                  Subnet         Subnet notes"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("1234    11.111.11.11     Plesk          Plesk - Unlimited Domain w/ Power Pack for VPS 17.8.11 Linux   ABCD.00000000.0000   44.444.44.44   test registration"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("12345   222.222.222.22   Plesk          Plesk - 30 Domain w/ Power Pack for VPS 18.x Windows           ABCD.11111111.0000                  -"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("VMware Licenses"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Name                            License key                     CPUs   Description                                     Manufacturer   Required User"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("vCenter                         ABCDE-00000-99999-88888-77777   4      vCenter Server Appliance 6.0                    VMware         administrator@vsphere.local"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Virtual SAN Advanced Tier III   ABCDE-11111-99999-88888-77777   1      VMware vSAN Advanced Tier III 64 - 124 TB 6.x   VMware         -"))
			})
			It("return account licenses in format json", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--output", "json")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Control Panel Licenses":`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Id": "1234","`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Software": "Plesk - Unlimited Domain w/ Power Pack for VPS 17.8.11 Linux","`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Ip_address": "11.111.11.11","`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"VMware Licenses":`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Name": "vCenter","`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"License key": "ABCDE-00000-99999-88888-77777","`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"CPUs": "4","`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Required User": "administrator@vsphere.local""`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`[`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`{`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`}`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`]`))
			})
		})
	})
})
